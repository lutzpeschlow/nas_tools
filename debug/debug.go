package debug

import (
	"fmt"
	"sort"

	"github.com/lutzpeschlow/nas_tools/objects"
)

func DebugPrintoutModelObj(obj *objects.Model) {

	fmt.Print("debug printout of control object: \n")

	keys := make([]int, 0, len(obj.NasCards))
	for id := range obj.NasCards {
		keys = append(keys, id)
	}
	sort.Ints(keys)

	for _, id := range keys {
		card := obj.NasCards[id]

		fmt.Print(id, len(card.Card), "\n")

		for index, value := range card.Card {
			fmt.Printf("  [%d] %s\n", index, value)
		}
	}

}
