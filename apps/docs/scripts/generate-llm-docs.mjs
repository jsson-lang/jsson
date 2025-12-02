import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const DOCS_DIR = path.join(__dirname, "../src/content/docs");
const OUTPUT_DIR = path.join(__dirname, "../public/llms.txt");

const MAX_CHUNK_TOKENS = 2000;

function countTokens(text) {
  return Math.ceil(text.split(/\s+/).length * 1.3);
}

function chunkText(text) {
  const words = text.split(" ");
  let chunks = [];
  let current = [];

  for (const word of words) {
    current.push(word);
    if (countTokens(current.join(" ")) > MAX_CHUNK_TOKENS) {
      chunks.push(current.join(" "));
      current = [];
    }
  }

  if (current.length > 0) chunks.push(current.join(" "));
  return chunks;
}

function processMDX(raw, fileName) {
  let txt = raw;

  // Remove frontmatter YAML
  txt = txt.replace(/^---[\s\S]*?---/, "");

  // Remove imports
  txt = txt.replace(/^import .*$/gm, "");

  // Remove HTML/JSX comments
  txt = txt.replace(/<!--[\s\S]*?-->/g, "");

  // Remove MDX JSX components fully
  txt = txt.replace(/<[A-Z][^>]*>[\s\S]*?<\/[A-Z][^>]*>/g, "");
  txt = txt.replace(/<[A-Z][^>]*\/>/g, "");

  // Convert headings ## ===> "1.2 Section"
  let h1 = 0,
    h2 = 0,
    h3 = 0;
  txt = txt.replace(/^(#{1,3})\s+(.*)$/gm, (_, hashes, title) => {
    if (hashes.length === 1) {
      h1++;
      h2 = h3 = 0;
      return `\n${h1}. ${title.toUpperCase()}\n${"-".repeat(title.length)}`;
    }
    if (hashes.length === 2) {
      h2++;
      h3 = 0;
      return `${h1}.${h2}. ${title}`;
    }
    if (hashes.length === 3) {
      h3++;
      return `${h1}.${h2}.${h3}. ${title}`;
    }
  });

  // Convert tables
  txt = txt.replace(/\|([^]+?)\|/g, (match) => {
    return match
      .split("\n")
      .map((row) => row.replace(/\|/g, " • ").replace(/-\s+-/g, ""))
      .join("\n");
  });

  // Convert bullet lists
  txt = txt.replace(/^\s*[-*+]\s+/gm, "• ");

  // Normalize inline markdown
  txt = txt.replace(/`([^`]+)`/g, "$1");
  txt = txt.replace(/\*\*([^*]+)\*\*/g, "$1");
  txt = txt.replace(/\*([^*]+)\*/g, "$1");

  // Collapse blank lines
  txt = txt.replace(/\n{3,}/g, "\n\n");

  txt = txt.trim();

  // Generate TOC
  const lines = txt.split("\n");
  const toc = lines
    .filter((l) => /^\d+(\.\d+)*\./.test(l))
    .map((l) => `- ${l}`)
    .join("\n");

  const meta = [
    `# JSSON Documentation — LLM Mode`,
    `Source: ${fileName}`,
    `Version: 0.0.5.2`,
    `Generated: ${new Date().toISOString()}`,
    "",
    `## Table of Contents`,
    toc,
    "",
    `## Content`,
  ].join("\n");

  return meta + "\n\n" + txt;
}

function processDirectory(sourceDir, outputDir) {
  const entries = fs.readdirSync(sourceDir, { withFileTypes: true });

  for (const entry of entries) {
    const sourcePath = path.join(sourceDir, entry.name);

    if (entry.isDirectory()) {
      const dirOutputPath = path.join(outputDir, entry.name);
      fs.mkdirSync(dirOutputPath, { recursive: true });
      processDirectory(sourcePath, dirOutputPath);
    } else if (entry.name.endsWith(".mdx")) {
      const raw = fs.readFileSync(sourcePath, "utf-8");
      const cleaned = processMDX(raw, entry.name);

      const chunks = chunkText(cleaned);
      const baseName = entry.name.replace(".mdx", "");

      if (chunks.length === 1) {
        const outputPath = path.join(outputDir, `${baseName}.txt`);
        fs.writeFileSync(outputPath, cleaned);
      } else {
        chunks.forEach((chunk, i) => {
          const outputPath = path.join(outputDir, `${baseName}.${i + 1}.txt`);
          fs.writeFileSync(outputPath, chunk);
        });
      }

      console.log(`✓ Processed ${entry.name}`);
    }
  }
}

function generateIndex() {
  let files = [];

  function scan(dir, prefix = "") {
    const entries = fs.readdirSync(dir, { withFileTypes: true });

    for (const e of entries) {
      const full = path.join(dir, e.name);
      const rel = path.join(prefix, e.name);

      if (e.isDirectory()) scan(full, rel);
      else if (e.name.endsWith(".txt")) files.push(rel.replace(/\\/g, "/"));
    }
  }

  scan(OUTPUT_DIR);

  const lines = [
    "# JSSON LLM Documentation Index",
    "Generated: " + new Date().toISOString(),
    "",
    "## Files",
    ...files.map((f) => `- /llms.txt/${f}`),
  ];

  fs.writeFileSync(path.join(OUTPUT_DIR, "index.txt"), lines.join("\n"));
}

console.log("🚀 Generating PRO-level LLM Docs...\n");

fs.rmSync(OUTPUT_DIR, { recursive: true, force: true });
fs.mkdirSync(OUTPUT_DIR);

processDirectory(DOCS_DIR, OUTPUT_DIR);
generateIndex();

console.log("\n✨ Done! LLM-ready docs generated at /public/llms.txt/");
