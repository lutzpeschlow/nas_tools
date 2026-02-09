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

	// switch for actions
	switch ctrl.Action {
	case "READ":
		write.WriteNasCards(ctrl, mod)
	case "STATS":
		read.GetNasCardsStatistics(mod)
		debug.DebugPrintoutNasCardStats(mod)
	case "SPLIT":
		write.WriteCardsToFiles(ctrl.OutputDir, mod)
	case "EXTRACT_ACC_LIST":
		modify.ExtractCardsAccordingList(ctrl, mod)
	//
	default:
		fmt.Println("unknown action: %s", ctrl.Action)
	}

	// return variable
	return nil
}
