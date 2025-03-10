#!/bin/bash
curl http://localhost:3000/api/users/login -X POST \
  -H "Content-Type: application/json" \
  -d '{
        "user": {
          "email": "ravi.me@hotmail.com",
          "password": "jumptty8000"
        }
      }' | jq
