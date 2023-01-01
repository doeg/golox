# golox

golox is a tree-walk interpreter for the [Lox programming language](https://github.com/munificent/craftinginterpreters).

## Local development

To start a REPL:

```
make repl
```

To build the binaries and make them globally available:

```
make install
golox
```

To build a binary into the local `build/` folder:

```
make build
./build/golox
```

To run the tests:

```
make test
```

To regenerate the generated files (e.g., `ast.go`): 

```
make generate
```
