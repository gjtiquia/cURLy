import { Vector2 } from "../vector2";

let gridElement: HTMLElement | undefined = undefined;

export function init(): { ok: boolean; error?: string } {
    const el = document.body.querySelector("[data-game-grid]");
    if (!el) return { ok: false, error: "cannot find [data-game-grid]!" };

    gridElement = el as HTMLElement;
    return { ok: true };
}

export function getSize(): Vector2 {
    // set for a pleasant game experience, that should also be supported in iPhone SE simulator
    // given that the game was designed for canvas 20x8 with border thickness 1 (21x9)
    // so (32x16) is a good size in powers of 2
    const size: Vector2 = { X: 32, Y: 12 };
    return size;
}

export function setText(text: string) {
    if (!gridElement) {
        console.error("game: gridElement undefined!");
        return;
    }

    gridElement.innerHTML = text;
}
