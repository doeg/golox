# golox

To build the binaries and make them globally available:

```
go install ./...
golox
```

To build a binary into the local `build/` folder:

```
go build -o build/ ./cmd/golox
./build/golox
```