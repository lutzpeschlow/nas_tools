package read

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
)

func ReadNasCards(filename string, obj *objects.Model) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if obj.NasCards == nil {
		obj.NasCards = make(map[int]*objects.NasCard)
	}

	scanner := bufio.NewScanner(file)
	var currentCard strings.Builder
	card_counter := 0
	parsingStarted := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if !parsingStarted {
			if strings.EqualFold(line, "BEGIN BULK") {
				parsingStarted = true
			}
			continue
		}

		// Kommentar-Zeilen überspringen ($)
		if strings.HasPrefix(line, "$") {
			continue
		}

		// Leere Zeilen = Kartentrenner
		if line == "" {
			if currentCard.Len() > 0 {
				// fertige Karte in Slice zerlegen und ins Model schreiben
				cardStr := strings.TrimSpace(currentCard.String())
				fields := strings.Fields(cardStr)
				fmt.Print(fields, "\n")

				obj.NasCards[card_counter] = &objects.NasCard{Card: fields}
				card_counter++

				currentCard.Reset()
			}
			continue
		}

		// Karte anhängen (free format: + für continuation)
		currentCard.WriteString(line + " ")

		// Optional: Prüfe auf Continuation-Marker
		if strings.HasSuffix(line, "+") {
			continue // Nächste Zeile anhängen
		}

		// Karte fertig → speichern
		cardStr := strings.TrimSpace(currentCard.String())
		fields := strings.Fields(cardStr)
		fmt.Print(fields, "\n")

		// obj.NasCards[card_counter] = &objects.NasCard{
		// 	card: fields,
		// }
		card_counter++

		currentCard.Reset()
	}

	// Letzte Karte prüfen
	if currentCard.Len() > 0 {
		cardStr := strings.TrimSpace(currentCard.String())
		fields := strings.Fields(cardStr)
		fmt.Print(fields, "\n")
		// obj.NasCards[card_counter] = &objects.NasCard{
		// 	card: fields,
		// }
	}

	return scanner.Err()

}

func ReadDat(filename string, obj *objects.Model) error {
	// open file, with defer as backup
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// scan file and assign data
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "GRID") {
			node, err := parseGRID(line)
			if err == nil {
				obj.Nodes[node.ID] = node
			}
		}
	}
	return scanner.Err()
}

func parseGRID(line string) (*objects.Node, error) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return nil, fmt.Errorf("invalid GRID line: %s", line)
	}

	node := &objects.Node{}

	fmt.Print(fields, "\n")

	// ID (Feld 1)
	node.ID, _ = strconv.Atoi(fields[1])

	// CP (optional, Feld 2 oder default 0)
	if len(fields) > 2 {
		cp, err := strconv.Atoi(fields[2])
		if err == nil {
			node.CP = cp
		}
	}

	// Koordinaten X,Y,Z (letzte 3 Felder)
	if len(fields) >= 7 {
		// Vollständig: ID,CP,X,Y,Z,CD,PS
		x, _ := strconv.ParseFloat(fields[3], 64)
		y, _ := strconv.ParseFloat(fields[4], 64)
		z, _ := strconv.ParseFloat(fields[5], 64)
		cd, _ := strconv.Atoi(fields[6])
		ps, _ := strconv.Atoi(fields[7])

		node.X, node.Y, node.Z = x, y, z
		node.CD, node.PS = cd, ps
	} else {
		// Dein File-Format: GRID ID X Y Z
		x, _ := strconv.ParseFloat(fields[2], 64)
		y, _ := strconv.ParseFloat(fields[3], 64)
		z, _ := strconv.ParseFloat(fields[4], 64)

		node.X, node.Y, node.Z = x, y, z
	}

	return node, nil
}
