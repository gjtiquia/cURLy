// be wary of functions with the same name, eg. init or main
// empty imports mean it will still be part of bundle while not having any conflicts
import {} from "./copy-button";
import { type Vector2 } from "./vector2";
import * as wasm from "./wasm";

async function initAsync() {
    // TODO : hardcode for now, should be set at runtime
    const size: Vector2 = { X: 4, Y: 4 };

    await wasm.initAsync(size);
}

initAsync();
