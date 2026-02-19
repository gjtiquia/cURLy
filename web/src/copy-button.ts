import { sleepAsync } from "./util";

listenToCopyButton()
export function listenToCopyButton() {
    // subscribe to the submit event itself, instead of having each element subscribe to the submit event
    // this is because, hx-boost does not do full page reloads, so <script> tags will not be reloaded, and it wont subscribe to new elements
    document.body.addEventListener("click", async (event) => {
        const button = event.target as HTMLElement;
        if (!button.matches("[data-copy-button]")) return;

        const codeElement = button.parentElement?.querySelector("code")
        if (codeElement) {
            if (button.innerHTML == "Copied!") {
                console.log("[data-copy-button] debounced!")
                return;
            }

            const code = codeElement.innerHTML
            await navigator.clipboard.writeText(code);

            console.log("[data-copy-button] copied ", code)

            const originalInnerHTML = button.innerHTML
            button.innerHTML = "Copied!"

            await sleepAsync(3000)

            button.innerHTML = originalInnerHTML
        }
    });
}
