build/cURLy:
	go run ./cmd/build

build/server:
	bunx @tailwindcss/cli -i ./web/input.css -o ./public/styles.css
	bunx tsc --noEmit
	bun build ./web/src/index.ts --outdir=./public
	GOOS=js GOARCH=wasm tinygo build -o ./public/main.wasm ./web/wasm
	go build -o ./bin/server ./cmd/server

