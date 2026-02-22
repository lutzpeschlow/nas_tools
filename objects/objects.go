package objects

type Control struct {
	// orig data stored in map
	Enable   map[string]bool `json:"enable"`
	Defaults struct {
		InputFile string `json:"input_file"`
		InputDir  string `json:"input_dir"`
	} `json:"defaults"`
	Actions map[string]interface{} `json:"actions"`
	// for better usage in functions, extract to single values
	Action     string
	InputFile  string
	InputDir   string
	OutputFile string
	OutputDir  string
	Option01   string
	Array01    []interface{}
	Input01    string
	// combined from previous
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
	NasFieldList []*ParsedCard
	NasCardStats map[string]int
}

// NasCard object
type NasCard struct {
	Card []string
}

type ParsedCard struct {
	Name   string
	Fields []string // Max 10 Felder
	Index  int      // Ursprungs-Index in NasCardList
}
