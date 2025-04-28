# Task Service

This project is a Go-based backend service that implements a simple task management API. It allows you to create tasks with a title, description, and automatically generated creation timestamp. Additionally, it integrates with Prometheus to expose metrics (total tasks created and task creation duration) and provides a Docker Compose setup to run the backend service along with PostgreSQL, Prometheus, and Grafana.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation and Setup](#installation-and-setup)
- [Running the Application](#running-the-application)
- [Testing the API](#testing-the-api)
- [Monitoring with Prometheus and Grafana](#monitoring-with-prometheus-and-grafana)
- [Environment Variables](#environment-variables)
- [Stopping the Application](#stopping-the-application)
- [License](#license)

## Features

- **Task API**:
  - `POST /tasks`: Create a new task with a title and description.
  - Automatically stores the task with an auto-generated `id` and `created_at` timestamp.
- **Database Integration**:
  - Uses PostgreSQL for data persistence.
  - Automatic migration using GORM.
- **Metrics**:
  - Exposes Prometheus metrics:
    - `tasks_created_total`: Counter for created tasks.
    - `task_creation_duration_seconds`: Histogram for task creation processing time.
- **Docker Compose Setup**:
  - Runs the backend service, PostgreSQL, Prometheus, and Grafana.
- **Clean Architecture**:
  - Organized into domain, repository, service, and delivery layers for maintainability.

## Project Structure

```
.  
├── cmd  
│   └── main.go  
├── internal  
│   ├── delivery  
│   │   └── rest  
│   │       └── handlers.go  
│   ├── domain  
│   │   └── task.go  
│   ├── repository  
│   │   └── task.go  
│   ├── server  
│   │   └── server.go  
│   └── service  
│       └── task.go  
├── pkg  
│   └── database  
│       └── postgres.go  
├── prometheus  
│   └── prometheus.yml  
├── .env  
├── docker-compose.yml  
├── Dockerfile  
├── go.mod  
├── go.sum  
└── README.md  
```

## Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- (Optional) [Go 1.20+](https://golang.org/dl/) for local development

## Installation and Setup

### 1. Clone the Repository

```bash
git clone https://github.com/valeraBerezovskij/testappDev.git
cd testappDev
```

### 2. Create a `.env` File

Create a `.env` file in the root of the project with the following contents:

```ini
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=tasks_db
```

### 3. Download Dependencies (if developing locally)

Initialize the module and download dependencies:

```bash
go mod tidy
```

## Running the Application

Use Docker Compose to build and start the entire environment:

```bash
docker-compose up --build
```

This command will build the Go backend, start PostgreSQL, Prometheus, and Grafana containers, and expose the following ports:

- **Backend**: [http://localhost:8080](http://localhost:8080)
- **PostgreSQL**: `localhost:5432`
- **Prometheus**: [http://localhost:9090](http://localhost:9090)
- **Grafana**: [http://localhost:3000](http://localhost:3000)

## Testing the API

You can create a new task by sending a `POST` request to the `/tasks` endpoint:

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Write Go test task",
    "description": "Implement a backend service with Prometheus metrics"
  }'
```

A successful request will return a JSON response with the task details, including an auto-generated `id` and `created_at` timestamp.

## Monitoring with Prometheus and Grafana

### Prometheus

#### Access Prometheus:

- Open [http://localhost:9090](http://localhost:9090) in your browser.

#### Check Scrape Targets:

- Go to **Status > Targets** and verify that the backend target is up.

#### Run Queries:

- Use the **Graph** tab to run queries like `tasks_created_total` or `task_creation_duration_seconds` to see your metrics.

### Grafana

#### Access Grafana:

- Open [http://localhost:3000](http://localhost:3000) in your browser.
- Log in with the default credentials (`admin/admin`).

#### Add Prometheus as a Data Source:

1. Click the gear icon in the sidebar and select **Data Sources**.
2. Click **Add data source** and choose **Prometheus**.
3. Set the URL to `http://prometheus:9090` (if using Docker networking) or `http://localhost:9090` (if accessing directly).
4. Click **Save & Test**.

#### Create a Dashboard:

1. Click on the **+** icon in the sidebar and select **Dashboard**.
2. Click **Add new panel** and choose your Prometheus data source.
3. Enter a query (for example, `tasks_created_total`) to create a graph panel.
4. Save your dashboard for ongoing monitoring.

## Environment Variables

The application uses environment variables for database configuration. You can define these in a `.env` file:

```ini
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=tasks_db
```

These values are used both by the Go application and the Docker Compose configuration.

## Stopping the Application

To stop the application and remove containers, run:

```bash
docker-compose down
```