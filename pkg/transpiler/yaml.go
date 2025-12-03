package transpiler

import (
	"encoding/json"
	"jsson/internal/ast"
	ie "jsson/internal/errors"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// TranspileToYAML converts the transpiled data to YAML format
func (t *Transpiler) TranspileToYAML() ([]byte, error) {
	// First, transpile to the internal representation
	root := make(map[string]interface{})

	for _, stmt := range t.program.Statements {
		switch s := stmt.(type) {
		case *ast.VariableDeclaration:
			// Variable declarations are stored in symbol table but not added to output
			val, err := t.evalExpression(s.Value, nil)
			if err != nil {
				return nil, err
			}
			t.symbolTable[s.Name.Value] = val
		case *ast.AssignmentStatement:
			key := s.Name.Value
			val, err := t.evalExpression(s.Value, nil)
			if err != nil {
				return nil, err
			}
			root[key] = val
		case *ast.IncludeStatement:
			// Handle includes (same logic as JSON transpiler)
			includePath := s.Path.Value
			var includeAbs string
			if filepath.IsAbs(includePath) {
				includeAbs = filepath.Clean(includePath)
			} else {
				includeAbs = filepath.Clean(filepath.Join(t.baseDir, includePath))
			}

			if t.inProgress[includeAbs] {
				return nil, t.errfNodeMsg(s, ie.CyclicInclude(includeAbs))
			}

			if cached, ok := t.includeCache[includeAbs]; ok {
				for k, v := range cached {
					if _, exists := root[k]; !exists {
						root[k] = v
					}
				}
				break
			}

			t.inProgress[includeAbs] = true

			data, err := os.ReadFile(includeAbs)
			if err != nil {
				t.inProgress[includeAbs] = false
				return nil, t.errfNode(s, "could not read include file %q — gremlin can't find it: %v", s.Path.Value, err)
			}

			l := lexer.New(string(data))
			l.SetSourceFile(includeAbs)
			p := parser.New(l)
			prog := p.ParseProgram()
			if len(p.Errors()) > 0 {
				t.inProgress[includeAbs] = false
				return nil, t.errfNode(s, "parser errors in included file %q — wizard got confused: %v", s.Path.Value, p.Errors())
			}

			incBase := filepath.Dir(includeAbs)
			incT := New(prog, incBase, t.mergeMode, includeAbs)
			incT.includeCache = t.includeCache
			incT.inProgress = t.inProgress

			incJSON, err := incT.Transpile()
			if err != nil {
				t.inProgress[includeAbs] = false
				return nil, t.errfNode(s, "transpile error in included file %q: %v", s.Path.Value, err)
			}

			var incRoot map[string]interface{}
			if err := json.Unmarshal(incJSON, &incRoot); err != nil {
				t.inProgress[includeAbs] = false
				return nil, t.errfNode(s, "invalid json from include %q: %v", s.Path.Value, err)
			}

			t.includeCache[includeAbs] = incRoot
			t.inProgress[includeAbs] = false

			for k, v := range incRoot {
				switch t.mergeMode {
				case "keep":
					if _, exists := root[k]; !exists {
						root[k] = v
					}
				case "overwrite":
					root[k] = v
				case "error":
					if _, exists := root[k]; exists {
						return nil, t.errfNode(s, "include merge conflict for key %q from %s", k, includeAbs)
					}
					root[k] = v
				default:
					if _, exists := root[k]; !exists {
						root[k] = v
					}
				}
			}
		}
	}

	// Marshal to YAML
	return yaml.Marshal(root)
}
