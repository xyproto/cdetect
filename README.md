# CDetect

Utility for figuring out which compiler and compiler version was used for compiling an executable file for Linux (in the ELF format).

* Supports detection of compiler name and version if an executable was built with one of these compilers:
  * GCC
  * Clang
  * FPC
  * OCaml
  * Go
  * TCC (compiler name only, executable does not include a version number)
* Works even with stripped executables.
* Should work for recent versions of the above compilers, but more testing is needed for supporting old versions.

### Installation (development version):

    go get github.com/xyproto/cdetect

### Example usage

    $ cdetect /bin/sh
    GCC 6.1.1

    $ cdetect /usr/bin/ls
    GCC 6.3.1

### General info

* Version: 0.3
