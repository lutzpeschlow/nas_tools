package read

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
)

func ReadNasCards(filename string, obj *objects.Model) error {
	fmt.Print("read nastran cards ............................................ \n")
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
	inCard := false
	var first_sign byte

	// extract each card in separate block
	// loop
	for scanner.Scan() {
		// trim line
		line := strings.TrimSpace(scanner.Text())
		// check   begin bulk  and use continue for next step
		if !parsingStarted {
			if strings.EqualFold(line, "BEGIN BULK") {
				parsingStarted = true
			}
			continue
		}

		if parsingStarted {
			fmt.Print(line, "\n")
		}

		// COMMENT - direct jump
		if strings.HasPrefix(line, "$") {
			continue
		}

		// ENDDATA - finisch, last current card should be saved before exit
		if strings.HasPrefix(line, "ENDDATA") {
			fmt.Print("found enddata ... \n")
			obj.NasCards[card_counter] = &objects.NasCard{Card: currentCard}
			continue
		}

		// check first sign of line
		if len(line) == 0 {
			first_sign = '-'
		} else {
			first_sign = line[0]
		}

		// CARD - first is a letter, new card alert, write existing buffer in object
		if (first_sign >= 'a' && first_sign <= 'z') || (first_sign >= 'A' && first_sign <= 'Z') {
			if len(currentCard) > 0 {
				// write existing card into NasCards and clean up currentCard

				obj.NasCards[card_counter] = &objects.NasCard{Card: currentCard}
				card_counter = card_counter + 1
				fmt.Print(" into obj: ", currentCard, "\n")
				// content
				fmt.Print("--- content --- \n")
				for id, card := range obj.NasCards {
					fmt.Printf("ID %d: %+v\n", id, card)
				}

				currentCard = currentCard[:0]
			} else {
				currentCard = append(currentCard, line)
				fmt.Print(" new or append: ", currentCard, " - ", len(currentCard), "\n")
				inCard = true
			}
			// end block - CARD
		} else {
			// CONT - any continuation line should be added to the current card
			currentCard = append(currentCard, line)
		}

	}
	return scanner.Err()
}
