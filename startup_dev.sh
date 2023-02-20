#!/bin/bash

rm -rf nav_service
git pull
go build -tags dev
./nav_service