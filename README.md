[![Go Reference](https://pkg.go.dev/badge/github.com/jeanmolossi/MaiGo.svg)](https://pkg.go.dev/github.com/jeanmolossi/MaiGo)
[![CI](https://github.com/jeanmolossi/MaiGo/actions/workflows/ci.yml/badge.svg)](https://github.com/jeanmolossi/MaiGo/actions/workflows/ci.yml)

# MaiGo

Se você procura um cliente HTTP simples, rápido e confiável para Go, o MaiGo é a escolha ideal. Inspirado na Mai Sakurajima, o pacote segue a filosofia de aparecer discretamente e resolver seu problema sem alarde.

![Mai Sakurajima](./docs/assets/mai-sakurajima-432-x-498.gif)

## Description
MaiGo oferece uma API baseada em *builders*, permitindo configurações de cabeçalhos, cookies, tentativas de *retry* e balanceamento de carga de forma fluida. Seu foco é produtividade sem deixar de lado a elegância – e alguns easter eggs para os fãs da Mai.

## Requirements
- Go 1.24 ou superior
- Módulos listados em `go.mod`

## Usage
Um exemplo de requisição básica pode ser encontrado em `examples/base_get_request`:

```go
client := maigo.DefaultClient("https://api.example.com")
resp, err := client.GET("/users").Send()
if err != nil {
    // Tratamento de erro
}

var users []User
if err := resp.Body().AsJSON(&users); err != nil {
    // Erro ao ler a resposta
}
```

Outros exemplos estão disponíveis na pasta `examples`, incluindo chamadas com cabeçalhos customizados, balanceamento de carga e *tracing* com OpenTelemetry. E não se surpreenda se surgir uma nova referência à Mai no meio dos logs.

## Contributing
Contribuições são muito bem-vindas! Para colaborar:
1. Fork este repositório.
2. Crie um branch com sua funcionalidade.
3. Envie um *pull request* detalhando as alterações.

Se tiver dúvidas, abra uma *issue*. Enquanto isso, curta a presença da Mai neste projeto!
