// Go WASM runtime from wasm_exec.js (loaded as script before this module).
declare const Go: new () => {
    importObject: WebAssembly.Imports;
    run(instance: WebAssembly.Instance): Promise<number>; // returns exit code
};

// Exports from our TinyGo WASM module (main.wasm).
interface WasmExports extends WebAssembly.Exports {
    memory: WebAssembly.Memory;
    getCanvasCellsAddr(): number;
}

let exports: WasmExports | undefined = undefined;

// TODO :
// - two-way communication using syscall/js in cmd/wasm/main.go
//   - most worth seeing is the Invoke example, but it sets a global function accessed via `window`, perhaps thats the way...?
// - research code-gen to improve DX when adding export/import functions in go
// - extract game logic into a shared package that can be used by cmd/tui and cmd/wasm

export async function initAsync() {
    const go = new Go();

    // import functions for main.wasm to use
    // TinyGo passes string args as (ptr, len) into linear memory, not as JS strings.
    go.importObject.env = {
        // TODO : work this out later
        getTermSize: function () {
            return { X: 10, Y: 10 };
        },
        notify: function (eventId: number) {
            console.log("notify:", eventId);

            if (exports) {
                const addr = exports.getCanvasCellsAddr();
                console.log("canvas cells addr:", exports.getCanvasCellsAddr());

                // TODO : hardcode for now
                const size = { X: 4, Y: 4 };
                const len = size.X * size.Y;

                const bytes = new Uint8Array(exports.memory.buffer, addr, len);
                console.log("canvas cells bytes:", bytes);
                console.log("canvas cells bytes[0]:", bytes[0]);
            }
        },
    };

    // polyfill if browsers do not support WebAssembly.instantiateStreaming
    if (!WebAssembly.instantiateStreaming) {
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }

    // fetch wasm and run main.wasm
    try {
        const result = await WebAssembly.instantiateStreaming(
            fetch("/public/main.wasm"),
            go.importObject,
        );

        const wasm = result.instance as WebAssembly.Instance & {
            exports: WasmExports;
        };

        exports = wasm.exports;

        console.log("running main.wasm...");
        const exitCode = await go.run(wasm); // runs main()
        console.log("main.wasm exit code:", exitCode);
    } catch (err) {
        console.error(err);
    }
}

function decodeString(ptr: number, len: number): string {
    if (!exports) return `<no memory: ${ptr}, ${len}>`;

    return new TextDecoder().decode(
        new Uint8Array(exports.memory.buffer, ptr, len),
    );
}
