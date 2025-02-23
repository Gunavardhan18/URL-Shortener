# URL Shortener

## Overview
This project is a URL shortener service built with **Go** and **Fiber**, using **PostgreSQL** for persistent storage and **Redis** for caching. It includes a **health check API** for monitoring service availability.

## Features
- Shorten long URLs into short codes
- Redirect short URLs to original long URLs
- Health check endpoint for service monitoring
- Database and cache integration
- Dockerized for easy deployment

## Tech Stack
- **Backend:** Go (Fiber framework)
- **Database:** PostgreSQL
- **Cache:** Redis
- **Containerization:** Docker, Kubernetes (optional)

## Setup Instructions

### Prerequisites
Ensure you have the following installed:
- Go 1.22+
- Docker & Docker Compose
- PostgreSQL
- Redis

### Clone the Repository
```sh
git clone https://github.com/Gunavardhan18/url-shortener.git
cd url-shortener
```

### Install Dependencies
```sh
go mod tidy
```

### Configure Environment Variables
Create a `.env` file:
```
DB_URL=postgres://user:password@localhost:5432/shortener
REDIS_URL=redis://localhost:6379
```

### Run the Service
```sh
go run cmd/server/main.go
```

### Run with Docker
```sh
docker-compose up --build
```

