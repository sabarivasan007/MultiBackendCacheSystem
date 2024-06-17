# Benchmark Tool - MultiBackend Cache System

## Overview
This tool allows you to run a benchmark for different cache operations (set, get, delete)
with configurable parameters such as the number of parallel clients, number of requests, payload size, 
and whether to use multithreading. Adjust the HTTP request logic to match your actual server endpoints and methods.

## System Requirements
- Docker
- Go (version 1.x or later)
- Cobra-cli library
- Redis server
- MemCacheServer

## Installation
```bash
git clone https://github.com/sabarivasan007/MultiCacheSystem.git

cd benchmark-tool

go build -o benchmark.exe
# This command will build the application to executable file

The base command is -> benchmark run 
Flags:
  -c, --clients int       Number of parallel connections (default 20)
  -h, --hostname string   Hostname (default "127.0.0.1")
  -n, --num int           Number of requests (default 1000)
  -p, --port int          Server post (default 8080)
  -s, --server string     Cache server to use (default "redis")
  -d, --size int          Data size in bytes (default 3)
  -t, --test string       Type of test to run (default "set")

Global Flags:
      --help   help for this command

Available options for flags

server  - redis, memcache, inmemory
req typ - get, set, delete
clients - 1 to 100
num     - 1 to 1000000
size    - 1 to 1000000
