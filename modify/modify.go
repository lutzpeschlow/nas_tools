package modify

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
	"github.com/lutzpeschlow/nas_tools/read"
)

func ExtractCardsAccordingList(ctrl *objects.Control, mod *objects.Model) error {
	fmt.Println("extract cards ...")
	fmt.Print("option: ", ctrl.Option01, " \n")
	fmt.Print("input: ", ctrl.Input01, " \n")
	fmt.Print("output file: ", ctrl.OutputFile, " \n")
	// id file
	idSet := make(map[string]bool)
	if ctrl.Input01 != "" {
		idFile, err := os.Open(ctrl.Input01)
		if err != nil {
			return fmt.Errorf("could not open ID file %q: %w", ctrl.Input01, err)
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
		fmt.Println("num of IDs: ", len(idSet), idSet)
	}
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
	// foundCount := 0
	for _, card := range mod.NasCardList {
		firstLine := card.Card[0]
		cardType := read.ExtractCardName(firstLine)

		if filter(cardType) {
			fields := strings.Fields(firstLine)
			cardID := strings.TrimSpace(fields[1])
			// if idSet[cardID] {
			fmt.Println(cardID)
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
