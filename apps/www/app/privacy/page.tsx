import { Metadata } from "next";
import Link from "next/link";
import Logo from "@/components/shared/logo";

export const metadata: Metadata = {
  title: "Privacy Policy - JSSON",
  description:
    "Privacy Policy for JSSON - JavaScript Simplified Object Notation",
};

export default function PrivacyPage() {
  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border">
        <div className="container mx-auto px-4 py-4">
          <Logo size="md" />
        </div>
      </header>

      <main className="container mx-auto px-4 py-12 max-w-4xl">
        <h1 className="text-4xl font-bold mb-2">Privacy Policy</h1>
        <p className="text-muted-foreground mb-8">
          Last updated: November 26, 2025
        </p>

        <div className="prose prose-neutral dark:prose-invert max-w-none">
          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">1. Introduction</h2>
            <p className="text-muted-foreground mb-4">
              JSSON ("we", "our", or "us") is committed to protecting your
              privacy. This Privacy Policy explains how we handle information
              when you use our website and services.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">
              2. Information We Don't Collect
            </h2>
            <p className="text-muted-foreground mb-4">
              JSSON is designed with privacy in mind:
            </p>
            <ul className="list-disc pl-6 text-muted-foreground space-y-2">
              <li>
                <strong>No code storage:</strong> All code written in the
                Playground runs locally in your browser via WebAssembly. We
                never see, store, or transmit your code.
              </li>
              <li>
                <strong>No user accounts:</strong> We don't require registration
                or collect personal information.
              </li>
              <li>
                <strong>No tracking:</strong> We don't use analytics, cookies,
                or tracking scripts.
              </li>
              <li>
                <strong>No server-side processing:</strong> Transpilation
                happens entirely in your browser.
              </li>
            </ul>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">
              3. Third-Party Services
            </h2>
            <p className="text-muted-foreground mb-4">
              Our website is hosted on Vercel. Vercel may collect standard
              server logs (IP addresses, browser types, timestamps) for
              infrastructure purposes. Please refer to{" "}
              <a
                href="https://vercel.com/legal/privacy-policy"
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary hover:underline"
              >
                Vercel's Privacy Policy
              </a>{" "}
              for more information.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">4. Open Source</h2>
            <p className="text-muted-foreground mb-4">
              JSSON is open-source software. You can review our entire codebase
              on{" "}
              <Link
                href="https://github.com/carlosedujs/jsson"
                className="text-primary hover:underline"
              >
                GitHub
              </Link>{" "}
              to verify our privacy claims.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">5. Local Storage</h2>
            <p className="text-muted-foreground mb-4">
              The Playground may use your browser's local storage to save your
              code locally on your device. This data never leaves your browser
              and can be cleared at any time through your browser settings.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">
              6. Children's Privacy
            </h2>
            <p className="text-muted-foreground mb-4">
              JSSON does not knowingly collect any information from anyone,
              including children under 13. Our service is designed to be
              privacy-first for all users.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">
              7. Changes to This Policy
            </h2>
            <p className="text-muted-foreground mb-4">
              We may update this Privacy Policy from time to time. We will
              notify users of any material changes by updating the "Last
              updated" date at the top of this policy.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">8. Contact Us</h2>
            <p className="text-muted-foreground mb-4">
              If you have any questions about this Privacy Policy, please
              contact us through our{" "}
              <Link
                href="https://github.com/carlosedujs/jsson"
                className="text-primary hover:underline"
              >
                GitHub repository
              </Link>
              .
            </p>
          </section>

          <div className="mt-8 p-6 bg-muted/50 rounded-lg border border-border">
            <p className="text-sm text-muted-foreground">
              <strong>TL;DR:</strong> We don't collect, store, or track
              anything. Your code stays in your browser. JSSON is privacy-first
              by design.
            </p>
          </div>
        </div>

        <div className="mt-12 pt-8 border-t border-border">
          <Link href="/" className="text-primary hover:underline">
            ← Back to Home
          </Link>
        </div>
      </main>
    </div>
  );
}
