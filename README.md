# swapenv

swap environment with breeze.

## Quick Start

- install `go install github.com/reduan2660/swapenv@latest` (binary coming soon)
- specify the environment in `.dev.env`, `.stage.env`, ...
- `swapenv load` to load the environments (to replace existing use --replace, otherwise it'll fast forward if already loaded)
- `swapenv to <environment-name>` to swap environments. e.g.: `swapenv to dev`
- `swapenv ls` to list all the available environments
- `swapenv` to show project staus or current active environment if any
- `swapenv spit` to write all the environment back to .*.env files (use --env to specify a single environment)

## Whats coming
- Sharing, and cloud sync

## Author
Alve Reduan - [hey@alvereduan.com](mailto:hey@alvereduan.com)

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
