#! /bin/bash
curl -X POST -H "Content-Type: application/json" -d '{"payload": {"email_address": "srt0422@yahoo.com", "password": "password"}}' -k http://localhost:8083/v1/authenticate-user