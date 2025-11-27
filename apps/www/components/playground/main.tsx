"use client";

import { useTranspiler } from "@/hooks/use-transpiler";
import { JSSONEditor } from "@/components/playground/editor";
import { OutputViewer } from "@/components/playground/output-viewer";

const DEFAULT_CODE = `// Welcome to the JSSON Playground!
// Try editing this code to see the magic happen.

server {
  port = 8080
  host = "localhost"
  debug = true
  
  // Database configuration
  database {
    type = "postgres"
    url = "postgres://user:pass@localhost:5432/db"
  }
}

// Generate some users
users [
  template { id, name, role }
  
  map (u) = {
    id = u.id
    name = u.name
    role = u.role
    active = true
  }
  
  1, "Alice", "admin"
  2, "Bob", "user"
  3, "Charlie", "user"
]`;

export default function MainPlayground() {
  const { code, setCode, error } =
    useTranspiler(DEFAULT_CODE);

  return (
    <main className="flex-1 flex gap-4 p-4 min-h-0">
      <div className="flex-1 min-w-0">
        <JSSONEditor value={code} onChange={(val) => setCode(val || "")} />
      </div>
      <div className="flex-1 min-w-0">
        <OutputViewer jssonCode={code} error={error} />
      </div>
    </main>
  );
}
