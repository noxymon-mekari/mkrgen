![logo](./public/logo.png)

<div style="text-align: center;">
  <h1>
    Production-Ready Microservices Boilerplate Generator
  </h1>
</div>

> **Acknowledgment**: This project is a fork of [go-blueprint](https://github.com/Melkeydev/go-blueprint/) that has been enhanced and adjusted to fulfill a broader use case of generating production-ready microservices architecture and patterns.

mkrgen is a powerful CLI tool that generates production-ready Go microservices with standardized architecture patterns. It provides comprehensive boilerplate code, advanced integrations, and enterprise-grade features to accelerate your microservices development from day one.

### Why Would I use this?

- Easy to set up and install
- Production-ready microservices architecture patterns
- Complete Go project structure with enterprise-grade organization
- Advanced HTTP server setup with popular frameworks
- Integrated database layers with multiple driver support
- Background job processing with Worker patterns
- Kafka consumer implementations for event-driven architecture
- Built-in observability and monitoring patterns
- Docker and containerization support
- CI/CD pipeline templates
- Advanced frontend integrations (React, HTMX)
- Focus on actual business logic instead of boilerplate setup

## Table of Contents

- [Install](#install)
- [Frameworks Supported](#frameworks-supported)
- [Database Support](#database-support)
- [Advanced Features](#advanced-features)
- [Usage Example](#usage-example)
- [License](#license)

<a id="install"></a>

<h2>
  <picture>
    <img src="./public/install.gif?raw=true" width="60px" style="margin-right: 1px;">
  </picture>
  Install
</h2>

```bash
go install github.com/noxymon-mekari/mkrgen@latest
```

This installs a go binary that will automatically bind to your $GOPATH

> if you’re using Zsh, you’ll need to add it manually to `~/.zshrc`.

```bash
GOPATH=$HOME/go  PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

don't forget to update

```bash
source ~/.zshrc
```

Then in a new terminal run:

```bash
mkrgen create
```

You can also use the provided flags to set up a project without interacting with the UI.

```bash
mkrgen create --name my-project --framework gin --driver postgres --git commit
```

See `mkrgen create -h` for all the options and shorthands.

<a id="frameworks-supported"></a>

<h2>
  <picture>
    <img src="./public/frameworks.gif?raw=true" width="60px" style="margin-right: 1px;">
  </picture>
  Frameworks Supported
</h2>

- [Chi](https://github.com/go-chi/chi)
- [Gin](https://github.com/gin-gonic/gin)
- [Fiber](https://github.com/gofiber/fiber)
- [HttpRouter](https://github.com/julienschmidt/httprouter)
- [Gorilla/mux](https://github.com/gorilla/mux)
- [Echo](https://github.com/labstack/echo)

<a id="database-support"></a>

<h2>
  <picture>
    <img src="./public/database.gif?raw=true" width="45px" style="margin-right: 15px;">
  </picture>
  Database Support
</h2>

mkrgen now offers enhanced database support, allowing you to choose your preferred database driver during project setup. Use the `--driver` or `-d` flag to specify the database driver you want to integrate into your project.

### Supported Database Drivers

Choose from a variety of supported database drivers:

- [Mysql](https://github.com/go-sql-driver/mysql)
- [Postgres](https://github.com/jackc/pgx/)
- [Sqlite](https://github.com/mattn/go-sqlite3)
- [Mongo](https://go.mongodb.org/mongo-driver)
- [ScyllaDB GoCQL](https://github.com/scylladb/gocql)

<a id="advanced-features"></a>

<h2>
  <picture>
    <img src="./public/advanced.gif?raw=true" width="70px" style="margin-right: 1px;">
  </picture>
  Advanced Features
</h2>

mkrgen provides enterprise-grade advanced features designed for production microservices. These features can be enabled individually or combined to create comprehensive, scalable applications.

You can use the `--advanced` flag when running the `create` command to access the following features. This is a multi-option prompt; one or more features can be used simultaneously:

- [HTMX](https://htmx.org/) support using [Templ](https://templ.guide/) for dynamic web interfaces
- CI/CD workflow setup using [Github Actions](https://docs.github.com/en/actions) with production-ready pipelines
- [Websocket](https://pkg.go.dev/github.com/coder/websocket) endpoints for real-time communication
- [Tailwind](https://tailwindcss.com/) CSS framework for modern styling
- Docker configuration with multi-stage builds and production optimization
- [React](https://react.dev/) frontend with TypeScript and complete backend integration
- [Worker](https://github.com/hibiken/asynq) background job processing with Redis and asynq
- [Kafka](https://github.com/segmentio/kafka-go) consumer implementation for event-driven architecture
- [Redis](https://github.com/redis/go-redis) integration for caching and session management
- [Swagger](https://swagger.io/) API documentation generation

Note: Selecting Tailwind option will automatically select HTMX unless React is explicitly selected

<a id="usage-example"></a>

<h2>
  <picture>
    <img src="./public/example.gif?raw=true" width="60px" style="margin-right: 1px;">
  </picture>
  Usage Example
</h2>

Here's an example of setting up a project with a specific database driver:

```bash
mkrgen create --name my-project --framework gin --driver postgres --git commit
```

<p align="center">
  <img src="./public/blueprint_1.png" alt="Starter Image" width="800"/>
</p>

Advanced features are accessible with the --advanced flag

```bash
mkrgen create --advanced
```

Advanced features can be enabled using the `--feature` flag along with the `--advanced` flag.

HTMX:

```bash
mkrgen create --advanced --feature htmx
```

CI/CD workflow:

```bash
mkrgen create --advanced --feature githubaction
```

Websocket:

```bash
mkrgen create --advanced --feature websocket
```

Tailwind:

```bash
mkrgen create --advanced --feature tailwind
```

Docker:

```bash
mkrgen create --advanced --feature docker
```

React:

```bash
mkrgen create --advanced --feature react
```

Worker (Background Jobs):

```bash
mkrgen create --advanced --feature worker
```

Kafka Consumer:

```bash
mkrgen create --advanced --feature kafka
```

Redis Integration:

```bash
mkrgen create --advanced --feature redis
```

Swagger Documentation:

```bash
mkrgen create --advanced --feature swagger
```

Or all features at once:

```bash
mkrgen create --name my-project --framework chi --driver mysql --advanced --feature htmx --feature githubaction --feature websocket --feature tailwind --feature docker --feature worker --feature kafka --feature redis --feature swagger --git commit --feature react
```

<p align="center">
  <img src="./public/blueprint_advanced.png" alt="Advanced Options" width="800"/>
</p>

**Visit [documentation](https://docs.mkrgen.dev) to learn more about mkrgen and its features.**

<a id="license"></a>

<h2>
  <picture>
    <img src="./public/license.gif?raw=true" width="50px" style="margin-right: 1px;">
  </picture>
  License
</h2>

Licensed under [MIT License](./LICENSE)
