# goduit

Golang implementation of conduit

## Queue Configuration

The application supports multiple queue backing services for article feed processing. You can configure which queue service to use via environment variables.

### Supported Queue Types

- **RabbitMQ** (default)
- **Redis**

### Configuration

Set the following environment variables:

- `QUEUE_TYPE`: The type of queue service to use (`rabbitmq` or `redis`). Default: `rabbitmq`
- `QUEUE_URL`: The connection URL for the queue service

#### RabbitMQ Configuration

```bash
QUEUE_TYPE=rabbitmq
QUEUE_URL=amqp://guest:guest@localhost:5672/
```

#### Redis Configuration

```bash
QUEUE_TYPE=redis
QUEUE_URL=redis://localhost:6379
```

For Redis with authentication:
```bash
QUEUE_URL=redis://user:password@localhost:6379/0
```

### Important Considerations

#### Message Reliability

**RabbitMQ**: Provides message acknowledgment semantics. Failed messages can be requeued for retry, ensuring at-least-once delivery.

**Redis**: Uses a fire-and-forget pattern with List operations (LPUSH/BLPOP). Messages are automatically removed when consumed, and cannot be requeued on failure. Failed messages are lost.

⚠️ **For critical workloads requiring guaranteed message delivery, use RabbitMQ.** If you must use Redis, consider implementing a dead-letter queue pattern or use Redis Streams instead of Lists.

### Migration Guide

#### Migrating from RabbitMQ to Redis

If you have an existing deployment using RabbitMQ and want to switch to Redis:

1. **Add Redis service** to your infrastructure
2. **Update environment variables**:
   ```bash
   QUEUE_TYPE=redis
   QUEUE_URL=redis://your-redis-host:6379
   ```
3. **Drain existing RabbitMQ queues** before switching to prevent message loss
4. **Restart services** to pick up the new configuration

**Note**: This is a one-time migration. Messages in RabbitMQ queues will not be automatically transferred to Redis.

#### Migrating from Redis to RabbitMQ

1. **Add RabbitMQ service** to your infrastructure
2. **Update environment variables**:
   ```bash
   QUEUE_TYPE=rabbitmq
   QUEUE_URL=amqp://guest:guest@rabbitmq-host:5672/
   ```
3. **Restart services** to pick up the new configuration

For docker-compose deployments, see the `.env.example` file for reference configuration.

## Running the Application

### With RabbitMQ (default)

```bash
make run
# or explicitly
make run QUEUE=rabbitmq
```

### With Redis

```bash
make run QUEUE=redis
```

### View Logs

```bash
# All services
make logs-all

# Specific services
make logs-api
make logs-worker
make logs-queue QUEUE=rabbitmq  # or QUEUE=redis
```

### Stop Services

```bash
make stop
```

## Testing

### Unit Tests

Run unit tests without database or queue dependencies:

```bash
go test $(go list ./... | grep -v integrationTests) -count=1
```

### Integration Tests

Integration tests require Docker containers to be running.

#### Quick Start (single queue)

```bash
# Test with RabbitMQ (default)
make run
make test-integration

# Or test with Redis
make run QUEUE=redis
make test-integration
```

#### Test Both Queue Types

Run integration tests against both RabbitMQ and Redis automatically:

```bash
make test-both-queues
```

This will:
1. Start services with RabbitMQ
2. Run all integration tests
3. Stop services
4. Start services with Redis
5. Run all integration tests
6. Stop services
7. Report results

#### Verbose Output

```bash
make test-integration-verbose
```

### CI/CD

The GitHub Actions workflow automatically runs integration tests with both RabbitMQ and Redis using a matrix strategy. See `.github/workflows/run-tests.yml` for details.

## Developer Workflow

### Typical Development Flow

```bash
# Start services (RabbitMQ by default)
make run

# Run unit tests during development
go test ./internal/queue/...

# Run integration tests
make test-integration

# View API logs
make logs-api

# Access API container
make bash

# Stop when done
make stop
```

### Testing Queue Implementations

```bash
# Test with RabbitMQ
make run QUEUE=rabbitmq
make test-integration
make stop

# Test with Redis
make run QUEUE=redis
make test-integration
make stop
```

### Switching Between Queues

Simply restart with a different queue parameter:

```bash
make stop
make run QUEUE=redis
```

The application will automatically use the specified queue backend.

### Testing and Mocks

This repository uses mockery for testing. See more [here](https://vektra.github.io/mockery/latest).
