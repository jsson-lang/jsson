"use client";

import { cn } from "@/lib/utils";
import Image from "next/image";
import Link from "next/link";
import { Badge } from "../ui/badge";

interface LogoProps {
  size?: "sm" | "md" | "lg" | "xl";
  variant?: "default" | "icon";
}

const textSizes = {
  sm: "text-sm",
  md: "text-md",
  lg: "text-lg",
  xl: "text-xl",
};

export default function Logo({ size = "xl", variant = "default" }: LogoProps) {
  return (
    <Link
      href="/"
      className={cn("flex items-center gap-2 font-bold", textSizes[size])}
    >
      <div className="flex h-8 w-8 items-center justify-center rounded-lg">
        <Image
          src="/logo.svg"
          alt="JSSON"
          width={
            size === "sm" ? 24 : size === "md" ? 32 : size === "lg" ? 48 : 64
          }
          height={
            size === "sm" ? 24 : size === "md" ? 32 : size === "lg" ? 48 : 64
          }
          className={cn("rounded-md", variant === "icon" && textSizes[size])}
        />
      </div>
      {variant === "default" && (
        <>
          <span className="text-primary">JSSON</span>
          <Badge variant={"secondary"}>V0.0.5.2</Badge>
        </>
      )}
    </Link>
  );
}
