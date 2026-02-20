export function getMaxCharPerLine(): number {
    const ruler = document.body.querySelector("[data-ruler]");
    if (!ruler) {
        console.error(
            "initAsync:",
            "cannot find [data-ruler]!",
            "returning 0...",
        );
        return 0;
    }

    const rect = ruler.getBoundingClientRect();
    const text = ruler.innerHTML;

    const charWidth = rect.width / text.length;

    // console.log("ruler.width:", rect.width);
    // console.log("ruler.units:", text.length);
    // console.log("charWidth:", charWidth);
    // console.log("window.width:", window.innerWidth);

    const maxCharCount = window.innerWidth / charWidth;
    // console.log("maxSizeX", maxSizeX);

    return maxCharCount;
}
