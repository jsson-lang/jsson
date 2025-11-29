export default function highlightLine(
    line: string,
    language: "jsson" | "json"
): string {
    if (!line) return "&nbsp;";

    // Escape HTML
    let escaped = line
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;");

    const placeholders: string[] = [];
    const ph = (cls: string, content: string) => {
        placeholders.push(`<span class="${cls}">${content}</span>`);
        return `###PH${placeholders.length - 1}###`;
    };

    /* ------------------------------------------------------------- */
    /* 🔹 COMMENTS                                                    */
    /* ------------------------------------------------------------- */
    if (escaped.trim().startsWith("//")) {
        return `<span class="text-muted-foreground">${escaped}</span>`;
    }

    /* ------------------------------------------------------------- */
    /* 🔹 STRINGS (must be early)                                    */
    /* ------------------------------------------------------------- */
    escaped = escaped.replace(/"([^"]*)"/g, (m) =>
        ph("text-green-400", m)
    );

    /* ------------------------------------------------------------- */
    /* 🔹 NUMBERS                                                     */
    /* ------------------------------------------------------------- */
    escaped = escaped.replace(/\b(\d+(\.\d+)?)\b/g, (m) =>
        ph("text-orange-400", m)
    );

    /* ------------------------------------------------------------- */
    /* 🔹 BOOLEANS                                                    */
    /* ------------------------------------------------------------- */
    escaped = escaped.replace(/\b(true|false|null)\b/g, (m) =>
        ph("text-blue-400", m)
    );

    /* ------------------------------------------------------------- */
    /* 🔹 SPECIAL CHARACTERS {}[]                                     */
    /* ------------------------------------------------------------- */
    escaped = escaped.replace(
        /([\{\}\[\]])/g,
        (m) => ph("text-yellow-400", m)
    );

    /* ------------------------------------------------------------- */
    /* 🔹 JSSON-SPECIFIC HIGHLIGHTS                                   */
    /* ------------------------------------------------------------- */
    if (language === "jsson") {
        // Keywords
        escaped = escaped.replace(
            /\b(template|map|include|server|api)\b/g,
            (m) => ph("text-purple-400 font-semibold", m)
        );

        // Identifiers before "=" (JSSON assignment)
        escaped = escaped.replace(
            /([a-zA-Z_][a-zA-Z0-9_]*)(?=\s*=)/g,
            (m) => ph("text-foreground", m)
        );

        // Operators
        escaped = escaped.replace(
            /(=|\+|-|\*|\/|\?|\:)/g,
            (m) => ph("text-pink-400", m)
        );
    }

    /* ------------------------------------------------------------- */
    /* 🔹 JSON Keys (only in JSON mode)                              */
    /* ------------------------------------------------------------- */
    if (language === "json") {
        escaped = escaped.replace(
            /"([^"]+)"(?=\s*:)/g,
            (m) => ph("text-foreground font-semibold", m)
        );
    }

    /* ------------------------------------------------------------- */
    /* 🔹 CLEANUP — restore placeholders                             */
    /* ------------------------------------------------------------- */
    return escaped.replace(
        /###PH(\d+)###/g,
        (_, i) => placeholders[parseInt(i)] ?? ""
    );
}