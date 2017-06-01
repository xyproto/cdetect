#!/bin/sh
for f in /usr/bin/*; do ./elfinfo $f; done | cut -d"," -f2 | grep -a compiler
