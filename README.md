# commentcov
pluggable comment coverage generator

## How to use

Place `.commentcov.yaml`, then run commentcov
```bash
$ cat .commentcov.yaml
plugins:
  - extension: .go
    install_command: go install github.com/terakoya76/commentcov-plugin-go@latest
    execute_command: commentcov-plugin-go
target_path: .
exclude_pathes:
  - ./vendor/**/**
mode: file_scope
```
