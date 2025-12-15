package objects

import (
	"fmt"

	"github.com/lutzpeschlow/nas_tools/objects"
)

// control object
type Control_Object struct {
	Action string
}

// Model object
// as main object to save model data
//
//	Nodes - hash map with integer key
type Model struct {
	NasCards map[int]*NasCard
}

// NasCard object
type NasCard struct {
	Card []string
}

func DebugPrintoutModelObj(obj *objects.Model) {
	fmt.Print("debug printout of control object: \n")
	fmt.Print(" Action:    ", obj.NasCards, "\n")
}
