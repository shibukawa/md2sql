/// <reference types="next" />
/// <reference types="next/image-types/global" />

// NOTE: This file should not be edited
// see https://nextjs.org/docs/basic-features/typescript for more information.

type sql = (src: string, dialect: string) => { ok: true, result: string} | {ok: false, message: string};
type f = (src: string) => { ok: true, result: string} | {ok: false, message: string};

declare var md2sql:{
    toSQL: sql,
    toMermaid: f,
    toPlantUML: f,
};