package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type NastranModel struct {
	Nodes map[int]*Node
}

type Node struct {
	ID      int
	CP      int
	X, Y, Z float64
	CD      int
	PS      int
}

func NewNastranModel() *NastranModel {
	return &NastranModel{
		Nodes: make(map[int]*Node),
	}
}

func (m *NastranModel) ReadDat(filename string) error {
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
	model := NewNastranModel()

	dir, e := os.Getwd()
	if e != nil {
		fmt.Printf("Fehler: %v\n", e)
		return
	}
	fmt.Println("Aktuelles Verzeichnis:", dir)

	err := model.ReadDat("./regression_tests/sol_103_meter.dat")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Gelesen: %d Knoten\n", len(model.Nodes))

	// Ausgabe der ersten paar Knoten
	for id, node := range model.Nodes {
		if id > 5 {
			break
		}
		fmt.Printf("Node %d: (%.3f, %.3f, %.3f)\n",
			node.ID, node.X, node.Y, node.Z)
	}
}

// func main() {
//
// 	type NastranModel struct {
// 		Nodes map[int]*Node // GRID: ID -> Node
// 		//	Elements   map[int]Element   // CBEAM, CTRIA3 etc.: ID -> Element (Interface)
// 		//	Properties map[int]Property  // PSHELL, PBAR etc.: ID -> Property (Interface)
// 		//	Materials  map[int]*Material // MAT1: ID -> Material
// 		//	Loads      map[int]*Load     // FORCE, SPC: ID -> Load
// 		// Weitere: Sets, Constraints etc.
// 	}
// 	type Node struct {
// 		ID      int
// 		CP      int // Coordinate system
// 		X, Y, Z float64
// 		CD      int // Displacement coord system
// 		PS      int // Permanent single-point constraints
// 	}
//
// 	var my_string string
// 	my_string = "nas_tools"
// 	fmt.Print("  ..  ", my_string, "\n")
//
// }

//	type Element interface {
//		ID() int
//		Pid() int // Property ID
//		NodeIDs() []int
//		Type() string // "CBEAM", "CTRIA3"
//	}
//
//	// Property-Interface
//	type Property interface {
//		ID() int
//		Mid() int // Material ID
//		Type() string
//	}
//	type BeamElement struct { // CBEAM
//		id     int
//		pid    int
//		n1, n2 int // Endknoten
//		orient int // Orientation vector/grid
//	}
//
//	type PShell struct {
//		Id   int
//		Mid1 int // Membran-Material
//		Mid2 int // Biegung (optional)
//		Mid3 int // Scherung (optional)
//		Mid4 int // Membran-Schub (optional)
//
//		T   float64 // Grunddicke
//		Nsm float64 // Nicht-strukturelle Masse
//		Z1  float64 // Obere Faserlage
//		Z2  float64 // Untere Faserlage
//
//		// optional: Dicken-Skalierung, Integration-Flag etc.
//	}
//
//	type ShellProperty struct { // PSHELL
//		id  int
//		mid int     // Material ID
//		t   float64 // Thickness
//		// Tflags, Ti für variable Dicken
//	}
//
//	type BarProperty struct { // PBAR: Querschnitt für Balken
//		id        int
//		mid       int
//		A, I1, I2 float64 // Area, Momente etc.
//		// Vollständig: 10+ Felder wie J, NSI etc.
//	}
//	type Material struct { // MAT1
//		ID      int
//		E       float64 // Young's modulus
//		G       float64 // Shear modulus
//		Nu      float64 // Poisson's ratio
//		Rho     float64 // Density
//		A, Tref float64
//	}
//
//	type Load struct { // FORCE/SPC
//		ID   int
//		Type string // "FORCE", "SPC"
//		Nid  int    // Node ID
//		Mag  float64
//		Dir  int // DOF 1-6
//		// Für FORCE: Scale, CID; für SPC: Components "123456"
//	}

// func (e *BeamElement) ID() int { return e.id }
// func (e *BeamElement) Pid() int { return e.pid }
// func (e *BeamElement) NodeIDs() []int { return []int{e.n1, e.n2} }
// func (e *BeamElement) Type() string { return "CBEAM" }
// func (p *PShell) ID() int   { return p.Id }
// func (p *PShell) Mid() int  { return p.Mid1 }
// func (p *PShell) Type() string { return "PSHELL" }
// func (p *ShellProperty) ID() int { return p.id }
// func (p *ShellProperty) Mid() int { return p.mid }
// func (p *ShellProperty) Type() string { return "PSHELL" }
