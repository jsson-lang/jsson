"use client";

import { buttonVariants } from "@/components/ui/button";
import { motion } from "motion/react";
import { ArrowRight, Terminal } from "lucide-react";
import Link from "next/link";
import { cn } from "@/lib/utils";
import TerminalCommand from "../shared/terminal-command";
import { CodeBlock } from "../shared/code-block";
import { Spotlight } from "@/components/ui/spotlight";

import {
  Announcement,
  AnnouncementTag,
  AnnouncementTitle,
} from "@/components/kibo-ui/announcement";
import { ArrowUpRightIcon } from "lucide-react";

export function Hero() {
  return (
    <section className="relative overflow-hidden py-24 lg:py-32 h-screen bg-linear-to-b from-background/10 to-foreground/10">
      <Spotlight />
      <div className="container  relative z-10 mx-auto px-4 md:px-6">
        <div className="flex flex-col items-center text-center">
          <Link href={"https://jsson-docs.vercel.app/changelog/#v005-2025-11-27"}>
            <Announcement variant="outline" themed className="max-w-4xl">
              <AnnouncementTag>Latest Update</AnnouncementTag>
              <AnnouncementTitle>
                v0.0.5 is now available
                <ArrowUpRightIcon className="ml-2 h-4 w-4" />
              </AnnouncementTitle>
            </Announcement>
          </Link>

          <motion.h1
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.1 }}
            className="mt-8 max-w-4xl text-5xl font-bold tracking-tight md:text-7xl lg:text-8xl bg-linear-to-b from-foreground to-foreground/70 bg-clip-text text-transparent"
          >
            Universal Config Format <br />
            <span className="text-primary">Supercharged</span>
          </motion.h1>

          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
            className="mt-6 max-w-2xl text-lg text-muted-foreground md:text-xl"
          >
            JSSON isn’t just another format — it’s a meta-format. Write your
            logic once, and let it spit out JSON, YAML, TOML, or fully typed
            TypeScript. No repetition. No boilerplate. No crying.
          </motion.p>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.3 }}
            className="mt-10 flex flex-col gap-4 sm:flex-row"
          >
            <Link
              href="/playground"
              className={cn(
                buttonVariants({ size: "lg" }),
                "h-12 px-8 text-base"
              )}
            >
              Start Coding <ArrowRight className="ml-2 h-4 w-4" />
            </Link>
            <Link
              href="https://jsson-docs.vercel.app"
              className={cn(
                buttonVariants({ size: "lg", variant: "outline" }),
                "h-12 px-8 text-base"
              )}
            >
              Documentation
            </Link>
          </motion.div>

          <TerminalCommand filename="terminal CLI">
            <CodeBlock code={`jsson --version`} />
          </TerminalCommand>
        </div>
      </div>

      {/* Background Gradients */}
      <div className="absolute top-0 -left-4 w-72 h-72 bg-primary/30 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob" />
      <div className="absolute top-0 -right-4 w-72 h-72 bg-purple-500/30 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-2000" />
      <div className="absolute -bottom-8 left-20 w-72 h-72 bg-pink-500/30 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-4000" />
    </section>
  );
}
