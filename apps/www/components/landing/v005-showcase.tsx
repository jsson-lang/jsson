"use client";

import { motion } from "motion/react";
import { Grid3x3, ShoppingCart, Sparkles } from "lucide-react";
import { CodeBlock } from "../shared/code-block";
import { cn } from "@/lib/utils";

const showcaseItems = [
  {
    icon: Grid3x3,
    title: "Nested Maps",
    subtitle: "Matrix Generation",
    badge: "25 numbers in 1 line",
    description: "Generate 2D matrices using nested map transformations",
    code: `// Multiplication table 5x5
rows := 5
cols := 5

table = (1..rows map (row) = (
  1..cols map (col) = row * col
))`,
    visual: (
      <div className="grid grid-cols-5 gap-1 p-4">
        {Array.from({ length: 25 }, (_, i) => {
          const row = Math.floor(i / 5) + 1;
          const col = (i % 5) + 1;
          return (
            <motion.div
              key={i}
              initial={{ opacity: 0, scale: 0 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.02 }}
              className="aspect-square flex items-center justify-center bg-primary/10 rounded text-xs font-mono text-primary font-semibold"
            >
              {row * col}
            </motion.div>
          );
        })}
      </div>
    ),
  },
  {
    icon: ShoppingCart,
    title: "Product Variants",
    subtitle: "E-commerce Power",
    badge: "6 variants automatically",
    description: "Generate all size/color combinations effortlessly",
    code: `products = (["S", "M", "L"] map (size) = (
  ["Red", "Blue"] map (color) = {
    sku = size + "-" + color
    price = 29.99
  }
))`,
    visual: (
      <div className="grid grid-cols-2 gap-2 p-4">
        {["S", "M", "L"].map((size, sizeIdx) =>
          ["Red", "Blue"].map((color, colorIdx) => {
            const index = sizeIdx * 2 + colorIdx;
            return (
              <motion.div
                key={`${size}-${color}`}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.1 }}
                className="p-3 rounded-lg border border-border bg-card"
              >
                <div className="text-xs font-mono text-muted-foreground mb-1">
                  SKU: {size}-{color}
                </div>
                <div className="flex items-center gap-2">
                  <div
                    className={cn(
                      "w-4 h-4 rounded-full",
                      color === "Red" ? "bg-red-500" : "bg-blue-500"
                    )}
                  />
                  <span className="text-sm font-medium">Size {size}</span>
                </div>
                <div className="text-xs text-muted-foreground mt-1">$29.99</div>
              </motion.div>
            );
          })
        )}
      </div>
    ),
  },
  {
    icon: Sparkles,
    title: "Large-Scale Generation",
    subtitle: "Massive Datasets",
    badge: "1,000 records instantly",
    description: "Generate thousands of records with patterns",
    code: `totalUsers := 1000
users = (0..(totalUsers - 1) map (id) = {
  id = id
  username = "user_" + id
  tier = id < 100 ? "bronze" : "silver"
})`,
    visual: (
      <div className="flex flex-col items-center justify-center p-8 space-y-4">
        <motion.div
          initial={{ scale: 0 }}
          whileInView={{ scale: 1 }}
          viewport={{ once: true }}
          transition={{ type: "spring", duration: 0.8 }}
          className="text-6xl font-bold bg-linear-to-r from-primary to-purple-500 bg-clip-text text-transparent"
        >
          1,000
        </motion.div>
        <div className="text-sm text-muted-foreground text-center">
          records generated
        </div>
        <div className="flex gap-2 mt-4">
          <div className="px-3 py-1 rounded-full bg-amber-500/10 text-amber-500 text-xs font-medium">
            100 bronze
          </div>
          <div className="px-3 py-1 rounded-full bg-gray-500/10 text-gray-400 text-xs font-medium">
            900 silver
          </div>
        </div>
      </div>
    ),
  },
];

export function V005Showcase() {
  return (
    <section className="py-24 bg-linear-to-b from-background/10 to-primary/5 relative overflow-hidden">
      <div className="absolute inset-0 bg-grid-white/10 bg-size-[20px_20px] mask-[radial-gradient(white,transparent_85%)]" />

      <div className="container mx-auto px-4 md:px-6 relative z-10">
        <div className="text-center mb-16">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            className="inline-block"
          >
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-primary/10 text-primary text-sm font-medium mb-4">
              <Sparkles className="h-4 w-4" />
              <span>New in v0.0.5</span>
            </div>
          </motion.div>

          <motion.h2
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.1 }}
            className="text-3xl font-bold tracking-tight md:text-4xl lg:text-5xl mb-4"
          >
            Data Generation <span className="text-primary">Powerhouse</span>
          </motion.h2>

          <motion.p
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.2 }}
            className="text-lg text-muted-foreground max-w-2xl mx-auto"
          >
            JSSON v0.0.5 introduces nested maps, nested arrays, and universal
            ranges. Turn complex data generation into simple, readable code.
          </motion.p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {showcaseItems.map((item, index) => (
            <motion.div
              key={index}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.2 }}
              className="group relative overflow-hidden rounded-2xl border border-border bg-card hover:shadow-xl transition-all duration-300"
            >
              <div className="p-6 border-b border-border">
                <div className="flex items-start justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-colors">
                      <item.icon className="h-5 w-5" />
                    </div>
                    <div>
                      <h3 className="font-semibold text-lg">{item.title}</h3>
                      <p className="text-xs text-muted-foreground">
                        {item.subtitle}
                      </p>
                    </div>
                  </div>
                  <div className="px-2 py-1 rounded-md bg-primary/10 text-primary text-xs font-medium whitespace-nowrap">
                    {item.badge}
                  </div>
                </div>
                <p className="text-sm text-muted-foreground">
                  {item.description}
                </p>
              </div>

              <div className="border-b border-border bg-muted/5">
                <div className="px-4 py-2 text-xs font-mono text-muted-foreground bg-muted/30">
                  input.jsson
                </div>
                <div className="max-h-[180px] overflow-auto">
                  <CodeBlock
                    code={item.code}
                    language="jsson"
                    className="p-4 text-sm"
                  />
                </div>
              </div>

              <div className="bg-card min-h-[200px]">
                <div className="px-4 py-2 text-xs font-mono text-muted-foreground border-b border-border bg-muted/10">
                  output
                </div>
                {item.visual}
              </div>

              <div className="absolute inset-0 bg-linear-to-br from-primary/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none" />
            </motion.div>
          ))}
        </div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.6 }}
          className="text-center mt-12"
        >
          <p className="text-muted-foreground mb-4">
            Ready to supercharge your data generation?
          </p>
          <a
            href="/playground"
            className="inline-flex items-center gap-2 px-6 py-3 rounded-lg bg-primary text-primary-foreground font-medium hover:bg-primary/90 transition-colors"
          >
            Try It in Playground
          </a>
        </motion.div>
      </div>
    </section>
  );
}
