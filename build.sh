#!/usr/bin/env bash

make linux
sudo docker build -t concurrent/ubuntu:v1 .
PASS "make linux"