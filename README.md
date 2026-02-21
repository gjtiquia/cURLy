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
- [TinyGo](https://tinygo.org/getting-started/install/)
- [Air](https://github.com/air-verse/air?tab=readme-ov-file#installation)
- [Bun](https://bun.sh/)

## tasks

```bash
# run cURLy tui
go run ./cmd/tui

# build cURLy tui executables
go run ./cmd/build
```

```bash
# install dependencies
bun install

# run web server on port 3000 with live reload
# also generates tailwind classes, bundles typescript, builds wasm
air

# build tailwind classes
bunx @tailwindcss/cli -i ./web/input.css -o ./public/styles.css

# bundle TypeScript scripts
bun build ./web/src/index.ts --outdir=./public

# build wasm
GOOS=js GOARCH=wasm tinygo build -o ./public/main.wasm ./web/wasm

# build web server
go build -o ./bin/server ./cmd/server

# run web server on specific port
PORT=4321 ./bin/server
```

## todos
- [ ] Makefile with build command, including TypeScript typecheck, and minified outputs
