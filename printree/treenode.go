package printree

import "fmt"

type (
	// IPrintableTreeNode is the interface should be implemented
	// by which node of the tree to be printed.
	IPrintableTreeNode interface {

		// ForEachPrintableChild is the method to iterating over direct child nodes
		ForEachPrintableChild(func(child IPrintableTreeNode, ind int) error) error

		// HierarchyBeginWord returns the information that needs to be printed
		// in the prefix for each line of the nodes directly children level
		HierarchyBeginWord() string

		// ChildrenCount returns the count of direct child nodes
		ChildrenCount() int
	}

	PrintableNodeSerializer func(node IPrintableTreeNode, ind int) string
)

func (nodeSerializer PrintableNodeSerializer) mkLine(st IPrintableTreeNode, prefix string, ind int) string {
	return nodeSerializer.mkMark(st, prefix, nodeSerializer(st, ind))
}

func (nodeSerializer PrintableNodeSerializer) mkMark(st IPrintableTreeNode, prefix string, mark string) string {
	head := st.HierarchyBeginWord() + prefix
	if head == "" {
		return mark
	}

	return fmt.Sprintf("%s %s", head, mark)
}

func (nodeSerializer PrintableNodeSerializer) RecursivePrint(st IPrintableTreeNode, opts ...Option) error {
	setting := (&SerializeSetting{}).pipe(opts...)
	canContinue, err := printRow(st, "", 0, nodeSerializer, setting)
	if err != nil {
		return err
	}
	if !canContinue {
		return nil
	}

	return st.ForEachPrintableChild(func(child IPrintableTreeNode, ind int) error {
		return recursivePrint(child, "", ind, st.ChildrenCount(), nodeSerializer, setting)
	})
}
