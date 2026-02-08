package write

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
	"github.com/lutzpeschlow/nas_tools/read"
)

func WriteNasCards(ctrl *objects.Control, obj *objects.Model) error {
	//
	filename := ctrl.FullOutputPath
	fmt.Print("write nas cards into file: ", filename, "\n")
	// assign file and get file object
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("... write file problems ")
		return fmt.Errorf("ERROR: %w", err)
	}
	defer f.Close()
	//
	fmt.Println("   write: ", filename)
	// sort keys of map
	keys := make([]int, 0, len(obj.NasCards))
	for id := range obj.NasCards {
		keys = append(keys, id)
	}
	sort.Ints(keys)
	// write loop
	for _, id := range keys {
		// get card object
		card := obj.NasCards[id]
		// loop though string slice in object
		for _, value := range card.Card {
			// write line
			if _, err := fmt.Fprintln(f, value); err != nil {
				return fmt.Errorf("write error: %w", err)
			}
		}
	}
	return nil
}

// --------------------------------------------------------------------------------------

func WriteCardsToFiles(dir string, obj *objects.Model) error {

	// directory handling
	if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("ERROR %s DEL: %w", dir, err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("ERROR %s CREATE: %w", dir, err)
	}

	// map for card types with loop through all cards
	cardTypes := make(map[string][][]string)

	for _, card := range obj.NasCardList {
		firstLine := card.Card[0]
		cardType := read.ExtractCardName(firstLine)
		cardTypes[cardType] = append(cardTypes[cardType], card.Card)
	}

	// create own file
	for cardType, cards := range cardTypes {
		//
		filename := strings.ToLower(cardType) + ".txt"
		filepath := filepath.Join(dir, filename)
		//
		var content strings.Builder
		for _, card := range cards {
			for _, line := range card {
				content.WriteString(line + "\n")
			}
			// content.WriteString("\n") // Leere Zeile zwischen Karten
		}
		//
		if err := os.WriteFile(filepath, []byte(content.String()), 0644); err != nil {
			return fmt.Errorf("ERROR: FILE %s WRITE: %w", filepath, err)
		}
		fmt.Printf("created: %s (%d cards)\n", filepath, len(cards))
	}

	return nil
}
