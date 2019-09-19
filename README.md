# :microscope: CDetect

Utility for figuring out which compiler and compiler version was used for compiling an executable file for Linux (in the ELF format).

### Installation

This requires Go 1.12 or later and will install the development version of `cdetect`:

    go get -u github.com/xyproto/cdetect

### Example usage

    $ cdetect /bin/sh
    GCC 8.1.1

    $ cdetect /usr/bin/ls
    GCC 8.2.0

    $ cdetect testdata/rust_hello
    Rust 1.27.0-nightly

    $ cdetect go
    Go 1.11.2

### Features and limitations

* Supports detection of compiler name and version if an executable was built with one of these compilers:
  * GCC
  * Clang
  * FPC
  * OCaml
  * Go
  * TCC (compiler name only, TCC does not store the version number in the executables)
  * Rust (for stripped executables, only the compiler name and GCC version used for linking)
  * GHC
* Works even with stripped executables.
* Should work for recent versions of all of the above compilers. Executables produced with old versions of the compilers may need more testing.

### Changelog

#### 0.5.3 to 0.5.4

* Add support for executables built with GCC 8 for 32-bit PowerPC.

#### 0.5.2 to 0.5.3

* Add detection of compiler name and version from executables built with `ghc` (Haskell).

#### 0.5.1 to 0.5.2

* Refactor out code to the [ainur](https://github.com/xyproto/ainur) module.
* Better support for 32-bit PowerPC ELF files.

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

* Version: 0.5.4
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
* License: MIT
