package cmd

import (
	"fmt"

	"github.com/lutzpeschlow/nas_tools/debug"
	"github.com/lutzpeschlow/nas_tools/f06_methods"
	"github.com/lutzpeschlow/nas_tools/nas_methods"
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
		dat_file := ctrl.FullInputPath
		_, _, err := read.ReadNasFile(dat_file, mod)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", dat_file, err)
		}
		err = write.WriteNasCards(ctrl, mod)
		if err != nil {
			return fmt.Errorf("WriteNasCards failed: %w", err)
		}
	// get statistics of nastran file
	case "STATS":
		dat_file := ctrl.FullInputPath
		_, _, err := read.ReadNasFile(dat_file, mod)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", dat_file, err)
		}
		_, err = read.GetNasCardsStatistics(mod)
		if err != nil {
			return fmt.Errorf("GetNasCardsStatistics failed: %w", err)
		}
		debug.DebugPrintoutNasCardStats(mod)
	// split nastran file into several files per card
	case "SPLIT":
		dat_file := ctrl.FullInputPath
		_, _, err := read.ReadNasFile(dat_file, mod)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", dat_file, err)
		}
		err = write.WriteCardsToFiles(ctrl.OutputDir, mod)
		if err != nil {
			return fmt.Errorf("WriteCardsToFiles failed: %w", err)
		}
	// extract entities according list
	case "EXTRACT_ACC_LIST":
		dat_file := ctrl.FullInputPath
		_, _, err := read.ReadNasFile(dat_file, mod)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", dat_file, err)
		}
		err = nas_methods.ExtractCardsAccordingList(ctrl, mod)
		if err != nil {
			return fmt.Errorf("ExtractCardsAccordingList failed: %w", err)
		}
	// get card entry
	case "GET_CARD_ENTRY":
		dat_file := ctrl.FullInputPath
		_, _, err := read.ReadNasFile(dat_file, mod)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", dat_file, err)
		}
		err, entry_list := nas_methods.GetCardEntries(ctrl, mod)
		if err != nil {
			return fmt.Errorf("GetCardEntry failed: %w", err)
		}
		fmt.Println("entry list length: ", len(entry_list))
	// read grounding forces
	case "READ_GROUNDING_FORCES":
		err := f06_methods.ReadGroundingForces(ctrl, mod)
		if err != nil {
			return fmt.Errorf("ReadGroundingForces failed: %w", err)
		}
	// read massless mechanisms
	case "READ_MASSLESS_MECH":
		err := f06_methods.ReadMasslessMechanisms(ctrl, mod)
		if err != nil {
			return fmt.Errorf("ReadMasslessMechanisms failed: %w", err)
		}

		// default
	default:
		return fmt.Errorf("unknown action: %s", ctrl.Action)
	}
	// return variable
	return nil
}
