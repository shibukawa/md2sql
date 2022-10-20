import React, { useCallback, useState, useRef, useEffect } from 'react'

import type { NextPage } from 'next'
import Head from 'next/head'
import Script from 'next/script'
import Image from 'next/image'

import SyntaxHighlighter from "react-syntax-highlighter";
import {tomorrow} from "react-syntax-highlighter/dist/cjs/styles/hljs";
import mermaid from "mermaid";

import gtihubImage from "../public/GitHub-Mark-Light-64px.png";
import { PlantUML } from '../components/plantuml'


// import styles from '../styles/Home.module.css'

const initialSrc = `# Sample Markdown

* table: User
    * ##id
    * name: string
    * #email: string
    * age: integer
    * jobs: *Job.id[]

* table: Job
    * ##id
    * name: string
`

const Home: NextPage = () => {
  const [format, setFormat] = useState("sql");
  const src = useRef(initialSrc);
  const [result, setResult] = useState("");
  const [tab, setTab] = useState("preview");

  const selectFormat = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setFormat(e.target.value);
  }, [])

  const modifySrc = useCallback((e: React.ChangeEvent<HTMLTextAreaElement>) => {
    src.current = e.target.value;
  }, [])

  useEffect(() => {
    if (result !== "") {
      generate();
    }
  }, [format]);

  useEffect(() => {
    if (format === "mermaid" && tab=== "preview" && result !== "") {
      mermaid.init({noteMargin: 10}, ".mermaid");
    }
  }, [format, result, tab])

  const generate = useCallback(() => {
    switch (format) {
      case "sql":
        const r1 = md2sql.toSQL(src.current);
        if (r1.ok) {
          setResult(r1.result);
        } else {
          console.error(r1.message);
        }
        break;
      case "plantuml":
        const r2 = md2sql.toPlantUML(src.current);
        if (r2.ok) {
          setResult(r2.result);
        } else {
          console.error(r2.message);
        }
        break;
      case "mermaid":
        const r3 = md2sql.toMermaid(src.current);
        if (r3.ok) {
          setResult(r3.result);
        } else {
          console.error(r3.message);
        }
        break;
    }
  }, [format])

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
          <textarea className="m-6 p-1 textarea textarea-bordered grow h-full" onInput={modifySrc} defaultValue={initialSrc}></textarea>
          <div className="flex">
            <label className="mx-2 label cursor-pointer">
              <span className="label-text">SQL</span> 
              <input type="radio" name="radio-6" className="mx-1 radio checked:bg-blue-500" onChange={selectFormat} value="sql" checked={format=="sql"}/>
            </label>
            <label className="mx-2 label cursor-pointer">
              <span className="label-text">Mermaid.js</span> 
              <input type="radio" name="radio-6" className="mx-1 radio checked:bg-blue-500" onChange={selectFormat} value="mermaid" checked={format=="mermaid"}/>
            </label>
            <label className="mx-2 label cursor-pointer">
              <span className="label-text">PlantUML</span> 
              <input type="radio" name="radio-6" className="mx-1 radio checked:bg-blue-500" onChange={selectFormat} value="plantuml" checked={format=="plantuml"}/>
            </label>
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
              <SyntaxHighlighter language={format} style={tomorrow} className="h-full">
                {result}
              </SyntaxHighlighter>
            </div>
            : format === "mermaid" ?
            <div className="grow m-2 mermaid" key={`${format}${result}${tab}`}>
              {result}
            </div>
            : format === "plantuml" ?
            <PlantUML src={result} className="grow m-2" key={`${format}${result}${tab}`}/>
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
