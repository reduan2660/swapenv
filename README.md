# swapenv

Swap environment with breeze.

## Quick Start

- install `go install github.com/reduan2660/swapenv@latest` (binary coming soon)
- specify the environment in `.dev.env`, `.stage.env`, ...
- `swapenv load` to load the environments (all of them will disappear from your current directory)
- `swapenv to <environment-name>` to swap environments. e.g.: `swapenv to dev`
- `swapenv ls` to list all the available environments
- `swapenv` to show project staus or current active environment if any

## Whats coming
- Sharing, and cloud sync

## Author
Alve Reduan - [hey@alvereduan.com](mailto:hey@alvereduan.com)
