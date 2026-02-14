import * as ANSI from "./ANSI";

process.on("exit", cleanup); // Regular exit
process.on("SIGINT", cleanupAndExit); // Ctrl-C, does not exit by default, need to manually exit
process.on("SIGTERM", cleanupAndExit); // Terminated by terminal

for (let i = 0; i < 5; i++) {
    const bufferArray = "     ".split("");
    bufferArray[i] = "x";

    const buffer = bufferArray.join("");

    clearAndDrawBuffer(buffer + "\n" + buffer);
    await Bun.sleep(500);
}

function clearAndDrawBuffer(buffer: string) {
    if (ANSI.isANSISupported) {
        ANSI.clearAndDrawBuffer(buffer)
    }
    else {
        console.clear();
        console.log(buffer)
    }
}

function cleanupAndExit() {
    cleanup();
    process.exit(0);
}

function cleanup() {
    if (ANSI.isANSISupported) {
        ANSI.cleanup();
    }
    else {
        console.clear();
    }
}

