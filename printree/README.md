# Printree - Pretty Tree Printer for Go

Printree is a versatile library for Go developers that allows for pretty-printing tree-like data structures in a format that's easy to read and understand. It's designed to be highly customizable and flexible, catering to a wide array of use cases where hierarchical data needs to be visualized.

## Features

- **Customizable Marks**: Configure the tree symbols to your preference (`LineItem`, `LastItem`, `Cross`, `Blank`).
- **Conditional Drill-Down**: Use predicates to selectively expand or collapse nodes in the tree.
- **Custom Printers**: Use your own print function if `fmt.Println` doesn't suit your needs.
- **Cross-Platform**: Works across different platforms with no dependencies on specific CLI utilities.

## Installation

```sh
go get github.com/bagaking/easycmd/printree
```

## Quick Start

Here is a quick example to print the file structure of a directory:

```go
package main

import (
    "fmt"
    "os"
	
    "github.com/bagaking/easycmd/printree"
)

func main() {
    // Construct the tree structure (fsNode) from your filesystem
    root, err := scanDir("/path/to/your/directory")
    if err != nil {
        panic(err)
    }

    // Print the filesystem tree
    printer := printree.PrintableNodeSerializer(func(node printree.IPrintableTreeNode, ind int) string {
        return node.(*fsNode).Content()
    })

    err = printer.RecursivePrint(root, printree.OptCustomPrinter(func(a ...interface{}) (n int, err error) {
        fmt.Println(a...)
        return len(a), nil
    }))
    if err != nil {
        panic(err)
    }
}
```

## Usage Scenarios

- **CLI Tools**: For command-line tools that need to present hierarchical data, such as task organizers or project structure visualizers. Example:

  Example from [./printree_fs_example_test.go](./printree_fs_example_test.go)
  ```go
  // This will print the file structure of the current directory
  rootPath, _ := os.Getwd()
  fsTree, _ := scanDir(rootPath)
  printTree(fsTree)
  ```

- **Debugging**: When you're trying to understand a complex nested data structure during development. Example:

  ```go
  // Example usage in a debugging scenario
  testTree := createTestTree() // function from printree_test.go
  printTree(testTree)
  ```

- **Testing Frameworks**: For a richer and more informative testing output when evaluating nested data structures. Example:

  ```go
  // Example usage in a testing framework
  TestPrintFSTree() // function from printree_fs_example_test.go
  ```

- **File Systems**: To visualize directory trees with permissions and types, similar to `tree` command on Unix. Example:

  ```go
  // Example from printree_fs_example_test.go
  // This prints a file system tree starting at the specified root path
  root, _ := scanDir("/path/to/root")
  printTree(root)
  ```

## API Documentation

For full API documentation and more detailed examples, visit the GoDoc page: [Printree GoDoc](https://pkg.go.dev/github.com/bagaking/easycmd/printree)

## Contribution

We encourage contributions! Feel free to submit pull requests or create issues for new features, bug fixes, or documentation enhancements.

## License

Printree is licensed under the MIT License. Detailed information can be found in the [LICENSE](LICENSE) file.

Happy tree printing!