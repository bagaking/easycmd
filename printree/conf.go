package printree

type (
	Println   func(a ...interface{}) (n int, err error)
	Predictor func(node IPrintableTreeNode) bool

	TreeMarks struct {
		LineItem string
		LastItem string
		Cross    string
		Blank    string
	}

	DrillDown struct {
		Predictor Predictor
		ShowMark  string
	}

	Config struct {
		*TreeMarks
		*DrillDown
		Println
	}

	Option func(conf *Config)
)

var DefaultMarks = TreeMarks{
	LineItem: "├─",
	LastItem: "└─",
	Cross:    "│ ",
	Blank:    "  ",
}

// OptSetMarks can set the tree marks
// Default Values:
// - TreeMarks.LineItem: "├─"
// - TreeMarks.LastItem: "└─",
// - TreeMarks.Cross:    "│ ",
// - TreeMarks.Blank:    "  ",
// You can change default values by set var DefaultMarks
func OptSetMarks(m TreeMarks) Option {
	return func(conf *Config) {
		conf.TreeMarks = &m
	}
}

// OptDrillDownPolicy can be used to set the DrillDown rules
// DrillDown.Predictor is a predicate function to determine if the node should be drill down
// DrillDown.ShowMark is the sentence to print when a not empty node are skipped
func OptDrillDownPolicy(dd *DrillDown) Option {
	return func(conf *Config) {
		conf.DrillDown = dd
	}
}

// OptCustomPrinter can be used to set a custom printer
// the default printer is `fmt.Println`
func OptCustomPrinter(printer Println) Option {
	return func(conf *Config) {
		conf.Println = printer
	}
}

func (m TreeMarks) HandlePrefix(LayerPrefix string, finalItem bool) (curItemPrefix, newLayerPrefix string) {
	newLayerPrefix = LayerPrefix
	curItemPrefix = LayerPrefix
	if finalItem {
		curItemPrefix += m.LastItem
		newLayerPrefix += m.Blank
		return
	}
	curItemPrefix += m.LineItem
	newLayerPrefix += m.Cross
	return
}
