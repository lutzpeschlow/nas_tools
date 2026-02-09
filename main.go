package main

// libraries
import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/lutzpeschlow/nas_tools/cmd"
	"github.com/lutzpeschlow/nas_tools/ctrl"
	"github.com/lutzpeschlow/nas_tools/objects"
	"github.com/lutzpeschlow/nas_tools/read"
)

// ============================================================================
// === main ===
// ============================================================================
func main() {
	// get os and current directory
	osName := runtime.GOOS
	current_dir, _ := os.Getwd()
	fmt.Println("current directory:", current_dir)
	// timer
	start := time.Now()
	fmt.Println(" start:", start.Format(time.RFC3339Nano))
	// ctrl_obj
	ctrl_obj := objects.Control{}
	err_ctrl := ctrl.ReadControlJsonFile("control.json", &ctrl_obj, osName)
	if err_ctrl != nil {
		fmt.Printf(" %v\n", err_ctrl)
		os.Exit(1)
	}
	ctrl.DebugPrintoutCtrlObj(&ctrl_obj)

	// model instance
	mod := objects.Model{}

	//
	// read input file
	dat_file := ctrl_obj.FullInputPath
	err := read.ReadNasFile(dat_file, &mod)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	// execute actaion according enabled value in config object
	if err := cmd.ExecuteAction(&ctrl_obj, &mod); err != nil {
		fmt.Printf("... %v\n", err)
		os.Exit(1)
	}
	// timer
	elapsed := time.Since(start)
	fmt.Println(" elapsed:", elapsed)

}
