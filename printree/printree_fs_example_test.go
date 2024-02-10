package printree_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bagaking/easycmd/printree"
)

type fsNode struct {
	path     string
	info     os.FileInfo
	children []*fsNode
}

func (n *fsNode) ForEachPrintableChild(fn func(child printree.IPrintableTreeNode, ind int) error) error {
	for i, child := range n.children {
		if err := fn(child, i); err != nil {
			return err
		}
	}
	return nil
}

func (n *fsNode) HierarchyBeginWord() string {
	if n.info.IsDir() {
		mode := n.info.Mode()
		return fmt.Sprintf("[D] %s ", mode.Perm()) // 格式化文件夹类型和权限信息
	} else {
		mode := n.info.Mode()
		return fmt.Sprintf("[F] %s ", mode.Perm()) // 格式化文件夹类型和权限信息
	}
	return ""
}

func (n *fsNode) Content() string {
	return filepath.Base(n.path) // 只展示文件或文件夹的基本名字
}

func (n *fsNode) ChildrenCount() int {
	return len(n.children)
}

// scanDir 递归扫描目录并构造 fsNode 树
func scanDir(path string) (*fsNode, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	node := &fsNode{
		path: path,
		info: fileInfo,
	}

	if fileInfo.IsDir() {
		files, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, fi := range files {
			childPath := filepath.Join(path, fi.Name())
			childNode, err := scanDir(childPath) // 递归构建子节点
			if err != nil {
				return nil, err
			}
			node.children = append(node.children, childNode)
		}
	}

	return node, nil
}

// TestPrintFSTree 测试使用 printree 打印文件系统树
func TestPrintFSTree(t *testing.T) {
	// 假设我们扫描当前目录
	rootPath, _ := os.Getwd()

	// 构建文件系统树
	fsTree, err := scanDir(rootPath)
	if err != nil {
		t.Fatal(err)
	}

	// 配置打印功能
	printFn := func(a ...interface{}) (n int, err error) {
		fmt.Println(a...)
		return len(a), nil
	}

	printer := printree.PrintableNodeSerializer(func(node printree.IPrintableTreeNode, ind int) string {
		return node.(*fsNode).Content()
	})

	// 使用 printree 打印文件系统树
	err = printer.RecursivePrint(fsTree, printree.OptCustomPrinter(printFn))
	if err != nil {
		t.Fatal(err)
	}
}
