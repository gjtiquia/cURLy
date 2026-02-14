export const isANSISupported = true;
// export const isANSISupported = false; // ANSI support is not guaranteed in Windows

const RESET = "\x1b[0m";
const HIDE_CURSOR = "\x1b[?25l";
const SHOW_CURSOR = "\x1b[?25h";
const CURSOR_HOME = "\x1b[H";
const CLEAR = "\x1b[2J";

export function clearAndDrawBuffer(buffer: string) {
    if (!isANSISupported)
        return;

    process.stdout.write(CLEAR + HIDE_CURSOR + CURSOR_HOME + buffer + RESET);
}

export function cleanup() {
    if (!isANSISupported)
        return;

    process.stdout.write(CLEAR + CURSOR_HOME + SHOW_CURSOR);
}

