package printree

import (
	"fmt"
)

type (
	// IPrintableTreeNode is the interface should be implemented
	// by which node of the tree to be printed.
	IPrintableTreeNode interface {

		// ForEachPrintableChild is the method to iterating over direct child nodes
		ForEachPrintableChild(func(child IPrintableTreeNode, ind int) error) error

		// PrintableTreeLinePrefix returns the information that needs to be printed
		// in the prefix for each line of the tree
		PrintableTreeLinePrefix() string

		// PrintableTreeLinePrefix returns the count of direct child nodes
		PrintableChildCount() int
	}

	PrintableNodeSerializer func(node IPrintableTreeNode, ind int) string
)

func (nodeSerializer PrintableNodeSerializer) RecursivePrint(st IPrintableTreeNode, opts ...Option) error {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	return st.ForEachPrintableChild(func(child IPrintableTreeNode, ind int) error {
		return recursivePrint(child, "", ind, st.PrintableChildCount()-1, nodeSerializer, cfg)
	})
}

func recursivePrint(
	st IPrintableTreeNode,
	prefix string,
	ind, maxInd int,
	nodeSerializer func(node IPrintableTreeNode, ind int) string,
	cfg *Config,
) error {
	if cfg.TreeMarks == nil {
		cfg.TreeMarks = &DefaultMarks
	}
	if cfg.Println == nil {
		cfg.Println = fmt.Println
	}

	curItemPrefix, newLayerPrefix := cfg.TreeMarks.HandlePrefix(prefix, ind == maxInd)

	if _, err := cfg.Println(fmt.Sprintf("%v %s%s", st.PrintableTreeLinePrefix(), curItemPrefix, nodeSerializer(st, ind))); err != nil {
		return err
	}

	if cfg.DrillDown != nil && cfg.Predictor != nil {
		if !cfg.Predictor(st) {
			if cfg.DrillDown.ShowMark != "" && st.PrintableChildCount() > 0 {
				if _, err := cfg.Println(fmt.Sprintf("%v %s%s", st.PrintableTreeLinePrefix(), newLayerPrefix, cfg.ShowMark)); err != nil {
					return err
				}
			}
			return nil
		}
	}

	if err := st.ForEachPrintableChild(func(child IPrintableTreeNode, ind int) error {
		if err := recursivePrint(child, newLayerPrefix, ind, st.PrintableChildCount()-1, nodeSerializer, cfg); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
