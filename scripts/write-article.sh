#!/bin/bash
curl http://localhost:3000/api/articles -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
        "article": {
          "title": "How to Use cURL for API Requests",
          "description": "A guide on making API requests with cURL.",
          "body": "This article explains how to use cURL to send POST requests to APIs...",
          "tagList": ["api", "curl", "http", "rest"]
        }
      }' | jq
