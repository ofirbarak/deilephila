#!/bin/bash

set -e

docker compose up -d
sleep 5  # Wait for mongo to running
docker exec -it mongo1 mongosh --eval 'rs.initiate({_id: "rs0",version: 1,members: [{ _id: 0, host: "localhost:27017" },]})'
