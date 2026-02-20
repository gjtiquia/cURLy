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

type Wasm = WebAssembly.Instance & { exports: WasmExports };

let wasm: Wasm | undefined = undefined;

// TODO :
// - two-way communication with go wasm
// - research code-gen to improve DX when adding export/import functions in go
// - extract game logic into a shared package that can be used by cmd/tui and cmd/wasm

const textDecoder = new TextDecoder();

export async function initAsync() {
    // TODO : hardcode for now, should pass in as arg
    const size = { X: 4, Y: 4 };

    const go = new Go();

    // import functions for main.wasm to use
    // TinyGo passes string args as (ptr, len) into linear memory, not as JS strings.
    go.importObject.env = {
        getTermSize: function (ptr: number) {
            if (!wasm) return;

            // directly setting Go's struct fields
            const view = new Int32Array(wasm.exports.memory.buffer, ptr, 2);
            view[0] = size.X;
            view[1] = size.Y;
        },
        notify: function (eventId: number) {
            if (!wasm) return;

            console.log("notify:", eventId);

            // getCanvasCellsAddr() returns the address of the Go slice *header*,
            // not the byte data. Slice header is [ptr: 4 bytes, len: 4 bytes, cap: 4 bytes].
            const sliceAddr = wasm.exports.getCanvasCellsAddr();
            const sliceDataView = new DataView(
                wasm.exports.memory.buffer,
                sliceAddr,
                12,
            );
            const ptr = sliceDataView.getUint32(0, true); // true = little-endian, the least significant byte is stored first, which Go's runtime uses
            const len = sliceDataView.getUint32(4, true);
            const cap = sliceDataView.getUint32(8, true);

            let out = "";
            for (let y = 0; y < size.Y; y++) {
                const rowBytes = new Uint8Array(
                    wasm.exports.memory.buffer,
                    ptr + y * size.X,
                    size.X,
                );
                out += textDecoder.decode(rowBytes);
                out += "\n";
            }
            console.log(out);
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

        wasm = result.instance as Wasm;

        console.log("running main.wasm...");
        const exitCode = await go.run(wasm); // runs main()
        console.log("main.wasm exit code:", exitCode);
    } catch (err) {
        console.error(err);
    }
}
