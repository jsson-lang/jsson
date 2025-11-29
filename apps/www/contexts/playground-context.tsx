"use client";

import { createContext, useContext, useState, ReactNode } from "react";

interface PlaygroundContextType {
  jssonCode: string;
  setJssonCode: (code: string) => void;
  output: string;
  setOutput: (output: string) => void;
  format: string;
  setFormat: (format: string) => void;
}

const PlaygroundContext = createContext<PlaygroundContextType | undefined>(
  undefined
);

export function PlaygroundProvider({ children }: { children: ReactNode }) {
  const [jssonCode, setJssonCode] = useState("");
  const [output, setOutput] = useState("");
  const [format, setFormat] = useState("json");

  return (
    <PlaygroundContext.Provider
      value={{ jssonCode, setJssonCode, output, setOutput, format, setFormat }}
    >
      {children}
    </PlaygroundContext.Provider>
  );
}

export function usePlaygroundContext() {
  const context = useContext(PlaygroundContext);
  if (context === undefined) {
    throw new Error(
      "usePlaygroundContext must be used within a PlaygroundProvider"
    );
  }
  return context;
}
