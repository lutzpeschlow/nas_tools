package cmd

import (
	"fmt"

	"github.com/lutzpeschlow/nas_tools/debug"
	"github.com/lutzpeschlow/nas_tools/modify"
	"github.com/lutzpeschlow/nas_tools/objects"
	"github.com/lutzpeschlow/nas_tools/read"
	"github.com/lutzpeschlow/nas_tools/write"
)

func ExecuteAction(ctrl *objects.Control, mod *objects.Model) error {
	// // action := config.Enable
	// fmt.Println(" - ", config.Enable, config.Defaults, config.Action)
	//
	// for actionName, enabled := range config.Enable {
	// 	if !enabled {
	// 		continue
	// 	}
	// 	switch actionName {
	// 	case "READ":
	// 		return write.WriteNasCards(config, mod)
	// 	case "STATS":
	// 		fmt.Println("stats")
	// 	case "SPLIT":
	// 		fmt.Println("split")
	// 	case "EXTRACT_ACC_LIST":
	// 		fmt.Println("extract")
	// 	default:
	// 		fmt.Printf("  unknown: %s\n", actionName)
	// 	}
	// }

	switch ctrl.Action {
	case "READ":
		return write.WriteNasCards(ctrl, mod)
	case "STATS":
		read.GetNasCardsStatistics(mod)
		debug.DebugPrintoutNasCardStats(mod)
		return nil
	case "SPLIT":
		return write.WriteCardsToFiles(ctrl.OutputDir, mod)
	case "EXTRACT_ACC_LIST":
		return modify.ExtractCardsAccordingList(ctrl, mod)
	//
	default:
		return fmt.Errorf("unknown action: %s", ctrl.Action)
	}
	return nil
}
