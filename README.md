<a href="https://github.com/xyproto/elfinfo"><img src="https://raw.githubusercontent.com/xyproto/elfinfo/master/web/elfinfo.png" style="margin-left: 2em" width="200px"></a>

# ELFinfo

Tiny program for emitting only the most basic information about an ELF file.

Can detect the compiler used to compile even a stripped binary for Go, GCC and FPC.

Installation (development version):

    go get github.com/xyproto/elfinfo

Example usage:

    $ elfinfo -c /bin/sh
    GCC 6.1.1

    $ elfinfo /usr/bin/ls
    /usr/bin/ls: stripped=true, compiler=GCC 6.3.1, byteorder=LittleEndian, machine=Advanced Micro Devices x86-64

* Version: 0.5
