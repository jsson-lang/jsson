// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";
import starlightThemeRapide from "starlight-theme-rapide";

// https://astro.build/config
export default defineConfig({
  site: "https://jsson-docs.vercel.app",
  integrations: [
    starlight({
      plugins: [starlightThemeRapide()],

      title: "JSSON Docs",
      description:
        "JSSON - The human-friendly way to write JSON. A modern syntax that transpiles to 100% valid JSON.",
      logo: {
        src: "./src/assets/logo.svg",
        alt: "JSSON Logo",
      },
      head: [
        {
          tag: "meta",
          attrs: {
            name: "google-site-verification",
            content: "Bco0meN-73Wimh1fOAJS4gtnEdGqooYR5zKQOfH0CkU",
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
              "JSON, JSSON, JSON syntax, transpiler, configuration files, human-friendly JSON, JSON generator, JSON template, infrastructure as code, Kubernetes config, API gateway configuration, i18n translations, feature flags, database seeding, geographic data, coordinate generation, multi-environment config, DevOps tools, configuration management, JSON alternative, YAML alternative, data generation, map transformation, conditional logic, range syntax, template arrays",
          },
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
            { label: "E-commerce Variants", slug: "real-world/ecommerce-variants" },
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
      ],
      customCss: [
        // Relative path to your custom CSS file
        "./src/styles/custom.css",
      ],
    }),
  ],
});
