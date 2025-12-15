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
	// scan object, str.Builder object and further variables
	scanner := bufio.NewScanner(file)
	var currentCard []string
	card_counter := 0
	parsingStarted := false
	var first_sign byte

	// extract each card in separate block
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

		if len(line) == 0 {
			first_sign = '-'
		} else {
			first_sign = line[0]
		}

		// first is a letter, new card alert, write existing buffer in object
		if (first_sign >= 'a' && first_sign <= 'z') || (first_sign >= 'A' && first_sign <= 'Z') {
			if len(currentCard) > 0 {
				// write existing card into NasCards and clean up currentCard
				fmt.Print("write card into objects ...", card_counter, " - ", len(currentCard), "\n")
				obj.NasCards[card_counter] = &objects.NasCard{Card: currentCard}
				card_counter = card_counter + 1

				currentCard = currentCard[:0]
			} else {
				fmt.Print("new or append ... \n")
				currentCard = append(currentCard, line)
				fmt.Print(currentCard, "\n")
			}

		}

		// Leere Zeilen = Kartentrenner
		// if line == "" {
		// 	if currentCard.Len() > 0 {

		//		cardStr := strings.TrimSpace(currentCard.String())
		//		fields := strings.Fields(cardStr)
		//		obj.NasCards[card_counter] = &objects.NasCard{Card: fields}
		//		card_counter++
		//		currentCard.Reset()
		//	}
		//	continue
		//}

		//currentCard.WriteString(line + " ")
		// Optional: Prüfe auf Continuation-Marker
		// if strings.HasSuffix(line, "+") {
		// 	continue // Nächste Zeile anhängen
		//}
		// Karte fertig → speichern
		// cardStr := strings.TrimSpace(currentCard.String())
		// fields := strings.Fields(cardStr)
		// obj.NasCards[card_counter] = &objects.NasCard{Card: fields}
		// card_counter++
		// currentCard.Reset()
	}
	// Letzte Karte prüfen
	// if currentCard.Len() > 0 {
	// 	cardStr := strings.TrimSpace(currentCard.String())
	// 	fields := strings.Fields(cardStr)
	// 	obj.NasCards[card_counter] = &objects.NasCard{Card: fields}
	// }
	return scanner.Err()
}
