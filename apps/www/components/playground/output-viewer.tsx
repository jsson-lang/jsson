"use client";

import { useState, useEffect, useMemo } from "react";
import Editor from "@monaco-editor/react";
import { Tabs, TabsList, TabsTrigger } from "../ui/tabs";
import { Button } from "../ui/button";
import { Copy, Download } from "lucide-react";
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
  output: string;
  error?: string | null;
  compilationTime?: number | null;
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
    value: "ts",
    label: "TypeScript",
  },
];

export function OutputViewer({
  output,
  error,
  compilationTime,
}: OutputViewerProps) {
  const {
    format,
    setFormat,
    setOutput: setContextOutput,
  } = usePlaygroundContext();

  function approxTokens(text: string) {
    if (!text) return 0;

    return Math.ceil(text.length / 3.5);
  }

  // Map format to Monaco language
  function getMonacoLanguage(format: string): string {
    const languageMap: Record<string, string> = {
      json: "json",
      yaml: "yaml",
      toml: "ini",
      ts: "typescript",
    };
    return languageMap[format] || "json";
  }

  useEffect(() => {
    setContextOutput(output);
  }, [output, setContextOutput]);

  const displayError = error;

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

  useEffect(() => {}, [format]);

  function generateDownload() {
    try {
      const blob = new Blob([output], { type: "text/plain" });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `output.${format}`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
      toastManager.add({
        title: "Downloaded!",
        description: "Output downloaded.",
        type: "success",
      });
    } catch (error) {
      toastManager.add({
        title: "Error",
        description: "Failed to download output.",
        type: "error",
      });
    }
  }

  function copyToClipboard() {
    navigator.clipboard.writeText(output);
    toastManager.add({
      title: "Copied!",
      description: "Output copied to clipboard.",
    });
  }

  function formatTime(ms: number): string {
    if (ms < 1) return `${(ms * 1000).toFixed(0)}µs`;
    if (ms < 1000) return `${ms.toFixed(2)}ms`;
    return `${(ms / 1000).toFixed(2)}s`;
  }

  return (
    <div className="h-full w-full flex flex-col border-l overflow-hidden shadow-sm">
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
          <Button
            variant="ghost"
            size="sm"
            onClick={copyToClipboard}
            disabled={!output}
          >
            <Copy className="h-4 w-4" />
            Copy
          </Button>
          <Button size="sm" disabled={!output} onClick={generateDownload}>
            <Download className="h-4 w-4" />
            Download File
          </Button>
        </div>
      </div>

      <div className="flex-1 min-h-0">
        {displayError ? (
          <div className="p-4 text-red-400 font-mono text-sm">
            Error: {displayError}
          </div>
        ) : (
          <Editor
            height="100%"
            language={getMonacoLanguage(format)}
            theme="jsson-dark"
            value={output}
            options={{
              readOnly: true,
              minimap: { enabled: false },
              fontSize: 14,
              fontFamily: "Geist Mono, monospace",
              padding: { top: 16 },
              scrollBeyondLastLine: false,
              automaticLayout: true,
              lineNumbers: "on",
              renderLineHighlight: "none",
              scrollbar: {
                vertical: "visible",
                horizontal: "visible",
              },
            }}
          />
        )}
      </div>

      <div className="flex items-center justify-between px-6 py-2 border-t border-border bg-muted/30">
        <div className="text-muted-foreground text-xs">
          <span className="font-semibold">{metrics.lines}</span> lines |
          <span className="font-semibold"> {metrics.chars}</span> chars |
          <span className="font-semibold"> ~ {approxTokens(output)}</span>{" "}
          tokens
        </div>
        {compilationTime !== null && (
          <div className="text-xs text-green-400 font-medium">
            ✓ Compiled in {formatTime(compilationTime!)}
          </div>
        )}
      </div>
    </div>
  );
}
