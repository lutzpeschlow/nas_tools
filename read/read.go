package read

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
)

func ReadDat(filename string, obj *objects.Model) error {
	// open file, with defer as backup
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// scan file and assign data
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "GRID") {
			node, err := parseGRID(line)
			if err == nil {
				obj.Nodes[node.ID] = node
			}
		}
	}
	return scanner.Err()
}

func parseGRID(line string) (*objects.Node, error) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return nil, fmt.Errorf("invalid GRID line: %s", line)
	}

	node := &objects.Node{}

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
