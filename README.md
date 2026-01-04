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

Under the hood, swapenv maintains a versioning, whenever we're loading / receiving new environment it increments the version. we can rename, select, rollback the vesions.
  - Each load creates a new version
  - Old versions auto-pruned (keeps latest N)
  - Named versions are protected from pruning

  - `swapenv version` - show current & latest version
  - `swapenv version <n>` - switch to version n
  - `swapenv version latest` - switch to latest
  - `swapenv version ls` - list all versions
  - `swapenv version rename <n> <name>` - name a version (protects from auto-delete)
  - `swapenv version rollback [steps]` - go back n versions (default 1)

  Flags:
  - `swapenv to <env> --version <n|name|latest>` - use specific version
  - `swapenv spit --version <n|name|latest>` - spit from specific version
  - `swapenv ls -v` - show versions alongside envs

  Config:
  - max_versions: 5 - how many versions to keep (default 5)

## Whats coming
- Sharing, and cloud sync

## Author
Alve Reduan - [hey@alvereduan.com](mailto:hey@alvereduan.com)

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
