package cmd

import (
	"fmt"

	"github.com/lutzpeschlow/nas_tools/debug"
	"github.com/lutzpeschlow/nas_tools/objects"
	"github.com/lutzpeschlow/nas_tools/read"
	"github.com/lutzpeschlow/nas_tools/write"
)

func ExecuteAction(ctrl_obj *objects.Control_Object, mod *objects.Model) error {
	switch ctrl_obj.Action {
	case "READ":
		return write.WriteNasCards(ctrl_obj.OutputFile, mod)
	case "STATS":
		read.GetNasCardsStatistics(mod)
		debug.DebugPrintoutNasCardStats(mod)
		return nil
	case "SPLIT":
		return write.WriteCardsToFiles(ctrl_obj.OutputDir, mod)
	default:
		return fmt.Errorf("unknown action: %s", ctrl_obj.Action)
	}
}
