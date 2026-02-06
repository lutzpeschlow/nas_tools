package main

// libraries
import (
	"fmt"
	"os"
	"runtime"

	"github.com/lutzpeschlow/nas_tools/cmd"
	"github.com/lutzpeschlow/nas_tools/ctrl"
	"github.com/lutzpeschlow/nas_tools/objects"
	"github.com/lutzpeschlow/nas_tools/read"
)

// ============================================================================
// === main ===
// ============================================================================
func main() {
	// ctrl_obj := objects.Control_Object{}
	config_obj := objects.Config{}
	osName := runtime.GOOS
	err_ctrl := ctrl.ReadControlJsonFile("control.json", &config_obj, osName)
	if err_ctrl != nil {
		fmt.Printf(" %v\n", err_ctrl)
		os.Exit(1)
	}
	// ctrl.DebugPrintoutCtrlObj(&ctrl_obj)

	// model instance
	mod := objects.Model{}
	// create map with key: int and value: *NasCard
	mod.NasCards = make(map[int]*objects.NasCard)
	// get current directory
	current_dir, _ := os.Getwd()
	fmt.Println("current directory:", current_dir)
	//
	// read input file
	dat_file := config_obj.FullInputPath

	err := read.ReadNasCards(dat_file, &mod)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	//
	if err := cmd.ExecuteAction(&config_obj, &mod); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

}
