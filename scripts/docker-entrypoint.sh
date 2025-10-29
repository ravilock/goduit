#!/bin/sh
set -e

# Load JWT keys as base64 environment variables if they exist
if [ -f "/app/jwtRS256.key" ]; then
    export JWT_PRIVATE_KEY_BASE64=$(base64 -w 0 /app/jwtRS256.key)
fi

if [ -f "/app/jwtRS256.key.pub" ]; then
    export JWT_PUBLIC_KEY_BASE64=$(base64 -w 0 /app/jwtRS256.key.pub)
fi

# Execute the original command
exec "$@"
