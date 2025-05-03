#!/bin/bash

echo "Builds app"
go build -o blog-corsa.bin

cd ./deploy

echo "build the zip package"
./deploy.bin -target service -outdir ~/app/go/igorrun/zips/
cd ~/app/go/igorrun/

echo "update the service"
./update-service.sh

echo "Ready to fly"