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
//	DebugPrintoutNasFieldEntries
//
// ----------------------------------------------------------------------------
func DebugPrintoutNasFieldEntries(obj *objects.Model) {
	for _, c := range obj.NasCardList {
		DebugPrintoutEntryLines(c.Fields)
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
//	DebugPrintoutEntries
//
// ----------------------------------------------------------------------------
func DebugPrintoutEntries(row []string) int {
	var debug_printout []string
	debug_printout = GetPrintoutEntries(row)
	fmt.Println(strings.Join(debug_printout, ""))
	return 0
}

// ----------------------------------------------------------------------------
//
//	DebugPrintoutEntryLines
//
// ----------------------------------------------------------------------------
func DebugPrintoutEntryLines(entry_lines [][]string) int {
	var debug_printout []string
	// fmt.Println(rows, len(entry_lines))
	fmt.Println(get_small_field_optic())
	for _, line := range entry_lines {
		debug_printout = GetPrintoutEntries(line)
		fmt.Println(strings.Join(debug_printout, ""))
	}
	return 0
}

// ----------------------------------------------------------------------------
//
//	GetPrintoutEntries
//
// ----------------------------------------------------------------------------
func GetPrintoutEntries(row []string) []string {
	var debug_printout []string
	isLargeField := false
	if len(row) > 0 && strings.Contains(row[0], "*") {
		isLargeField = true
	}
	// (1)
	// small field output
	if !isLargeField {
		for _, field := range row {
			debug_printout = append(debug_printout, fmt.Sprintf("%-8s", field))
		}
		// (2)
		// large field output
	} else {
		debug_printout = append(debug_printout, fmt.Sprintf("%-8s", row[0]))
		for i := 1; i < 5 && i < len(row); i++ {
			debug_printout = append(debug_printout, fmt.Sprintf("%-16s", row[i]))
		}
		if len(row) >= 5 {
			debug_printout = append(debug_printout, fmt.Sprintf("%-8s", row[5]))
		} else {
			debug_printout = append(debug_printout, "")
		}
	}
	// return value as slice of strings
	return debug_printout
}

// ----------------------------------------------------------------------------
//
//	get_small_field_optic
//
// ----------------------------------------------------------------------------
func get_small_field_optic() string {
	return "+-------+-------+-------+-------+-------+-------+-------+-------+-------+-------+"
}

// ----------------------------------------------------------------------------
//
//	get_large_field_optic
//
// ----------------------------------------------------------------------------
func get_large_field_optic() string {
	return "+-------+---------------+---------------+---------------+---------------+-------+"
}
