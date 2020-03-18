#!/usr/bin/env bash
. ./scripts/ci-helper.sh
make linux
sudo docker build -t concurrent/ubuntu:v1 .
PASS "make linux"
