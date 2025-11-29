"use client";

import { motion } from "motion/react";
import {
  Zap,
  Copy,
  Repeat,
  Calculator,
  FileJson,
  ArrowRightLeft,
  Layers,
  Table2,
  Variable,
} from "lucide-react";

const features = [
  {
    icon: ArrowRightLeft,
    title: "Universal Multi-Format Output",
    description:
      "Write once, export everywhere. Convert JSSON into JSON, YAML, TOML, or TypeScript — your entire config ecosystem from a single source of truth.",
  },
  {
    icon: Variable,
    title: "Declare Variables",
    description:
      "Declare variables once and reuse them throughout your config. No more copy-pasting values.",
  },
  {
    icon: Copy,
    title: "Reusable Templates",
    description:
      "Define structured data models once and instantiate them with spreadsheet-style rows. Zero duplication, maximum clarity.",
  },
  {
    icon: ArrowRightLeft,
    title: "Advanced Maps",
    description:
      "Transform values declaratively with map(). Generate derived fields, dynamic objects and multidimensional data pipelines with nested mappings.",
  },
  {
    icon: Repeat,
    title: "Powerful Ranges",
    description:
      "Generate massive datasets with 1..10000, nested ranges, and range-driven maps. Ideal for geodata, schedules, grids and simulations.",
  },
  {
    icon: Calculator,
    title: "Expressions & Logic",
    description:
      "Full math, comparisons, modulo, ternaries and computed properties right inside the DSL — no scripting needed.",
  },
  {
    icon: Layers,
    title: "Nested Arrays & Complex Structures",
    description:
      "Naturally express deep arrays, nested maps, multi-level expansions and matrix-like outputs with clean declarative syntax.",
  },
  {
    icon: FileJson,
    title: "Includes & Composition",
    description:
      "Break large configs into logical parts and merge them with include. Supports keep, overwrite, and strict error modes.",
  },
  {
    icon: Zap,
    title: "Native Types, No Boilerplate",
    description:
      "Strings, ints, floats, booleans, objects and arrays — without JSON’s quoting hell. Clean syntax built for humans.",
  },
  {
    icon: Table2,
    title: "High-Scale Data Generation",
    description:
      "Create thousands or millions of records with a few lines. Perfect for seeds, mocks, simulations, grids, and AI datasets.",
  },
  {
    icon: Zap,
    title: "Inline-Safe & LLM-Friendly",
    description:
      "Minifiable, stable under compression, and easy for LLMs to parse. JSSON wasn’t built for AI, but AI absolutely loves it.",
  },
];

export function Features() {
  return (
    <section className="py-24 bg-muted/30">
      <div className="container mx-auto px-4 md:px-6">
        <div className="text-center mb-16">
          <motion.h2
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5 }}
            className="text-3xl font-bold tracking-tight md:text-4xl lg:text-5xl"
          >
            Why JSSON?
          </motion.h2>
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5, delay: 0.1 }}
            className="mt-4 text-lg text-muted-foreground max-w-2xl mx-auto"
          >
            Stop repeating yourself in JSON. Upgrade to a language designed for
            modern application configuration.
          </motion.p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {features.map((feature, index) => (
            <motion.div
              key={index}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5, delay: index * 0.1 }}
              className="group relative overflow-hidden rounded-2xl border border-border bg-card p-8 hover:shadow-lg transition-all duration-300 hover:-translate-y-1"
            >
              <div className="absolute inset-0 bg-linear-to-br from-primary/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
              <div className="relative z-10">
                <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-colors duration-300">
                  <feature.icon className="h-6 w-6" />
                </div>
                <h3 className="text-xl font-semibold mb-2">{feature.title}</h3>
                <p className="text-muted-foreground">{feature.description}</p>
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
