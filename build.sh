#!/bin/bash
docker run --rm -it -v "$PWD"/docker:/docker -v "$PWD"/bin:/go/bin golang /docker/build-script.sh
