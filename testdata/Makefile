.PHONY:all clean

CC = tcc
CXX = clang++
CXXFLAGS ?= -O2

all: clang_hello tcc_hello rust_hello rust_hello_stripped

clang_hello: clang_hello.cpp
	${CXX} ${CXXFLAGS} $< -o $@
	chmod -x $@

tcc_hello: tcc_hello.c
	${CC} ${CXXFLAGS} $< -o $@
	chmod -x $@

rust_hello: rust_hello.rs
	rustc $<
	chmod -x $@

rust_hello_stripped: rust_hello.rs
	rustc $< -o $@
	strip $@
	chmod -x $@

clean:
	rm -f clang_hello tcc_hello rust_hello
