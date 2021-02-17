# msh

A mini Unix shell.

## Development

```
fswatch -Ee "msh/msh|msh/parser/main.go|msh/coverage.out|msh/.git|msh/README.md" -or . parser | xargs -n1 -I {} ./make {}
```

## Known Issues

*   You can't run `msh` from inside `msh`. It hangs.
*   When you Ctrl-C a running process, it prints two instances of '^C' on the screen.
*   When I try to wrap os.Stdin in a noop closer implementation, `msh` pauses
    for user input after each time it runs a process. (Not sure about builtin
    commands.)
*   If I don't connect os.Stdin to the stdin of subprocesses, then grep doesn't
    behave as expected (grep doesn't detect stdin, so it doesn't hang.) If I do
    connect os.Stdin to the stdin of the subprocesses, then `ls` doesn't behave
    as expected (it makes 1 long list of files, instead of columns.)
*   The `cd` stack sometimes contains relative directories (if you `cd` to a
    relative directory.) Use `cd #` to see a listing of the directories.
