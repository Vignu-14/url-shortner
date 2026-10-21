# High-Throughput URL Shortener in Go

A high-performance, concurrent URL shortener written in Go using PostgreSQL for persistent storage and Redis for in-memory caching. Built with standard library routing, a thread-safe Cache-Aside decorator, custom middleware (rate limiting & latency logging), and benchmarked across different concurrency tiers.

---

## Architecture Overview

The system uses a **Decorator Pattern** to decouple the caching logic from the core business service. The HTTP service only talks to a clean `Store` interface, while `CachedStore` transparently orchestrates Redis lookups, cache misses, and PostgreSQL fallbacks.

```
                          HTTP Request
                               │
                               ▼
                   ┌───────────────────────┐
                   │   Logger Middleware   │
                   └───────────┬───────────┘
                               │
                               ▼
                   ┌───────────────────────┐
                   │ RateLimiter Middleware│
                   └───────────┬───────────┘
                               │
                               ▼
                   ┌───────────────────────┐
                   │   Handler & Service   │
                   │ (SHA-256 + Base62)    │
                   └───────────┬───────────┘
                               │
                               ▼
                   ┌───────────────────────┐
                   │ CachedStore Decorator │
                   └───────────┬───────────┘
                              / \
                 Cache Hit   /   \   Cache Miss / Write
                            /     \
                           ▼       ▼
                    ┌─────────┐ ┌────────────┐
                    │  Redis  │ │ PostgreSQL │
                    │ (Cache) │ │ (Database) │
                    └─────────┘ └────────────┘
```

---

## Key Features

- **Standard Library Routing:** Built using Go 1.22+ `http.ServeMux` routing (`POST /`, `GET /{shortUrl}`) with zero third-party web frameworks.
- **Cache-Aside Pattern:** Reads check Redis first. On a cache miss, data is fetched from PostgreSQL, written back into Redis with a TTL, and returned.
- **Custom Middleware Pipeline:**
  - **Latency Logger:** Measures and logs method, path, and request duration in real-time.
  - **Rate Limiter:** Thread-safe Fixed Window Counter protecting endpoints against abuse (e.g., 30 POSTs/min, 500 GETs/min).
- **Persistent Relational Storage:** PostgreSQL schema with a B-tree primary key index on `short_url` and `ON CONFLICT` upsert handling.
- **Connection Pooling:** Tuned pools using `jackc/pgx/v5/pgxpool` and `redis/go-redis/v9`.

---

## Performance & Benchmark Analysis

Tested using `wrk` on Linux x86_64 (8 CPU cores: Intel Core i5-11320H @ 3.20GHz, 7.0 GiB RAM):

| Storage Engine | Concurrency Tier | Operation | Throughput | Median Latency (p50) |
| :--- | :--- | :--- | :--- | :--- |
| **In-Memory Go Map** (`sync.RWMutex`) | 4 threads, 50 conns | Read (`GET`) | **~255,590 req/s** | 113 µs |
| **Redis Cache** (`go-redis/v9`) | 4 threads, 50 conns | Read (`GET`) | **~44,300 req/s** | 1.51 ms |
| **PostgreSQL** (`pgxpool`) | 4 threads, 50 conns | Read (`GET`) | **~23,800 req/s** | 1.81 ms |
| **PostgreSQL** (`pgxpool`) | 4 threads, 50 conns | Write (`POST` unique) | **~9,200 req/s** | 4.91 ms |

### Systems Insights:
1. **In-Process RAM vs. Out-of-Process Redis:** The pure in-memory Go map reached ~255k req/s because memory reads are direct CPU cache pointer lookups (~100 ns). Redis introduces inter-process loopback TCP socket syscalls and RESP protocol serialization.
2. **Redis CPU Profiling (`pidstat`):** Profiling Redis during load showed it consuming ~76% of its single CPU core, with **57% of CPU time spent in the Linux kernel (`%system`)** handling TCP socket buffers and context switches.
3. **Cache Efficiency:** Adding Redis to PostgreSQL nearly **doubled read throughput** (~23.8k $\rightarrow$ ~44.3k req/s) while offloading query load from the relational database.

---

## Getting Started

### 1. Prerequisites
- [Go](https://go.dev/) (version 1.22 or higher)
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)

### 2. Database Setup
Connect to PostgreSQL and create the database and table:

```sql
CREATE DATABASE urlshortner;

\c urlshortner

CREATE TABLE short_urls (
    short_url VARCHAR(20) PRIMARY KEY,
    long_url TEXT NOT NULL
);
```

### 3. Start Redis
```bash
sudo systemctl start redis-server
```

### 4. Run the API Server
```bash
go run cmd/api/main.go
```
The server will start listening on `http://localhost:3000`.

---

## API Endpoints

### 1. Shorten a URL
- **Endpoint:** `POST /`
- **Headers:** `Content-Type: application/json`
- **Request Body:**
  ```json
  {
    "url": "https://github.com/vignesh"
  }
  ```
- **Example cURL:**
  ```bash
  curl -X POST http://localhost:3000/ \
    -H "Content-Type: application/json" \
    -d '{"url":"https://github.com/vignesh"}'
  ```
- **Response (`201 Created`):**
  ```json
  {
    "status": "success",
    "shortUrl": "http://localhost:3000/hgZp4UlDoOn"
  }
  ```

---

### 2. Redirect to Original URL
- **Endpoint:** `GET /{shortUrl}`
- **Example cURL:**
  ```bash
  curl -i http://localhost:3000/hgZp4UlDoOn
  ```
- **Response (`308 Permanent Redirect`):**
  ```text
  HTTP/1.1 308 Permanent Redirect
  Location: https://github.com/vignesh
  ```

---

## Project Structure

```text
├── cmd/
│   └── api/
│       └── main.go               # Server initialization, DB/Redis connections, route wiring
├── internal/
│   ├── handlers/
│   │   └── generate_url.go       # HTTP handlers for POST and GET
│   ├── middleware/
│   │   ├── logger.go             # Request latency and access logger
│   │   └── rate_limiter.go       # Thread-safe Fixed Window Counter rate limiter
│   ├── models/
│   │   └── url.go                # JSON request/response models
│   ├── services/
│   │   ├── short_url.go          # SHA-256 + Base62 hashing algorithm
│   │   └── url_service.go        # Business logic resolving and creating URLs
│   └── storage/
│       ├── db-storage.go         # PostgreSQL implementation via pgxpool
│       ├── memory.go             # Thread-safe in-memory store (RWMutex)
│       ├── redis.go              # Redis client store with TTL
│       └── store.go              # Store interface & CachedStore decorator
├── benchmarks/                   # Initial memory benchmarks and lua scripts
└── benchmarks-postgres/          # PostgreSQL performance results and scripts
```
