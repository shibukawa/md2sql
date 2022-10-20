import { useMemo, forwardRef } from "react"; 
import pako from "pako";
import { encode64 } from "../lib/encode64"

type ImageProps = JSX.IntrinsicElements['img'];

export const PlantUML = forwardRef<HTMLImageElement, ImageProps>(function PlantUML(plops, ref) {
    const { src, ...remained } = plops;

    const base64 = useMemo(() => {
        if (!src) {
            return "";
        }
        const bin = pako.deflateRaw(unescape(encodeURIComponent(src)));
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
    console.log({base64})
    
    return src ? <img {...remained} ref={ref} src={`http://www.plantuml.com/plantuml/svg/${base64}`} /> : null;
});

