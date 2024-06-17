# MultiBackend Cache System

## Overview
The MultiBackend Cache System is designed to facilitate efficient data caching by utilizing multiple caching strategies. 
This system supports Memcache and Redis implementations, allowing for flexible, scalable caching solutions suitable for a variety of applications.

## Features
- **Multiple Cache Backends**: Supports Memcache and Redis.
- **Docker Integration**: Easily deployable in a containerized environment with Docker.
- **Concurrency Safe**: Thread-safe implementations ensuring data integrity.

System Architecture
The system is designed with modularity in mind, featuring distinct components for each functionality:

Interface: Defines the cache interface that all backend implementations must adhere to.
MemCache: An in-memory cache implementation.
Redis: A Redis-based cache implementation.
Server: Manages the HTTP server for cache operations.
Utilities: Helper functions for various tasks.
Main: The entry point of the application.

## System Requirements
- Docker
- Go (version 1.x or later)
- Redis server
- MemCacheServer

## Installation

### Clone the repository
```bash
git clone https://github.com/sabarivasan007/MultiCacheSystem.git
cd go-cache-server

#Using Docker
To deploy the cache system using Docker, run the following commands:

docker-compose up --build
This command will set up the necessary containers for Redis, Memcache, and the application server.


APIs Interact with the cache:
examples:
Set a cache entry:
curl -X POST -d '{"key": "example", "value": "123"}' http://localhost:8080/cache?cache=redis


Get a cache entry:
curl -X GET http://localhost:8080/cache/example?cache=redis


Delete a cache entry:
curl -X DELETE http://localhost:8080/cache/example?cache=redis
