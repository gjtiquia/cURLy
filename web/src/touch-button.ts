import { getInputActionId, InputAction } from "./game/input";
import { wasm } from "./wasm/wasm";

export function init() {
    document.body.addEventListener("touchstart", async (event) => {
        const button = event.target as HTMLElement;
        if (!button.matches("[data-touch-button]")) return;

        const action = button.getAttribute("data-touch-button") as InputAction;
        const actionId = getInputActionId(action);

        if (wasm && action != "none") {
            // Check if API is supported
            if ("vibrate" in navigator) {
                navigator.vibrate(20); // ms
            }

            wasm.exports.onInputAction(actionId);
        }
    });
}

init();
