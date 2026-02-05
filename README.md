# swapenv

swap environment with breeze.

## quick start

- install `go install github.com/reduan2660/swapenv@latest` (binary coming soon)
- specify the environment in `.dev.env`, `.stage.env`, ...
- `swapenv load` to load the environments (to replace existing use --replace, otherwise it'll fast forward if already loaded)
- `swapenv to <environment-name>` to swap environments. e.g.: `swapenv to dev`
- `swapenv ls` to list all the available environments
- `swapenv` to show project staus or current active environment if any
- `swapenv spit` to write all the environment back to .*.env files (use --env to specify a single environment)


## share/receive

e2e encyprted share and sync. the server only carries receiver's public key and encyprted payload.

- `swapenv share` - share current project (all envs, latest version)
  - `--project` - specific project
  - `--env` - specific env only
  - `--version` - specific version (default: latest)

- `swapenv receive` - receive shared environment
  - existing project → new version
  - new project → created (no localPath)

- `swapenv map <project>` - assign current directory to project

### flow
1. device A: `swapenv share` → stream code shown
2. device B: `swapenv receive` → receives & saves
3. device B: `swapenv map myproject` → links to directory
4. device B: `swapenv to dev` → activates env


under the hood, swapenv maintains a versioning, whenever we're loading / receiving new environment it increments the version. we can rename, select, rollback the vesions.
  - each load creates a new version
  - old versions auto-pruned (keeps latest N)
  - named versions are protected from pruning

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

## Integrations

### oh-my-posh

Add to your config's `segments` array:
```json
{
    "type": "command",
        "style": "plain",
        "foreground": "#7dcfff",
        "template": "{{ if .Output }} {{ .Output }}{{ end }}",
        "properties": {
            "command": "swapenv info --format plain --env-only"
        }
}

```

### server
- default: `app.swapenv.sh`
- self-host: [github.com/reduan2660/swapenv-server](github.com/reduan2660/swapenv-server)
- override: `--server <url>` or set `server` in config

## whats coming
- cloud sync

## author
alve reduan - [iam.reduan@gmail.com](mailto:iam.reduan@gmail.com)

## license
this project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
