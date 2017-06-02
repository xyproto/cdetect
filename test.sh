#!/bin/sh
for f in /usr/bin/*; do
  echo compiler=$f; echo "$f: " $(./elfinfo "$f" | cut -d, -f2 || echo "$f")
done | grep -ia " compiler="
