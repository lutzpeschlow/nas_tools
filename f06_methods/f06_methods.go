package f06_methods

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
	"github.com/lutzpeschlow/nas_tools/utils"
)

// ----------------------------------------------------------------------------
//
//		ReadGroundingForces
//
//	                                                                                                   DIRECTION        4
//	                          G R O U N D   C H E C K   F O R C E S  ( N + A UT O - S E T )
//
//	     POINT ID.   TYPE          T1             T2             T3             R1             R2             R3
//	     10800080      G      0.0            0.0            0.0            2.958194E+11   0.0            0.0
//	     10800100      G      0.0           -5.570115E+11   5.321358E+11   0.0            0.0            0.0
//	     10800270      G      0.0            9.622509E+11   0.0            2.724761E+11   0.0           -8.686390E+11
//	     10800280      G      0.0            0.0           -2.535309E+11   0.0            0.0           -2.648989E+11
//	     10800290      G      0.0            0.0           -2.943361E+11   0.0            0.0            0.0
//	     10800320      G      0.0           -4.100117E+11  -3.341411E+11   0.0            0.0            0.0
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
	err = utils.WriteTxtFile(ctrl.FullOutputPath, foundIDs)
	if err != nil {
		return err
	}
	err = utils.WriteSessionFile("grounding_nodes", "node", ctrl.FullOutputPath, foundIDs)
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
//		ReadMasslessMechanisms
//
//	   15201301 T1  1.00000E+00  15201301 T2  1.00000E+00  15201301 T3  1.00000E+00  15201301 R1  1.00000E+00  15201301 R2  1.00000E+00
//	   15201301 R3  1.00000E+00  15201302 T1  1.00000E+00  15201302 T2  1.00000E+00  15201302 T3  1.00000E+00
//	   15201302 R2  1.00000E+00  15201302 R3  1.00000E+00  15201303 T1  1.00000E+00  15201303 T2  1.00000E+00  15201303 T3  1.00000E+00
//	   15201301 T1  1.00000E+00  15201301 T2  1.00000E+00  15201301 T3  1.00000E+00
//	   15201301 R3  1.00000E+00  15201302 T1  1.00000E+00  15201302 T2  1.00000E+00  15201302 T3  1.00000E+00  15201302 R1  1.00000E+00
//	   15201302 R2  1.00000E+00  15201302 R3  1.00000E+00
//	   15201301 T1  1.00000E+00  15201301 T2  1.00000E+00  15201301 T3  1.00000E+00  15201301 R1  1.00000E+00  15201301 R2  1.00000E+00
//	   15201301 R3  1.00000E+00  15201302 T1  1.00000E+00  15201302 T2  1.00000E+00  15201302 T3  1.00000E+00  15201302 R1  1.00000E+00
//	   15201302 R2  1.00000E+00
//
// ----------------------------------------------------------------------------
func ReadMasslessMechanisms(ctrl *objects.Control, mod *objects.Model) error {
	fmt.Println("reading massless mechanisms ...")
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
	allowedDirs := []string{"T1", "T2", "T3", "R1", "R2", "R3"}
	// loop over input file lines
	for scanner.Scan() {
		line := scanner.Text()
		// VAXW table recognized
		if strings.Contains(line, "VAXW") {
			readFlag = true
			continue
		}
		if !readFlag {
			continue
		}
		// split line
		fields := strings.Fields(line)
		// check block length of 3,6,9,12,15
		if len(fields)%3 != 0 {
			continue
		}
		// extract node ids from blocks
		for i := 0; i < len(fields); i += 3 {
			// check of T1,T2,T3,R1,R2,R3 content
			if i+1 >= len(fields) {
				continue
			}
			dir := fields[i+1]
			okDir := false
			for _, allowed := range allowedDirs {
				if strings.Contains(dir, allowed) {
					okDir = true
					break
				}
			}
			if !okDir {
				continue
			}
			// put node id into foundIDs array
			if _, err := strconv.Atoi(fields[i]); err == nil {
				foundIDs = append(foundIDs, fields[i])
			}
		}
	}
	// remove duplicate entries
	foundIDs = utils.RemoveDuplicateEntries(foundIDs)
	// write text file and session file
	err = utils.WriteTxtFile(ctrl.FullOutputPath, foundIDs)
	if err != nil {
		return err
	}
	err = utils.WriteSessionFile("massless_nodes", "node", ctrl.FullOutputPath, foundIDs)
	if err != nil {
		return err
	}
	// reporting
	fmt.Println("number of massless nodes: ", len(foundIDs))
	// return value
	return nil
}
