/** Go WASM runtime from wasm_exec.js (loaded as script before this module). */
declare const Go: new () => {
    importObject: WebAssembly.Imports;
    run(instance: WebAssembly.Instance): Promise<number>; // returns exit code
};

/** Exports from our TinyGo WASM module (main.wasm). */
interface WasmExports extends WebAssembly.Exports {
    multiply(a: number, b: number): number;
}

export function init() {

    console.log("wasm.init")

    const go = new Go();

    // Providing the environment object, used in WebAssembly.instantiateStreaming.
    // This part goes after "const go = new Go();" declaration.
    go.importObject.env = {
        add: function(x: number, y: number) {
            return x + y;
        },
    };

    // polyfill if browsers do not support WebAssembly.instantiateStreaming
    if (!WebAssembly.instantiateStreaming) {
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }

    // fetch wasm
    WebAssembly.instantiateStreaming(fetch("/public/main.wasm"), go.importObject).then((result) => {
        const wasm = result.instance as WebAssembly.Instance & { exports: WasmExports };
        go.run(wasm); // runs main()

        // Calling the multiply function:
        console.log('multiplied two numbers:', wasm.exports.multiply(5, 3));
    }).catch((err) => {
        console.error(err);
    });
}

init()
