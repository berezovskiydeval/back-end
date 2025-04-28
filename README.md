
# Back-End Applications

A monorepo with several independent Go services that showcase different architectural styles and integrations.

---

## Repository Structure

| Folder | Purpose | Tech Stack |
|--------|---------|-----------|
| **`crud-banana`**     | REST API with JWT auth for the *Banana* entity + log producer to RabbitMQ | Go 1.23, PostgreSQL, RabbitMQ, Viper, Logrus |
| **`bn-logger-mongo`** | RabbitMQ consumer that writes audit logs to MongoDB | Go 1.23, MongoDB Driver, Envconfig |
| **`rpc-application`** | gRPC CRUD example (*Product*) with client and server | Go 1.23, gRPC, MongoDB |
| **`crud-task`**    | Task‑management API with Prometheus metrics + Docker Compose stack | Go 1.23, PostgreSQL, Prometheus, Grafana, GORM |

<details>
<summary>Quick interaction diagram</summary>

```text
crud-banana ──▶ RabbitMQ ──▶ bn-logger-mongo ──▶ MongoDB
          │
          └──▶ PostgreSQL

crud-task ──▶ PostgreSQL
          └──▶ Prometheus ──▶ Grafana

rpc-application (pure gRPC) ──▶ MongoDB
```
</details>

---

## Quick Start (local)

```bash
# 1. Clone the repo
git clone https://github.com/berezovskiydeval/back-end.git
cd back-end
```

### Option A — run each service manually

Spin up infrastructure (Postgres, Mongo, RabbitMQ) any way you like, then:

```bash
# REST API
cd crud-banana && go run ./cmd

# Logger
cd ../bn-logger-mongo && go run ./cmd

# gRPC server
cd ../rpc-application/rpc-server && go run ./cmd
```

### Option B — Task Service stack (Docker Compose)

```bash
cd task-service
cp .env.example .env        # tweak if needed
docker compose up --build   # backend 8080, PG 5432, Prom 9090, Grafana 3000
```

---

## Service Details

### `crud-banana` — REST Service

#### Main Endpoints

| Method | URI                | Description                                |
|--------|--------------------|--------------------------------------------|
| `POST` | `/api/auth/signup` | user registration                          |
| `POST` | `/api/auth/signin` | login, returns **access + refresh** tokens |
| `POST` | `/api/auth/refresh`| refresh tokens                             |
| `GET`  | `/api/items/`      | list bananas                               |
| `GET`  | `/api/items/{id}`  | get one banana                             |
| `POST` | `/api/items/`      | create                                     |
| `PUT`  | `/api/items/{id}`  | update                                     |
| `DELETE`| `/api/items/{id}` | delete                                     |

> All `/api/items/**` routes require a **Bearer** token (middleware).

#### Environment Variables

```env
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=bananas
DB_SSLMODE=disable

RABBIT_URL=amqp://guest:guest@localhost:5672/
RABBIT_QUEUE=logs
```

HTTP port is configured in `configs/config.yml` (default **8080**).

---

### `bn-logger-mongo` — Audit Log Service

* Consumes messages from RabbitMQ, deserialises `LogItem`, and inserts a document into the **audit** collection.  
* Configured via environment variables with `DB_` and `SERVER_` prefixes.  
* Dependencies are declared in its `go.mod`.

---

### `rpc-application` — gRPC Example

* **Server:** `rpc-application/rpc-server`  
* **Client:** `rpc-application/rpc-client`  

Protocol defined in `proto/product.proto`  
(generated files `product.pb.go`, `product_grpc.pb.go`).

The server stores data in **MongoDB**; connection parameters come from `.env`.

---

### `task-service` — Task Management + Metrics

* **API** — `POST /tasks` to create a task (`id`, `title`, `description`, `created_at`).  
* **Database** — PostgreSQL with auto-migration (GORM).  
* **Metrics** — Prometheus counters & histograms (`tasks_created_total`, `task_creation_duration_seconds`).  
* **Docker Compose** — backend, PostgreSQL, Prometheus and Grafana.  
* **Clean Architecture** — domain / repository / service / delivery layers.

See [`task-service/README.md`](task-service/README.md) for full instructions.

---
## License

**MIT** — free to use for study, demos or pet projects.
