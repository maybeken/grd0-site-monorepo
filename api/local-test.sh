#!/bin/bash
docker build -t grd0/api:latest . && \
docker run --rm -p 80:80 \
	-v ./data/api.db:/opt/app/data/api.db \
	-v ./data/files/files.json:/opt/app/data/files/files.json \
	--env-file ./api.env \
	grd0/api
