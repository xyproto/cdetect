#!/bin/sh
if [ ! -x ./cdetect ]; then
  echo 'Could not find executable: cdetect'
  echo 'Has it been built with "go build"?'
  exit 1
fi
for f in /usr/bin/*; do echo -n "$f: "; ./cdetect "$f"; done
