package modify

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
	"github.com/lutzpeschlow/nas_tools/read"
)

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
