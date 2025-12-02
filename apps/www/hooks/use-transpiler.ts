"use client";

import { useState, useEffect, useCallback, useRef } from "react";

// Declare the global function type
declare global {
    interface Window {
        Go: any;
        transpileJSSON: (input: string, format?: string) => { output?: string; error?: string };
    }
}

export function useTranspiler(initialCode: string = "") {
    const [code, setCode] = useState(initialCode);
    const [output, setOutput] = useState("");
    const [error, setError] = useState<string | null>(null);
    const [isTranspiling, setIsTranspiling] = useState(false);
    const [isWasmLoaded, setIsWasmLoaded] = useState(false);
    const [compilationTime, setCompilationTime] = useState<number | null>(null);
    const goRef = useRef<any>(null);

    // Load WASM
    useEffect(() => {
        const loadWasm = async () => {
            try {
                if (!window.Go) {
                    // Load wasm_exec.js dynamically if not present
                    const script = document.createElement("script");
                    script.src = "/wasm_exec.js";
                    script.async = true;
                    script.onload = async () => {
                        await initGo();
                    };
                    document.body.appendChild(script);
                } else {
                    await initGo();
                }
            } catch (err) {
                console.error("Failed to load WASM:", err);
                setError("Failed to load JSSON compiler");
            }
        };

        const initGo = async () => {
            if (goRef.current) return;

            const go = new window.Go();
            goRef.current = go;

            const result = await WebAssembly.instantiateStreaming(
                fetch("/jsson.wasm"),
                go.importObject
            );

            go.run(result.instance);
            setIsWasmLoaded(true);
        };

        loadWasm();
    }, []);

    const transpile = useCallback((sourceCode: string, format: string = "json") => {
        if (!isWasmLoaded) return;

        setIsTranspiling(true);
        setError(null);

        try {
            setTimeout(() => {
                if (window.transpileJSSON) {
                    const startTime = performance.now();
                    const result = window.transpileJSSON(sourceCode, format);
                    const endTime = performance.now();
                    const elapsed = endTime - startTime;

                    setCompilationTime(elapsed);

                    if (result.error) {
                        setError(result.error);
                    } else {
                        setOutput(result.output || "");
                    }
                }
                setIsTranspiling(false);
            }, 10);
        } catch (err: any) {
            setError(err.message || "Transpilation failed");
            setIsTranspiling(false);
        }
    }, [isWasmLoaded]);


    function runCode(format: string = "json") {
        transpile(code, format);
    }

    return {
        code,
        setCode,
        output,
        error,
        isTranspiling,
        isWasmLoaded,
        compilationTime,
        runCode,
    };
}
