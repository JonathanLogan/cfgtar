#!/bin/sh

tar -c testData/ | go run main.go cfg.data | tar -x -C tmp/
