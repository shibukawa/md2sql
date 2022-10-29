import { useMemo, createContext, useContext } from "react";
import pako from "pako";
import { encode64 } from "../lib/encode64"

type ImageProps = JSX.IntrinsicElements['img'];

const encoder = new TextEncoder();

const PlantUMLContext = createContext("http://www.plantuml.com/plantuml");

export const PlantUMLProvider = PlantUMLContext.Provider;

export function PlantUML(plops: ImageProps) {
    const { src, ...remained } = plops;

    const serverUrl = useContext(PlantUMLContext) || "http://www.plantuml.com/plantuml";

    const base64 = useMemo(() => {
        if (!src) {
            return "";
        }
        const bin = pako.deflateRaw(encoder.encode(src));
        // https://stackoverflow.com/a/21214792
        const CHUNK_SIZE = 0x8000;
        let index = 0;
        const length = bin.length;
        const strs: string[] = [];
        while (index < length) {
          const slice = bin.subarray(index, Math.min(index + CHUNK_SIZE, length));
          // @ts-ignore 
          strs.push(String.fromCharCode.apply(null, slice));
          index += CHUNK_SIZE;
        }
        return encode64(strs.join(''));
    }, [src])
    
    return src ? <img {...remained} src={`${serverUrl}/svg/${base64}`} /> : null;
};

