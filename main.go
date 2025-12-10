package main

// libraries
import (
	"fmt"
	"os"
	"runtime"

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

func main() {
	ctrl_obj := ctrl.Control_Object{}
	osName := runtime.GOOS
	err_ctrl := ctrl.ReadControlFile("control.txt", &ctrl_obj, osName)
	if err_ctrl != nil {
		fmt.Printf(" %v\n", err_ctrl)
		os.Exit(1)
	}
	ctrl.DebugPrintoutCtrlObj(&ctrl_obj)

	// model instance
	mod := &Model{}
	// create map with key: int and value: *Node
	mod.Nodes = make(map[int]*Node)
	// get current directory
	current_dir, _ := os.Getwd()
	fmt.Println("current directory:", current_dir)
	// read input file
	dat_file := "./regression_tests/sol_103_meter.dat"
	err := read.ReadDat(dat_file, &mod)
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
