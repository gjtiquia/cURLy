import * as ANSI from "./ANSI";

// TODO : for rendering on web, <pre> seems faster than canvas

process.on("exit", cleanup); // Regular exit
process.on("SIGINT", cleanupAndExit); // Ctrl-C, does not exit by default, need to manually exit
process.on("SIGTERM", cleanupAndExit); // Terminated by terminal

// environment config
const MAX_WIDTH = process.stdout.columns;
const MAX_HEIGHT = process.stdout.rows;

// game config
const FPS = 10;
const DELTA_TIME_MS = 1000 / FPS

const WIDTH = 40;
const HEIGHT = 10;

const PADDING_CHAR = " ";
const BG_CHAR = " ";
const SNAKE_CHAR = "x";

// init game logic
const snakeHeadPos = [0, HEIGHT / 2]; // x,y

// init draw
const canvas = createCanvas();

while (true) {

    // game logic
    snakeHeadPos[0] = (snakeHeadPos[0]! + 1) % WIDTH

    // draw
    resetCanvas(canvas);
    drawChar(canvas, snakeHeadPos[0], snakeHeadPos[1]!, SNAKE_CHAR);

    // render
    const buffer = canvasToStringBuffer(canvas);
    clearAndDrawBuffer(buffer);

    // frame
    // TODO : can see multiplayer book on their suggested architecture
    await Bun.sleep(DELTA_TIME_MS);
}

function createCanvas() {
    const canvas: string[][] = [];
    for (let y = 0; y < HEIGHT; y++) {
        const row: string[] = [];
        for (let x = 0; x < WIDTH; x++) {
            row.push(BG_CHAR)
        }
        canvas.push(row);
    }
    return canvas;
}

function resetCanvas(canvas: string[][]) {
    for (let y = 0; y < canvas.length; y++) {
        for (let x = 0; x < canvas[y]!.length; x++) {
            canvas[y]![x] = BG_CHAR
        }
    }
    return canvas;
}

function drawChar(canvas: string[][], x: number, y: number, char: string) {
    canvas[y]![x] = char;
}

function canvasToStringBuffer(canvas: string[][]) {
    let buffer = "";

    const yPadding = Math.floor((MAX_HEIGHT - canvas.length) / 2);
    const xPadding = Math.floor((MAX_WIDTH - canvas[0]!.length) / 2);

    // upper padding
    for (let y = 0; y < yPadding; y++) {
        for (let x = 0; x < MAX_WIDTH; x++) {
            buffer += PADDING_CHAR;
        }
        buffer += "\n"
    }

    for (let row of canvas) {

        // left padding
        for (let x = 0; x < xPadding; x++)
            buffer += PADDING_CHAR

        // render
        buffer += row.join("");

        // right padding
        for (let x = 0; x < xPadding; x++)
            buffer += PADDING_CHAR

        buffer += "\n"
    }

    // lower padding
    for (let y = 0; y < yPadding; y++) {
        for (let x = 0; x < MAX_WIDTH; x++) {
            buffer += PADDING_CHAR;
        }
        buffer += "\n"
    }

    return buffer
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

