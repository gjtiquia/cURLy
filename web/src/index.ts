// be wary of functions with the same name, eg. init or main
// empty imports mean it will still be part of bundle while not having any conflicts
import {} from "./copy-button";
import { getMaxCharPerLine } from "./ruler";
import { type Vector2 } from "./vector2";
import * as wasm from "./wasm";

async function initAsync() {
    // set for a pleasant game experience, that should also be supported in iPhone SE simulator
    // given that the game was designed for canvas 20x8 with border thickness 1 (21x9)
    // so (32x16) is a good size in powers of 2
    const size: Vector2 = { X: 32, Y: 16 };

    const maxSizeX = getMaxCharPerLine();
    if (size.X >= maxSizeX) {
        return;
    }

    await wasm.initAsync(size);
}

initAsync();
