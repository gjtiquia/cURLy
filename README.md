# cURLy

snake game with cURL

```bash
# for Linux/Mac
curl -fsSL cURLy.gjt.io/install.sh | bash

# for Windows 
powershell -c "irm curly.gjt.io/install.ps1 | iex"
```

## pre-requisites

- [Go](https://go.dev/doc/install)
- [Air](https://github.com/air-verse/air?tab=readme-ov-file#installation)

## tasks

```bash
# run cURLy tui
go run ./cmd/tui

# build cURLy tui executables
go run ./cmd/build
```

```bash
# run web server (port 3000 by default)
go run ./cmd/server

# run web server on specific port
PORT=4321 go run ./cmd/server

# run web server on port 3000 with live reload
air

# build web server
go build -o ./bin/server ./cmd/server

# run web server on specific port
PORT=4321 ./bin/server
```
