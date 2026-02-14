// https://bun.com/docs/bundler/executables

await Bun.$`rm -rf ./dist`

console.log("Compiling for different platforms...")

const targets = ["bun-linux-x64", "bun-linux-arm64", "bun-windows-x64", "bun-darwin-x64", "bun-darwin-arm64", "bun-linux-x64-musl", "bun-linux-arm64-musl"] as const

for (let target of targets) {
    console.log("Compiling... target:", target)

    const outfile = `./dist/cURLy_${target.replace("bun-", "")}`;
    await Bun.build({
        entrypoints: ["./src/index.ts"],
        compile: { target, outfile },
    });

    console.log("Success! outfile:", outfile)
}

