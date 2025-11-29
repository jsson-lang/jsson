"use client";

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { motion } from "motion/react";

const faqs = [
  {
    question: "Is JSSON a replacement for JSON?",
    answer:
      "No. JSSON doesn’t try to replace JSON — it replaces the pain of writing JSON. You write JSSON, and it compiles into standard JSON (or YAML, TOML, TS...). Your existing tools still work exactly the same.",
  },
  {
    question: "How is JSSON different from TOON, TONL, YAML, or TOML?",
    answer:
      "Those formats are final outputs. JSSON is a meta-format: a source language that adds logic, templates, ranges, and transformations, and then generates JSON/YAML/TOML/TypeScript from it. It’s not competing — it’s upstream.",
  },
  {
    question: "Does JSSON support loops, ranges, and dynamic data?",
    answer:
      "Yes. Ranges, arithmetic, template blocks, nested maps, and inline conditionals are all built into the language. JSSON can generate thousands of records from a few lines of declarative logic.",
  },
  {
    question: "Can JSSON generate TypeScript with types?",
    answer:
      "Yes. JSSON can output TypeScript objects with `as const` for immutability and automatically infer types via `typeof`. No duplicated definitions, no drifting configs — full type safety.",
  },
  {
    question: "Is JSSON whitespace-dependent?",
    answer:
      "No. Unlike YAML or TOON, JSSON has clear delimiters and is not indentation-sensitive. You can minify or inline it however you want — it always works.",
  },
  {
    question: "Does the playground run entirely in the browser?",
    answer:
      "Yes. Thanks to WebAssembly, the JSSON compiler runs 100% client-side in your browser. Nothing leaves your machine.",
  },
  {
    question: "Is JSSON good for LLM prompts?",
    answer:
      "Yes. It removes unnecessary tokens (quotes, commas, redundancy) while keeping structure clear for AI models. Ideal for prompts that need typed, structured, token-efficient data.",
  },
  {
    question: "Can I extend JSSON with my own transpilers?",
    answer:
      "Yes. The compiler is modular, and you can add custom generators to output any format or language your project requires.",
  },
  {
    question: "Why does JSSON exist?",
    answer:
      "Because writing large JSON/YAML/TOML configs by hand is painful and repetitive. JSSON adds logic, structure, and reusability — without sacrificing compatibility with existing ecosystems.",
  },
];


export function FAQ() {
  return (
    <section className="py-24 bg-muted/30">
      <div className="container mx-auto px-4 md:px-6 max-w-3xl">
        <div className="text-center mb-16">
          <motion.h2
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            className="text-3xl font-bold tracking-tight md:text-4xl"
          >
            Frequently Asked Questions
          </motion.h2>
        </div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.2 }}
        >
          <Accordion className="w-full">
            {faqs.map((faq, index) => (
              <AccordionItem key={index} value={`item-${index}`}>
                <AccordionTrigger className="text-left text-lg font-medium">
                  {faq.question}
                </AccordionTrigger>
                <AccordionContent className="text-muted-foreground text-base">
                  {faq.answer}
                </AccordionContent>
              </AccordionItem>
            ))}
          </Accordion>
        </motion.div>
      </div>
    </section>
  );
}
