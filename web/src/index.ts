// be wary of functions with the same name, eg. init or main
// empty imports mean it will still be part of bundle while not having any conflicts
import {} from "./copy-button";
import {} from "./toggle-touch-controls";
import { getMaxCharPerLine } from "./ruler";
import { type Vector2 } from "./vector2";
import * as game from "./game";
import * as wasm from "./wasm";

async function initAsync() {
    const { ok, error } = game.init();
    if (!ok) {
        console.error("initAsync:", error);
        return;
    }

    const size = game.getSize();

    const maxSizeX = getMaxCharPerLine();
    if (size.X >= maxSizeX) {
        console.error(
            "initAsync:",
            "game.size.X",
            size.X,
            ">= maxCharPerLine",
            maxSizeX,
        );
        return;
    }

    await wasm.initAsync(size);
}

initAsync();
