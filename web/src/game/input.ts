import { wasm } from "../wasm/wasm";

type InputAction = "none" | "up" | "down" | "left" | "right" | "restart";

export function subscribeToKeyDownEvent() {
    document.addEventListener("keydown", (e) => {
        const action = mapCodeToInputAction(e.code);
        const actionId = getInputActionId(action);

        if (wasm && action != "none") {
            // console.log("js actionId", actionId);
            wasm.exports.onInputAction(actionId);
        }
    });
}

function mapCodeToInputAction(code: string): InputAction {
    switch (code) {
        case "KeyW":
        case "ArrowUp":
        case "KeyK":
            return "up";

        case "KeyS":
        case "ArrowDown":
        case "KeyJ":
            return "down";

        case "KeyA":
        case "ArrowLeft":
        case "KeyH":
            return "left";

        case "KeyD":
        case "ArrowRight":
        case "KeyL":
            return "right";

        case "KeyR":
            return "restart";

        default:
            return "none";
    }
}

function getInputActionId(action: InputAction): number {
    switch (action) {
        case "up":
            return 1;
        case "down":
            return 2;
        case "left":
            return 3;
        case "right":
            return 4;
        case "restart":
            return 5;
        case "none":
            return 0;
    }
}
