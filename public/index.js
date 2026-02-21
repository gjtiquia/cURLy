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

// web/src/ruler.ts
function getMaxCharPerLine() {
  const ruler = document.body.querySelector("[data-ruler]");
  if (!ruler) {
    console.error("initAsync:", "cannot find [data-ruler]!", "returning 0...");
    return 0;
  }
  const rect = ruler.getBoundingClientRect();
  const text = ruler.innerHTML;
  const charWidth = rect.width / text.length;
  const maxCharCount = window.innerWidth / charWidth;
  return maxCharCount;
}

// web/src/wasm/exports.ts
var textDecoder = new TextDecoder;
function createExports(size) {
  return {
    getTermSize: function(ptr) {
      if (!wasm)
        return;
      const view = new Int32Array(wasm.exports.memory.buffer, ptr, 2);
      view[0] = size.X;
      view[1] = size.Y;
    },
    notify: function(eventId) {
      if (!wasm)
        return;
      const slicePtr = wasm.exports.getCanvasCellsPtr();
      const sliceHeader = new Uint32Array(wasm.exports.memory.buffer, slicePtr, 3);
      const ptr = sliceHeader[0];
      const len = sliceHeader[1];
      const cap = sliceHeader[2];
      let text = "";
      for (let y = 0;y < size.Y; y++) {
        const rowBytes = new Uint8Array(wasm.exports.memory.buffer, ptr + y * size.X, size.X);
        text += textDecoder.decode(rowBytes);
        text += `
`;
      }
      setText(text);
    }
  };
}

// web/src/wasm/wasm.ts
var wasm = undefined;
async function initAsync(size) {
  const go = new Go;
  go.importObject.env = createExports(size);
  if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
      const source = await (await resp).arrayBuffer();
      return await WebAssembly.instantiate(source, importObject);
    };
  }
  try {
    const result = await WebAssembly.instantiateStreaming(fetch("/public/main.wasm"), go.importObject);
    wasm = result.instance;
    console.log("running main.wasm...");
    const exitCode = await go.run(wasm);
    console.log("main.wasm exit code:", exitCode);
  } catch (err) {
    console.error("wasm.initAsync: error");
    console.error(err);
  }
}

// web/src/game/input.ts
function subscribeToKeyDownEvent() {
  document.addEventListener("keydown", (e) => {
    const action = mapCodeToInputAction(e.code);
    const actionId = getInputActionId(action);
    if (wasm && action != "none") {
      wasm.exports.onInputAction(actionId);
    }
  });
}
function mapCodeToInputAction(code) {
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
function getInputActionId(action) {
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

// web/src/game/index.ts
var gridElement = undefined;
function init2() {
  const el = document.body.querySelector("[data-game-grid]");
  if (!el)
    return { ok: false, error: "cannot find [data-game-grid]!" };
  gridElement = el;
  subscribeToKeyDownEvent();
  return { ok: true };
}
function getSize() {
  const size = { X: 32, Y: 12 };
  return size;
}
function setText(text) {
  if (!gridElement) {
    console.error("game: gridElement undefined!");
    return;
  }
  gridElement.innerHTML = text;
}
// web/src/index.ts
async function initAsync2() {
  const { ok, error } = init2();
  if (!ok) {
    console.error("initAsync:", error);
    return;
  }
  const size = getSize();
  const maxSizeX = getMaxCharPerLine();
  if (size.X >= maxSizeX) {
    console.error("initAsync:", "game.size.X", size.X, ">= maxCharPerLine", maxSizeX);
    return;
  }
  await initAsync(size);
}
initAsync2();
