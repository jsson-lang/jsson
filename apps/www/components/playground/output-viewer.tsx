"use client";

import { useState, useEffect, useMemo } from "react";
import { CodeBlock } from "../shared/code-block";
import { Tabs, TabsList, TabsTrigger } from "../ui/tabs";
import { Button } from "../ui/button";
import { Copy } from "lucide-react";
import { toastManager } from "@/components/ui/toast";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectGroupLabel,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";

import { usePlaygroundContext } from "@/contexts/playground-context";

interface OutputViewerProps {
  jssonCode: string;
  error?: string | null;
}

const structuredFormats = [
  {
    value: "json",
    label: "JSON",
  },
  {
    value: "yaml",
    label: "YAML",
  },
  {
    value: "toml",
    label: "TOML",
  },
];
const typedFormats = [
  {
    value: "typescript",
    label: "TypeScript",
  },
];

export function OutputViewer({ jssonCode, error }: OutputViewerProps) {
  const {
    format,
    setFormat,
    setOutput: setContextOutput,
  } = usePlaygroundContext();
  const [output, setOutput] = useState("");
  const [transpileError, setTranspileError] = useState<string | null>(null);

  function approxTokens(text: string) {
    if (!text) return 0;

    return Math.ceil(text.length / 3.5);
  }

  useEffect(() => {
    if (!jssonCode || error) {
      setOutput("");
      setContextOutput("");
      return;
    }

    if (window.transpileJSSON) {
      const result = window.transpileJSSON(jssonCode, format);

      if (result.error) {
        setTranspileError(result.error);
        setOutput("");
        setContextOutput("");
      } else {
        setTranspileError(null);
        setOutput(result.output || "");
        setContextOutput(result.output || "");
      }
    }
  }, [jssonCode, format, error, setContextOutput]);

  const displayError = error || transpileError;

  const metrics = useMemo(() => {
    if (!output) {
      return {
        lines: 0,
        chars: 0,
        tokens: 0,
      };
    }

    return {
      lines: output.split("\n").length,
      chars: output.length,
      tokens: approxTokens(output),
    };
  }, [output]);

  function copyToClipboard() {
    navigator.clipboard.writeText(output);
    toastManager.add({
      title: "Copied!",
      description: "Output copied to clipboard.",
    });
  }

  return (
    <div className="h-full w-full flex flex-col  overflow-hidden border-l bg-card/50 backdrop-blur-sm shadow-sm">
      <div className="flex items-center justify-between px-6 py-2 border-b border-border bg-muted/30">
        <div className="flex items-center gap-2">
          <span className="text-xs font-medium text-muted-foreground">
            output.{format}
          </span>
          <Select
            aria-label="Select Output Format"
            value={format}
            items={[...structuredFormats, ...typedFormats]}
            onValueChange={(value) => value !== null && setFormat(value)}
          >
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectGroupLabel>Structured Formats</SelectGroupLabel>
                {structuredFormats.map((fmt) => (
                  <SelectItem key={fmt.value} value={fmt.value}>
                    {fmt.label}
                  </SelectItem>
                ))}
              </SelectGroup>
              <SelectGroup>
                <SelectGroupLabel>Typed Formats</SelectGroupLabel>
                {typedFormats.map((fmt) => (
                  <SelectItem key={fmt.value} value={fmt.value}>
                    {fmt.label}
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>

        <div className="flex items-center gap-3">
          <div className="text-muted-foreground text-xs border-r pr-3">
            <span className="font-semibold">{metrics.lines}</span> lines |
            <span className="font-semibold"> {metrics.chars}</span> chars |
            <span className="font-semibold"> ~ {approxTokens(output)}</span>{" "}
            tokens
          </div>

          <Button
            variant="ghost"
            size="sm"
            onClick={copyToClipboard}
            disabled={!output}
          >
            <Copy className="h-4 w-4" />
            Copy
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-card/50 p-0">
        {displayError ? (
          <div className="p-4 text-red-400 font-mono text-sm">
            Error: {displayError}
          </div>
        ) : (
          <CodeBlock
            code={output}
            language={format === "json" ? "json" : "jsson"}
            className="h-full rounded-none border-none bg-transparent"
            showLineNumbers
          />
        )}
      </div>
    </div>
  );
}
