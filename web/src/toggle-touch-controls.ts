export function init() {
    document.body.addEventListener("click", async (event) => {
        const button = event.target as HTMLElement;
        if (!button.matches("[data-toggle-touch-controls]")) return;

        const touchControls = document.body.querySelector(
            "[data-show-on-touch]",
        );
        if (!touchControls) return;

        const isHidden = touchControls.classList.contains("hidden");

        if (isHidden) {
            touchControls.classList.remove("hidden");
            button.innerText = "Hide Touch Controls";
        } else {
            touchControls.classList.add("hidden");
            button.innerText = "Show Touch Controls";
        }

        const elementsToHide = document.body.querySelectorAll(
            "[data-hide-on-touch]",
        );
        for (let el of elementsToHide) {
            if (!isHidden && el.classList.contains("hidden")) {
                el.classList.remove("hidden");
            } else if (!el.classList.contains("hidden")) {
                el.classList.add("hidden");
            }
        }

        // Check if API is supported
        if ("vibrate" in navigator) {
            navigator.vibrate(20); // ms
        }
    });
}

init();
