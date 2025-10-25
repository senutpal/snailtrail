# Snailtrail 

A collection of rate limiting implementations in Go, demonstrating different approaches to throttling HTTP requests.

## Overview

This repository showcases three different rate limiting strategies for HTTP APIs:

1. **tollbooth** - Third-party library implementation using [didip/tollbooth](https://github.com/didip/tollbooth)
2. **tokenbucket** - Custom fixed-window rate limiter implementation
3. **clientlim** - Per-client IP-based rate limiter with automatic cleanup

## Features

- Simple HTTP server examples with `/ping` endpoint
- JSON response formatting
- Configurable rate limits and time windows
- HTTP 429 (Too Many Requests) responses when limits exceeded

## Getting Started

### Prerequisites

- Go 1.24.4 or later

### Installation

```bash
git clone https://github.com/senutpal/snailtrail.git
cd snailtrail
go mod download
```

## Running the Examples
### Tollbooth (Library-based):
```bash
cd tollbooth
go run .
```
### Token Bucket (Simple Rate Limiter):
```bash
cd tokenbucket
go run .
```
### Client Limiter (Per-IP Rate Limiting):
```bash
cd clientlim
go run .
```
Each server runs on http://localhost:8080. Test the /ping endpoint:
```bash
curl http://localhost:8080/ping
```
# Configuration
- tokenbucket: 2 requests per second
- clientlim: 2 requests per second per IP address
- tollbooth: 1 request per second