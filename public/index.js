// web/src/copy-button.ts
listenToCopyButton();
function listenToCopyButton() {
  document.body.addEventListener("click", async (event) => {
    const target = event.target;
    if (!target.matches("[data-copy-button]"))
      return;
    const codeElement = target.parentElement?.querySelector("code");
    if (codeElement) {
      const code = codeElement.innerHTML;
      await navigator.clipboard.writeText(code);
      console.log("[data-copy-button] copied ", code);
    }
  });
}
export {
  listenToCopyButton
};
