import { Metadata } from "next";
import Link from "next/link";
import Logo from "@/components/shared/logo";

export const metadata: Metadata = {
  title: "Terms of Service - JSSON",
  description:
    "Terms of Service for JSSON - JavaScript Simplified Object Notation",
};

export default function TermsPage() {
  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border">
        <div className="container mx-auto px-4 py-4">
          <Logo size="md" />
        </div>
      </header>

      <main className="container mx-auto px-4 py-12 max-w-4xl">
        <h1 className="text-4xl font-bold mb-2">Terms of Service</h1>
        <p className="text-muted-foreground mb-8">
          Last updated: November 26, 2025
        </p>

        <div className="prose prose-neutral dark:prose-invert max-w-none">
          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">
              1. Acceptance of Terms
            </h2>
            <p className="text-muted-foreground mb-4">
              By accessing and using JSSON (the "Service"), you accept and agree
              to be bound by the terms and provision of this agreement.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">2. Use License</h2>
            <p className="text-muted-foreground mb-4">
              JSSON is provided as open-source software under the MIT License.
              You are free to:
            </p>
            <ul className="list-disc pl-6 text-muted-foreground space-y-2">
              <li>
                Use the software for any purpose, including commercial
                applications
              </li>
              <li>Modify the source code</li>
              <li>Distribute copies of the software</li>
              <li>Sublicense and sell copies of the software</li>
            </ul>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">3. Disclaimer</h2>
            <p className="text-muted-foreground mb-4">
              The Service is provided "as is" without warranty of any kind,
              express or implied. We do not guarantee that:
            </p>
            <ul className="list-disc pl-6 text-muted-foreground space-y-2">
              <li>The Service will meet your specific requirements</li>
              <li>
                The Service will be uninterrupted, timely, secure, or error-free
              </li>
              <li>
                The results obtained from using the Service will be accurate or
                reliable
              </li>
            </ul>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">4. Limitations</h2>
            <p className="text-muted-foreground mb-4">
              In no event shall JSSON or its contributors be liable for any
              damages (including, without limitation, damages for loss of data
              or profit, or due to business interruption) arising out of the use
              or inability to use the Service.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">5. Playground Usage</h2>
            <p className="text-muted-foreground mb-4">
              The JSSON Playground runs entirely in your browser using
              WebAssembly. We do not collect, store, or transmit any code you
              write in the playground. All transpilation happens locally on your
              device.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">6. Modifications</h2>
            <p className="text-muted-foreground mb-4">
              We reserve the right to modify these terms at any time. Continued
              use of the Service after changes constitutes acceptance of the new
              terms.
            </p>
          </section>

          <section className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">7. Contact</h2>
            <p className="text-muted-foreground mb-4">
              For questions about these Terms of Service, please contact us
              through our{" "}
              <Link
                href="https://github.com/carlosedujs/jsson"
                className="text-primary hover:underline"
              >
                GitHub repository
              </Link>
              .
            </p>
          </section>
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
