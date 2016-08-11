#!/bin/bash
docker run --rm -it -p 8080:8080 -v "$PWD"/docker:/docker golang /docker/run.sh
