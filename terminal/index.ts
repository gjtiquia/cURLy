const isANSISupported = true;
// const isANSISupported = false; // ANSI support is not guaranteed in Windows

const RESET = "\x1b[0m";
const HIDE_CURSOR = "\x1b[?25l";
const SHOW_CURSOR = "\x1b[?25h";
const CURSOR_HOME = "\x1b[H";
const CLEAR = "\x1b[2J";

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
    if (isANSISupported) {
        // hides the cursor and does a single write, preventing flicker
        process.stdout.write(CLEAR + HIDE_CURSOR + CURSOR_HOME + buffer + RESET);
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
    if (isANSISupported) {
        process.stdout.write(CLEAR + CURSOR_HOME + SHOW_CURSOR);
    }
    else {
        console.clear();
    }
}

