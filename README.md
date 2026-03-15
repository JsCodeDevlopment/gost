<div align="center">

# Gost 🚀

**A powerful, opinionated Go boilerplate engineered with a NestJS-like architecture.**

[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue.svg)](https://golang.org/dl/)
[![Gin](https://img.shields.io/badge/framework-Gin-brightgreen.svg)](https://github.com/gin-gonic/gin)
[![GORM](https://img.shields.io/badge/orm-GORM-red.svg)](https://gorm.io/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

<div align="center">
   <img src="https://img.lightshot.app/I94-6gInQIOFtWawF839Ww.png" />
</div>

## 📖 About The Project

**Gost** is a robust boilerplate and starter kit designed for building scalable backend applications in Golang. By drawing heavy inspiration from the modular architecture and design patterns of features of **NestJS**, Gost brings structure, order, and developer ergonomics to your Go web applications.

It provides a ready-to-use environment completely configured with a powerful HTTP framework, an ORM, caching, decoupled validation, centralized error handling, and file uploads. Everything is orchestrated in a familiar way to developers transitioning from the Node/NestJS ecosystem, bringing the strong typing and high performance of Go.

---

## 🌟 Key Features

- **NestJS-like Architecture**: Logical separation of concerns through Modules, Controllers, Services, and Repositories.
- **Dependency Injection pattern**: Clean and manual wiring of dependencies keeping the codebase tightly coupled only where it needs to be.
- **Built-in Validation**: Class-validator style validation using Go Generics and struct tags (`Pipes`).
- **Global Error Handling**: Centralized exception filtering to avoid leaking panics and standardizing API error JSON responses (`Filters`).
- **Middleware Abstractions**: Simple interfaces for `Interceptors` (request logging/modification) and `Guards` (authentication/authorization).
- **Advanced Security**: Integrated JWT with Access/Refresh tokens, Redis-based blacklisting, and RBAC (Role-Based Access Control).
- **Communication & Real-time**: Fully integrated RabbitMQ for async processing, Websockets (Hub/Client) for real-time interaction, and secure Webhooks with HMAC signatures and auto-retries.
- **Data Protection**: Hardened password hashing with Bcrypt and AES-256-GCM field-level encryption.
- **Resilience**: Redis-powered Rate Limiting to prevent DDoS and brute-force attacks.
- **CORS Configured**: Out-of-the-box support for frontend consumers (SPA-friendly). Dynamically configured via the `ALLOWED_CORS` environment variable logic.
- **File Upload Support**: Built-in utility for handling `multipart/form-data` uploads.
- **Database & Cache Ready**: Pre-configured with PostgreSQL (via GORM) and Redis, easily testable via Docker Compose.

---

## 🛠️ Technologies Used

- **Web Framework:** [Gin Web Framework](https://github.com/gin-gonic/gin)
- **Database ORM:** [GORM](https://gorm.io/)
- **Messaging:** [RabbitMQ (AMQP)](https://www.rabbitmq.com/)
- **Caching & State:** [Go-Redis](https://github.com/go-redis/redis)
- **Security:** JWT, Bcrypt, AES-256-GCM
- **Containerization:** Docker & Docker Compose

---

## 🏗️ Architecture & Core Concepts

Gost recreates the building blocks of modern backend frameworks utilizing Go's native constructs and the Gin framework context.

### 1. Modules (`src/modules`)

Modules group related domain entities, logic, and networking (e.g., Users, Products, Orders) into a cohesive block. Each module exposes an `InitModule` function, acting as the module's wiring board (similar to the `@Module` decorator), setting up and injecting dependencies.

### 2. Controllers, Services & Repositories

- **Controllers** (`*.controller.go`): Handle incoming HTTP requests, extract parameters/body, and format responses. They delegate logic processing.
- **Services** (`*.service.go`): Protect the core business logic. Agnostic of HTTP rules.
- **Repositories** (`*.repository.go`): The persistence layer. Handles direct database interactions (GORM), abstracting the DB logic from the Service.

### 3. Interceptors (`src/common/interceptors`)

Middlewares that wrap route handlers. Use them for request logging, mapping payloads, tracking execution time, or even mutating the response context (e.g., `LoggerInterceptor`).

### 4. Guards (`src/common/guards`)

Middlewares dedicated exclusively to authorization and authentication fluxes. The provided `AuthGuard` checks for valid tokens before allowing the request's execution pipeline to proceed.

### 5. Filters (`src/common/filters`)

Global exception filters. If a controller encounters an error, it shouldn't need to format the error manually. By bubbling it up via `c.Error(err)`, the globally attached `ErrorHandler` intercepts it and formats an elegant JSON response identical to NestJS's `HttpException`.

### 6. Pipes (`src/common/pipes`)

Used for input payload serialization and syntax validation. Gost leverages Go Generics in `pipes.ValidateBody[DTO](c)` to parse JSON bodies directly into typed DTOs and validate them strictly based on Gin's binding tags.

### 7. Connectivity & Real-time (`src/modules/ws`, `src/common/messaging`)

Gost provides out-of-the-box support for:
- **Websockets**: Persistent bidirectional communication using a central Hub.
- **RabbitMQ**: Asynchronous message production and consumption (Scaffolding ready).
- **Webhooks**: Reliable event dispatching with HMAC signatures and exponential backoff retries.

### 🔐 8. Security & Protection (`src/common/security`, `src/modules/auth`)

Security is baked into the framework core:
- **JWT Auth**: Access tokens and Refresh tokens managed via Redis.
- **RBAC**: Protect your routes using `Guards.RolesGuard("admin")`.
- **Rate Limit**: Stop brute force attacks with the built-in Redis rate limiter.
- **Encryption**: Built-in utilities for Bcrypt hashing and AES-256 encryption.

---

## Used Tecnologies

<div align="center">
  <img src="https://img.shields.io/badge/Go-000?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Gin-000?style=for-the-badge&logo=gin&logoColor=white" />
  <img src="https://img.shields.io/badge/GORM-000?style=for-the-badge&logo=gorm&logoColor=white" />
  <img src="https://img.shields.io/badge/PostgreSQL-000?style=for-the-badge&logo=postgresql&logoColor=white" />
  <img src="https://img.shields.io/badge/Redis-000?style=for-the-badge&logo=redis&logoColor=white" />
  <img src="https://img.shields.io/badge/Docker-000?style=for-the-badge&logo=docker&logoColor=white" />
</div>

---

## 🚦 Prerequisites

To run and develop on this project, ensure you have installed:

- [Go](https://go.dev/dl/) >= 1.21
- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- Git

---

## 🚀 Installation & Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yourusername/gost.git my-app
   cd my-app
   ```

2. **Start the Infrastructure (PostgreSQL & Redis):**

   ```bash
   docker-compose up -d
   ```

   _(Wait a few seconds for the database containers to fully initialize)._

3. **Configure Environment Variables:**

   Copy the example environment file and define your database/redis configuration:

   ```bash
   cp .env.example .env
   ```

   The `.env` file manages configurations for the DB, Redis, and application port.

4. **Install Go Dependencies:**

   ```bash
   go mod tidy
   ```

5. **Run the Application:**
   ```bash
   go run main.go
   ```
   The API will boot up and be accessible at `http://localhost:3000/api/v1`.

---

## 📖 Usage Guide

### Directory Structure Overview

```text
gost/
├── main.go                     // Application entry point (Bootstrap)
├── docker-compose.yml          // Infrastructure definitions (Postgres, Redis, RabbitMQ)
├── src/
│   ├── app/
│   │   └── app.module.go       // Main registrar (mounts routes, configs, middlewares)
│   ├── common/
│   │   ├── filters/            // Global Error Handling
│   │   ├── guards/             // Authentication, JWT, and RBAC Middlewares
│   │   ├── interceptors/       // Request flow modifications, Logging, Rate Limiting
│   │   ├── messaging/          // RabbitMQ Producers, Consumers, and Webhook Workers
│   │   ├── pipes/              // Payload Validations logic
│   │   ├── security/           // Cryptographic utils (Bcrypt, AES-256)
│   │   └── utils/              // Utilities (File Upload, Webhook Dispatcher)
│   ├── config/
│   │   ├── database.go         // Database configuration
│   │   ├── rabbitmq.go         // Messaging broker setup
│   │   └── redis.go            // Cache configuration
│   └── modules/
│       ├── auth/               // JWT Identity management (Login, Refresh, Logout)
│       ├── ws/                 // Websocket Hub and Client management
│       └── users/              // [Example] Domain Module
│           ├── dto/            // Payload validation and input schemas
│           ├── entities/       // Database models
│           ├── users.controller.go 
│           ├── users.module.go     
│           ├── users.repository.go 
│           └── users.service.go    
```

### Creating a New Module

Say you want to create a `Products` module:

1. Create a folder `src/modules/products`.
2. Following the NestJS pattern, segregate your files:
   - `entities/`: Define your generic entity models (e.g., `product.entity.go`).
   - `dto/`: Put your request payload structs here.
   - `exceptions/`, `presenters/`, `queries/`, and `tests/`: Scaffold these directories to keep concerns separated as the module grows.
3. At the root of the module folder (`src/modules/products/`):
   - Scaffold the core layers: `products.repository.go`, `products.service.go`, `products.controller.go` (and `products.consumer.go` if parsing messages from queues).
   - Wire them inside `products.module.go` mapping from DB to Repo, Service to Controller.
4. Create `products.module.go` containing `func InitModule(router *gin.RouterGroup)` to manually wire these layers together. Register your `POST`, `GET` handlers here.
5. Finally, register the new module in `src/app/app.module.go`: `products.InitModule(apiGroup)`.

### Handing Incoming Validations (Using Pipes)

Create a DTO with struct tags for automated validation:

```go
type CreateProductDto struct {
    Name  string  `json:"name" binding:"required,min=3"`
    Price float64 `json:"price" binding:"required,gt=0"`
}
```

Use the Pipe in your controller:

```go
func (ctrl *ProductController) Create(c *gin.Context) {
    // Throws a beautifully handled 400 Bad Request if fields are invalid
    dto, err := pipes.ValidateBody[CreateProductDto](c)
    if err != nil {
        return
    }

    // dto is strongly typed as *CreateProductDto
    product, err := ctrl.service.Create(*dto)
    // ...
}
```

---

### 📡 Included API Overview

#### Auth Module
- `POST /api/v1/auth/login` - Authenticate and receive Access & Refresh tokens
- `POST /api/v1/auth/logout` - Invalidate current session (Redis Blacklist)

#### Users Module
- `GET /api/v1/users` - List all users (Example of Repository pattern)
- `GET /api/v1/users/:id` - Fetch user details
- `POST /api/v1/users` - Create a new user (Validates Email & Name size)
- `PUT /api/v1/users/:id` - Update user details
- `DELETE /api/v1/users/:id` - Delete a user
- `POST /api/v1/users/:id/avatar` - Upload a user avatar image (multipart/form-data)

---

## 📚 Masterclass Documentation

To explore the full potential of the library, we've created a directory with explanatory guides teaching step-by-step the inner workings behind Gost. If you want to learn how to extract 100% from every file, keep reading:

1. [01 - Introduction and Architectural Philosophy](./docs/01-introduction-and-architecture.md)
2. [02 - Bootstrap and Server Configurations](./docs/02-bootstrap-and-server-configurations.md)
3. [03 - Building Modules, Injection and Domain-Driven](./docs/03-building-modules-injection-and-domain-driven.md)
4. [04 - DTOs, Security with Pipes and ORM (Entities)](./docs/04-dtos-security-with-pipes-and-orm.md)
5. [05 - Shields: Filters, Auth Guards and Interceptors](./docs/05-shields-filters-auth-guards-and-interceptors.md)
6. [06 - Utilities: Micro-Caching and File Uploads](./docs/06-utilities-micro-caching-and-file-uploads.md)
7. [07 - Deployment and Hosting Strategy](./docs/07-deployment-and-hosting-strategy.md)
8. [08 - Testing Strategies (Unit & E2E)](./docs/08-testing-strategies.md)
9. [09 - Security Deep Dive: Authenticity and Protection](./docs/09-security-deep-dive-authenticity-and-protection.md)
10. [10 - Communication and Connectivity (RabbitMQ, WS, Webhooks)](./docs/10-communication-and-connectivity.md)

---

## 🤝 Contributing

Contributions make the open-source community an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**. If you have a fix or suggestion, please open a pull request.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'feat: Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Contribution Guidelines

- Formatting: Follow standard Go formatting (`go fmt` / `gofumpt`).
- Architecture: Respect the established decoupled design and layers boundaries. Controllers shouldn't call logic directly.
- Commits: Try using [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

---

## Desenvolvedor

| Foto                                                                                                                             | Nome                                                 | Cargo                                                                   |
| -------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------- | ----------------------------------------------------------------------- |
| <img src="https://avatars.githubusercontent.com/u/100796752?s=400&u=ae99bd456c6b274cd934d85a374a44340140e222&v=4" width="100" /> | [Jonatas Silva](https://github.com/JsCodeDevlopment) | Senior Software Engineer / CTO at [PokerNetic](https://pokernetic.com/) |

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---

<div align="center">
  <i>Built with ❤️ by and for <a href="https://www.linkedin.com/in/jscodedevelopment/">Jonatas Silva</a></i>
</div>
