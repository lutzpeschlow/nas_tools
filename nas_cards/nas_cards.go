package nas_cards

import (
	"fmt"

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

func CreateNasBlocksFromTxt(ctrl *objects.Control, mod *objects.Model) error {
	fmt.Println(" nas blocks ... \n")

	return nil
}
