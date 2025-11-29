"use client";

import { StatsData } from "./stats-utils";
import {
  PieChart,
  Pie,
  Cell,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import {
  ArrowRight,
  Box,
  FileText,
  Hash,
  Layers,
  Type,
  Zap,
  Activity,
} from "lucide-react";

interface StatsOverviewProps {
  stats: StatsData;
  format: string;
}

export function StatsOverview({ stats, format }: StatsOverviewProps) {
  // Data for donut chart
  const donutData = [
    {
      name: "Input",
      value: stats.input.chars,
      percentage: stats.ratios.inputPercentage,
    },
    {
      name: format.toUpperCase(),
      value: stats.output.chars,
      percentage: stats.ratios.outputPercentage,
    },
  ];

  // Consolidated data for comparison
  const comparisonData = [
    {
      name: "Lines",
      Input: stats.input.lines,
      Output: stats.output.lines,
      ratio: stats.expansion.lines,
    },
    {
      name: "Chars",
      Input: stats.input.chars,
      Output: stats.output.chars,
      ratio: stats.expansion.chars,
    },
    {
      name: "Tokens",
      Input: stats.input.tokens,
      Output: stats.output.tokens,
      ratio: stats.expansion.tokens,
    },
  ];

  const COLORS = {
    input: "#8b5cf6", // purple-500
    output: "#10b981", // emerald-500
    grid: "#27272a", // zinc-800
    text: "#a1a1aa", // zinc-400
  };

  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      return (
        <div className="rounded-lg border border-border bg-popover p-3 shadow-md">
          <p className="mb-2 font-medium text-popover-foreground">{label}</p>
          {payload.map((entry: any, index: number) => (
            <div key={index} className="flex items-center gap-2 text-sm">
              <div
                className="h-2 w-2 rounded-full"
                style={{ backgroundColor: entry.color }}
              />
              <span className="text-muted-foreground">{entry.name}:</span>
              <span className="font-mono font-medium text-foreground">
                {entry.value.toLocaleString()}
              </span>
            </div>
          ))}
        </div>
      );
    }
    return null;
  };

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      {/* Header Section */}
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h2 className="text-2xl font-bold tracking-tight flex items-center gap-2">
            <Activity className="h-6 w-6 text-primary" />
            Stats Overview
          </h2>
          <p className="text-sm text-muted-foreground">
            Real-time analysis of your JSSON transformation
          </p>
        </div>
        {stats.expansion.lines > 1 && (
          <div className="hidden sm:flex items-center gap-2 rounded-full bg-primary px-4 py-1.5 text-sm font-medium text-background">
            <Zap className="h-4 w-4" />
            <span className="text-background">
              {stats.expansion.lines.toFixed(1)}x Efficiency Boost
            </span>
          </div>
        )}
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        <div className="rounded-xl bg-card p-4 shadow-sm">
          <div className="flex items-center justify-between pb-2">
            <p className="text-sm font-medium text-muted-foreground">Input</p>
            <FileText className="h-4 w-4 text-purple-500" />
          </div>
          <div className="space-y-1">
            <div className="text-2xl font-bold">{stats.input.lines}</div>
            <p className="text-xs text-muted-foreground">lines of code</p>
          </div>
          <div className="mt-3 flex items-center gap-3 text-xs text-muted-foreground">
            <div className="flex items-center gap-1">
              <Type className="h-3 w-3" /> {stats.input.chars}
            </div>
            <div className="flex items-center gap-1">
              <Hash className="h-3 w-3" /> ~{stats.input.tokens}
            </div>
          </div>
        </div>

        <div className="rounded-xl bg-card p-4 shadow-sm">
          <div className="flex items-center justify-between pb-2">
            <p className="text-sm font-medium text-muted-foreground">
              {format.toUpperCase()} Output
            </p>
            <Box className="h-4 w-4 text-emerald-500" />
          </div>
          <div className="space-y-1">
            <div className="text-2xl font-bold">{stats.output.lines}</div>
            <p className="text-xs text-muted-foreground">lines generated</p>
          </div>
          <div className="mt-3 flex items-center gap-3 text-xs text-muted-foreground">
            <div className="flex items-center gap-1">
              <Type className="h-3 w-3" /> {stats.output.chars}
            </div>
            <div className="flex items-center gap-1">
              <Hash className="h-3 w-3" /> ~{stats.output.tokens}
            </div>
          </div>
        </div>

        <div className="rounded-xl bg-card p-4 shadow-sm relative overflow-hidden group">
          <div className="absolute inset-0 bg-linear-to-br from-background/10 to-foreground/10 opacity-0 group-hover:opacity-100 transition-opacity" />
          <div className="flex items-center justify-between pb-2 relative">
            <p className="text-sm font-medium text-muted-foreground">
              Expansion
            </p>
            <Layers className="h-4 w-4 text-blue-500" />
          </div>
          <div className="space-y-1 relative">
            <div className="text-2xl font-bold text-foreground">
              {stats.expansion.lines.toFixed(1)}x
            </div>
            <p className="text-xs text-emerald-500 font-medium flex items-center">
              <ArrowRight className="h-3 w-3 mr-1" />
              Multiplication Factor
            </p>
          </div>
          <div className="mt-3 text-xs text-muted-foreground relative">
            {stats.input.lines} lines → {stats.output.lines} lines
          </div>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-7">
        <div className="col-span-4 rounded-xl bg-card p-6 shadow-sm">
          <div className="mb-6">
            <h3 className="font-semibold">Growth Metrics</h3>
            <p className="text-sm text-muted-foreground">
              Side-by-side comparison of structure size
            </p>
          </div>
          <div className="h-[250px] w-full">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart
                data={comparisonData}
                margin={{ top: 10, right: 10, left: -20, bottom: 0 }}
                barSize={32}
              >
                <CartesianGrid
                  strokeDasharray="3 3"
                  vertical={false}
                  stroke={COLORS.grid}
                  opacity={0.4}
                />
                <XAxis
                  dataKey="name"
                  axisLine={false}
                  tickLine={false}
                  tick={{ fill: COLORS.text, fontSize: 12 }}
                  dy={10}
                />
                <YAxis
                  axisLine={false}
                  tickLine={false}
                  tick={{ fill: COLORS.text, fontSize: 12 }}
                />
                <Tooltip
                  cursor={{ fill: "transparent" }}
                  content={CustomTooltip}
                />
                <Legend
                  wrapperStyle={{ paddingTop: "20px" }}
                  iconType="circle"
                  iconSize={8}
                />
                <Bar
                  dataKey="Input"
                  fill={COLORS.input}
                  radius={[4, 4, 0, 0]}
                  animationDuration={1000}
                />
                <Bar
                  dataKey="Output"
                  fill={COLORS.output}
                  radius={[4, 4, 0, 0]}
                  animationDuration={1000}
                />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="col-span-3 rounded-xl bg-card p-6 shadow-sm">
          <div className="mb-6">
            <h3 className="font-semibold">Volume Distribution</h3>
            <p className="text-sm text-muted-foreground">
              Character count breakdown
            </p>
          </div>
          <div className="h-[250px] w-full relative">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={donutData}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={85}
                  paddingAngle={4}
                  dataKey="value"
                  stroke="none"
                >
                  <Cell fill={COLORS.input} />
                  <Cell fill={COLORS.output} />
                </Pie>
                <Tooltip content={CustomTooltip} />
                <Legend
                  verticalAlign="bottom"
                  height={36}
                  iconType="circle"
                  iconSize={8}
                />
              </PieChart>
            </ResponsiveContainer>

            <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none pb-8">
              <span className="text-2xl font-bold">
                {stats.expansion.chars.toFixed(1)}x
              </span>
              <span className="text-xs text-muted-foreground">Growth</span>
            </div>
          </div>
        </div>
      </div>

      {stats.expansion.lines > 1 && (
        <div className="mt-2 rounded-lg border-purple-500/20 bg-purple-500/5 p-4 flex items-start gap-3">
          <div className="rounded-full bg-purple-500/10 p-1.5 mt-0.5">
            <Zap className="h-4 w-4 text-purple-500" />
          </div>
          <div className="space-y-1">
            <p className="text-sm font-medium text-foreground">
              High Efficiency Detected
            </p>
            <p className="text-xs text-muted-foreground leading-relaxed">
              You generated{" "}
              <span className="font-medium text-foreground">
                {stats.output.lines} lines
              </span>{" "}
              of output from just{" "}
              <span className="font-medium text-foreground">
                {stats.input.lines} lines
              </span>{" "}
              of JSSON. That's a{" "}
              <span className="font-medium text-purple-500">
                {stats.expansion.lines.toFixed(1)}x
              </span>{" "}
              productivity boost compared to writing raw {format.toUpperCase()}.
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
