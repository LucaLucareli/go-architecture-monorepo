# api-golang

A api está seguindo princípios de **Clean Architecture** e **DDD (Domain Driven Design)**.

## Estrutura do Projeto

O projeto é organizado como um monorepo utilizando `go.work`, dividido em `apps` (microserviços) e um diretório `shared` que contém a lógica de domínio e infraestrutura comum.

### Diretório `shared/`

O core da aplicação reside aqui, dividido em camadas para garantir separação de preocupações:

1.  **`domain/`**: Camada mais interna. Contém as regras de negócio puras.
    *   `entities/`: Estruturas de dados fundamentais (ex: `User`, `Business`).
    *   `repositories/`: Interfaces de persistência que a aplicação espera.
    *   `enums/`: Constantes de negócio.
    *   `types/`: Tipos de dados globais de domínio.

2.  **`application/`**: Casos de uso e serviços que orquestram o domínio.
    *   `auth/`: `AuthService` que implementa lógica de login e tokens.
    *   `interfaces/`: Interfaces de comunicação da aplicação (ex: `HttpResponse`).

3.  **`infrastructure/`**: Implementações concretas de tecnologias externas.
    *   `persistence/`: Implementação de repositórios usando **Ent ORM** (Postgres).
    *   `queue/`: Implementação de filas usando **Asynq** (Redis).
    *   `security/`: Implementação de segurança (JWT).
    *   `report/`: Geradores de relatórios (Excel, CSV).
    *   `container/`: Injeção de dependências e configuração de módulos.

4.  **`pkg/`**: Utilitários e helpers transversais (cross-cutting concerns).
    *   `logger/`, `validation/`, `middlewares/`, `helpers/`, `interceptors/`.

### Diretório `apps/`

Contém os serviços executáveis:

*   **`auth-api/`**: Responsável pela autenticação e gestão de tokens.
*   **`employee-api/`**: Responsável pela gestão de funcionários e geração de relatórios.
    *   `worker/`: Processador de tarefas assíncronas (background jobs).

## Destaques Técnicos & Performance

Este projeto implementa padrões avançados para garantir alta disponibilidade e baixo consumo de recursos:

### 1. Lazy Loading DI Container
Implementado o carregamento tardio de dependências utilizando o padrão **Singleton Seguro** com `sync.Once`.
- **Vantagem:** Conexões pesadas (Postgres, Redis) só são estabelecidas no momento do primeiro uso real.
- **Resultado:** Startup quase instantâneo dos serviços e economia de memória em ambientes escaláveis.

### 2. Auditoria Assíncrona (Non-blocking)
Middleware de auditoria que utiliza **Goroutines** e **Channels com buffer**.
- **Mecânica:** Os logs de acesso são enviados para um buffer em memória. Um worker em background consome esses dados e os processa sem adicionar milissegundos à latência da requisição do usuário.
- **Resiliência:** Inclui mecanismos de queda controlada (drop entry) caso o buffer atinja o limite, protegendo a estabilidade da aplicação principal.

### 3. Resiliência e Proteção (Middlewares)
- **Timeout (Context Deadline):** Controle rígido do ciclo de vida da request. Se uma query ou processamento travar, o contexto é cancelado automaticamente para evitar processos zumbis.
- **Rate Limiting (Token Bucket):** Proteção contra ataques de força bruta e DoS no nível de IP, garantindo que o throughput da aplicação seja respeitado.

## Tecnologias Utilizadas

*   **Go (Golang)** 1.24+
*   **Ent ORM**: Entity framework para Go.
*   **Echo**: Web framework minimalista.
*   **Asynq**: Redis-based task queue para processamento assíncrono.
*   **Redis**: Cache e storage para as filas.
*   **PostgreSQL**: Banco de dados relacional.
*   **Docker & Docker Compose**: Orquestração de infraestrutura local.

## Padrões de Projeto Aplicados

*   **Repository Pattern**: Desacopla a lógica de negócio do banco de dados.
*   **Dependency Injection**: Facilita testes e modularidade.
*   **DTO (Data Transfer Objects)**: Define contratos claros entre camadas e APIs.
*   **Strategy Pattern**: Utilizado na fábrica de geradores de relatórios.
*   **Singleton with sync.Once**: Utilizado para injeção de dependência preguiçosa (Lazy Loading).
