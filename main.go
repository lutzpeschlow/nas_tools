package main

// libraries
import (
	"fmt"
	"os"
	"runtime"
	"sort"

	"github.com/lutzpeschlow/nas_tools/ctrl"
	"github.com/lutzpeschlow/nas_tools/objects"
	"github.com/lutzpeschlow/nas_tools/read"
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
		fmt.Print(id, len(card.Card), card.Card, "\n")
		// for index, value := range card.Card {
		// 	fmt.Printf("  [%d] %s\n", index, value)
		// }
	}

}

// ============================================================================
// === main ===
// ============================================================================
func main() {
	ctrl_obj := objects.Control_Object{}
	osName := runtime.GOOS
	err_ctrl := ctrl.ReadControlFile("control.txt", &ctrl_obj, osName)
	if err_ctrl != nil {
		fmt.Printf(" %v\n", err_ctrl)
		os.Exit(1)
	}
	ctrl.DebugPrintoutCtrlObj(&ctrl_obj)

	// model instance
	mod := objects.Model{}
	// create map with key: int and value: *NasCard
	mod.NasCards = make(map[int]*objects.NasCard)
	// get current directory
	current_dir, _ := os.Getwd()
	fmt.Println("current directory:", current_dir)
	// read input file
	dat_file := "./regression_tests/nast_card_test_short.dat"
	err := read.ReadNasCards(dat_file, &mod)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	DebugPrintoutModelObj(&mod)

}
