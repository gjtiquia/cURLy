import * as ansi from "./ansi";

// TODO : for rendering on web, <pre> seems faster than canvas

// setup
setup();

// environment config
const MAX_WIDTH = process.stdout.columns;
const MAX_HEIGHT = process.stdout.rows;

// game config
const DEBUG_MODE = (process.env.DEBUG ?? "0") == "1"; // TODO : ehhh might be better that logs save to a log file...? and can see in real time the tail of the logs?

const FPS = 10;
const DELTA_TIME_MS = 1000 / FPS

const BORDER_X_THICKNESS = 1;
const BORDER_Y_THICKNESS = 1;

const CANVAS_WIDTH = 40;
const CANVAS_HEIGHT = 10;

const PADDING_Y = Math.floor((MAX_HEIGHT - CANVAS_HEIGHT) / 2);
const PADDING_X = Math.floor((MAX_WIDTH - CANVAS_WIDTH) / 2);

// game config - input // TODO : change to WASD default, but allow config override
const UP_KEY_VAR_1 = "e"
const DOWN_KEY_VAR_1 = "d"
const LEFT_KEY_VAR_1 = "s"
const RIGHT_KEY_VAR_1 = "f"

const UP_KEY_VAR_2 = "\u001b[A" // arrow up
const DOWN_KEY_VAR_2 = "\u001b[B" // arrow down
const LEFT_KEY_VAR_2 = "\u001b[D" // arrow left
const RIGHT_KEY_VAR_2 = "\u001b[C" // arrow right

type InputAction = "up" | "down" | "left" | "right"

const UP_ACTION: InputAction = "up"
const DOWN_ACTION: InputAction = "down"
const LEFT_ACTION: InputAction = "left"
const RIGHT_ACTION: InputAction = "right"

const INPUT_MAP = new Map<string, InputAction>([
    [UP_KEY_VAR_1, UP_ACTION],
    [DOWN_KEY_VAR_1, DOWN_ACTION],
    [LEFT_KEY_VAR_1, LEFT_ACTION],
    [RIGHT_KEY_VAR_1, RIGHT_ACTION],
    [UP_KEY_VAR_2, UP_ACTION],
    [DOWN_KEY_VAR_2, DOWN_ACTION],
    [LEFT_KEY_VAR_2, LEFT_ACTION],
    [RIGHT_KEY_VAR_2, RIGHT_ACTION],
])

// game config - display
const PADDING_CHAR = " ";
const BORDER_X_CHAR = "-";
const BORDER_Y_CHAR = "|";
const BG_CHAR = " ";
const SNAKE_CHAR = "x";

// init game logic
const inputActionBuffer: InputAction[] = [];
const snakeHeadPos = [0, CANVAS_HEIGHT / 2]; // x,y
const snakeDirection = [1, 0]

function onInput(key: string) {
    // debug("Key:" + JSON.stringify(key));

    // map key to action
    if (INPUT_MAP.has(key)) {
        const action = INPUT_MAP.get(key)!
        inputActionBuffer.push(action);

        debug("Action:" + action);
    }
}

// init draw
const canvas = createCanvas();
debug("canvas width:" + canvas[0]!.length.toString())
debug("canvas height:" + canvas.length.toString())

while (true) {

    // TODO : improve with input buffer logic (eg. quick up/left succession)
    // poll input
    const lastAction = inputActionBuffer.pop()

    // game logic - OnUpdate
    if (lastAction === "up") {
        snakeDirection[0] = 0;
        snakeDirection[1] = 1;
    }
    else if (lastAction === "down") {
        snakeDirection[0] = 0;
        snakeDirection[1] = -1;
    }
    else if (lastAction === "left") {
        snakeDirection[0] = -1;
        snakeDirection[1] = 0;
    }
    else if (lastAction === "right") {
        snakeDirection[0] = 1;
        snakeDirection[1] = 0;
    }

    snakeHeadPos[0] = (snakeHeadPos[0]! + snakeDirection[0]!) % CANVAS_WIDTH
    snakeHeadPos[1] = (snakeHeadPos[1]! + snakeDirection[1]!) % CANVAS_HEIGHT

    if (snakeHeadPos[0] < 0) snakeHeadPos[0] += CANVAS_WIDTH
    if (snakeHeadPos[1] < 0) snakeHeadPos[1] += CANVAS_HEIGHT

    // game logic - OnAfterUpdate
    inputActionBuffer.length = 0

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
    // canvas is drawn from top to bottom but game coordinates is from bottom to top
    canvas[MAX_HEIGHT - (PADDING_Y + y) - 2]![PADDING_X + x] = char;
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

    if (ansi.isANSISupported) {
        ansi.clearAndDrawBuffer(buffer)
    }
    else {
        console.clear();
        console.log(buffer)
    }
}

function setup() {
    process.stdin.setRawMode(true)
    process.stdin.resume(); // necessary or else "data" event wont fire
    process.stdin.setEncoding("utf8") // so can do string comparison on received keypresses

    // cleanup listeners
    process.on("exit", cleanup); // Regular exit on program end
    process.on("SIGINT", cleanupAndExit); // Ctrl-C, does not exit by default, need to manually exit
    process.on("SIGTERM", cleanupAndExit); // Terminated by terminal

    // input listeners
    process.stdin.on("data", (key: string) => {
        // Ctrl+C sends character code 3
        if (key === "\u0003") {
            process.kill(process.pid, "SIGINT");
            return;
        }

        onInput(key);
    });
}

function cleanup() {
    process.stdin.setRawMode(false)

    if (DEBUG_MODE)
        return;

    if (ansi.isANSISupported) {
        ansi.cleanup();
    }
    else {
        console.clear();
    }
}

function cleanupAndExit() {
    cleanup();
    process.exit(0);
}

function debug(message: string) {
    if (DEBUG_MODE)
        console.log(message)
}
