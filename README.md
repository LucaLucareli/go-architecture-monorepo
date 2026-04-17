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
