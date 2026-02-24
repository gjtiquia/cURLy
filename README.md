# cURLy

snake game with cURL

```bash
# for Linux/Mac
curl -fsSL cURLy.gjt.io/install.sh | bash

# for Windows 
powershell -c "irm curly.gjt.io/install.ps1 | iex"
```

also available on the web (desktop and mobile) at [curly.gjt.io](https://curly.gjt.io)

## pre-requisites

- [Go](https://go.dev/doc/install)
- [TinyGo](https://tinygo.org/getting-started/install/)
- [Air](https://github.com/air-verse/air?tab=readme-ov-file#installation)
- [Bun](https://bun.sh/)
- [make](https://en.wikipedia.org/wiki/Make_(software))

## tasks

```bash
# run cURLy tui
go run ./cmd/tui

# build cURLy tui executables
make build/tui
```

```bash
# install dependencies
bun install

# run web server on port 3000 with live reload
# also generates tailwind classes, bundles typescript, builds wasm
air

# build web server
make build/server

# run web server on specific port
PORT=4321 ./bin/server
```
