package read

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
)

func ReadNasCards(filename string, obj *objects.Model) error {
	fmt.Print("read nastran cards ... \n")
	// get file object and close with defer
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// if obj.NasCards == nil {
	// 	obj.NasCards = make(map[int]*objects.NasCard)
	// }

	// scan object, str.Builder object and further variables
	scanner := bufio.NewScanner(file)
	var currentCard strings.Builder
	card_counter := 0
	parsingStarted := false
	// loop
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// check   begin bulk  and use continue for next step
		if !parsingStarted {
			if strings.EqualFold(line, "BEGIN BULK") {
				parsingStarted = true
			}
			continue
		}
		// comment line - jump
		if strings.HasPrefix(line, "$") {
			continue
		}

		// Leere Zeilen = Kartentrenner
		if line == "" {
			if currentCard.Len() > 0 {
				// fertige Karte in Slice zerlegen und ins Model schreiben
				cardStr := strings.TrimSpace(currentCard.String())
				fields := strings.Fields(cardStr)
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
		obj.NasCards[card_counter] = &objects.NasCard{Card: fields}
		card_counter++
		currentCard.Reset()
	}
	// Letzte Karte prüfen
	if currentCard.Len() > 0 {
		cardStr := strings.TrimSpace(currentCard.String())
		fields := strings.Fields(cardStr)
		obj.NasCards[card_counter] = &objects.NasCard{Card: fields}
	}
	return scanner.Err()
}
