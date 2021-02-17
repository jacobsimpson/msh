# msh

A mini Unix shell.

## Examples

To see the built in commands:

```
/Users/jsimpson/src/msh % help
msh, version 20210217.034931
These shell commands are defined internally. Type `help` to see this list.

cd         cd [dir]
exit       exit [status]
export     export
help       help [pattern]
pwd        pwd
```

Some specialty things about `cd`.

```
/Users/jsimpson/src/msh % cd
/Users/jsimpson % cd src
/Users/jsimpson/src % cd #
  /Users/jsimpson/src/msh
  /Users/jsimpson
* src
/Users/jsimpson/src % cd --
/Users/jsimpson/src/msh % cd #
* /Users/jsimpson/src/msh
  /Users/jsimpson
  src
/Users/jsimpson/src/msh % cd +
/Users/jsimpson % cd #
  /Users/jsimpson/src/msh
* /Users/jsimpson
  src
/Users/jsimpson %
```

A couple basic redirects (`>`, `>&`) and piping (`|`) work.

```
/Users/jsimpson % echo 1 > t
/Users/jsimpson % echo 2 >> t
/Users/jsimpson % echo 3 >> t
/Users/jsimpson % cat t
1
2
3
/Users/jsimpson % cat t | grep 3
3
/Users/jsimpson %
```

## Development

Handy little command line that rebuilds the project on each change.

```
fswatch -Ee "$(basename `pwd`)/$(basename `pwd`)|$(basename `pwd`)/parser/grammar.go|$(basename `pwd`)/coverage.out|$(basename `pwd`)/.git|$(basename `pwd`)/README.md" -or . parser | xargs -n1 -I {} ./make {}
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
