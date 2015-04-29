#!/bin/bash

set -e;

docker build -t ipx800_watcher_image .;

docker run -d --name ipx800_watcher ipx800_watcher_image;