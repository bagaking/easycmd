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

``bash
go get github.com/bagaking/easycmd
``

## Quick Start

Below is a concise example displaying how you can set up a CLI application that manages inventory with commands to sell different items:

```go
package main

import (
    "github.com/bagaking/easycmd"
    "github.com/urfave/cli/v2"
)

// Define actions for different commands
func SellGem(c *cli.Context) error {
// Your code to sell a gem.
return nil
}

func SellWood(c *cli.Context) error {
// Your code to sell wood.
return nil
}

func main() {
    app := easycmd.New("inventory").
        Child("sell").Set.
        Alias("s").Usage("Commands to sell items").End.
        Base().Child("gem").Set.
        Alias("g").Usage("Sell some gems").End.Action(SellGem).
        Base().Child("wood").Set.
        Alias("w").Usage("Sell some wood").End.Action(SellWood).Flags(flagsWood).
		Base().Child("idle").Handler(mainIdleHandler).
		RunBaseAsApp(); 

    if err := app.RunBaseAsApp(); err != nil {
        panic(err)
    }
}
```

In this example, we created a CLI with an `inventory` command followed by `sell` command having subcommands `gem` and `wood` each with their specific actions and usages.

## Documentation

For more detailed information, visit urfave/cli documentation as easycmd is a complementary wrapper to simplify using urfave/cli functionalities.

## Contributions

Contributions are welcome! Feel free to submit pull requests to help improve easycmd or create an issue for any bugs or enhancements.