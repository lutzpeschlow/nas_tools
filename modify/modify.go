package modify

import (
	"fmt"
	"os"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
	"github.com/lutzpeschlow/nas_tools/read"
)

func ExtractCardsAccordingList(ctrl *objects.Control_Object, mod *objects.Model) error {
	fmt.Println("extract cards ...")
	fmt.Print("option: ", ctrl.Option01, " \n")
	fmt.Print("output file: ", ctrl.OutputFile, " \n")
	// output file
	file, err := os.Create(ctrl.OutputFile)
	if err != nil {
		return fmt.Errorf("could not create output file %q: %w", ctrl.OutputFile, err)
	}
	defer file.Close()

	// pre-filter function
	var filter func(cardType string) bool

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
		return fmt.Errorf("unknown Option01: %q", ctrl.Option01)
	}

	// loop over card list
	// filter is now used instead of ctrl.option01
	for _, card := range mod.NasCardList {
		firstLine := card.Card[0]
		cardType := read.ExtractCardName(firstLine)

		if filter(cardType) {
			for _, line := range card.Card {
				_, err := file.WriteString(line + "\n")
				if err != nil {
					return fmt.Errorf("could not write to output file %q: %w", ctrl.OutputFile, err)
				}
			}
		}
	}

	// return value
	return nil
}
