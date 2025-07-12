# Redis Integration

This guide explains how to use Redis as an advanced feature in your mkrgen project.

## Overview

The Redis integration provides a ready-to-use Redis client setup with configuration management, interface abstractions, and example implementations.

## What's Included

When you select Redis as an advanced feature, mkrgen will generate:

- **Redis Client**: A configured Redis client using `github.com/redis/go-redis/v9`
- **Interface Abstraction**: Clean interface definitions for Redis operations
- **Configuration Management**: Environment-based Redis configuration
- **Example Service**: Sample implementation showing how to use Redis in your application

## Generated Files

```
pkg/redis/
├── interface.go       # Redis interface definitions
├── client.go         # Redis client implementation
└── config.go         # Configuration management

internal/example_service/
└── redis_example.go  # Example service using Redis
```

## Usage

### 1. Environment Configuration

The Redis client is configured through environment variables:

```bash
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 2. Using the Redis Client

```go
// Initialize Redis client
redisClient := redis.NewClient()

// Use in your service
exampleService := example_service.NewRedisExampleService(redisClient)
```

### 3. Available Operations

The generated interface includes common Redis operations:

- `Set(key, value string) error`
- `Get(key string) (string, error)`
- `Delete(key string) error`
- `Exists(key string) (bool, error)`

## Example Implementation

The generated example service demonstrates:

- Setting and getting values
- Error handling
- Best practices for Redis integration

## Dependencies

The Redis integration automatically adds:

- `github.com/redis/go-redis/v9` - Redis client for Go

## Best Practices

1. **Environment Variables**: Always use environment variables for Redis configuration
2. **Connection Pooling**: The client automatically handles connection pooling
3. **Error Handling**: Always handle Redis errors appropriately
4. **Key Naming**: Use consistent key naming conventions
5. **Serialization**: Consider using JSON for complex data structures

## Testing

The integration includes example tests showing how to:

- Mock Redis operations
- Test Redis-dependent services
- Use Redis for caching strategies

For more advanced Redis usage patterns, refer to the [Redis Go client documentation](https://github.com/redis/go-redis).
