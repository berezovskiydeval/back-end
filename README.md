# Back-End Applications

---

## Repository Structure

| Folder | Purpose | Tech Stack |
|--------|---------|-----------|
| **`crud-banana`** | REST API with JWT authentication for the *Banana* entity + log producer to RabbitMQ | Go 1.23, PostgreSQL, RabbitMQ, Viper, Logrus |
| **`bn-logger-mongo`** | RabbitMQ consumer that writes audit logs to MongoDB | Go 1.23, MongoDB Driver, Envconfig |
| **`rpc-application`** | gRPC CRUD example (*Product*) with client and server | Go 1.23, gRPC, MongoDB |

<details>
<summary>Quick interaction diagram</summary>

The REST service **`crud-banana`** → publishes *LogItem* to ✉︎ RabbitMQ →  
consumer **`bn-logger-mongo`** → writes to 🗄 MongoDB  
<br/>The **`rpc-application`** runs separately (pure gRPC, no queue).
</details>

---

## Quick Start (local)

```bash
# 1. Clone the repo
git clone https://github.com/berezovskiydeval/back-end.git
cd back-end

# 2. Start the infrastructure (Postgres + Mongo + RabbitMQ)
docker run -d --name pg -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:16
docker run -d --name mongo -p 27017:27017 mongo:7
docker run -d --name rabbit -p 5672:5672 -p 15672:15672 rabbitmq:3-management

# 3. Copy example env files and tweak if needed
cp crud-banana/.env.example crud-banana/.env
cp bn-logger-mongo/.env.example bn-logger-mongo/.env
cp rpc-application/rpc-server/.env.example rpc-application/rpc-server/.env

# 4. Build & run the services

# REST API
cd crud-banana
go run ./cmd

# Logger
cd ../bn-logger-mongo
go run ./cmd
```
### `crud-banana` — REST Service

#### Main Endpoints

| Method | URI                | Description                               |
|--------|--------------------|-------------------------------------------|
| `POST` | `/api/auth/signup` | user registration                         |
| `POST` | `/api/auth/signin` | login, returns **access + refresh** tokens |
| `POST` | `/api/auth/refresh`| refresh tokens                            |
| `GET`  | `/api/items/`      | list bananas                              |
| `GET`  | `/api/items/{id}`  | get one banana                            |
| `POST` | `/api/items/`      | create                                    |
| `PUT`  | `/api/items/{id}`  | update                                    |
| `DELETE`| `/api/items/{id}` | delete                                    |

> All `/api/items/**` routes require a **Bearer** token (checked by middleware).

#### Environment Variables

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
HTTP port is set in `configs/config.yml` (default **8080**).

---

### `bn-logger-mongo` — Audit Log Service

* Subscribes to the RabbitMQ queue, deserialises `LogItem`, and inserts a document into the **audit** collection.  
* Configured via environment variables with `DB_` and `SERVER_` prefixes.  
* All dependencies are declared in its `go.mod`.

---

### `rpc-application` — gRPC Example

* **Server:** `rpc-application/rpc-server`  
* **Client:** `rpc-application/rpc-client`

The protocol is defined in `proto/product.proto`  
(resulting generated files: `product.pb.go`, `product_grpc.pb.go`).

The server persists data in **MongoDB**; connection parameters are read from `.env`.

cd ../rpc-application/rpc-server
go run ./cmd
