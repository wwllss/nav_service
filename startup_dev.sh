#!/bin/bash

rm -rf nav_service
git pull
go build -tags dev
nohup ./nav_service conf/config.yaml > nav_service.log 2>&1 &