# boilerplate-compose

Minimal Golang CLI for reading a boilerplate compose file with include/extends and variable files.

Usage:

```sh
boilerplate-compose -f boilerplate-compose.yaml --var-file vars.yml --var Title=Boilerplate --print-config
```

Compose file example (`boilerplate-compose.yaml`):

```yaml
include:
  - common.yml

extends:
  - base.yml

vars:
  app: myapp

templates:
  app:
    template_url: ./path/to/templates
    output_folder: ./out
    missing_key_action: error
    vars:
      Title: "My App"
      ShowLogo: false
```

Notes:
- include and extends accept local file paths. They are loaded before the current file and merged with last-wins semantics.
- Variables from `--var-file` and `--var` are merged into top-level `vars` after file resolution; later values override earlier ones.
- Use `--print-config` to output the fully resolved YAML.