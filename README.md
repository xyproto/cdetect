# CDetect

Utility for figuring out which compiler and compiler version was used for compiling an executable file for Linux (in the ELF format).

* Can detect which compiler was used for Go, GCC, FPC and OCaml.
* Works even with stripped executables.

### Installation (development version):

    go get github.com/xyproto/cdetect

### Example usage

    $ cdetect /bin/sh
    GCC 6.1.1

    $ cdetect /usr/bin/ls
    GCC 6.3.1

### General info

* Version: 0.3
