#!/bin/bash
# Re-build PMS

set -xe

mkdir -p build
cd build
cmake ..
make clean
make
