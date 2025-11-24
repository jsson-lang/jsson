// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";
import starlightThemeRapide from "starlight-theme-rapide";

// https://astro.build/config
export default defineConfig({
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
              "JSON, JSSON, JSON syntax, transpiler, configuration, human-friendly JSON",
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
        { label: "FAQ", slug: "faq" },
      ],
      customCss: [
        // Relative path to your custom CSS file
        "./src/styles/custom.css",
      ],
    }),
  ],
});
