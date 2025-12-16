package main

// libraries
import (
	"fmt"
	"os"
	"runtime"

	"github.com/lutzpeschlow/nas_tools/ctrl"
	"github.com/lutzpeschlow/nas_tools/debug"
	"github.com/lutzpeschlow/nas_tools/objects"
	"github.com/lutzpeschlow/nas_tools/read"
	"github.com/lutzpeschlow/nas_tools/write"
)

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

	// debug printout
	debug.DebugPrintoutModelObj(&mod)

	// write content file
	write.WriteNasCards("result.txt", &mod)

}

// func extractCardName(line string) string {
//     // Nastran: Erste 8 Zeichen als Kartenname
//     if len(line) < 8 {
//         return ""
//     }
//     name := strings.TrimSpace(line[:8])
//
//     // Entferne + * , am Ende
//     name = strings.TrimRight(name, "+*, ")
//
//     // Zu Großbuchstaben
//     return strings.ToUpper(name)
// }
// func (m *Model) FillStats() {
//     m.NasCardStats = make(map[string]int)
//
//     for _, card := range m.NasCards {
//         for _, line := range card.Card {
//             if len(line) == 0 { continue }
//
//             cardType := extractCardName(line)
//             if cardType != "" {
//                 m.NasCardStats[cardType]++
//             }
//         }
//     }
// }
