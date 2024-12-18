#!/bin/bash

rm -rf nav_service
git pull
go build -tags release
nohup ./nav_service > nav_service.log &