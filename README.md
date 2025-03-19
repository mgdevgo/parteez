<h1 align="center"> 
    Parteez
</h1>
<p align="center">A modern event management platform built with Go</p>

## Getting Started

```shell
make run
```

## Development

Install project dependencies

```shell
go install tool
```

## Requirements

| Component | Version |
|-----------|---------|
| Go | 1.24.0 |
| PostgreSQL | Latest |
| Make | GNU Make 4.4.1 |
| Docker | Docker Engine 24.00 |
| Docker Compose | 3.9 |

## Tech Stack

- **Go** - Core programming language
- **Fiber** - Fast HTTP web framework
- **PGX** - PostgreSQL driver and toolkit
- **Koanf** - Configuration management
- **golang-migrate** - Database migration tool
- **slog** - Structured logging
- **testify** - Testing framework
- **swaggo** - Swagger/OpenAPI documentation generator

## Troubleshooting

### Remove warnings in proto files

Open `settings.json` in vscode and paste this options

```json
"protoc": {
    "options": [
        "--proto_path=${userHome}/.easyp/"
    ]
}
```
