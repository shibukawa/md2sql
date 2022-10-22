import baseX from "base-x";

const base62 = baseX('0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ');
const encoder = new TextEncoder();
const decoder = new TextDecoder();

export function encode62(src: string) {
    return base62.encode(encoder.encode(src))
}

export function decode62(src: string) {
    return decoder.decode(base62.decode(src))
}