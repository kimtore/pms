#!/bin/bash
# Re-build PMS

set -e

intltoolize --force
aclocal -I m4
autoreconf --force --install --verbose
./configure
make clean
make

echo
echo "Build finished."
echo
