// web/src/util.ts
function sleepAsync(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

// web/src/copy-button.ts
function init() {
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
init();
// web/src/wasm.ts
function init2() {
  console.log("wasm.init");
  const go = new Go;
  go.importObject.env = {
    add: function(x, y) {
      return x + y;
    }
  };
  if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
      const source = await (await resp).arrayBuffer();
      return await WebAssembly.instantiate(source, importObject);
    };
  }
  WebAssembly.instantiateStreaming(fetch("/public/main.wasm"), go.importObject).then((result) => {
    const wasm = result.instance;
    go.run(wasm);
    console.log("multiplied two numbers:", wasm.exports.multiply(5, 3));
  }).catch((err) => {
    console.error(err);
  });
}
init2();
