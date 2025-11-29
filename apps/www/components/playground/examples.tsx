"use client";

import { useState } from "react";
import { Button } from "../ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetPanel,
  SheetTitle,
  SheetTrigger,
} from "../ui/sheet";
import {
  ChevronRight,
  BookOpen,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { EXAMPLES } from "./examples-data";

interface ExamplesSheetProps {
  onSelect: (code: string) => void;
}

export default function ExamplesSheet({ onSelect }: ExamplesSheetProps) {
  const [open, setOpen] = useState(false);

  const handleSelect = (code: string) => {
    onSelect(code);
    setOpen(false);
  };

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger
        render={
          <Button variant="outline" className="gap-2">
            <BookOpen className="h-4 w-4" />
            Examples
          </Button>
        }
      ></SheetTrigger>

      <SheetContent className="sm:max-w-xl overflow-y-auto">
        <SheetHeader className="mb-6">
          <SheetTitle className="text-2xl font-bold tracking-tight">
            Examples Library
          </SheetTitle>
          <SheetDescription>
            Explore real-world use cases and patterns. Click to load into the
            editor.
          </SheetDescription>
        </SheetHeader>

        <SheetPanel>
          <div className="grid gap-4">
            {EXAMPLES.map((example, index) => (
              <button
                key={index}
                onClick={() => handleSelect(example.code)}
                className={cn(
                  "flex items-start gap-4 rounded-xl border p-4 text-left transition-all hover:bg-accent hover:border-accent-foreground/20 group",
                  "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                )}
              >
                <div className="rounded-lg bg-primary/10 p-2 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-colors">
                  <example.icon className="h-5 w-5" />
                </div>
                <div className="flex-1 space-y-1">
                  <div className="flex items-center justify-between">
                    <span className="font-semibold tracking-tight">
                      {example.title}
                    </span>
                    <span className="text-[10px] uppercase tracking-wider font-medium text-muted-foreground bg-muted px-2 py-0.5 rounded-full">
                      {example.category}
                    </span>
                  </div>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    {example.description}
                  </p>
                </div>
                <ChevronRight className="h-4 w-4 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity self-center" />
              </button>
            ))}
          </div>
        </SheetPanel>
      </SheetContent>
    </Sheet>
  );
}
