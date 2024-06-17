
# Multi-Backend Cache System

## Overview
The MultiBackend Cache System is designed to facilitate efficient data caching by utilizing multiple caching strategies. 
This system supports Inmemory, Memcache and Redis implementations, allowing for flexible, scalable caching solutions suitable for a variety of applications.
 
## Features
- **Multiple Cache Backends**: Supports Inmemory, Memcache and Redis.
- **Docker Integration**: Easily deployable in a containerized environment with Docker.
- **Concurrency Safe**: Thread-safe implementations ensuring data integrity.

## Tenent Feature (only for inmemmory)
- **In Memory** - Included Tenant Partition
If you want to store the cache data into a specific tenant then have to change config IsTenantBased to true and provide tenantNames (max 3)
There will be a slight changes in Api Endpoints
## Table of Contents

1. [Project Structure](#project-structure)
2. [Getting Started](#getting-started)
3. [Configuration](#configuration)
4. [Deployment](#deployment)
5. [Usage](#usage)
6. [Testing](#testing)
7. [Contributing](#contributing)
8. [License](#license)

## Project Structure

\`\`\`
multi-backend-cache/
├── docker-compose.yml
├── dockerFile
├── go.mod
├── go.sum
├── prometheus.yml
├── router.go
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── Internal/
│   ├── cache/
│   │   ├── inmemory.go
│   │   ├── interface.go
│   │   ├── memcache.go
│   │   ├── mock_cache.go
│   │   ├── redis.go
│   │   └── tenantCacheManager.go
│   ├── config/
│   │   ├── config.go
│   │   └── config.yaml
│   ├── Handler/
│   │   ├── handler.go
│   │   └── middleware.go
│   └── metrices/
│       └── prometheus_model_info.go
├── kuber-deployment/
│   ├── go-service-deployment.yaml
│   ├── mem-statefulset.yaml
│   ├── redis-configmap.yaml
│   ├── redis-statefulset.yaml
│   ├── twem-config.yaml
│   └── twemproxy.yaml
├── packageUtils/Utils/
│   └── utils.go
└── test/
    ├── inmemory_test.go
    ├── server_test.go
    ├── test_mock.go
    └── test_server_test.go
\`\`\`

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/)
- [Docker](https://www.docker.com/products/docker-desktop)

### Installation

1. Clone the repository:
   \`\`\`sh
   git clone https://github.com/sabarivasan007/MultiBackendCacheSystem
   cd multi-backend-cache
   \`\`\`

2. Build the Docker image:
   \`\`\`sh
   docker build -t multi-backend-cache .
   \`\`\`

3. Start the services using Docker Compose:
   \`\`\`sh
   docker-compose up
   \`\`\`

## Configuration

The configuration file \`config.yaml\` is located in the \`Internal/config/\` directory. It includes settings for the different cache backends and other application configurations.

## Deployment

### Docker Compose

To deploy the application using Docker Compose, run:
\`\`\`sh
docker-compose up
\`\`\`


## Usage

The application routes and handlers are defined in \`router.go\` and the \`Internal/Handler/\` directory. API documentation is available in the \`docs/\` directory and can be accessed via Swagger UI.



## System Options:
- redis
- memcache
- inmemory

## APIs Interact with the cache:
examples:

### Set a cache entry:
## POST - http://34.234.207.91:8080/cache?system=inmemory
## with Tenant details  http://34.234.207.91:8080/cache?system=inmemory&tenantId=tenant1
\`\`\`json
{
    "key": "example", 
    "value": "123",
    "ttl": 10
}
\`\`\`

### Get a cache entry:
GET - http://34.234.207.91:8080/cache/exampleKey?system=inmemory
## with Tenant details  http://34.234.207.91:8080/cache/tenant1/exampleKey?system=inmemory&tenantId=tenant1

### Delete a cache entry:
DELETE - http://34.234.207.91:8080/cache/dhoni?system=inmemory
## with Tenant details  http://34.234.207.91:8080/cache?system=inmemory&tenantId=tenant1


### Clear all cache entries:
PUT - http://34.234.207.91:8080/cache/clear?system=inmemory
## with Tenant details  http://34.234.207.91:8080/cache/clear?system=inmemory&tenantId=tenant1



