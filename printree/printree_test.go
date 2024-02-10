package printree_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bagaking/easycmd/printree"
)

// mockPrintableNode 是 IPrintableTreeNode 的一个模拟实现，用于测试
type mockPrintableNode struct {
	prefix     string
	children   []printree.IPrintableTreeNode
	linePrefix string
	content    string
}

func (m *mockPrintableNode) ForEachPrintableChild(fn func(child printree.IPrintableTreeNode, ind int) error) error {
	for i, child := range m.children {
		if err := fn(child, i); err != nil {
			return err
		}
	}
	return nil
}

func (m *mockPrintableNode) HierarchyBeginWord() string {
	return m.linePrefix
}

func (m *mockPrintableNode) Content() string {
	return m.content
}

func (m *mockPrintableNode) ChildrenCount() int {
	return len(m.children)
}

// createTestTree 创建一个用于测试的树结构
func createTestTree() printree.IPrintableTreeNode {
	return &mockPrintableNode{
		content: "Root",
		children: []printree.IPrintableTreeNode{
			&mockPrintableNode{
				content: "Child1",
				children: []printree.IPrintableTreeNode{
					&mockPrintableNode{
						content:  "Grandchild1",
						children: []printree.IPrintableTreeNode{},
					},
				},
			},
			&mockPrintableNode{
				content: "Child2",
				children: []printree.IPrintableTreeNode{
					&mockPrintableNode{
						content:  "Grandchild2",
						children: []printree.IPrintableTreeNode{},
					},
				},
			},
		},
	}
}

func TestPrintTree(t *testing.T) {
	root := createTestTree()

	var output strings.Builder
	printer := printree.PrintableNodeSerializer(func(node printree.IPrintableTreeNode, ind int) string {
		return node.(*mockPrintableNode).Content()
	})

	printFn := func(a ...interface{}) (n int, err error) {
		return output.WriteString(a[0].(string) + "\n")
	}
	marks := printree.TreeMarks{
		LineItem: "├─",
		LastItem: "└─",
		Cross:    "│ ",
		Blank:    "  ",
	}

	// 执行打印函数
	err := printer.RecursivePrint(root,
		printree.OptSetMarks(marks),
		printree.OptCustomPrinter(printFn),
	)
	result := output.String()

	assert.Nil(t, err, "error should be nil")

	// 确定每行的输出
	expectedOutput := strings.Join([]string{
		"Root",
		"├─ Child1",
		"│ └─ Grandchild1",
		"└─ Child2",
		"  └─ Grandchild2",
		"",
	}, "\n")

	assert.Equal(t, expectedOutput, result, "output should be equal")
}
