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
	start := time.Now()
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("elapsed:", time.Since(start))
}

// ----------------------------------------------------------------------------
// function run()
//
//	containing main logic
//	clean specification between CLI entrpoint and business logic
//
// in run following actions:
//
//	first information as OS and working directory
//	control object
func run() error {
	// get os and current directory
	osName := runtime.GOOS
	current_dir, _ := os.Getwd()
	fmt.Println(osName, current_dir)
	// ctrl_obj
	ctrl_obj := objects.Control{}
	if err := ctrl.ReadControlJsonFile("control.json", &ctrl_obj, osName); err != nil {
		return fmt.Errorf("failed to read control.json: %w", err)
	}
	ctrl.DebugPrintoutCtrlObj(&ctrl_obj)
	// model instance
	mod := objects.Model{}
	// read input file
	dat_file := ctrl_obj.FullInputPath
	_, _, err := read.ReadNasFile(dat_file, &mod)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", dat_file, err)
	}
	// execute main action according enabled value in control object
	if err := cmd.ExecuteAction(&ctrl_obj, &mod); err != nil {
		return fmt.Errorf("execute action failed: %w", err)
	}
	// explicit nil return value
	return nil
}
