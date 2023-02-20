#!/bin/bash

rm -rf nav_service
git pull
go build -tags release
./nav_service