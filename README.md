# URL Shortener

Um projeto simples de encurtador de URLs feito em Go.

## Funcionalidades

- Criar URLs encurtadas
- Redirecionar para a URL original usando o slug
- Health check da API
- Observabilidade (logs, métricas, tracing)
- Configuração via Docker e variáveis de ambiente

## Como rodar

1. Clone o repositório:
   ```
   git clone <url-do-repo>
   cd url-shortener
   ```

2. Configure o banco de dados no arquivo `.env` (veja `.env.example`).

3. Suba os serviços com Docker Compose:
   ```
   docker-compose up -d
   ```

## Rotas principais

- `POST /url`  
  Cria uma URL encurtada.  
  Exemplo:
  ```
  curl -X POST http://localhost:8080/url \
    -H "Content-Type: application/json" \
    -d '{"original_url": "https://exemplo.com", "slug": "meuslug"}'
  ```

- `GET /{slug}`  
  Redireciona para a URL original.

- `GET /health`  
  Verifica se a API está online.

## Observabilidade

O projeto já possui métricas e tracing integrados.

<img width="1358" height="523" alt="image" src="https://github.com/user-attachments/assets/8601f239-328a-491b-9501-84cc2851bec6" />
