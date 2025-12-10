package ctrl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// control object
type Control_Object struct {
	Action string
}

// ReadControlFile function to read a control file
//
// input:
//
// output:
//   - error: if read or parse fails, put back error, else nil
func ReadControlFile(path string, obj *Control_Object, osName string) error {
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
			}
		}
	}
	// return value is the error interface value of the scanner
	return scanner.Err()
}

func DebugPrintoutCtrlObj(obj *Control_Object) {
	fmt.Print("debug printout of control object: \n")
	fmt.Print(" Action:    ", obj.Action, "\n")
}
