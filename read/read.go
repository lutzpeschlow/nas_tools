package read

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
)

// ----------------------------------------------------------------------------
//
//	ReadNasFile
//
// ----------------------------------------------------------------------------
func ReadNasFile(filename string, obj *objects.Model) (int, int, error) {
	fmt.Println("... read nastran cards: ", filename)
	// open nas file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("... problem reading file: ", filename)
		return 0, 0, err
	}
	defer file.Close()
	// start parser
	return ParseNasFromReader(file, obj)
}

// ----------------------------------------------------------------------------
//
//	ParseNasFromReader
//
// ----------------------------------------------------------------------------
func ParseNasFromReader(r io.Reader, obj *objects.Model) (int, int, error) {
	// PASS 1:
	// search for BEGIN BULK
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, 0, err
	}
	text := string(data)
	scanner1 := bufio.NewScanner(strings.NewReader(text))
	// scan through complete file and set hasBulk accordingly
	hasBulk := false
	for scanner1.Scan() {
		line := strings.TrimLeft(scanner1.Text(), " \t") // Spaces + Tabs
		if strings.HasPrefix(line, "BEGIN BULK") {
			hasBulk = true
			break
		}
	}
	fmt.Println("INFO:  BEGIN BULK - ", hasBulk)
	// PASS 2
	// loop over lines
	scanner := bufio.NewScanner(strings.NewReader(text))

	// set true if there was no BEGIN BULK in PASS 1
	// set false if there was any BEGIN BULK, and wait till BEGIN BULK
	inBulk := !hasBulk
	// create map with   key: int ;  value: *NasCard
	obj.NasCards = make(map[int]*objects.NasCard)
	var currentCard []string
	inCard := false
	var firstSign byte
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
			firstSign = '-'
		} else {
			firstSign = line[0]
		}
		// (2) CARD - first is a letter
		// new card alert, write existing buffer in object and setup for new card
		if (firstSign >= 'a' && firstSign <= 'z') || (firstSign >= 'A' && firstSign <= 'Z') {
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
	// stats of file
	fmt.Println("lines/cards: ", lineCount, len(obj.NasCards), len(obj.NasCardList))
	// return scanner error
	return len(obj.NasCards), len(obj.NasCardList), scanner.Err()
}

// ----------------------------------------------------------------------------
//
//	GetNasCardsStatistics
//
// ----------------------------------------------------------------------------
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

// ----------------------------------------------------------------------------
//
//	ExtractCardName
//
// ----------------------------------------------------------------------------
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

// ----------------------------------------------------------------------------
//
//	ExtractCardID
//
// ----------------------------------------------------------------------------
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
