#!/bin/bash
docker-compose run trigger-instance-service dep ensure
docker-compose run trigger-instance-service go run migrate.go up