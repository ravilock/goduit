#!/bin/bash
curl http://localhost:3000/api/users -X POST \
  -H "Content-Type: application/json" \
  -d '{
        "user": {
          "username": "raylok",
          "email": "ravi.me@hotmail.com",
          "password": "jumptty8000"
        }
      }' | jq
