package debug

import (
	"fmt"
	"sort"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
)

// ----------------------------------------------------------------------------
//
//	DebugPrintoutModelObj
//
// ----------------------------------------------------------------------------
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

// ----------------------------------------------------------------------------
//
//	DebugPrintoutNasCardStats
//
// ----------------------------------------------------------------------------
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

// ----------------------------------------------------------------------------
//
//	DeburgPrintoutEntries
//
// ----------------------------------------------------------------------------
func DeburgPrintoutEntries(row []string) int {
	isLargeField := false
	if len(row) > 0 && strings.Contains(row[0], "*") {
		isLargeField = true
	}
	// (1)
	// small field output
	if !isLargeField {
		// optic for small field: +--------+--------+--------+
		fmt.Print("+")
		for range row {
			fmt.Print(strings.Repeat("-", 7) + "+")
		}
		fmt.Println()
		// data
		for _, field := range row {
			fmt.Printf("%-8s", field)
		}
		fmt.Println("")
		// (2)
		// large field output
	} else {
		// optic for large field: +-------+---------------+---------------+
		fmt.Print("+" + strings.Repeat("-", 7) + "+")
		for i := 0; i < 4; i++ {
			fmt.Print(strings.Repeat("-", 15) + "+")
		}
		fmt.Print(strings.Repeat("-", 7) + "+")
		fmt.Println()
		// data
		fmt.Println(row, len(row))
		fmt.Printf("%-8s", row[0])
		for i := 1; i < 5 && i < len(row); i++ {
			fmt.Printf("%-16s", row[i])
		}
		if len(row) >= 5 {
			fmt.Printf("%-8s", row[5])
			fmt.Printf("\n")
		} else {
			fmt.Println("")
		}

	}

	return 0
}
