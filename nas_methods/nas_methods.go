package nas_methods

import (
	"bufio"
	"fmt"
	"io/ioutil"
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

// ----------------------------------------------------------------------------
//
//	GetCardEntry
//
// ----------------------------------------------------------------------------
func GetCardEntries(ctrl *objects.Control, mod *objects.Model) (error, []string) {
	// variables
	var card_name string
	var line int
	var entry int
	// return list
	var entry_list []string
	// var adapted_line int
	// (0) check length of input array
	if len(ctrl.Array01) < 3 {
		return fmt.Errorf("ERROR: input array is expecting 3 entries: %d", len(ctrl.Array01)), entry_list
	}
	// (1) assign card name
	if cn, ok := ctrl.Array01[0].(string); ok {
		card_name = cn
	} else {
		return fmt.Errorf("Array01[0] no string: %T=%v", ctrl.Array01[0], ctrl.Array01[0]), entry_list
	}
	// (2) assign line number
	if l, ok := ctrl.Array01[1].(float64); ok {
		line = int(l)
	} else {
		return fmt.Errorf("Array01[1] no float64: %T=%v", ctrl.Array01[1], ctrl.Array01[1]), entry_list
	}
	// (3) assign entry number
	if e, ok := ctrl.Array01[2].(float64); ok {
		entry = int(e)
	} else {
		return fmt.Errorf("Array01[2] no float: %T=%v", ctrl.Array01[2], ctrl.Array01[2]), entry_list
	}
	//
	pos := (line-1)*10 + entry
	fmt.Println(card_name, line, entry, " - ", pos)

	for _, field := range mod.NasCardList {
		current_name := read.ExtractCardName(field.Fields[0][0])
		if current_name == card_name {
			// fmt.Println(i, field.Fields, current_name)
			one_liner := read.GetOneLiner(field.Fields)
			entry_list = append(entry_list, one_liner[pos-1])

		}
	}

	data_to_file := strings.Join(entry_list, "\n") + "\n"
	ioutil.WriteFile("entry_list.txt", []byte(data_to_file), 0644)

	// two return variables: error value and entry list
	return nil, entry_list
}
