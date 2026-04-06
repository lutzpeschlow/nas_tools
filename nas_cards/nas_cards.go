package nas_cards

import "github.com/lutzpeschlow/nas_tools/objects"

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
