import { Metadata } from "next";
import { Button } from "@/components/ui/button";
import MainPlayground from "@/components/playground/main";
import Logo from "@/components/shared/logo";
import Link from "next/link";
import { Stats } from "@/components/playground/stats";
import { PlaygroundProvider } from "@/contexts/playground-context";

export const metadata: Metadata = {
  title: "Playground - JSSON",
  description: "Playground for JSSON - JavaScript Simplified Object Notation",
};

export default function PlaygroundPage() {
  return (
    <PlaygroundProvider>
      <div className="flex flex-col h-screen overflow-hidden">
        <header className="flex items-center justify-between px-6 py-3 border-b border-border bg-background/80 backdrop-blur-md z-10">
          <div className="flex items-center gap-4">
            <Logo size="md" />
            <div className="h-6 w-px bg-border" />
            <span className="text-sm font-medium text-muted-foreground">
              Playground
            </span>
          </div>
          <Stats />
          <Link href={"https://docs.jssonlang.tech/"}>
            <Button size="sm" variant={"ghost"}>
              Go to docs
            </Button>
          </Link>
        </header>

        <MainPlayground />
      </div>
    </PlaygroundProvider>
  );
}
