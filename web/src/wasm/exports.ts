import { wasm } from "./wasm";

const textDecoder = new TextDecoder();

// TODO : probably make a Vector2 type

// functions exported to Go
export function createExports(size: { X: number; Y: number }) {
    return {
        getTermSize: function (ptr: number) {
            if (!wasm) return;

            // directly setting Go's struct fields
            const view = new Int32Array(wasm.exports.memory.buffer, ptr, 2);
            view[0] = size.X;
            view[1] = size.Y;
        },
        notify: function (eventId: number) {
            if (!wasm) return;

            console.log("notify:", eventId);

            // getCanvasCellsAddr() returns the address of the Go slice *header*,
            // not the byte data. Slice header is [ptr: 4 bytes, len: 4 bytes, cap: 4 bytes].
            const sliceAddr = wasm.exports.getCanvasCellsAddr();
            const sliceDataView = new DataView(
                wasm.exports.memory.buffer,
                sliceAddr,
                12,
            );
            const ptr = sliceDataView.getUint32(0, true); // true = little-endian, the least significant byte is stored first, which Go's runtime uses
            const len = sliceDataView.getUint32(4, true);
            const cap = sliceDataView.getUint32(8, true);

            let out = "";
            for (let y = 0; y < size.Y; y++) {
                const rowBytes = new Uint8Array(
                    wasm.exports.memory.buffer,
                    ptr + y * size.X,
                    size.X,
                );
                out += textDecoder.decode(rowBytes);
                out += "\n";
            }
            console.log(out);
        },
    };
}
