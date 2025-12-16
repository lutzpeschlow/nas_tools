package write

import (
	"fmt"
	"os"
	"sort"

	"github.com/lutzpeschlow/nas_tools/objects"
)

func WriteNasCards(filename string, obj *objects.Model) error {
	// assign file and get file object
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ERROR: %w", err)
	}
	defer f.Close()
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
