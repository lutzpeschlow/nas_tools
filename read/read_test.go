package read

import (
	"strings"
	"testing"

	"github.com/lutzpeschlow/nas_tools/objects"
)

// TEST
// parsing test without any IO, the content of file
// is delivered here via input
func TestParseNasFromReader(t *testing.T) {
	input := `
$ comment
BEGIN BULK
GRID    1       0.0     0.0     0.0
GRID    2       1.0     0.0     0.0
ENDDATA
`
	r := strings.NewReader(input)

	var m objects.Model
	nCards, nList, err := ParseNasFromReader(r, &m)

	if err != nil {
		t.Fatal(err)
	}
	if nCards != 2 || nList != 2 {
		t.Fatalf("expected 2/2, got %d/%d", nCards, nList)
	}
	// weitere Asserts auf m.NasCardList ...
}
