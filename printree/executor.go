package printree

import "fmt"

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

	SerializeSetting struct {
		*TreeMarks
		*DrillDown
		Println
	}

	Option func(conf *SerializeSetting)
)

var DefaultMarks = TreeMarks{
	LineItem: "├─",
	LastItem: "└─",
	Cross:    "│ ",
	Blank:    "  ",
}

func (e *SerializeSetting) pipe(opts ...Option) *SerializeSetting {
	if e == nil {
		e = &SerializeSetting{}
	}
	for _, opt := range opts {
		opt(e)
	}
	if e.TreeMarks == nil {
		e.TreeMarks = &DefaultMarks
	}
	if e.Println == nil {
		e.Println = fmt.Println
	}
	return e
}

// OptSetMarks can set the tree marks
// Default Values:
// - TreeMarks.LineItem: "├─"
// - TreeMarks.LastItem: "└─",
// - TreeMarks.Cross:    "│ ",
// - TreeMarks.Blank:    "  ",
// You can change default values by set var DefaultMarks
func OptSetMarks(m TreeMarks) Option {
	return func(conf *SerializeSetting) {
		conf.TreeMarks = &m
	}
}

// OptDrillDownPolicy can be used to set the DrillDown rules
// DrillDown.Predictor is a predicate function to determine if the node should be drill down
// DrillDown.ShowMark is the sentence to print when a not empty node are skipped
func OptDrillDownPolicy(dd *DrillDown) Option {
	return func(conf *SerializeSetting) {
		conf.DrillDown = dd
	}
}

// OptCustomPrinter can be used to set a custom printer
// the default printer is `fmt.Println`
func OptCustomPrinter(printer Println) Option {
	return func(conf *SerializeSetting) {
		conf.Println = printer
	}
}

// BuildPrefixes 为当前节点及其子节点构建前缀字符串
//
// @params
// - isLastNode 参数标记着是否是当前层次的最后一个元素
//
// @returns
// - currentPrefix 是加上了当前层次前缀之后的字符串
// - nextLevelPrefix 是用来构建下一层的前缀字符串
func (m TreeMarks) BuildPrefixes(currentLevelPrefix string, isLastNode bool) (currentPrefix, nextLevelPrefix string) {
	currentPrefix = currentLevelPrefix // 使用当前层的前缀
	if isLastNode {
		currentPrefix += m.LastItem                    // 如果是最后一个节点，使用 LastItem 标记
		nextLevelPrefix = currentLevelPrefix + m.Blank // 为子节点打印准备不包含竖线的新层次前缀
	} else {
		currentPrefix += m.LineItem                    // 如果不是最后一个节点，使用 LineItem 标记
		nextLevelPrefix = currentLevelPrefix + m.Cross // 为子节点打印准备包含竖线的新层次前缀
	}
	return
}
