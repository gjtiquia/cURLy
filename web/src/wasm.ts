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

const textDecoder = new TextDecoder();

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
                // TODO : hardcode for now
                const size = { X: 4, Y: 4 };

                // getCanvasCellsAddr() returns the address of the Go slice *header*,
                // not the byte data. Slice header is [ptr: 4 bytes, len: 4 bytes, cap: 4 bytes].
                const sliceAddr = exports.getCanvasCellsAddr();
                const sliceDataView = new DataView(
                    exports.memory.buffer,
                    sliceAddr,
                    12,
                );
                const ptr = sliceDataView.getUint32(0, true); // true = little-endian, the least significant byte is stored first, which Go's runtime uses
                const len = sliceDataView.getUint32(4, true);
                const cap = sliceDataView.getUint32(8, true);

                let out = "";
                for (let y = 0; y < size.Y; y++) {
                    const rowBytes = new Uint8Array(
                        exports.memory.buffer,
                        ptr + y * size.X,
                        size.X,
                    );
                    out += textDecoder.decode(rowBytes);
                    out += "\n";
                }
                console.log(out);

                // console.log(decodeString(ptr, len));
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
