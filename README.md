# Centralized Logging System

This project is a centralized logging solution built with Go, MongoDB, and Docker. It provides a scalable way to collect, store, and query logs from multiple services in a distributed environment.

## Features

- Centralized log collection from multiple clients
- RESTful API for log ingestion and querying
- MongoDB backend for efficient storage and retrieval
- Dockerized deployment for easy setup

---

## Project Structure

```plaintext
centeralised-logging/
├── centeral-logging/   # Main logging service (API, DB integration)
│   ├── main.go
│   ├── handler/
│   ├── database/
│   ├── model/
│   └── Dockerfile
├── client/             # Example client for sending logs
│   ├── client.go
│   └── Dockerfile
├── server/             # Example server for parsing/validating logs (Log-Controller)
│   ├── server.go
│   └── Dockerfile
├── docker-compose.yml  # Multi-service orchestration
├── README.md           # Project documentation
└── mongo-data;C/       # MongoDB data volume
```

---

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)

### Quick Start

1. Clone the repository:

 ```powershell
 git clone https://github.com/gunjanmistry08/centeralised-logging.git
 cd centeralised-logging
 ```

1. Start all services using Docker Compose:

 ```powershell
 docker-compose up --build
 ```

1. The logging API will be available at `http://localhost:8080` (default).

---

## Usage

### Sending Logs

1. The client would send logs automaticly to the server with a random time interval of 1/2 seconds

### Querying Logs

Send HTTP GET requests to `/logs` endpoint with query parameters:

```bash
curl "http://localhost:8080/logs?limit=5"
```

### Metric

To get metric data of logs

```bash
curl "http://localhost:8080/metrics"
```

---

## Configuration

Environment variables and configuration options can be set in the Docker Compose file or passed to containers as needed.

---

## Imporvements

1. Using another collection to keep track of user who are blacklisted and feed it to redis for faster lookup
2. Using kafka/rabbitmq to  pass logs to centeral logging server from log-controller
