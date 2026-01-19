package objects

// control object
type Control_Object struct {
	Action     string
	InputFile  string
	InputDir   string
	OutputFile string
	OutputDir  string
	Option01   string
	//
	FullInputPath  string
	FullOutputPath string
}

// Model object
// as main object to save model data
//
//	Nodes - hash map with integer key
type Model struct {
	NasCards     map[int]*NasCard
	NasCardList  []*NasCard
	NasCardStats map[string]int
}

// NasCard object
type NasCard struct {
	Card []string
}
