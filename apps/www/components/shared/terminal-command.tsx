import { cn } from "@/lib/utils";
import { Terminal } from "lucide-react";
import { motion } from "motion/react";

interface TerminalCommandProps {
  filename: string;
  children: React.ReactNode;
  className?: string;
}

export default function TerminalCommand({
  filename,
  children,
  className,
}: TerminalCommandProps) {
  return (
    <motion.div 
      initial={{ opacity: 0, y: 40 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.7, delay: 0.4 }}
      className={cn("mt-20 w-fit", className)}
    >
      <div className="relative rounded-xl border border-border bg-card/50 shadow-2xl backdrop-blur-sm">
        <div className="flex items-center gap-2 border-b border-border px-4 py-3">
          <div className="flex gap-1.5">
            <div className="h-3 w-3 rounded-full bg-red-500/20" />
            <div className="h-3 w-3 rounded-full bg-yellow-500/20" />
            <div className="h-3 w-3 rounded-full bg-green-500/20" />
          </div>

          <div className="ml-4 flex items-center gap-2 rounded-md bg-muted/50 px-3 py-1 text-xs text-muted-foreground">
            <Terminal className="h-3 w-3" />
            {filename}
          </div>
        </div>

        <div className="p-6 text-left">{children}</div>
      </div>
    </motion.div>
  );
}
