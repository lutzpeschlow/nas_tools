package cmd

import (
	"fmt"

	"github.com/lutzpeschlow/nas_tools/debug"
	"github.com/lutzpeschlow/nas_tools/modify"
	"github.com/lutzpeschlow/nas_tools/objects"
	"github.com/lutzpeschlow/nas_tools/read"
	"github.com/lutzpeschlow/nas_tools/write"
)

// ----------------------------------------------------------------------------
//
//	ExecuteAction
//
// ----------------------------------------------------------------------------
func ExecuteAction(ctrl *objects.Control, mod *objects.Model) error {
	// switch for actions
	switch ctrl.Action {
	// read nastran file and write to new file
	case "READ":
		err := write.WriteNasCards(ctrl, mod)
		if err != nil {
			return fmt.Errorf("WriteNasCards failed: %w", err)
		}
	// get statistics of nastran file
	case "STATS":
		_, err := read.GetNasCardsStatistics(mod)
		if err != nil {
			return fmt.Errorf("GetNasCardsStatistics failed: %w", err)
		}
		debug.DebugPrintoutNasCardStats(mod)
	// split nastran file into several files per card
	case "SPLIT":
		err := write.WriteCardsToFiles(ctrl.OutputDir, mod)
		if err != nil {
			return fmt.Errorf("WriteCardsToFiles failed: %w", err)
		}
	// extract entities according list
	case "EXTRACT_ACC_LIST":
		err := modify.ExtractCardsAccordingList(ctrl, mod)
		if err != nil {
			return fmt.Errorf("ExtractCardsAccordingList failed: %w", err)
		}
	// get card entry
	case "GET_CARD_ENTRY":
		err := read.GetCardEntry(ctrl, mod)
		if err != nil {
			return fmt.Errorf("GetCardEntry failed: %w", err)
		}
	// unknown action
	default:
		return fmt.Errorf("unknown action: %s", ctrl.Action)
	}
	// return variable
	return nil
}
