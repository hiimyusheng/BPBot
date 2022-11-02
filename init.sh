#!/bin/bash
docker pull mongo:4.4
docker run --name mongo4 -v $(pwd)/data:/data/db -d -p 27017:27017 --rm mongo:4.4
