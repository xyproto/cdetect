.PHONY:all clean nox 

CC = tcc
CXX = clang++
CXXFLAGS ?= -O2

all: clang_hello tcc_hello nox

clang_hello: clang_hello.cpp
	${CXX} ${CXXFLAGS} $< -o $@

tcc_hello: tcc_hello.c
	${CC} ${CXXFLAGS} $< -o $@

nox:
	chmod -x clang_hello
	chmod -x tcc_hello

clean:
	rm -f clang_hello tcc_hello
