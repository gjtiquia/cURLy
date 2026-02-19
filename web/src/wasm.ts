// Go WASM runtime from wasm_exec.js (loaded as script before this module).
declare const Go: new () => {
    importObject: WebAssembly.Imports;
    run(instance: WebAssembly.Instance): Promise<number>; // returns exit code
};

// Exports from our TinyGo WASM module (main.wasm).
interface WasmExports extends WebAssembly.Exports {
    multiply(a: number, b: number): number;
}

export async function init() {
    const go = new Go();

    // import functions for main.wasm to use
    go.importObject.env = {
        add: function(x: number, y: number) {
            return x + y;
        },
    };

    // TODO : how to run exported functions if wasm needs to wait to fetch, and wasm is long running?

    // polyfill if browsers do not support WebAssembly.instantiateStreaming
    if (!WebAssembly.instantiateStreaming) {
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }

    // fetch wasm and run main.wasm
    try {
        const result = await WebAssembly.instantiateStreaming(fetch("/public/main.wasm"), go.importObject)
        const wasm = result.instance as WebAssembly.Instance & { exports: WasmExports };

        console.log("running main.wasm...")
        const exitCode = await go.run(wasm); // runs main()
        console.log("main.wasm exit code:", exitCode)
    } catch (err) {
        console.error(err);
    }
}

init()
