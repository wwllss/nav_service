#!/bin/bash

rm -rf nav_service
git pull
go build -tags release
nohup ./nav_service conf/config.yaml > nav_service.log 2>&1 &