#!/bin/bash

# Simple bash script to build basic lnd tools for all the platforms
# we support with the golang cross-compiler.
#
# Copyright (c) 2016 Company 0, LLC.
# Use of this source code is governed by the ISC
# license.

# If no tag specified, use date + version otherwise use tag.
if [[ $1x = x ]]; then
    DATE=`date +%Y%m%d`
    VERSION=`git describe --tags`
    TAG=$DATE-$VERSION
else
    TAG=$1
    VERSION=$1
fi

PACKAGE=sweetd
MAINDIR=$PACKAGE-$TAG
mkdir -p $MAINDIR
cd $MAINDIR

SYS="linux-amd64 linux-armv6 linux-armv7 linux-arm64"

# Use the first element of $GOPATH in the case where GOPATH is a list
# (something that is totally allowed).
GPATH=$(echo $GOPATH | cut -f1 -d:)
COMMITFLAGS="-X main.Commit=$(git rev-parse HEAD) -X main.Version=$VERSION -X main.Date=$(date +%Y-%m-%d)"

for i in $SYS; do
    OS=$(echo $i | cut -f1 -d-)
    ARCH=$(echo $i | cut -f2 -d-)
    ARM=
    
    if [[ $ARCH = "armv6" ]]; then
      ARCH=arm
      ARM=6
    elif [[ $ARCH = "armv7" ]]; then
      ARCH=arm
      ARM=7
    fi
    
    mkdir $PACKAGE-$i-$TAG
    cd $PACKAGE-$i-$TAG
    
    echo "Building:" $OS $ARCH $ARM
    env GOOS=$OS GOARCH=$ARCH GOARM=$ARM GO111MODULE=on go build -v -ldflags "$COMMITFLAGS" github.com/the-lightning-land/sweetd
    cd ..
    
    cp ../README.md ../LICENSE $PACKAGE-$i-$TAG/
    tar cvzf $PACKAGE-$i-$TAG.tar.gz $PACKAGE-$i-$TAG
    rm -r $PACKAGE-$i-$TAG
done

shasum -a 256 * > manifest-$TAG.txt