# :microscope: CDetect

Utility for figuring out which compiler and compiler version was used for compiling an executable file for Linux (in the ELF format).

### Installation (development version):

    go get github.com/xyproto/cdetect

### Example usage

    $ cdetect /bin/sh
    GCC 7.3.0

    $ cdetect /usr/bin/ls
    GCC 7.2.1

    $ cdetect testdata/rust_hello
    Rust 1.27.0-nightly

    $ cdetect go
    Go 1.10.2

### Features and limitations

* Supports detection of compiler name and version if an executable was built with one of these compilers:
  * GCC
  * Clang
  * FPC
  * OCaml
  * Go
  * TCC (compiler name only, executable does not include a version number)
  * Rust (for stripped executables, only the compiler name and GCC version used for linking)
* Works even with stripped executables.
* Should work for recent versions of the above compilers, but more testing is needed for supporting all old versions.

### Changelog

#### 0.5 to 0.5.1

* Fix an issue with version detection for Rust.

#### 0.4 to 0.5

* Add support for detecting executables compiled with Rust.
* Will now look for the given filename in PATH, if not found.

#### 0.3 to 0.4

* Add support for detecting executables compiled with Clang or TCC.

#### 0.2 to 0.3

* Fix issue #1, detection of executables compiled with GCC on Void Linux.

#### 0.1 to 0.2

* Rename the utility to `cdetect`.

#### 0.1

* Support for detecting various compilers and compiler version numbers.

### General info

* Version: 0.5.1
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
* License: MIT
