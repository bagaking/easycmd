# easycmd

easycmd is a tool for making cmd's, based on urfave/cli, weakening some concepts and providing some scaffolding capabilities.

## Usage

```go
import (
    "github.com/bagaking/easycmd"
)

func main() {
    if err := easycmd.New("sell").
        Child("gem").Set.
            Alias("g").Usage("sell some gem").End.Action(SellSomeGem).
        Base().Child("wood").Set.
            Alias("w").Usage("sell some wood").End.Action(SellSomeWood).Flags(flagsWood).
        Base().Child("idle").Handler(mainIdleHandler).
        RunBaseAsApp(); err != nil {
        panic(err)
    }
}
```

