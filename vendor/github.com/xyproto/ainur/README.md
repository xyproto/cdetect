# Ainur

Go module for figuring out which compiler and compiler version was used for compiling an executable file for Linux (in the ELF format).

### Features and limitations

* Supports detection of compiler name and version if an executable was built with one of these compilers:
  * GCC
  * Clang
  * FPC
  * OCaml
  * Go
  * TCC (compiler name only, TCC does not store the version number in the executables)
  * Rust (for stripped executables, only the compiler name and GCC version used for linking)
* Works even with stripped executables.
* Should work for recent versions of all of the above compilers. Executables produced with old versions of the compilers may need more testing.

### General info

* Version: 1.0.0
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
* License: MIT
