# JSSON Examples - Planned Features (Roadmap)

**ATENÇÃO**: Os exemplos neste diretório usam **sintaxe planejada** que ainda **NÃO está implementada**.

Estes arquivos servem como **especificação** e **testes futuros** para features do roadmap.

## Arquivos

### Database Configuration

- **`database_config.jsson`** - Configuração de banco de dados multi-ambiente
  - Usa: List comprehension `[| |] ->`
  - Usa: `@use` directive
  - Usa: String filters `{var | uppercase}`

### API Configuration

- **`api_config.jsson`** - Gerador de endpoints CRUD
  - Usa: List comprehension `[| |] ->`
  - Usa: Flatten operator `| flatten`
  - Usa: Template strings avançadas

## Roadmap

Estas features estão planejadas para as próximas versões:

### v0.0.7 - `@use` Directive

```jsson
@preset "defaults" { timeout = 30 }
api = @use "defaults" { timeout = 60 }
```

### v0.0.8 - List Comprehension

```jsson
servers = [| (env) "dev", "staging", "prod" |] -> (e) {
    name = "{e}-server"
}
```

### v0.0.9 - String Filters

```jsson
password = "${DB_PASSWORD_{env | uppercase}}"
slug = "{name | kebab-case}"
```

### v0.1.0 - Flatten Operator

```jsson
data = [...nested arrays...] | flatten
```

## Documentação Completa

Veja o roadmap completo em: `../../ROADMAP.md`

## ⚠️ Não Tente Executar

Estes arquivos **NÃO vão compilar** com a versão atual do JSSON:

```bash
# Isto vai falhar
.\jsson.exe -i .\examples\planned\database_config.jsson
# Error: unknown expression type: <nil>
```

## 🤝 Contribuindo

Quer implementar alguma dessas features? Veja:

- `ROADMAP.md` - Detalhes de implementação
- Issues com tag `roadmap` e `enhancement`
