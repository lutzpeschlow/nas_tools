package nas_methods

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
	"github.com/lutzpeschlow/nas_tools/read"
)

// ----------------------------------------------------------------------------
//
//	ExtractCardsAccordingList
//
// ----------------------------------------------------------------------------
func ExtractCardsAccordingList(ctrl *objects.Control, mod *objects.Model) error {
	fmt.Println("extract cards ...", ctrl.Option01, ctrl.Input01, ctrl.OutputFile)
	// (1) read id file
	// final result is a map: idSet containing required ids as string
	FullInput01 := filepath.Join(ctrl.InputDir, ctrl.Input01)
	idSet := make(map[string]bool)
	idFile, err := os.Open(FullInput01)
	if err != nil {
		return fmt.Errorf(FullInput01, err)
	}
	defer idFile.Close()
	scanner := bufio.NewScanner(idFile)
	for scanner.Scan() {
		id := strings.TrimSpace(scanner.Text())
		if id != "" {
			idSet[id] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading ID file %q: %w", ctrl.Input01, err)
	}
	fmt.Println("num of IDs: ", len(idSet))
	// (2) output
	// (2.1) prepare output file
	file, err := os.Create(ctrl.OutputFile)
	if err != nil {
		return fmt.Errorf("could not create output file %q: %w", ctrl.OutputFile, err)
	}
	defer file.Close()
	// (2.2) pre-filter function
	// variable filter which contains a function that will be called later
	// depending on option01, the filter function is adapted and delivers later boolean
	var filter func(cardType string) bool
	// adapt filter function according selected option
	switch ctrl.Option01 {
	case "NOD":
		filter = func(cardType string) bool {
			return strings.HasPrefix(cardType, "GRID")
		}
	case "ELM":
		filter = func(cardType string) bool {
			return strings.HasPrefix(cardType, "C") && !strings.HasPrefix(cardType, "CORD")
		}
	case "MPC":
		filter = func(cardType string) bool {
			return cardType == "RBE" || cardType == "MPC"
		}
	default:
		targetCard := strings.ToUpper(strings.TrimSpace(ctrl.Option01))
		filter = func(cardType string) bool {
			return cardType == targetCard
		}
	}
	// loop over card list
	for _, card := range mod.NasCardList {
		// get card name
		firstLine := card.Card[0]
		cardType := read.ExtractCardName(firstLine)
		// filter function as gate keeper
		if filter(cardType) {
			cardId := read.ExtractCardID(firstLine)

			if idSet[cardId] {

				for _, line := range card.Card {
					_, err := file.WriteString(line + "\n")
					if err != nil {
						return fmt.Errorf("could not write to output file %q: %w", ctrl.OutputFile, err)
					}
				}
			}
		}
	}

	// return value
	return nil
}

func getPosition(row, entry int) int {
	if row == 1 {
		return entry
	}
	return 9 + 8*(row-2) + entry
}

// ----------------------------------------------------------------------------
//
//	GetCardEntry
//
// ----------------------------------------------------------------------------
func GetCardEntry(ctrl *objects.Control, mod *objects.Model) error {
	// variables
	var card_name string
	var line int
	var entry int
	// var adapted_line int
	// (0) check length of input array
	if len(ctrl.Array01) < 3 {
		return fmt.Errorf("ERROR: input array is expecting 3 entries: %d", len(ctrl.Array01))
	}
	// (1) assign card name
	if cn, ok := ctrl.Array01[0].(string); ok {
		card_name = cn
	} else {
		return fmt.Errorf("Array01[0] no string: %T=%v", ctrl.Array01[0], ctrl.Array01[0])
	}
	// (2) assign line number
	if l, ok := ctrl.Array01[1].(float64); ok {
		line = int(l)
	} else {
		return fmt.Errorf("Array01[1] no float64: %T=%v", ctrl.Array01[1], ctrl.Array01[1])
	}
	// (3) assign entry number
	if e, ok := ctrl.Array01[2].(float64); ok {
		entry = int(e)
	} else {
		return fmt.Errorf("Array01[2] no float: %T=%v", ctrl.Array01[2], ctrl.Array01[2])
	}
	//
	// !!!
	// !!!
	one_line := make([]string, 0)
	fmt.Println(card_name, line, entry)
	// loop through all cards to get requested entries
	for _, nas_card := range mod.NasCardList {
		// card name in fields vs. requested card name
		current_name := read.ExtractCardName(nas_card.Fields[0][0])
		if current_name == card_name {
			for i, array := range nas_card.Fields {
				fmt.Println(i)

				array = array[:len(array)-1]

				if i > 0 {
					array = array[1:]
				}

				fmt.Println(">", array, len(array))
				one_line = append(one_line, array...)
			}
			fmt.Println(one_line)
			// fmt.Println(nas_card.Fields, len(nas_card.Fields))

			//
			// 	// small field or large field
			// 	isLargeField := strings.Contains(nas_card.Fields[0][0], "*")
			// 	adapted_line = 0
			// 	// adapt line counter depending on format of previous lines
			// 	var line_add int
			// 	if line >= 2 {
			// 		for i, l := range nas_card.Fields {
			// 			if i+1 < line {
			// 				fmt.Println(i+1, l)
			// 				if strings.ContainsAny(nas_card.Fields[i][0], "*") {
			// 					line_add = line_add + 1
			// 					fmt.Println(nas_card.Fields[i][0], line_add)
			// 				}
			// 			}
			// 		}
			//
			// 		adapted_line = line + line_add
			// 		fmt.Println("> ", line, line_add, adapted_line)
			// 	}
			// 	// (1) LARGE
			// 	if isLargeField {
			// 		// entry > 5 go to next slice, 6-9 to index 0-3
			// 		if entry > 5 {
			// 			fieldIndex := entry - 6
			// 			fmt.Println(nas_card.Fields[adapted_line-1][fieldIndex])
			// 			// regular line
			// 		} else {
			// 			fmt.Println(nas_card.Fields[adapted_line-1][entry-1])
			// 		}
			// 		// (2) SMALL
			// 	} else {
			// 		// Small Field: normaler Zugriff
			// 		fmt.Println(nas_card.Fields[adapted_line-1][entry-1])
			// 	}

		}

		//
	}
	return nil
}
