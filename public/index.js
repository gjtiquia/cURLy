// web/src/util.ts
function sleepAsync(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

// web/src/copy-button.ts
listenToCopyButton();
function listenToCopyButton() {
  document.body.addEventListener("click", async (event) => {
    const button = event.target;
    if (!button.matches("[data-copy-button]"))
      return;
    const codeElement = button.parentElement?.querySelector("code");
    if (codeElement) {
      if (button.innerHTML == "Copied!") {
        console.log("[data-copy-button] debounced!");
        return;
      }
      const code = codeElement.innerHTML;
      await navigator.clipboard.writeText(code);
      console.log("[data-copy-button] copied ", code);
      const originalInnerHTML = button.innerHTML;
      button.innerHTML = "Copied!";
      await sleepAsync(3000);
      button.innerHTML = originalInnerHTML;
    }
  });
}
export {
  listenToCopyButton
};
