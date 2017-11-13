#!/bin/sh
if [ ! -x ./compiledwith ]; then
  echo 'Could not find executable: compiledwith'
  echo 'Has it been built with "go build"?'
  exit 1
fi
for f in /usr/bin/*; do echo -n "$f: "; ./compiledwith "$f"; done
