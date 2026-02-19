// be wary of functions with the same name, eg. init or main
// empty imports mean it will still be part of bundle while not having any conflicts
import {} from "./copy-button";
import * as wasm from "./wasm";

async function initAsync() {
    await wasm.initAsync();
}

initAsync();
