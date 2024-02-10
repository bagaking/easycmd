# easycmd

easycmd is a comprehensive command-line interface (CLI) toolkit 
simplifying the development of structured command-line applications. 
Relying on urfave/cli, easycmd streamlines CLI creation by abstracting 
away boilerplate code and providing easy-to-use scaffolding facilities.

Whether you are crafting simple tools or complex multi-command apps, 
easycmd enables you to focus on core functionality.

## Features

- Simplifies the definition of command hierarchies.
- Supports aliases and custom flag handling.
- Facilitates middleware usage for common setup functions.
- Extensible through custom builders for individual command actions.
- Well-defined error handling and custom output formatting.

## Installation

To install easycmd, simply run:

```bash
go get github.com/bagaking/easycmd
```

## Quick Start

Below is a small, complete example that defines a `hello` subcommand:

```go
package main

import (
	"fmt"
	"os"

	"github.com/bagaking/easycmd"
	"github.com/urfave/cli/v2"
)

func main() {
	err := easycmd.New("example").
		Set.Usage("Small easycmd example").End.
		Child("hello").
		Set.Alias("hi").Usage("Print a greeting").End.
		Action(func(c *cli.Context) error {
			name := c.Args().First()
			if name == "" {
				name = "world"
			}
			fmt.Printf("hello, %s\n", name)
			return nil
		}).
		RunBaseAsApp()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

Run it with `go run . hello Alice` or the alias `go run . hi Alice`.

## Local validation

```bash
go test ./...
```

## Documentation

For more detailed information, visit urfave/cli documentation as easycmd is a complementary wrapper to simplify using urfave/cli functionalities.

## Contributions

Contributions are welcome! Feel free to submit pull requests to help improve easycmd or create an issue for any bugs or enhancements.
