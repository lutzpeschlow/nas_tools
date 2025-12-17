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

func DebugPrintoutNasCardStats(obj *objects.Model) {
	fmt.Print("debug printout of nas card stats: \n")
	if len(obj.NasCardStats) == 0 {
		fmt.Println("no stats available")
		return
	}
	// get keys from map and sort the keys
	keys := make([]string, 0, len(obj.NasCardStats))
	for typ := range obj.NasCardStats {
		keys = append(keys, typ)
	}
	sort.Strings(keys)
	// output
	maxLen := 10
	total := 0
	for _, typ := range keys {
		count := obj.NasCardStats[typ]
		fmt.Printf("%-*s: %d\n", maxLen, typ, count)
		total += count
	}
	fmt.Print("-----------\n")
	fmt.Printf("%-*s: %d\n", maxLen, "TOTAL", total)
	fmt.Print("===========\n")
}
