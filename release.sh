#!/bin/sh
ver=$(git describe --tags)
mkdir -p "cdetect-$ver"
cp -rv main.go LICENSE README.md "cdetect-$ver"
tar Jcvf "cdetect-$ver.tar.xz" "cdetect-$ver"
