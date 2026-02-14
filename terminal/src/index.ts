import * as ANSI from "./ANSI";

// TODO : for rendering on web, <pre> seems faster than canvas

process.on("exit", cleanup); // Regular exit
process.on("SIGINT", cleanupAndExit); // Ctrl-C, does not exit by default, need to manually exit
process.on("SIGTERM", cleanupAndExit); // Terminated by terminal

// environment config
const MAX_WIDTH = process.stdout.columns;
const MAX_HEIGHT = process.stdout.rows;

// game config
const DEBUG_MODE = (process.env.DEBUG ?? "0") == "1";

const FPS = 10;
const DELTA_TIME_MS = 1000 / FPS

const BORDER_X_THICKNESS = 1;
const BORDER_Y_THICKNESS = 1;

const CANVAS_WIDTH = 40;
const CANVAS_HEIGHT = 10;

const PADDING_Y = Math.floor((MAX_HEIGHT - CANVAS_HEIGHT) / 2);
const PADDING_X = Math.floor((MAX_WIDTH - CANVAS_WIDTH) / 2);

const PADDING_CHAR = " ";
const BORDER_X_CHAR = "-";
const BORDER_Y_CHAR = "|";
const BG_CHAR = " ";
const SNAKE_CHAR = "x";

// init game logic
const snakeHeadPos = [0, CANVAS_HEIGHT / 2]; // x,y

// init draw
const canvas = createCanvas();
debug("canvas width:" + canvas[0]!.length.toString())
debug("canvas height:" + canvas.length.toString())

while (true) {

    // game logic
    snakeHeadPos[0] = (snakeHeadPos[0]! + 1) % CANVAS_WIDTH

    // draw
    resetCanvas(canvas);
    drawChar(canvas, snakeHeadPos[0], snakeHeadPos[1]!, SNAKE_CHAR);

    // render
    const buffer = canvasToStringBuffer(canvas);
    clearAndDrawBuffer(buffer);

    // frame
    // TODO : can see multiplayer book on their suggested architecture
    await Bun.sleep(DELTA_TIME_MS);

    if (DEBUG_MODE)
        break;
}

function createCanvas() {
    const canvas: string[][] = [];

    // upper padding
    for (let y = 0; y < PADDING_Y - BORDER_Y_THICKNESS; y++) {
        const row: string[] = [];
        for (let x = 0; x < MAX_WIDTH; x++) {
            row.push(PADDING_CHAR);
        }
        canvas.push(row);
    }

    // upper border
    for (let y = 0; y < BORDER_Y_THICKNESS; y++) {
        const row: string[] = [];

        // left padding
        for (let x = 0; x < PADDING_X - BORDER_X_THICKNESS; x++)
            row.push(PADDING_CHAR)

        // render border
        for (let x = 0; x < BORDER_X_THICKNESS + CANVAS_WIDTH + BORDER_X_THICKNESS; x++)
            row.push(BORDER_X_CHAR)

        // right padding
        for (let x = 0; x < PADDING_X - BORDER_X_THICKNESS; x++)
            row.push(PADDING_CHAR)

        canvas.push(row);
    }

    // render canvas background
    for (let y = 0; y < CANVAS_HEIGHT; y++) {
        const row: string[] = [];

        // left padding
        for (let x = 0; x < PADDING_X - BORDER_X_THICKNESS; x++)
            row.push(PADDING_CHAR)

        // left border
        for (let x = 0; x < BORDER_X_THICKNESS; x++)
            row.push(BORDER_Y_CHAR)

        // render
        for (let x = 0; x < CANVAS_WIDTH; x++)
            row.push(BG_CHAR);

        // right border
        for (let x = 0; x < BORDER_X_THICKNESS; x++)
            row.push(BORDER_Y_CHAR)

        // right padding
        for (let x = 0; x < PADDING_X - BORDER_X_THICKNESS; x++)
            row.push(PADDING_CHAR)

        canvas.push(row);
    }

    // lower border
    for (let y = 0; y < BORDER_Y_THICKNESS; y++) {
        const row: string[] = [];

        // left padding
        for (let x = 0; x < PADDING_X - BORDER_X_THICKNESS; x++)
            row.push(PADDING_CHAR)

        // render border
        for (let x = 0; x < BORDER_X_THICKNESS + CANVAS_WIDTH + BORDER_X_THICKNESS; x++)
            row.push(BORDER_X_CHAR)

        // right padding
        for (let x = 0; x < PADDING_X - BORDER_X_THICKNESS; x++)
            row.push(PADDING_CHAR)

        canvas.push(row);
    }

    // lower padding
    for (let y = 0; y < PADDING_Y; y++) {
        const row: string[] = [];

        for (let x = 0; x < MAX_WIDTH; x++) {
            row.push(PADDING_CHAR);
        }
        canvas.push(row);
    }

    return canvas;
}

function resetCanvas(canvas: string[][]) {
    for (let y = 0; y < CANVAS_HEIGHT; y++) {
        for (let x = 0; x < CANVAS_WIDTH; x++) {
            drawChar(canvas, x, y, BG_CHAR)
        }
    }
    return canvas;
}

function drawChar(canvas: string[][], x: number, y: number, char: string) {
    canvas[PADDING_Y + y]![PADDING_X + x] = char;
}

function canvasToStringBuffer(canvas: string[][]) {
    let buffer = "";

    for (let row of canvas) {
        buffer += row.join("");
        buffer += "\n"
    }

    return buffer
}

function clearAndDrawBuffer(buffer: string) {
    if (DEBUG_MODE)
        return;

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
    if (DEBUG_MODE)
        return;

    if (ANSI.isANSISupported) {
        ANSI.cleanup();
    }
    else {
        console.clear();
    }
}

function debug(message: string) {
    if (DEBUG_MODE)
        console.log(message)
}
