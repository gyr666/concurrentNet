#!/usr/bin/env bash

pipeline -e

function PASS() {
    echo -e "$*\t\033[38;5;32m[PASS]\033[0m"
}

function FAIL() {
    echo -e "$*\t\033[38;5;124m[FAIL]\033[0m"
}