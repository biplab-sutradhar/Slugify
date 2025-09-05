# Slugify 

## Overview

This project provides a **scalable URL shortener** system with both core features and an admin dashboard. The URL shortening service allows users to shorten long URLs, track analytics on link usage, and manage links programmatically via API. It also features an intuitive web-based admin interface for link creation, analytics, and API key management.

## Core Features

### URL Shortening
- Accepts long URLs and returns a unique shortened URL.
  
### URL Resolution
- Redirects short links to the original URLs, ensuring proper HTTP status codes.

### Link Management
- Allows creation, viewing, activation/deactivation, and deletion of shortened URLs through API endpoints.

### Analytics Tracking
- Tracks and stores click events with timestamps, user agents, referrers, and geolocation data.

### API Key Authentication
- Supports API key-based authentication with scoped access and rate limiting per key.

## Admin Dashboard Features

- **Link Creation Interface**: A web form that validates and generates short links instantly.
- **Link Listing**: Paginated view with search, filtering, and sorting capabilities.
- **Analytics Visualization**: Charts displaying click counts, geographic distribution, and time-series data.
- **API Key Management**: Interface to create, revoke, and monitor API key usage.

## Data Management

- **Link Expiration**: Support for optional TTL (Time-to-Live) on links.
- **Bulk Operations**: Import and export functionality for link backups and management.
- **Click Aggregation**: Aggregated click data for detailed analytics.

## Performance and Scalability

- Sub-10ms cache hits using Redis for popular URLs.
- 1000+ requests per second throughput under load.
- Horizontal scaling to 8 API containers without state issues.
- 95th percentile latency of under 100ms for API responses.

## Security

- IP and API key-based rate limiting to prevent abuse.
- URL validation and sanitization to ensure no malicious inputs.
- IP address hashing and minimal personally identifiable information (PII) storage.
- Secure API key generation, validation, and storage.

## Tech Stack

- **Backend**: Go (Gin web framework)
- **Database**: PostgreSQL (via Neon free tier)
- **Cache**: Redis (via Upstash free tier)
- **Frontend**: React with TanStack Query
- **Containerization**: Docker (multi-stage builds, docker-compose)
  
## Non-Functional Requirements

- **Uptime Target**: 99.9% availability during load testing periods.
- **Graceful Degradation**: Ability to fall back to the database when Redis is unavailable.

## Testing

- **Load Testing**: Simulate concurrent users with k6.
- **Performance Testing**: Validate Redis hit rates, rate-limiting, and autoscaling of the database.
- **Functional Testing**: Extensive unit tests, integration tests, and React component tests.

## Deployment

- **Docker Compose**: Full local environment setup for API, database, cache, and load balancer.
- **CI/CD**: GitHub Actions for automated testing, builds, and performance benchmarks.
- **Monitoring**: Prometheus integration for latency, throughput, and error rates.

## Performance Metrics

- Linear throughput scaling: 200 requests per second per instance.
- Sub-100ms response time for cache hits.
- Redis hit rate: 80%+ under realistic traffic.

## How to Run Locally

### Prerequisites

- Docker and Docker Compose
- Go 1.21+
- Node.js (for the frontend)

### Setup

1. Clone the repository:

```bash
git clone https://github.com/yourusername/url-shortener.git
cd url-shortener
