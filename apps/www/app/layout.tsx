import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import { AnchoredToastProvider, ToastProvider } from "@/components/ui/toast"
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "JSSON — JavaScript Simplified Object Notation",
  description:
    "JSSON is a meta-format that generates JSON, YAML, TOML and TypeScript using templates, logic and smart ranges.",
  keywords: [
    "JSSON",
    "JSON",
    "YAML",
    "TOML",
    "TypeScript",
    "config",
    "transpiler",
    "meta-format",
  ],
  verification: {
    google: "Bco0meN-73Wimh1fOAJS4gtnEdGqooYR5zKQOfH0CkU"
  },
  robots: {
    index: true,
    follow: true,
  },
  openGraph: {
    title: "JSSON — The Config Meta-Language",
    description:
      "Write once, generate everywhere. Templates, maps, logic, ranges and multi-format output.",
    url: "https://jsson.vercel.app",
    siteName: "JSSON",
    images: [
      {
        url: "/og-image.png",
        width: 1200,
        height: 630,
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    creator: "@jssonlang",
    images: ["/og-image.png"],
  },
  icons: {
    icon: "/favicon.svg",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <ToastProvider position="top-center">
          <AnchoredToastProvider>
            {children}
          </AnchoredToastProvider>
        </ToastProvider>
      </body>
    </html>
  );
}
