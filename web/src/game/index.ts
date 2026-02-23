import { Vector2 } from "../vector2";
import { subscribeToKeyDownEvent } from "./input";

let gridElement: HTMLElement | undefined = undefined;

export function init(): { ok: boolean; error?: string } {
    const el = document.body.querySelector("[data-game-grid]");
    if (!el) return { ok: false, error: "cannot find [data-game-grid]!" };

    gridElement = el as HTMLElement;

    subscribeToKeyDownEvent();

    return { ok: true };
}

export function getSize(): Vector2 {
    // set for a pleasant game experience, that should also be supported in iPhone SE simulator
    // given that the game was designed for canvas 20x8 with border thickness 1 (21x10), and a header + message + footer (21x13)
    // so 32 is chosen as width for power of 2, 14 is chosen just a bit more than the min height (13)
    const size: Vector2 = { X: 32, Y: 14 };
    return size;
}

export function setText(text: string) {
    if (!gridElement) {
        console.error("game: gridElement undefined!");
        return;
    }

    gridElement.innerHTML = text;
}
