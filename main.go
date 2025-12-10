package main

// libraries
import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/lutzpeschlow/nas_tools/ctrl"
)

// Model object
// as main object to save model data
//
//	Nodes - hash map with integer key
type Model struct {
	Nodes map[int]*Node
}

// Node object
type Node struct {
	ID      int
	CP      int
	X, Y, Z float64
	CD      int
	PS      int
}

func (m *Model) ReadDat(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "GRID") {
			node, err := parseGRID(line)
			if err == nil {
				m.Nodes[node.ID] = node
			}
		}
	}
	return scanner.Err()
}

func parseGRID(line string) (*Node, error) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return nil, fmt.Errorf("invalid GRID line: %s", line)
	}

	node := &Node{}

	fmt.Print(fields, "\n")

	// ID (Feld 1)
	node.ID, _ = strconv.Atoi(fields[1])

	// CP (optional, Feld 2 oder default 0)
	if len(fields) > 2 {
		cp, err := strconv.Atoi(fields[2])
		if err == nil {
			node.CP = cp
		}
	}

	// Koordinaten X,Y,Z (letzte 3 Felder)
	if len(fields) >= 7 {
		// Vollständig: ID,CP,X,Y,Z,CD,PS
		x, _ := strconv.ParseFloat(fields[3], 64)
		y, _ := strconv.ParseFloat(fields[4], 64)
		z, _ := strconv.ParseFloat(fields[5], 64)
		cd, _ := strconv.Atoi(fields[6])
		ps, _ := strconv.Atoi(fields[7])

		node.X, node.Y, node.Z = x, y, z
		node.CD, node.PS = cd, ps
	} else {
		// Dein File-Format: GRID ID X Y Z
		x, _ := strconv.ParseFloat(fields[2], 64)
		y, _ := strconv.ParseFloat(fields[3], 64)
		z, _ := strconv.ParseFloat(fields[4], 64)

		node.X, node.Y, node.Z = x, y, z
	}

	return node, nil
}

func main() {
	ctrl_obj := ctrl.Control_Object{}
	osName := runtime.GOOS

	// model instance
	mod := &Model{}
	// create map with key: int and value: *Node
	mod.Nodes = make(map[int]*Node)
	// get current directory
	current_dir, _ := os.Getwd()
	fmt.Println("current directory:", current_dir)
	// read input file
	err := mod.ReadDat("./regression_tests/sol_103_meter.dat")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	// number of nodes
	fmt.Print("num nodes: ", len(mod.Nodes), "\n")
	// node list
	for id, node := range mod.Nodes {
		if id < 5 {
			fmt.Print("", node.ID, node.X, node.Y, node.Z, "\n")
		}
	}
}
