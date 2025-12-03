# JSSON Examples - Current Features

Este diretório contém exemplos que **funcionam completamente** com a versão atual do JSSON.

## Estrutura

### Validation Examples

- **`simple_config.jsson`** - Configuração básica com strings, números, booleans e arrays
- **`invalid_config.jsson`** - Exemplo de configuração inválida (para testes de validação futura)

## ✅ Features Demonstradas

Todos os exemplos neste diretório usam apenas features **implementadas e funcionando**:

- ✅ Variáveis com `:=`
- ✅ Objetos e arrays
- ✅ Strings, números, booleans
- ✅ Comentários
- ✅ Sintaxe limpa (sem quotes nas keys, sem vírgulas)

## Como Testar

```bash
# Build do JSSON
go build -o jsson.exe ./cmd/jsson

# Testar exemplo
.\jsson.exe -i .\examples\current\simple_config.jsson

# Output em diferentes formatos
.\jsson.exe -i .\examples\current\simple_config.jsson -f yaml
.\jsson.exe -i .\examples\current\simple_config.jsson -f toml
.\jsson.exe -i .\examples\current\simple_config.jsson -f ts
```

## Mais Exemplos

Para exemplos mais avançados (ranges, maps, templates), veja o diretório raiz `examples/`.

Para features planejadas (roadmap), veja `examples/planned/`.
