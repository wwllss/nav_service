#!/bin/bash

rm -rf nav_service
git pull
go build -tags dev
nohup ./nav_service > nav_service.log &