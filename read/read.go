package read

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
)

// function: ReadNasCards
//
// description:
//
// input :  file name,  model object
// output : error/return value
func ReadNasFile(filename string, obj *objects.Model) (len1 int, len2 int, err error) {
	fmt.Println("... read nastran cards: ", filename)
	// PASS 1
	// searching for BEGIN BULK
	// hasBulk is set
	file1, err := os.Open(filename)
	if err != nil {
		fmt.Println("... problem reading file: ", filename)
		return len1, len2, err
	}
	defer file1.Close()

	scanner1 := bufio.NewScanner(file1)
	hasBulk := false
	for scanner1.Scan() {
		if strings.Contains(scanner1.Text(), "BEGIN BULK") {
			hasBulk = true
			break
		}
	}
	if !hasBulk {
		fmt.Println("WARN: No BEGIN BULK")
	}
	// PASS 2
	// now regular parsing
	// inBulk will be set depending on hasBulk
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("... problem reading file: ", filename)
		return len1, len2, err
	}
	defer file.Close()
	// set true if there was no BEGIN BULK in PASS 1
	// set false if there was any BEGIN BULK, and wait till BEGIN BULK
	inBulk := !hasBulk
	// scan object, str.Builder object and further variables
	scanner := bufio.NewScanner(file)
	// create map with key: int and value: *NasCard
	obj.NasCards = make(map[int]*objects.NasCard)
	var currentCard []string
	inCard := false
	var first_sign byte
	lineCount := 0
	// extract each card in separate block
	// loop over lines
	for scanner.Scan() {
		lineCount = lineCount + 1
		line := scanner.Text()

		// are we in BEGIN BULK, set inBulk true
		if strings.Contains(line, "BEGIN BULK") {
			inBulk = true
			fmt.Println("found BEGIN BULK ...")
			continue
		}
		// in case of ENDDATA, break loop routine
		if strings.Contains(line, "ENDDATA") {
			inBulk = false
			fmt.Println("found ENDDATA ...")
			break
		}
		// inside BULK block, we are allowed to parse
		if !inBulk {
			continue
		}
		// (1) COMMENT - direct jump
		if strings.HasPrefix(line, "$") {
			continue
		}
		// check first sign of line
		if len(line) == 0 {
			first_sign = '-'
		} else {
			first_sign = line[0]
		}
		// (2) CARD - first is a letter
		// new card alert, write existing buffer in object and setup for new card
		if (first_sign >= 'a' && first_sign <= 'z') || (first_sign >= 'A' && first_sign <= 'Z') {
			// (2.1) there is content in current card that needs to be saved in object
			if len(currentCard) > 0 {
				// write existing card into NasCards and NasCardList
				// create new string slice, initialize nil, add content
				nextID := len(obj.NasCards)
				newCard := &objects.NasCard{Card: append([]string(nil), currentCard...)}
				obj.NasCards[nextID] = newCard
				obj.NasCardList = append(obj.NasCardList, newCard)
				// clean up current card and assign new data
				currentCard = currentCard[:0]
				currentCard = append(currentCard, line)
				inCard = true
				// (2.2) simply add line to current card, no action with previous card
			} else {
				currentCard = append(currentCard, line)
				inCard = true
			}
			// (3) CONT - in card any continuation line
		} else {
			if inCard {
				currentCard = append(currentCard, line)
			}
		}
	}
	// cover last line in buffer ...
	if len(currentCard) > 0 {
		nextID := len(obj.NasCards)
		newCard := &objects.NasCard{Card: append([]string(nil), currentCard...)}
		obj.NasCards[nextID] = newCard
		obj.NasCardList = append(obj.NasCardList, newCard)
	}
	// read file stats
	fmt.Println("lines/cards: ", lineCount, len(obj.NasCards), len(obj.NasCardList))
	// return scanner error
	return len(obj.NasCards), len(obj.NasCardList), scanner.Err()
	// end of function - ReadNasFile
}

