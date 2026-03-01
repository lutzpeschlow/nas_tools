package nas_deck

import (
	"fmt"
	"strings"

	"github.com/lutzpeschlow/nas_tools/objects"
)

// ----------------------------------------------------------------------------
//
//	p
//
// ----------------------------------------------------------------------------

func ParseAllCards(model *objects.Model) error {
	var parsed []*objects.ParsedCard
	i := 0

	for i < len(model.NasCardList) {
		fmt.Println(model.NasCardList[i].Card[0])
		switch {
		// case strings.Contains(model.NasCardList[i].Card[0], ","):
		// p := extractFree(model.NasCardList[i], i) // ← Index übergeben
		// if p != nil {
		// 	parsed = append(parsed, p)
		// }
		// i++

		// case strings.Contains(model.NasCardList[i].Card[0], "*") && i+1 < len(model.NasCardList):
		// p := extractLarge(model.NasCardList[i], model.NasCardList[i+1], i) // ← Index
		// if p != nil {
		// 	parsed = append(parsed, p)
		// }
		// i += 2

		default:
			p, advance := extractSmall(model.NasCardList, i) // ← start=i
			if p != nil {
				parsed = append(parsed, p)
			}
			i += advance
			// i = i + 1
		}
	}
	model.NasFieldList = parsed
	return nil
}

// Free Field: "grid,11,,100."
// func extractFree(card *objects.NasCard, index int) *objects.ParsedCard {
// 	full := strings.Join(card.Card, " ")
// 	if idx := strings.Index(full, "$"); idx != -1 {
// 		full = full[:idx]
// 	}
// 	parts := strings.Split(full, ",")
// 	if len(parts) < 2 {
// 		return nil
// 	}
//
// 	name := strings.TrimSpace(strings.ToUpper(parts[0]))
// 	fields := make([]string, 0, 10)
// 	for _, part := range parts[1:] {
// 		clean := strings.TrimSpace(part)
// 		if clean != "" && len(fields) < 10 {
// 			fields = append(fields, clean)
// 		}
// 	}
//
// 	return &objects.ParsedCard{Name: name, Fields: fields, Index: index}
// }

// Large Field: GRID* + *Zeile
// func extractLarge(card1, card2 *objects.NasCard, index int) *objects.ParsedCard {
// 	full1 := strings.Join(card1.Card, "")
// 	full2 := strings.Join(card2.Card, "")
// 	content := full1 + full2
//
// 	// Name: GRID* → GRID
// 	re := regexp.MustCompile(`^([A-Z]{1,8})\*`)
// 	matches := re.FindStringSubmatch(content)
// 	name := ""
// 	if len(matches) > 1 {
// 		name = strings.ToUpper(matches[1])
// 	}
//
// 	// 16-Char-Felder extrahieren
// 	fields := make([]string, 0, 10)
// 	for j := 8; j < len(content)-16 && len(fields) < 10; j += 16 {
// 		field := content[j : j+16]
// 		clean := strings.TrimSpace(field)
// 		if clean != "" && !strings.HasPrefix(clean, "*") {
// 			fields = append(fields, clean)
// 		}
// 	}
//
// 	return &objects.ParsedCard{Name: name, Fields: fields, Index: index}
// }

// Small Field + Kontinuationen
func extractSmall(cards []*objects.NasCard, start int) (*objects.ParsedCard, int) {
	var fields []string
	name := ""
	i := start
	for i < len(cards) && len(fields) < 10 {
		card := cards[i].Card
		if len(card) == 0 {
			i++
			continue
		}
		// Name aus erster Karte
		if name == "" {
			name = strings.TrimSpace(strings.ToUpper(card[0]))
		}
		// Felder (Name überspringen)
		for j, f := range card {
			if name != "" && j == 0 {
				continue
			}
			clean := cleanField(f)
			if clean != "" {
				fields = append(fields, clean)
			}
		}
		if !isContinuation(cards[i+1].Card) {
			break
		}
		i++
	}
	if name == "" {
		return nil, i - start + 1
	}

	return &objects.ParsedCard{Name: name, Fields: fields, Index: start}, i - start + 1
}

func cleanField(field string) string {
	padded := padField(field, 8)
	return strings.TrimFunc(padded, func(r rune) bool {
		return r == ' ' || r == '+' || r == '*'
	})
}

func padField(field string, width int) string {
	if len(field) < width {
		return field + strings.Repeat(" ", width-len(field))
	}
	return field[:width]
}

func isContinuation(card []string) bool {
	if len(card) == 0 {
		return false
	}
	first := strings.TrimSpace(card[0])
	return first == "" || first == "+" || strings.HasPrefix(first, "*") || strings.HasPrefix(first, ",")
}
