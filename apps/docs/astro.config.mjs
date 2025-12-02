// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";
import starlightThemeRapide from "starlight-theme-rapide";

// https://astro.build/config
export default defineConfig({
  site: "https://docs.jssonlang.tech",
  integrations: [
    starlight({
      plugins: [starlightThemeRapide()],

      title: "JSSON Documentation",
      description:
        "Official JSSON documentation. Learn how to use JSSON, a powerful meta-language and transpiler for generating JSON, YAML, TOML, and TypeScript. Features templates, ranges, maps, and streaming support.",
      logo: {
        src: "./src/assets/logo.svg",
        alt: "JSSON Logo",
      },
      head: [
        {
          tag: "meta",
          attrs: {
            name: "google-site-verification",
            content: "6qeyW-HDbBat5RyhEyy94NNXlTpn_gyBGuX9kCi8EBw",
          },
        },
        {
          tag: "meta",
          attrs: {
            property: "og:image",
            content: "/og-image.png",
          },
        },
        {
          tag: "meta",
          attrs: {
            property: "og:type",
            content: "website",
          },
        },
        {
          tag: "meta",
          attrs: {
            name: "twitter:card",
            content: "summary_large_image",
          },
        },
        {
          tag: "meta",
          attrs: {
            name: "keywords",
            content:
              "JSSON, JSSON language, JSSON documentation, JSSON transpiler, JSSON syntax, config generator, multi-format transpiler, YAML generator, TOML generator, TypeScript generator, template syntax, range expressions, map transformations, streaming support, configuration files, infrastructure as code, Kubernetes config, API gateway, i18n translations, feature flags, database seeding, DevOps tools, configuration management, data generation tool",
          },
        },
        {
          tag: "script",
          attrs: {
            type: "application/ld+json",
          },
          content: JSON.stringify({
            "@context": "https://schema.org",
            "@type": "TechArticle",
            headline: "JSSON Documentation",
            description:
              "Complete documentation for JSSON - a meta-language and transpiler for generating JSON, YAML, TOML, and TypeScript",
            url: "https://docs.jssonlang.tech",
            author: {
              "@type": "Organization",
              name: "JSSON",
            },
            keywords: "JSSON, transpiler, documentation, config generator",
          }),
        },
      ],
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/carlosedujs/jsson",
        },
      ],
      sidebar: [
        {
          label: "Getting Started",
          items: [
            { label: "Introduction", slug: "guides/getting-started" },
            { label: "Basic Syntax", slug: "guides/basic-syntax" },
          ],
        },
        {
          label: "Guides",
          items: [
            { label: "Multi-Format Output", slug: "guides/multi-format" },
            { label: "Templates & Arrays", slug: "guides/templates" },
            { label: "Includes & Modules", slug: "guides/include-modules" },
            { label: "CLI Usage", slug: "guides/cli" },
            { label: "Transpiler Usage", slug: "guides/transpiler" },
            { label: "Advanced Patterns", slug: "guides/advanced-patterns" },
          ],
        },
        {
          label: "Reference",
          items: [
            { label: "Syntax Reference", slug: "reference/syntax" },
            { label: "Go API Reference", slug: "api/transpiler" },
            { label: "Errors & Debugging", slug: "reference/errors" },
            { label: "AST Reference", slug: "reference/ast" },
          ],
        },
        {
          label: "Examples",
          items: [
            { label: "Demo", slug: "examples/demo" },
            { label: "Templates", slug: "examples/template" },
          ],
        },
        {
          label: "Real-World Use Cases",
          items: [
            { label: "Overview", slug: "real-world/overview" },
            { label: "Geographic Data", slug: "real-world/geographic-data" },
            {
              label: "E-commerce Variants",
              slug: "real-world/ecommerce-variants",
            },
            { label: "Scheduling Matrix", slug: "real-world/scheduling" },
            { label: "Kubernetes Config", slug: "real-world/kubernetes" },
            { label: "API Gateway", slug: "real-world/api-gateway" },
            { label: "i18n Translations", slug: "real-world/i18n" },
            { label: "Feature Flags", slug: "real-world/feature-flags" },
            { label: "Database Seeding", slug: "real-world/database-seed" },
          ],
        },
        { label: "FAQ", slug: "faq" },
        { label: "Changelog", slug: "changelog" },
        { label: "LLM Docs", link: "/llms.txt/index.txt" },
        { label: "Playground", link: "https://jssonlang.tech/playground" },
      ],
      customCss: [
        // Relative path to your custom CSS file
        "./src/styles/custom.css",
      ],
    }),
  ],
});
