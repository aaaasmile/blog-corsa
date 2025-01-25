#!/bin/bash

echo "Builds app"
go build -o comment-blog.bin

cd ./deploy

echo "build the zip package"
./deploy.bin -target service -outdir ~/app/go/comment-blog/zips/
cd ~/app/go/comment-blog/

echo "update the service"
./update-service.sh

echo "Ready to fly"