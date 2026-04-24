# Guia de Desenvolvimento - Goodwe Backend

Este guia contém todos os comandos necessários para configurar e rodar o projeto localmente.

## 1. Pré-requisitos

*   **Go 1.24+**
*   **Docker e Docker Compose**
*   **Postgres (via Docker)**

---

## 2. Configuração Inicial

Primeiro, prepare o arquivo de ambiente:

```bash
cp .env.example .env
```

*Certifique-se de ajustar as variáveis no `.env` se necessário (ex: `DATABASE_URL`).*

---

## 3. Infraestrutura (Banco de Dados e Redis)

Suba os containers do Postgres e Redis:

```bash
docker-compose up -d
```

---

## 4. Banco de Dados e Migrações

O projeto utiliza o **Ent ORM**. As migrações (criação de tabelas) ocorrem automaticamente ao iniciar qualquer API ou o script de seed.

Se você alterar o schema (`shared/infrastructure/persistence/postgres/ent/schema/`), regenere o código do Ent:

```bash
cd shared/infrastructure/persistence/postgres/ent/ && go generate ./...
```

---

## 5. Seeding (Dados Iniciais)

Para popular o banco com usuários de teste (`admin@goodwe.com`, `employee@goodwe.com`), grupos de acesso e empresas:

```bash
cd tools/seed && go run main.go
```
*Senha padrão para todos os usuários: `password123`*

---

## 6. Executando as APIs

Você precisará de terminais separados para cada serviço:

### Auth API (Porta 3001)
Responsável por login e geração de tokens JWT.
```bash
cd apps/auth-api && go run main.go
```

### Employee API (Porta 3003)
Responsável por gestão de usuários e relatórios.
```bash
cd apps/employee-api && go run main.go
```

### Worker (Processamento em Background)
Responsável por processar tarefas assíncronas (como geração de Excel/CSV).
```bash
cd apps/employee-api/worker && go run main.go
```

---

## 7. Comandos de Manutenção

### Limpar Dependências
```bash
go work sync
```

### Resetar Banco de Dados (Cuidado!)
```bash
docker-compose down -v
docker-compose up -d
# Após isso, rode o seed novamente
```