// function: GetNasCardsStatistics
//
// description:
//
// input : model object
// output : error/return value
func GetNasCardsStatistics(obj *objects.Model) (int, error) {
	fmt.Println("... get stats", len(obj.NasCards))
	// init map for stats
	obj.NasCardStats = make(map[string]int)
	// loop through all cards
	for _, card := range obj.NasCards {
		if len(card.Card) == 0 {
			continue
		}
		// first line of card
		firstLine := card.Card[0]
		cardType := ExtractCardName(firstLine)
		// add counter to according card in the statistics object
		if cardType != "" {
			obj.NasCardStats[cardType]++
		}
	}
	return len(obj.NasCardStats), nil
}

// function: extractCardName
//
// description:
//
// input : string
// output : string
func ExtractCardName(line string) string {
	if len(line) < 4 {
		return ""
	}
	name := strings.TrimSpace(line[:8])
	sepIndex := strings.IndexAny(name, ",+ *")

	if sepIndex > 0 {
		name = name[:sepIndex] // Bis zum Trennzeichen abschneiden
	}

	return strings.ToUpper(strings.TrimSpace(name))

}

func ExtractCardID(line string) string {
	// (1) Free Field
	if strings.Contains(line[:10], ",") {
		fields := strings.Split(line, ",")
		if len(fields) >= 2 {
			return strings.TrimSpace(fields[1])
		}
		return ""
	}
	// (2) FIXED FIELD
	if strings.Contains(line[:8], "*") {
		// (2.1) Large Field
		id := strings.TrimSpace(line[8:24])
		return strings.TrimRight(id, " ")
	} else {
		// (2.2) Small Field
		id := strings.TrimSpace(line[8:16])
		return strings.TrimRight(id, " ")
	}
}

// GetCardEntry extracts the entry at the specified position from a NASTRAN card.
// line is 0-based index of the starting line in the card slice.
// entry is 1-based field number (1-10).
// card is slice of strings, each representing one line of the card.
func GetCardEntry(line, entry int, card []string) string {
	if line < 0 || entry < 1 || entry > 10 || line >= len(card) {
		return ""
	}

	// Detect format from first line of card
	firstLine := strings.TrimLeft(card[line], " \t")
	isLongField := strings.Contains(firstLine, "*")

	fieldSize := 8
	fieldsPerLine := 10
	if isLongField {
		fieldSize = 16
		fieldsPerLine = 6 // First field 8 chars, then 4x16 chars = 72 chars total
	}

	// Calculate which line and position the entry is on
	targetLine := line + (entry-1)/fieldsPerLine
	if targetLine >= len(card) {
		return ""
	}

	fieldPos := (entry - 1) % fieldsPerLine
	lineStr := card[targetLine]

	// Handle continuation lines (start with spaces, +, ,, or *)
	contPrefix := ""
	if targetLine > line {
		trimmed := strings.TrimLeft(lineStr, " \t")
		if len(trimmed) > 0 && (trimmed[0] == '+' || trimmed[0] == ',' || trimmed[0] == '*') {
			contPrefix = trimmed[:1]
			lineStr = strings.TrimLeft(lineStr, " \t"+contPrefix)
		}
	}

	// Extract field using fixed width
	startPos := fieldPos * fieldSize
	if startPos >= len(lineStr) {
		return ""
	}

	endPos := startPos + fieldSize
	if endPos > len(lineStr) {
		endPos = len(lineStr)
	}

	field := strings.TrimSpace(lineStr[startPos:endPos])

	// Remove continuation identifier if it's in the last field of continuation line
	if fieldPos == fieldsPerLine-1 && targetLine > line && len(field) > 0 {
		field = strings.TrimSpace(regexp.MustCompile(`^[+\s,*$]+|[+\s,*$]+$`).ReplaceAllString(field, ""))
	}

	return field
}

func Add(a, b int) int {
	return a + b
}
