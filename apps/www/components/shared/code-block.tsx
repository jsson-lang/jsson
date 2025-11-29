"use client";

import { cn } from "@/lib/utils";
import { useMemo } from "react";
import highlightLine from "@/lib/highlightLines";

interface CodeBlockProps {
  code: string;
  language?: "jsson" | "json";
  className?: string;
  showLineNumbers?: boolean;
}

export function CodeBlock({
  code,
  language = "jsson",
  className,
  showLineNumbers = false,
}: CodeBlockProps) {
  const highlighted = useMemo(() => {
    const lines = code.split("\n");

    // Remove empty first line
    if (lines[0]?.trim() === "") lines.shift();

    // Detect indentation
    const minIndent = lines
      .filter((l) => l.trim() !== "")
      .reduce((min, l) => {
        const indent = l.match(/^\s*/)?.[0].length ?? 0;
        return Math.min(min, indent);
      }, Infinity);

    return lines.map((line, i) => {
      const clean =
        minIndent !== Infinity ? line.slice(minIndent) : line;

      return (
        <div
          key={i}
          className="leading-relaxed whitespace-pre font-mono text-sm"
        >
          {showLineNumbers && (
            <span className="inline-block w-8 select-none text-muted-foreground/30 text-right mr-4">
              {i + 1}
            </span>
          )}

          <span
            dangerouslySetInnerHTML={{
              __html: highlightLine(clean, language),
            }}
          />
        </div>
      );
    });
  }, [code, language, showLineNumbers]);

  return (
    <div
      className={cn(
        "font-mono text-sm overflow-x-auto rounded-lg",
        className
      )}
    >
      {highlighted}
    </div>
  );
}

/* ------------------------------------------------------------------------- */
/* 🔥 HIGH-END HIGHLIGHTER — FAST + CONSISTENT                               */
/* ------------------------------------------------------------------------- */


