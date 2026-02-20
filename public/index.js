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
var exports = undefined;
async function initAsync() {
  const go = new Go;
  go.importObject.env = {
    getTermSize: function() {
      return { X: 10, Y: 10 };
    },
    notify: function(eventId) {
      console.log("notify:", eventId);
      if (exports) {
        const addr = exports.getCanvasCellsAddr();
        console.log("canvas cells addr:", exports.getCanvasCellsAddr());
        const size = { X: 4, Y: 4 };
        const len = size.X * size.Y;
        const bytes = new Uint8Array(exports.memory.buffer, addr, len);
        console.log("canvas cells bytes:", bytes);
        console.log("canvas cells bytes[0]:", bytes[0]);
      }
    }
  };
  if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
      const source = await (await resp).arrayBuffer();
      return await WebAssembly.instantiate(source, importObject);
    };
  }
  try {
    const result = await WebAssembly.instantiateStreaming(fetch("/public/main.wasm"), go.importObject);
    const wasm = result.instance;
    exports = wasm.exports;
    console.log("running main.wasm...");
    const exitCode = await go.run(wasm);
    console.log("main.wasm exit code:", exitCode);
  } catch (err) {
    console.error(err);
  }
}

// web/src/index.ts
async function initAsync2() {
  await initAsync();
}
initAsync2();
