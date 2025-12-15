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
	NasCards map[int]*NasCard
}

// NasCard object
type NasCard struct {
	Card []string
}
