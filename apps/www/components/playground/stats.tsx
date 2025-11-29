"use client";

import { ChartCandlestick } from "lucide-react";
import { Button } from "../ui/button";
import {
  Dialog,
  DialogTrigger,
  DialogPopup,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogPanel,
} from "../ui/dialog";
import { StatsOverview } from "./stats-overview";
import { calculateStats } from "./stats-utils";
import { usePlaygroundContext } from "@/contexts/playground-context";

export function Stats() {
  const { jssonCode, output, format } = usePlaygroundContext();
  const stats = calculateStats(jssonCode, output);
  const hasData = jssonCode.length > 0 || output.length > 0;

  return (
    <Dialog>
      <DialogTrigger render={<Button />}>
        <ChartCandlestick />
        Productivity Overview
      </DialogTrigger>
      <DialogPopup className="md:max-w-5xl">
        <DialogHeader>
          <DialogTitle className={"text-2xl font-bold tracking-tight"}>
            JSSON Productivity Overview
          </DialogTitle>
          <DialogDescription className={"text-sm text-muted-foreground"}>
            Comprehensive statistics about your JSSON transpilation efficiency
          </DialogDescription>
        </DialogHeader>
        <DialogPanel>
          {hasData ? (
            <StatsOverview stats={stats} format={format} />
          ) : (
            <div className="py-12 text-center text-muted-foreground">
              <ChartCandlestick className="mx-auto h-12 w-12 mb-4 opacity-50" />
              <p>Write some JSSON code to see productivity stats</p>
            </div>
          )}
        </DialogPanel>
      </DialogPopup>
    </Dialog>
  );
}
