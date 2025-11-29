"use client";

import Editor, { useMonaco, OnMount } from "@monaco-editor/react";
import { useEffect } from "react";
import { Button } from "../ui/button";
import ExamplesSheet from "./examples";

interface JSSONEditorProps {
  value: string;
  onChange: (value: string | undefined) => void;
}

export function JSSONEditor({ value, onChange }: JSSONEditorProps) {
  const monaco = useMonaco();

  useEffect(() => {
    if (monaco) {
      monaco.languages.register({ id: "jsson" });

      monaco.languages.setMonarchTokensProvider("jsson", {
        keywords: ["include", "template", "map", "step"],
        constants: ["true", "false", "null"],

        tokenizer: {
          root: [
            // Comments
            [/\/\/.*$/, "comment"],

            // Keywords
            [/\b(include|template|map|step)\b/, "keyword"],

            // Booleans and null
            [/\b(true|false|null)\b/, "constant.language"],

            // Strings
            [/"([^"\\]|\\.)*$/, "string.invalid"], // non-terminated string
            [/"/, { token: "string.quote", bracket: "@open", next: "@string" }],

            // Numbers (floats first, then integers)
            [/\d+\.\d+/, "number.float"],
            [/\d+/, "number"],

            // Range operator
            [/\.\./, "keyword.operator.range"],

            // Member access
            [/\./, "keyword.operator.member"],

            // Arithmetic operators
            [/[+\-*/%]/, "keyword.operator.arithmetic"],

            // Comparison operators
            [/[<>]=?|[!=]=/, "keyword.operator.comparison"],

            // Ternary operators
            [/[?:]/, "keyword.operator.ternary"],

            // Assignment
            [/=/, "keyword.operator.assignment"],

            // Identifiers (variables and keys)
            [/[a-zA-Z_][a-zA-Z0-9_]*/, "identifier"],

            // Delimiters
            [/[{}]/, "@brackets"],
            [/[\[\]]/, "@brackets"],
            [/[()]/, "@brackets"],
            [/,/, "delimiter.comma"],
          ],

          string: [
            [/[^\\"]+/, "string"],
            [/\\./, "string.escape"],
            [/"/, { token: "string.quote", bracket: "@close", next: "@pop" }],
          ],
        },
      });

      monaco.editor.defineTheme("jsson-dark", {
        base: "vs-dark",
        inherit: true,
        rules: [
          { token: "comment", foreground: "6A9955", fontStyle: "italic" },
          { token: "keyword", foreground: "C586C0", fontStyle: "bold" },
          { token: "constant.language", foreground: "569CD6" },
          { token: "string", foreground: "CE9178" },
          { token: "string.escape", foreground: "D7BA7D" },
          { token: "number", foreground: "B5CEA8" },
          { token: "number.float", foreground: "B5CEA8" },
          { token: "keyword.operator", foreground: "D4D4D4" },
          { token: "identifier", foreground: "9CDCFE" },
          { token: "delimiter", foreground: "D4D4D4" },
        ],
        colors: {
          "editor.background": "#1e1e1e",
          "editor.foreground": "#D4D4D4",
          "editorLineNumber.foreground": "#858585",
          "editor.selectionBackground": "#264F78",
          "editor.inactiveSelectionBackground": "#3A3D41",
        },
      });

      // Set theme
      monaco.editor.setTheme("jsson-dark");
    }
  }, [monaco]);

  function approximateTokens(text: string) {
    if (!text) return 0;

    return Math.ceil(text.length / 3.5);
  }

  return (
    <div className="h-full w-full overflow-hidden bg-card/50 backdrop-blur-sm shadow-sm">
      <div className="flex items-center justify-between px-6 py-2 border-b border-border bg-muted/30">
        <span className="text-xs font-medium text-muted-foreground">
          input.jsson
        </span>

        <div className="flex items-center gap-2">
          <h3 className="text-muted-foreground text-xs border-r pr-2">
            {value.split("\n").length} lines | {value.length} chars | ~{" "}
            {approximateTokens(value)} tokens
          </h3>

          <Button variant={"destructive-outline"} onClick={() => onChange("")}>
            Clear
          </Button>
          <ExamplesSheet onSelect={(code) => onChange(code)} />
        </div>
      </div>
      <Editor
        height="100%"
        defaultLanguage="jsson"
        theme="jsson-dark"
        value={value}
        onChange={onChange}
        options={{
          minimap: { enabled: false },
          fontSize: 14,
          fontFamily: "Geist Mono, monospace",
          padding: { top: 16 },
          scrollBeyondLastLine: false,
          automaticLayout: true,
          lineNumbers: "on",
          renderLineHighlight: "all",
          scrollbar: {
            vertical: "visible",
            horizontal: "visible",
          },
        }}
      />
    </div>
  );
}
