package nas_cards

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
)

// ----------------------------------------------------------------------------
//
//	CreateCard
//
// ----------------------------------------------------------------------------
func CreateCard(cardLine string) *objects.NasCard {
	return &objects.NasCard{
		Card: []string{cardLine},
	}
}

// ----------------------------------------------------------------------------
//
//	CreateNasBlocksFromTxt
//
// ----------------------------------------------------------------------------
func CreateNasBlocksFromTxt(ctrl *objects.Control) error {
	// variables
	var ent_ids []int
	// open id file
	id_file, err := os.Open(ctrl.FullInputPath)
	if err != nil {
		return fmt.Errorf("input file error: %w", err)
	}
	defer id_file.Close()
	// read id file
	scanner := bufio.NewScanner(id_file)
	for scanner.Scan() {
		// make a string strip
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// range case
		// example: 4:5   or 7:10
		if strings.Contains(line, ":") {
			range_values := strings.Split(line, ":")
			start_value, _ := strconv.Atoi(range_values[0])
			end_value, _ := strconv.Atoi(range_values[1])
			for id := start_value; id <= end_value; id++ {
				ent_ids = append(ent_ids, id)
			}
			// go direclty to next line
			continue
		}
		// normal case - only one entity
		id, _ := strconv.Atoi(line)
		ent_ids = append(ent_ids, id)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading ID file %q: %w", ctrl.Input01, err)
	}
	fmt.Println("number of ids: ", len(ent_ids))
	// output
	// output file
	out_file, err := os.Create(ctrl.OutputFile)
	if err != nil {
		return fmt.Errorf("Ausgabedatei erstellen: %w", err)
	}
	defer out_file.Close()
	w := bufio.NewWriter(out_file)
	defer w.Flush()
	// variables for writing loop
	const valuesPerLine = 8
	const fieldWidth = 8
	// write entity ids as blocks
	for i := 0; i < len(ent_ids); i += valuesPerLine {
		// consider variable end value of array
		end := min(i+valuesPerLine, len(ent_ids))
		// setup 8 values as inner loop
		var sb strings.Builder
		for _, id := range ent_ids[i:end] {
			sb.WriteString(fmt.Sprintf("%*d", fieldWidth, id))
		}
		// write into buffer
		if _, err := fmt.Fprintln(w, sb.String()); err != nil {
			return fmt.Errorf("output error: %w", err)
		}
	}
	// return value
	return nil
}
