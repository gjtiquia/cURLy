# cURLy

snake game with cURL

```
 cURLy.gjt.io          
 ----------------------
 |                    |
 |                    |
 |                    |
 |             *      |
 |                    |
 |             0      |
 |   o         o      |
 |   ooooooooooo      |
 ----------------------
 Score: 130            
 Move: WASD; Restart: R
```

```bash
# for Linux/Mac
curl -fsSL cURLy.gjt.io/install.sh | bash

# for Windows 
powershell -c "irm curly.gjt.io/install.ps1 | iex"
```

also available on the web (desktop and mobile) at [cURLy.gjt.io](https://cURLy.gjt.io)

```bash
# can also be installed via go install
go install github.com/gjtiquia/cURLy/cmd/cURLy@latest

# then run the commmand
cURLy
```

## pre-requisites

- [Go](https://go.dev/doc/install)
- [TinyGo](https://tinygo.org/getting-started/install/)
- [Air](https://github.com/air-verse/air?tab=readme-ov-file#installation)
- [Bun](https://bun.sh/)
- [make](https://en.wikipedia.org/wiki/Make_(software))

## tasks

```bash
# run cURLy
go run ./cmd/cURLy

# build cURLy executables
make build/cURly
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

## publish

- build executables
- create new release (and new tag) on GitHub releases
- upload all executables
- publish new release
- install scripts will then download from the latest release
- `go install` command will also download from the latest release tag
