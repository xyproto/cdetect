.PHONY:all clean nox 

CC = tcc
CXX = clang++
CXXFLAGS ?= -O2

all: clang_hello tcc_hello rust_hello nox

clang_hello: clang_hello.cpp
	${CXX} ${CXXFLAGS} $< -o $@

tcc_hello: tcc_hello.c
	${CC} ${CXXFLAGS} $< -o $@

rust_hello: rust_hello.rs
	rustc $<
	strip $@

nox:
	chmod -x clang_hello
	chmod -x tcc_hello
	chmod -x rust_hello

clean:
	rm -f clang_hello tcc_hello rust_hello
