package objects

// control object
type Control_Object struct {
	Action string
}

// Model object
// as main object to save model data
//
//	Nodes - hash map with integer key
type Model struct {
	Nodes    map[int]*Node
	NasCards map[int]*NasCard
}

// NasCard object
type NasCard struct {
	Card []string
}

// Node object
type Node struct {
	ID      int
	CP      int
	X, Y, Z float64
	CD      int
	PS      int
}
