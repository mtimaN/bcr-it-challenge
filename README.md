# bcr-it-challenge
# Contributors

* [Radu Marin](https://github.com/radum157)

* [Filip Dumitrascu](https://github.com/filipdumitrascu)

* [Luca Botez](https://github.com/lucabotez)

* [Matei Mantu](https://github.com/mtimaN)

# Secure User Authentication System

A full-stack application featuring user authentication, rate limiting, monitoring, and machine learning capabilities. Built with Go backend, React frontend, and includes comprehensive monitoring and analytics.

## Architecture

- **Backend**: Go HTTP server with JWT authentication and rate limiting
- **Frontend**: React dashboard application
- **Database**: Cassandra for persistent storage, Redis for caching
- **Monitoring**: Prometheus metrics collection and Grafana dashboards
- **ML**: K-means clustering model for user analysis
- **Security**: TLS encryption and comprehensive authentication

## Features

- **User Management**: Registration, authentication with secure password hashing
- **JWT Authentication**: Secure token-based authentication system
- **Rate Limiting**: IP-based request throttling (100 requests/minute)
- **Caching**: Redis-backed user data caching for improved performance
- **CORS Support**: Configured for frontend integration
- **Monitoring**: Real-time metrics and alerting
- **Machine Learning**: User behavior analysis with K-means clustering

## Tech Stack

### Backend
- **Go**: HTTP server with standard library
- **Cassandra**: Primary database for user data
- **Redis**: Caching layer for improved performance
- **JWT**: Token-based authentication
- **Prometheus**: Metrics collection

### Frontend
- **React**: User interface dashboard
- **Vite**: Build tool and development server

### Infrastructure
- **Docker**: Containerized services
- **TLS**: Encrypted connections
- **Grafana**: Monitoring dashboards

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.23+ (for local development)
- Node.js 22+ (for frontend development)
- TLS certificates (see Setup section)

### 1. Clone and Setup
```bash
git clone <repository-url>
cd <project-directory>
```

### 2. Generate TLS Certificates
Create the `certs/` directory and generate your TLS certificates:
```bash
mkdir certs
# Generate your certificates (server.crt and server.key)
# Place them in the certs/ directory
```

### 3. Start Infrastructure Services
```bash
docker-compose up -d
```

This starts:
- Cassandra (port 9042)
- Redis (port 6379)
- Prometheus (port 9090)
- Grafana (port 3000)

### 4. Setup Backend
```bash
# Install dependencies
go mod tidy

# Run the backend server
go run main.go
```

### 5. Setup Frontend
```bash
cd dashboard/
# Follow instructions in dashboard/README.md
```

### 6. Access Services
- **Frontend**: http://localhost:5173
- **Grafana**: http://localhost:3000 (admin/GFPass0319)
- **Prometheus**: http://localhost:9090

## Configuration

### Environment Variables

The application supports the following environment variables with sensible defaults:

```bash
# Database Configuration
CASS_USERNAME=backend          # Cassandra username
CASS_PASSWORD=BPass0319        # Cassandra password
CASS_KEYSPACE=cass_keyspace    # Cassandra keyspace

# Cache Configuration
REDIS_PASSWORD=RPass0319       # Redis password

# Security Configuration
JWT_SECRET=some_secret         # JWT signing secret

# TLS Configuration
TLS_CERT_PATH=certs/server.crt # TLS certificate path
TLS_KEY_PATH=certs/server.key  # TLS private key path
```

### Docker Services Configuration

The `docker-compose.yml` configures:
- **Cassandra**: Authentication enabled, custom keyspace initialization
- **Redis**: Password-protected instance
- **Prometheus**: Configured with custom config and TLS support
- **Grafana**: Pre-configured dashboards and data sources

## Monitoring

### Prometheus Metrics
The application exposes metrics for:
- Request counts and latency
- Rate limiting events
- Database operation metrics
- Authentication success/failure rates

### Grafana Dashboards
Pre-configured dashboards include:
- Application performance metrics
- Database health monitoring
- Rate limiting analytics
- User authentication patterns

Access Grafana at http://localhost:3000 with credentials:
- Username: `admin`
- Password: `GFPass0319`

## Security Features

### Authentication
- Secure password hashing using bcrypt
- JWT tokens with configurable expiration
- Bearer token validation for protected endpoints

### Rate Limiting
- IP-based rate limiting (100 requests/minute)
- Configurable limits per endpoint
- Automatic rate limit violation logging

### Network Security
- TLS encryption for all communications
- CORS configuration for frontend integration
- Secure headers and response handling

## Machine Learning

The project includes a Jupyter notebook with a K-means clustering model for user behavior analysis. The model helps identify user patterns and categories for better service personalization.

To run the ML analysis:
```bash
jupyter notebook analysis.ipynb
```

## Database Schema

### Cassandra Tables
The database schema is automatically initialized from `internal/schema.cql` when the container starts. Key tables include:
- User management tables
- Authentication logs
- Rate limiting records

### Redis Caching
Redis is used for:
- User session caching
- Rate limiting counters
- Temporary data storage

## Development

### Backend Development
```bash
# Run with hot reload (if using air)
air

# Run tests
go test ./...

# Build for production
go build -o app main.go
```

### Frontend Development
```bash
cd dashboard/
# See dashboard/README.md for detailed instructions
npm run dev
```

### Database Management
```bash
# Connect to Cassandra
docker exec -it cassandraDB cqlsh -u admin -p CPass0319

# Connect to Redis
docker exec -it redisDB redis-cli -a RPass0319
```

## Troubleshooting

### Common Issues

1. **Database Connection Errors**
   - Ensure Cassandra and Redis containers are healthy
   - Check network connectivity between services
   - Verify credentials match environment variables

2. **TLS Certificate Issues**
   - Ensure certificates are properly generated and placed in `certs/`
   - Check certificate permissions and paths
   - Verify certificate validity dates

3. **Rate Limiting Issues**
   - Check IP detection logic for proxy/load balancer setups
   - Adjust rate limiting parameters if needed
   - Monitor rate limiting metrics in Grafana

### Health Checks
```bash
# Check service health
docker-compose ps

# View service logs
docker-compose logs [service-name]

# Test backend endpoint
curl -k https://localhost:8080/health
```

## Performance Considerations

- Redis caching reduces database load
- Rate limiting prevents abuse
- Connection pooling for database efficiency
- Prometheus metrics help identify bottlenecks
- TLS termination at application level

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

[Add your license information here]

## Additional Resources

- [Internal Setup Guide](internal/README.md)
- [Frontend Documentation](dashboard/README.md)
- [Cassandra Documentation](https://cassandra.apache.org/doc/)
- [Redis Documentation](https://redis.io/documentation)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
