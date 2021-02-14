# msh

A mini Unix shell.

## Development

```
fswatch -Ee "msh/msh|msh/parser/main.go|msh/coverage.out|msh/.git|msh/README.md" -or . parser | xargs -n1 -I {} ./make {}
```
