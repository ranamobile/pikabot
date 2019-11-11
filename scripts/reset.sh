#!/bin/bash
#
#
docker-compose down
docker-compose build
docker system prune -f --volumes

docker-compose up -d
docker-compose logs -f pikabot
