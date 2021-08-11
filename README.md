# maschinenkonzept
Bot for vk.com

## How to write modules?
1. Create new go package for your module
2. Import core  
```go
import "github.com/beshenkaD/maschinenkonzept/core"
```
3. Write your command/hooks/ticks
4. Register it in init func
```go
func init() {
    core.RegisterCommand(...)
    core.RegisterHook(...)
    core.RegisterTick(...)
}
```
5. Import your module in main.go
```go
import (
    _ path/to/module
)
```
## Dependencies
For `quote` module

Droid font
```sh
sudo pacman -S ttf-droid              # arch
sudo apt install fonts-droid-fallback # debian
```
