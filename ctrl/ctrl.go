package ctrl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
)

func ReadControlJsonFile(path string, obj *objects.Config, osName string) error {
	// read json control file
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}
	fmt.Print("... scanning control json file  \n")
	// loop through enabled actions
	for actionName, enabled := range obj.Enable {
		if !enabled {
			continue
		}
		fmt.Printf("ACTION: '%s' \n", actionName)
		// create map with key as integer but with flexible values - interface
		actionParams := map[string]interface{}{}
		// definition of input file and input dir
		actionParams["input_file"] = obj.Defaults.InputFile
		actionParams["input_dir"] = obj.Defaults.InputDir
		// further parameters
		actionData, exists := obj.Actions[actionName]
		fmt.Println(actionData, " - ", exists)
		//
		for k, v := range actionData.(map[string]interface{}) {
			actionParams[k] = v
		}
		// full input path
		if obj.Defaults.InputDir != "" && obj.Defaults.InputFile != "" {
			obj.FullInputPath = filepath.Join(obj.Defaults.InputDir, obj.Defaults.InputFile)
		} else if obj.Defaults.InputFile != "" {
			obj.FullInputPath = obj.Defaults.InputFile
		} else if obj.Defaults.InputDir != "" {
			obj.FullInputPath = obj.Defaults.InputDir
		}
		// full output path
		//   output dir
		outputDir := ""
		if val, exists := actionParams["output_dir"]; exists {
			if dirStr, ok := val.(string); ok {
				outputDir = dirStr
			}
		}
		//    output file
		outputFile := ""
		if val, exists := actionParams["output_file"]; exists {
			if fileStr, ok := val.(string); ok {
				outputFile = fileStr
			}
		}
		// build full path from both - dir and file
		if outputDir != "" && outputFile != "" {
			obj.FullOutputPath = filepath.Join(outputDir, outputFile)
		} else if outputFile != "" {
			obj.FullOutputPath = outputFile
		} else if outputDir != "" {
			obj.FullOutputPath = outputDir
		}

	}
	return err
}

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
	// full input path
	if obj.InputDir != "" && obj.InputFile != "" {
		obj.FullInputPath = filepath.Join(obj.InputDir, obj.InputFile)
	} else if obj.InputFile != "" {
		obj.FullInputPath = obj.InputFile
	} else if obj.InputDir != "" {
		obj.FullInputPath = obj.InputDir
	}
	// full output path
	if obj.OutputDir != "" && obj.OutputFile != "" {
		obj.FullOutputPath = filepath.Join(obj.OutputDir, obj.OutputFile)
	} else if obj.OutputFile != "" {
		obj.FullOutputPath = obj.OutputFile
	} else if obj.OutputDir != "" {
		obj.FullOutputPath = obj.OutputDir
	}
	// return value is the error interface value of the scanner
	return scanner.Err()
}

func DebugPrintoutCtrlObj(obj *objects.Control_Object) {
	fmt.Print("debug printout of control object: \n")
	fmt.Print("   Action:       ", obj.Action, "\n")
	fmt.Print("   InputFile:    ", obj.InputFile, "\n")
	fmt.Print("   InputDir:     ", obj.InputDir, "\n")
	fmt.Print("   OutputFile:   ", obj.OutputFile, "\n")
	fmt.Print("   OutputDir:    ", obj.OutputDir, "\n")
	fmt.Print("   Option01:     ", obj.Option01, "\n")
	fmt.Print("      FullInputPath::     ", obj.FullInputPath, "\n")
	fmt.Print("      FullOutpuPath::     ", obj.FullOutputPath, "\n")

}
