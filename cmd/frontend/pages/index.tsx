import React, { useCallback, useState, useRef, useEffect } from 'react'

import type { NextPage } from 'next'
import Head from 'next/head'
import Script from 'next/script'
import Image from 'next/image'
import { useRouter } from 'next/router'

import SyntaxHighlighter from "react-syntax-highlighter";
import {tomorrow} from "react-syntax-highlighter/dist/cjs/styles/hljs";

import gtihubImage from "../public/GitHub-Mark-Light-64px.png";
import { PlantUML } from '../components/plantuml'
import { decode62, encode62 } from '../lib/base62'
import { Mermaid } from '../components/mermaid'

const defaultSrc = `# Sample Markdown

* table: User
    * @id
    * name: string
    * $email: string
    * age: integer
    * jobs: *Job.id[]

* table: Job
    * @id
    * name: string
`;

function useSelect(defaultValue: string): [string, (e: React.ChangeEvent<HTMLSelectElement>)=>void, (v: string)=>void] {
  const [value, setValue] = useState(defaultValue);
  const onSelect = useCallback(function selectFormat(e: React.ChangeEvent<HTMLSelectElement>) {
    setValue(e.target.value);
  }, [])
  return [value, onSelect, setValue];
}

const Home: NextPage = () => {
  const router = useRouter();

  const [format, onSelectFormat, setFormat] = useSelect("sql");
  const [dialect, onSelectDialect, setDialect] = useSelect("postgres");
  const [initialSrc, setInitialSrc] = useState(""); // init after loading in browser
  const [tab, setTab] = useState("preview");

  const src = useRef("");                           // textarea is uncontrolled form. this keeps the value
  const [result, setResult] = useState("");

  useEffect(function initializeStatusFromQueryParameter() {
    setFormat(router.query["f"] as string || "sql");
    src.current = router.query["s"] ? decode62(router.query["s"] as string) : defaultSrc;
    setInitialSrc(src.current);
    setTab(router.query["t"] as string || "preview");
    setDialect(router.query["d"] as string || "postgres")
  }, [router.isReady])

  const modifySrc = useCallback(function modifySrcWhenEdit(e: React.ChangeEvent<HTMLTextAreaElement>) {
    src.current = e.target.value;
  }, [])

  useEffect(function regenerateWhenFormatIsChanged(){
    if (result !== "") {
      generate();
    }
  }, [format]);

  const generate = useCallback(() => {
    let result;
    switch (format) {
      case "sql":
        result = md2sql.toSQL(src.current, dialect);
        break;
      case "plantuml":
        result = md2sql.toPlantUML(src.current);
        break;
      default: //"mermaid":
        result = md2sql.toMermaid(src.current);
        break;
    }
    if (result.ok) {
      setResult(result.result);
      router.replace({
        query: {
          f: format,
          s: encode62(src.current),
          t: tab,
          d: dialect,
        },        
      }, undefined, { shallow: true});
      } else {
      console.error(result.message);
    }
  }, [format, dialect])

  const copyToClipboard = useCallback(() => {
    navigator.clipboard.writeText(result);
  }, [result])

  return (
    <div className="flex flex-col h-full">
      <Head>
        <title>md2sql</title>
        <meta name="description" content="Generate SQL/ERD from Markdown" />
        <link rel="icon" href="/md2sql/favicon.ico" />
      </Head>
      { /* Load web assembly */ }
      <Script id="exec-wasm" src="/md2sql/wasm_exec.js" onLoad={() => {
        // @ts-ignore
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("/md2sql/md2sql.wasm"), go.importObject).then((result) => {
          go.run(result.instance);
        });
      }}/>

      <header className="navbar mb-2 shadow-lg bg-neutral text-neutral-content rounded-box">
        <div className="flex-1 px-2 mx-2">
          <span className="text-lg font-bold">
            md2sql
          </span>
        </div>
        <div className="flex-none">
          <a className="btn btn-square btn-ghost" href="https://github.com/shibukawa/md2sql">
            <Image src={gtihubImage} alt="github"/>
          </a>
        </div>
      </header>

      <main className="main flex flex-col w-full grow lg:flex-row items-stretch">
        <div className="flex flex-grow h-full card rounded-box shadow-2xl p-6">
          <h2 className="grow-0 font-medium leading-tight text-4xl mt-0 mb-2 text-blue-600">Markdown Source</h2>
          { initialSrc
             ? <textarea className="m-6 p-1 textarea textarea-bordered grow h-full" onInput={modifySrc} defaultValue={initialSrc}></textarea>
             : <textarea className="m-6 p-1 textarea textarea-bordered grow h-full"></textarea>
          }
          <div className="flex">
            <div className="form-control mx-1">
              <label className="label">
                <span className="label-text">Export Format</span>
              </label>
              <select className="select select-bordered" value={format} onChange={onSelectFormat}>
                <option value="sql">SQL</option>
                <option value="mermaid">Mermaid.js</option>
                <option value="plantuml">PlantUML</option>
              </select>
            </div>
            {
              format === "sql" ? 
                <div className="form-control mx-1">
                  <label className="label">
                    <span className="label-text">Dialect</span>
                  </label>
                  <select className="select select-bordered" value={dialect} onChange={onSelectDialect}>
                    <option value="postgres">PostgreSQL</option>
                    <option value="mysql">MySQL</option>
                    <option value="sqlite">SQLite</option>
                  </select>
                </div>
              : null
            }
          </div>
        </div> 
        <div className="divider lg:divider-vertical"></div> 
        <div className="flex flex-grow h-full card rounded-box shadow-2xl p-6">
          <h2 className="grow-0 font-medium leading-tight text-4xl mt-0 mb-2 text-blue-600">
            Result
            { format !== "sql" ?
            <div className="tabs">
              <a className={`tab tab-lifted ${tab==="preview" ? "tab-active" : ""}`} onClick={() => { setTab("preview")}}>Preview</a> 
              <a className={`tab tab-lifted ${tab==="src" ? "tab-active" : ""}`} onClick={() => { setTab("src")}}>Source</a> 
            </div> : null}
          </h2>

          { format === "sql" || tab === "src" ? 
            <div className="grow m-2">
              <SyntaxHighlighter language={format} style={tomorrow} className="h-full rounded-md border border-slate-300 text-xs">
                {result}
              </SyntaxHighlighter>
            </div>
            : format === "mermaid" ?
            <Mermaid className="grow m-2" src={result} />
            : format === "plantuml" ?
            <PlantUML src={result} className="grow m-2" key={`${format}${result}${tab}`} />
            : null
          }
          <div className="grow-0 flex w-full">
            <button className="btn m-1" onClick={generate}>Generate</button>
            <button className="btn m-1" disabled={result === ""} onClick={copyToClipboard}>Copy To Clipboard</button>
          </div>
        </div>
      </main>
    </div>
  );
}

export default Home
