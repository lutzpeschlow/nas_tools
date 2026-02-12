package ctrl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lutzpeschlow/nas_tools/objects"
)

func ReadControlJsonFile(path string, obj *objects.Control, osName string) error {
	// (1) read json control file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read control file: %w", err)
	}
	// (2) json parsen
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("parse config json %s: %w", path, err)
	}
	// (3) manage and process actions
	//    store data from json file in several maps and structs
	//    enable, defaults, actions
	//    loop through enabled actions
	for actionName, enabled := range obj.Enable {
		if !enabled {
			continue
		}
		// does exist action in object
		actionData, exists := obj.Actions[actionName]
		if !exists {
			return fmt.Errorf("action %s not found in config", actionName)
		}
		// create map with key as integer but with flexible values - interface
		actionParams := map[string]interface{}{}
		// definition of input file and input dir
		actionParams["input_file"] = obj.Defaults.InputFile
		actionParams["input_dir"] = obj.Defaults.InputDir
		// further parameters
		// actionData, _ := obj.Actions[actionName]
		// create map of map with data depending on action
		for k, v := range actionData.(map[string]interface{}) {
			actionParams[k] = v
		}
		// (2)
		// extract data to single values for better usage
		obj.Action = actionName
		obj.InputFile = actionParams["input_file"].(string)
		obj.InputDir = actionParams["input_dir"].(string)
		if val, exists := actionParams["output_file"]; exists {
			obj.OutputFile = val.(string)
		}
		if val, exists := actionParams["output_dir"]; exists {
			obj.OutputDir = val.(string)
		}
		if val, exists := actionParams["option_01"]; exists {
			obj.Option01 = val.(string)
		}
		if val, exists := actionParams["input_01"]; exists {
			obj.Input01 = val.(string)
		}
		// combined from previous
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
	}
	return nil
}

func DebugPrintoutCtrlObj(obj *objects.Control) {
	fmt.Print("debug printout of control object: \n")
	fmt.Print("   Action:       ", obj.Action, "\n")
	fmt.Print("   InputFile:    ", obj.InputFile, "\n")
	fmt.Print("   InputDir:     ", obj.InputDir, "\n")
	fmt.Print("   OutputFile:   ", obj.OutputFile, "\n")
	fmt.Print("   OutputDir:    ", obj.OutputDir, "\n")
	fmt.Print("   Option01:     ", obj.Option01, "\n")
	fmt.Print("   Input01:     ", obj.Input01, "\n")
	fmt.Print("      FullInputPath:     ", obj.FullInputPath, "\n")
	fmt.Print("      FullOutpuPath:     ", obj.FullOutputPath, "\n")

}
