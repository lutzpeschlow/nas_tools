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
	// card_counter := 0
	parsingStarted := false
	inCard := false
	var first_sign byte

	// extract each card in separate block
	// loop
	for scanner.Scan() {
		// trim line
		// line := strings.TrimSpace(scanner.Text())
		line := scanner.Text()
		// check   begin bulk  and use continue for next step
		if !parsingStarted {
			if strings.EqualFold(line, "BEGIN BULK") {
				parsingStarted = true
			}
			continue
		}

		// COMMENT - direct jump
		if strings.HasPrefix(line, "$") {
			continue
		}

		// ENDDATA - finisch, last current card should be saved before exit
		if strings.HasPrefix(line, "ENDDATA") {
			fmt.Print("found enddata, save last card ... \n")
			// write existing card into NasCards
			nextID := len(obj.NasCards)
			obj.NasCards[nextID] = &objects.NasCard{Card: append([]string(nil), currentCard...)}
			continue
		}

		// check first sign of line
		if len(line) == 0 {
			first_sign = '-'
		} else {
			first_sign = line[0]
		}

		// CARD - first is a letter, new card alert
		// write existing buffer in object and setup for new card
		if (first_sign >= 'a' && first_sign <= 'z') || (first_sign >= 'A' && first_sign <= 'Z') {
			// there is content in current card that needs to be saved
			if len(currentCard) > 0 {
				// write existing card into NasCards
				// create new string slice, initialize nil, add content
				nextID := len(obj.NasCards)
				obj.NasCards[nextID] = &objects.NasCard{Card: append([]string(nil), currentCard...)}
				// clean up current card and assign new data
				currentCard = currentCard[:0]
				currentCard = append(currentCard, line)
				inCard = true
				// no content in current card, simply add line
			} else {
				currentCard = append(currentCard, line)
				inCard = true
			}
			// end block - CARD
		} else {
			if inCard {
				// CONT - any continuation line should be added to the current card
				currentCard = append(currentCard, line)
			}
		}
	}
	return scanner.Err()
}

// content
//fmt.Print("    --- content --- \n")
//for id, card := range obj.NasCards {
//	fmt.Printf("   ID %d: %+v\n", id, card)
//}
