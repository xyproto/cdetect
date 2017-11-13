# CompiledWith

Utility for figuring out which compiler and compiler version was sued for compiling an executable file for Linux (in the ELF format).

* Can detect which compiler was used for Go, GCC, FPC and OCaml.
* Works even with stripped executables.

Installation (development version):

    go get github.com/xyproto/compiledwith

Example usage:

    $ compiledwith /bin/sh
    GCC 6.1.1

    $ compiledwith /usr/bin/ls
    GCC 6.3.1

* Version: 0.1
