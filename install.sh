#!/bin/sh

# Infer GOBIN from GOPATH if necessary and able to
[ "$GOBIN" == "" ] && [ "$GOPATH" != "" ]  &&
  echo "[warning] missing \$GOBIN, using $GOPATH/bin" &&
  GOBIN="$GOPATH/bin"

# Make sure we know where to copy to
[ "$GOBIN" == "" ] &&
  echo '[error] $GOBIN not set' &&
  exit 1

SOURCE="$(pwd)/build/pms"
DESTINATION="$GOBIN/pms"

# Make sure we have something to copy
[ ! -f "$SOURCE" ] &&
  echo "[error] $SOURCE not found" &&
  exit 1

# Make room for the binary
[ -f "$DESTINATION" ] &&
  echo "[warning] removing existing binary $DESTINATION" &&
  rm -f "$DESTINATION"

# Check if the user wanted to link instead of copy
[ "$INSTALL_TYPE" == "link" ] &&
  echo "[info] linking $SOURCE to $DESTINATION" &&
  ln -sf "$SOURCE" "$DESTINATION" &&
  exit 0

echo "[info] copying $SOURCE to $DESTINATION"
cp "$SOURCE" "$DESTINATION"
