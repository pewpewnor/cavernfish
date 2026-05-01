export function rootFontSizePx(): number {
    return parseFloat(getComputedStyle(document.documentElement).fontSize) || 16;
}

export function pxToRem(px: number): number {
    return px / rootFontSizePx();
}

export function remToPx(rem: number): number {
    return rem * rootFontSizePx();
}
