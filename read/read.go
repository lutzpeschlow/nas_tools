package read

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
)

// function: ReadNasCards
//
// description:
//
// input :  file name,  model object
// output : error/return value
func ReadNasFile(filename string, obj *objects.Model) error {
	fmt.Println("... read nastran cards: ", filename)
	// get file object and close with defer
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("... problem reading file: ", filename)
		return err
	}
	defer file.Close()
	// scan object, str.Builder object and further variables
	scanner := bufio.NewScanner(file)
	// create map with key: int and value: *NasCard
	obj.NasCards = make(map[int]*objects.NasCard)
	var currentCard []string
	inCard := false
	var first_sign byte
	lineCount := 0
	inBulk := false
	// extract each card in separate block
	// loop
	for scanner.Scan() {
		lineCount = lineCount + 1
		line := scanner.Text()

		// are we in BEGIN BULK, set inBulk true
		if strings.Contains(line, "BEGIN BULK") {
			inBulk = true
			fmt.Println("found BEGIN BULK ...")
			continue
		}
		// in case of ENDDATA, break loop routine
		if strings.Contains(line, "ENDDATA") {
			inBulk = false
			fmt.Println("found ENDDATA ...")
			break
		}
		// inside BULK block, we are allowed to parse
		if !inBulk {
			continue
		}
		// (1) COMMENT - direct jump
		if strings.HasPrefix(line, "$") {
			continue
		}
		// check first sign of line
		if len(line) == 0 {
			first_sign = '-'
		} else {
			first_sign = line[0]
		}
		// (2) CARD - first is a letter
		// new card alert, write existing buffer in object and setup for new card
		if (first_sign >= 'a' && first_sign <= 'z') || (first_sign >= 'A' && first_sign <= 'Z') {
			// (2.1) there is content in current card that needs to be saved in object
			if len(currentCard) > 0 {
				// write existing card into NasCards and NasCardList
				// create new string slice, initialize nil, add content
				nextID := len(obj.NasCards)
				newCard := &objects.NasCard{Card: append([]string(nil), currentCard...)}
				obj.NasCards[nextID] = newCard
				obj.NasCardList = append(obj.NasCardList, newCard)
				// clean up current card and assign new data
				currentCard = currentCard[:0]
				currentCard = append(currentCard, line)
				inCard = true
				// (2.2) simply add line to current card, no action with previous card
			} else {
				currentCard = append(currentCard, line)
				inCard = true
			}
			// (3) CONT - in card any continuation line
		} else {
			if inCard {
				currentCard = append(currentCard, line)
			}
		}
	}
	// cover last line in buffer ...
	if len(currentCard) > 0 {
		nextID := len(obj.NasCards)
		newCard := &objects.NasCard{Card: append([]string(nil), currentCard...)}
		obj.NasCards[nextID] = newCard
		obj.NasCardList = append(obj.NasCardList, newCard)
	}
	// read file stats
	fmt.Println("lines/cards: ", lineCount, len(obj.NasCards), len(obj.NasCardList))
	// return scanner error
	return scanner.Err()
}

// function: GetNasCardsStatistics
//
// description:
//
// input : model object
// output : error/return value
func GetNasCardsStatistics(obj *objects.Model) error {
	//
	obj.NasCardStats = make(map[string]int)
	// loop through all cards
	for _, card := range obj.NasCards {
		if len(card.Card) == 0 {
			continue
		}
		// first line of card
		firstLine := card.Card[0]
		cardType := ExtractCardName(firstLine)
		// add counter to according card in the statistics object
		if cardType != "" {
			obj.NasCardStats[cardType]++
		}
	}
	return nil
}

// function: extractCardName
//
// description:
//
// input : string
// output : string
func ExtractCardName(line string) string {
	if len(line) < 4 {
		return ""
	}
	name := strings.TrimSpace(line[:8])
	sepIndex := strings.IndexAny(name, ",+ *")

	if sepIndex > 0 {
		name = name[:sepIndex] // Bis zum Trennzeichen abschneiden
	}

	return strings.ToUpper(strings.TrimSpace(name))

}

// content
//fmt.Print("    --- content --- \n")
//for id, card := range obj.NasCards {
//	fmt.Printf("   ID %d: %+v\n", id, card)
//}

// func extractCardName(line string) string {
//     // Nastran: Erste 8 Zeichen als Kartenname
//     if len(line) < 8 {
//         return ""
//     }
//     name := strings.TrimSpace(line[:8])
//
//     // Entferne + * , am Ende
//     name = strings.TrimRight(name, "+*, ")
//
//     // Zu Großbuchstaben
//     return strings.ToUpper(name)
// }
// func (m *Model) FillStats() {
//     m.NasCardStats = make(map[string]int)
//
//     for _, card := range m.NasCards {
//         for _, line := range card.Card {
//             if len(line) == 0 { continue }
//
//             cardType := extractCardName(line)
//             if cardType != "" {
//                 m.NasCardStats[cardType]++
//             }
//         }
//     }
// }
