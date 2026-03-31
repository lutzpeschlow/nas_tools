package f06_methods

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
)

// ----------------------------------------------------------------------------
//
//	ReadGroundingForces
//
// ----------------------------------------------------------------------------
func ReadGroundingForces(ctrl *objects.Control, mod *objects.Model) error {
	fmt.Println("reading grounding forces ...")
	// set input file
	f06_file := ctrl.FullInputPath
	f, err := os.Open(f06_file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// prepare input file reader
	scanner := bufio.NewScanner(f)
	readFlag := false
	foundIDs := []string{}
	// loop over input file lines
	for scanner.Scan() {
		line := scanner.Text()
		// grounding force flag
		if strings.Contains(line, "G R O U N D   C H E C K   F O R C E S") {
			readFlag = true
			continue
		}
		// jump to next
		if !readFlag {
			continue
		}
		// create field from line
		//    check length of 8 members and "G" entry in second member
		//    check first member as integer value
		fields := strings.Fields(line)
		if len(fields) != 8 || fields[1] != "G" {
			continue
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			continue
		}
		// assign node id and calculate max absolute value from 6 force values
		maxAbs := 0.0
		for i := 2; i < 8; i++ {
			v, err := strconv.ParseFloat(fields[i], 64)
			if err != nil {
				continue
			}
			av := math.Abs(v)
			if av > maxAbs {
				maxAbs = av
			}
		}
		// compare with limit size to write into result file
		if maxAbs > ctrl.LimitSize {
			foundIDs = append(foundIDs, fields[0])
		}
		// fmt.Printf("Node %s -> max abs value: %.6e\n", id, maxAbs)
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("scan input file %s: %w", f06_file, err)
		}

	}
	// write text file and session file
	err = WriteTxtFile(ctrl.FullOutputPath, foundIDs)
	if err != nil {
		return err
	}
	err = WriteSessionFile("g", "node", ctrl.FullOutputPath, foundIDs)
	if err != nil {
		return err
	}
	// reporting
	fmt.Println("number of grounding nodes: ", len(foundIDs))
	// return value
	return nil
}

// ----------------------------------------------------------------------------
//
//	WriteTxtFile
//
// ----------------------------------------------------------------------------
func WriteTxtFile(filePath string, lines []string) error {
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create output file %s: %w", filePath, err)
	}
	defer out.Close()
	for _, line := range lines {
		if _, err := fmt.Fprintln(out, line); err != nil {
			return fmt.Errorf("write output file %s: %w", filePath, err)
		}
	}
	return nil
}

// ----------------------------------------------------------------------------
//
//	WriteSessionFile
//
// ----------------------------------------------------------------------------
func WriteSessionFile(groupName, entityType, outputPath string, ids []string) error {
	sesPath := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".ses"
	f, err := os.Create(sesPath)
	if err != nil {
		return fmt.Errorf("create session file %s: %w", sesPath, err)
	}
	defer f.Close()
	if _, err := fmt.Fprintf(f, "ga_group_create(%q)\n", groupName); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(f, "ga_group_entity_add(%q, %q // @ \n", groupName, entityType); err != nil {
		return err
	}
	for i, id := range ids {
		if i == len(ids)-1 {
			if _, err := fmt.Fprintf(f, "\" %s \" )\n", id); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(f, "\" %s \" // @ \n", id); err != nil {
				return err
			}
		}
	}
	return nil
}

// ----------------------------------------------------------------------------
//
//	ReadMasslessMechanisms
//
// ----------------------------------------------------------------------------
func ReadMasslessMechanisms(ctrl *objects.Control, mod *objects.Model) error {
	// variables
	fmt.Println(" .. .. ..")
	// return list

	return nil
}
