"use client";

import { motion } from "motion/react";
import { Check } from "lucide-react";
import { useState } from "react";
import { cn } from "@/lib/utils";
import { CodeBlock } from "../../shared/code-block";
import { codeExamples } from "./examples";

export function CodeDemo() {
  const [activeTab, setActiveTab] = useState(0);
  const [outputFormat, setOutputFormat] = useState<
    "json" | "yaml" | "toml" | "typescript"
  >("json");

  const getOutput = (example: (typeof codeExamples)[0]) => {
    switch (outputFormat) {
      case "yaml":
        return example.yaml || example.json;
      case "toml":
        return example.toml || example.json;
      case "typescript":
        return example.typescript || example.json;
      default:
        return example.json;
    }
  };

  return (
    <section className="py-24 overflow-hidden">
      <div className="container mx-auto px-4 md:px-6">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          <div>
            <motion.h2
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              className="text-3xl font-bold tracking-tight md:text-4xl lg:text-5xl mb-6"
            >
              Write Logic, <br />
              <span className="text-primary">
                Get JSON, YAML, TOML, TypeScript.
              </span>
            </motion.h2>
            <motion.p
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: 0.1 }}
              className="text-lg text-muted-foreground mb-8"
            >
              JSSON brings the power of a real programming language to your
              configuration files. No more copy-pasting or manual error-prone
              editing.
            </motion.p>

            <div className="space-y-4">
              {[
                "Native variables & constants",
                "Arithmetic & Conditional Logic",
                "Templates & Maps for complex arrays",
                "Smart Ranges with steps",
                "Modular configuration (include)",
              ].map((item, i) => (
                <motion.div
                  key={i}
                  initial={{ opacity: 0, x: -20 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: 0.2 + i * 0.1 }}
                  className="flex items-center gap-3"
                >
                  <div className="flex h-6 w-6 items-center justify-center rounded-full bg-primary/10 text-primary">
                    <Check className="h-3.5 w-3.5" />
                  </div>
                  <span className="text-muted-foreground">{item}</span>
                </motion.div>
              ))}
            </div>
          </div>

          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            whileInView={{ opacity: 1, scale: 1 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5 }}
            className="relative"
          >
            <div className="relative rounded-xl border border-border bg-card overflow-hidden">
              <div className="flex border-b border-border bg-muted/30">
                {codeExamples.map((example, i) => (
                  <button
                    key={i}
                    onClick={() => setActiveTab(i)}
                    className={cn(
                      "px-6 py-3 font-medium transition-colors border-r text-xs border-border",
                      activeTab === i
                        ? "bg-card text-foreground"
                        : "bg-transparent text-muted-foreground hover:text-foreground"
                    )}
                  >
                    {example.title}
                  </button>
                ))}
              </div>

              <div className="grid grid-cols-2 divide-x divide-border">
                <div className="p-0">
                  <div className="px-4 py-2 text-xs font-mono text-muted-foreground border-b border-border bg-muted/10 flex justify-between">
                    <span>input.jsson</span>
                    <span className="text-primary">JSSON</span>
                  </div>
                  <div className="h-[400px] overflow-hidden bg-card">
                    <CodeBlock
                      code={codeExamples[activeTab].jsson}
                      language="jsson"
                      className="p-4 h-full overflow-auto"
                    />
                  </div>
                </div>
                <div className="p-0 bg-muted/5">
                  <div className="px-4 py-2 text-xs font-mono text-muted-foreground border-b border-border bg-muted/10 flex justify-between items-center">
                    <span>
                      output.
                      {outputFormat === "typescript" ? "ts" : outputFormat}
                    </span>
                    <div className="flex gap-1">
                      {(["json", "yaml", "toml", "typescript"] as const).map(
                        (format) => (
                          <button
                            key={format}
                            onClick={() => setOutputFormat(format)}
                            className={cn(
                              "px-2 py-0.5 text-[10px] rounded transition-colors uppercase font-semibold",
                              outputFormat === format
                                ? "bg-primary text-primary-foreground"
                                : "bg-muted/50 text-muted-foreground hover:bg-muted hover:text-foreground"
                            )}
                          >
                            {format === "typescript" ? "TS" : format}
                          </button>
                        )
                      )}
                    </div>
                  </div>
                  <div className="h-[400px] overflow-hidden">
                    <CodeBlock
                      code={getOutput(codeExamples[activeTab])}
                      language="json"
                      className="p-4 h-full overflow-auto"
                    />
                  </div>
                </div>
              </div>
            </div>
          </motion.div>
        </div>
      </div>
    </section>
  );
}
