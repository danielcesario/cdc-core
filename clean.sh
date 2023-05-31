eal#!/bin/bash
docker-compose down --rmi all -v
docker rmi $(docker images -q)