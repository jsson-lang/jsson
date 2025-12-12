# JSSON Language Server

Language Server Protocol (LSP) implementation for JSSON.

## Features

- **Real-time Diagnostics**: Syntax errors displayed as you type
- **Intelligent Autocomplete**: 
  - Keywords (`template`, `map`, `zip`, etc.)
  - Variables declared in the document
  - Code snippets
  - Property suggestions
- **Hover Documentation**: Detailed help for keywords and operators
- **Full LSP Compliance**: Works with any LSP-compatible editor

## Building

```bash
go build -o jsson-lsp ./cmd/jsson-lsp
```

## Running

The Language Server communicates via stdio:

```bash
./jsson-lsp
```

## Integration with VS Code

The JSSON VS Code extension automatically uses this Language Server. To configure manually:

1. Build the Language Server
2. Update the extension to point to the executable
3. Restart VS Code

## Architecture

```
internal/lsp/
├── server.go       # Main LSP server and message handling
├── diagnostics.go  # Error detection and reporting
├── completion.go   # Autocomplete functionality
└── hover.go        # Hover documentation
```

The LSP reuses JSSON's existing:
- `internal/lexer` - Tokenization
- `internal/parser` - Syntax analysis
- `internal/ast` - Abstract syntax tree

## Supported LSP Methods

- `initialize` - Server capabilities
- `textDocument/didOpen` - Document opened
- `textDocument/didChange` - Document changed
- `textDocument/didClose` - Document closed
- `textDocument/completion` - Autocomplete
- `textDocument/hover` - Hover information
- `textDocument/publishDiagnostics` - Error reporting
- `shutdown` - Graceful shutdown
- `exit` - Server exit

## Logging

The server logs to `jsson-lsp.log` in the current directory for debugging.

## Future Enhancements

- Go to definition
- Find references
- Rename symbol
- Document formatting
- Code actions (quick fixes)
- Semantic tokens (better syntax highlighting)
