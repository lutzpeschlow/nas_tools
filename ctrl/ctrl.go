package ctrl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"path/filepath"

	"github.com/lutzpeschlow/nas_tools/objects"
)

// ReadControlFile function to read a control file
//
// input:
//
// output:
//   - error: if read or parse fails, put back error, else nil
func ReadControlFile(path string, obj *objects.Control_Object, osName string) error {
	// defaults
	obj.Action = "READ"
	// pointer to file for later opening, err as interface value
	file, err := os.Open(path)
	// if we have an error go out with returning err value
	if err != nil {
		return err
	}
	// close file at the end of the function
	defer file.Close()
	// read content from file object and scan
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// trim and split the line, parts as slice of strings (dynamical array)
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			switch parts[0] {
			case "ACTION":
				obj.Action = parts[1]
			case "INPUT_FILE":
				obj.InputFile = parts[1]
			case "INPUT_DIR":
				obj.InputDir = parts[1]
			case "OUTPUT_FILE":
				obj.OutputFile = parts[1]
			case "OUTPUT_DIR":
				obj.OutputDir = parts[1]
			case "OPTION_01":
				obj.Option01 = parts[1]
			}
		}
	}
	//
	if obj.InputDir != "" && obj.InputFile != "" {
        obj.FullInputPath = filepath.Join(obj.InputDir, obj.InputFile)
    } else if obj.InputFile != "" {
        obj.FullInputPath = obj.InputFile  
    } else if obj.InputDir != "" {
        obj.FullInputPath = obj.InputDir   
    }
	// return value is the error interface value of the scanner
	return scanner.Err()
}

func DebugPrintoutCtrlObj(obj *objects.Control_Object) {
	fmt.Print("debug printout of control object: \n")
	fmt.Print("  Action:       ", obj.Action, "\n")
	fmt.Print("  InputFile:    ", obj.InputFile, "\n")
	fmt.Print("  InputDir:     ", obj.InputDir, "\n")
	fmt.Print("  OutputFile:   ", obj.OutputFile, "\n")
	fmt.Print("  OutputDir:    ", obj.OutputDir, "\n")
	fmt.Print("  Option01:     ", obj.Option01, "\n")
	fmt.Print("     FullInputPath::     ", obj.FullInputPath, "\n")

}
