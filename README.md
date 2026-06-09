# pm-auth-service

Microservicio de autenticación en **Go**, parte del backend de **Pure Minerals**, un e-commerce construido sobre una arquitectura de microservicios. Expone una API **gRPC** que el API Gateway consume para registro e inicio de sesión, persiste credenciales en **PostgreSQL** y emite tokens **JWT** para la comunicación entre servicios.

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat-square&logo=go&logoColor=white)
![gRPC](https://img.shields.io/badge/gRPC-enabled-4CAF50?style=flat-square)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=flat-square&logo=postgresql&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-HS256-F59E0B?style=flat-square)
![Docker](https://img.shields.io/badge/Docker-compose-2496ED?style=flat-square&logo=docker&logoColor=white)

---

## ¿Qué hace?

| Funcionalidad | Detalle |
|---|---|
| **Registro** | Valida email único, hashea password con bcrypt, persiste en PostgreSQL y retorna JWT |
| **Login** | Busca credencial por email, compara hash, retorna JWT |
| **Token** | JWT firmado con HS256, claims: `email`, `userID`, `exp`, `iat` |

## Arquitectura

Organizado en capas siguiendo principios de **Clean / Hexagonal Architecture**:

```
cmd/           → Punto de entrada
config/        → Variables de entorno
internal/
  domain/      → Entidades y errores de dominio
  application/ → Puertos y casos de uso
  infrastructure/ → PostgreSQL (GORM), bcrypt, JWT
  interfaces/grpc/ → Adaptador gRPC
```

## Inicio rápido

### Opción A — Todo con Docker Compose (recomendado)

Levanta PostgreSQL **y** el servicio de autenticación, ya containerizados:

```bash
docker compose -f Docker-compose.yml up --build
```

> **Podman:** el proyecto funciona igual con Podman (que expone `docker` como alias).
> Si usas el plugin nativo: `podman compose -f Docker-compose.yml up --build`.
> El flag `-f` es necesario porque el archivo se llama `Docker-compose.yml` (con D mayúscula).

Esto:
- Construye la imagen de la app a partir del `Dockerfile` (build multi-stage, binario estático sobre `alpine`).
- Arranca `postgres` con un `healthcheck`; la app espera (`depends_on: service_healthy`) a que la base esté lista antes de conectarse.
- Expone el servidor gRPC en `localhost:50051`.

Para detener y limpiar:

```bash
docker compose -f Docker-compose.yml down        # agrega -v para borrar también el volumen de datos
```

### Opción B — Solo la base en Docker, la app en local

Útil para desarrollo con hot-reload manual:

```bash
# 1. Levantar PostgreSQL
docker compose -f Docker-compose.yml up -d postgres

# 2. Instalar dependencias y correr el servicio
go mod download && go run ./cmd/main.go
```

El servidor queda disponible en `localhost:50051`.

### Probar con grpcurl

```bash
# Registro
grpcurl -plaintext \
  -d '{"email":"user@example.com","password":"pass123"}' \
  localhost:50051 auth.AuthService/Register

# Login
grpcurl -plaintext \
  -d '{"email":"user@example.com","password":"pass123"}' \
  localhost:50051 auth.AuthService/Login
```

## Variables de entorno

| Variable | Default | Descripción |
|---|---|---|
| `GRPC_PORT` | `50051` | Puerto del servidor gRPC |
| `DB_HOST` | `localhost` | Host de PostgreSQL |
| `DB_USER` | `postgres` | Usuario de base de datos |
| `DB_PASSWORD` | `postgres` | Password de base de datos |
| `DB_PORT` | `5432` | Puerto de PostgreSQL |
| `DB_NAME` | `pm_auth` | Nombre de la base de datos |
| `JWT_SECRET` | `secret` | Clave para firmar tokens |
| `JWT_EXPIRATION` | `24h` | Duración del access token |
| `JWT_REFRESH_EXPIRATION` | `168h` | Duración del refresh token |

> En Docker Compose estas variables ya vienen configuradas en el servicio `app` (con `DB_HOST: postgres`). Para correr en local se leen desde el archivo `.env`.

## Stack

- **Go 1.26** · **gRPC** · **PostgreSQL 16**
- **GORM** — ORM con AutoMigrate al inicio
- **bcrypt** — hashing de passwords
- **golang-jwt/jwt v5** — emisión y verificación de tokens
- **Docker Compose** — orquestación local

## Estado del proyecto

> Servicio en desarrollo activo como parte de un sistema de microservicios.

- [x] Flujo de registro y login completo
- [x] Persistencia con GORM y AutoMigrate
- [x] API gRPC con reflection habilitado
- [x] Dockerfile (build multi-stage) y Docker Compose con servicio `app`
- [x] Refresh token
- [x] Tests unitarios
- [ ] Códigos de error gRPC más específicos

---

*Parte de una arquitectura de microservicios en Go.*