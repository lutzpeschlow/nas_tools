package nas_methods

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lutzpeschlow/nas_tools/nas_cards"
	"github.com/lutzpeschlow/nas_tools/objects"
	"github.com/lutzpeschlow/nas_tools/read"
	"github.com/lutzpeschlow/nas_tools/write"
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

// ----------------------------------------------------------------------------
//
//	MpcToCbush
//
// $  1   ||  2   ||  3   ||  4   ||  5   ||  6   ||  7   ||  8   ||  9   ||  10  |
// MPC     9000000050000001       1     -1.  600620       1      1.
// MPC     9000000050000001       2     -1.  600620       2      1.
// MPC     9000000050000001       3     -1.  600620       3      1.
// MPC     9000000050000001       4     -1.  600620       4      1.
// MPC     9000000050000001       5     -1.  600620       5      1.
// MPC     9000000050000001       6     -1.  600620       6      1.
// MPC     9000000050000002       1     -1.  610620       1      1.
// MPC     9000000050000002       2     -1.  610620       2      1.
// MPC     9000000050000002       3     -1.  610620       3      1.
//
// $  1   ||  2   ||  3   ||  4   ||  5   ||  6   ||  7   ||  8   ||  9   ||  10  |
// PBUSH   302     K               10000.
// CBUSH   302     302     501     502     1.0     1.0     0.0
// CBUSH   cbushid pbushid a       b                                0
// ----------------------------------------------------------------------------
func MpcToCbush(ctrl *objects.Control, mod *objects.Model) error {
	// current settins for used IDs
	fmt.Println(ctrl.Action, ctrl.PbushID, ctrl.CbushID)
	// variables
	// node pairing struct for node a and node b
	type Pair struct {
		A string
		B string
	}
	seen := make(map[Pair]int)
	lines := make([]string, 0)
	cbush_id := ctrl.CbushID
	pbush_id := ctrl.PbushID
	// loop through nas cards in memory
	for _, field := range mod.NasCardList {
		// get card name
		current_name := read.ExtractCardName(field.Fields[0][0])
		// found MPC entry
		if current_name == "MPC" {
			// assign to struct
			a := field.Fields[0][2]
			b := field.Fields[0][5]
			p := Pair{A: a, B: b}
			// assign struct to map and count the duplicates
			seen[p]++
			line := fmt.Sprintf("%-8s%-8d%-8d%-8s%-8s%-8s%-8s%-8s%-8s%-8s",
				"CBUSH", cbush_id, pbush_id, a, b, "", "", "", "0", "")
			newCard := nas_cards.CreateCard(line)
			mod.NewCardList = append(mod.NewCardList, newCard)
			lines = append(lines, line)
			cbush_id = cbush_id + 1
		}
	}
	// reporting
	fmt.Println("lines: ", len(lines))
	for p, count := range seen {
		fmt.Println("  ", p.A, p.B, count)
	}
	// write to file
	err := write.WriteNewCards(ctrl.FullOutputPath, mod)
	if err != nil {
		return err
	}
	// return value
	return nil
}

// ----------------------------------------------------------------------------
// 0..100000 - 4164  with last id: 21024
// 100000..200000 - 0
// 200000..300000 - 66  with last id: 200166
// ----------------------------------------------------------------------------
func GetIdRangeTables(ctrl *objects.Control, mod *objects.Model) error {

	fmt.Println(ctrl.OutputDir, ctrl.Option01)

	return nil
}
