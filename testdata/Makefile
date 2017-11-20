.PHONY: clean

SRC = main.cpp
EXE = clang_hello
CXX = clang++
CXXFLAGS ?= -O2

${EXE}: ${SRC}
	${CXX} ${CXXFLAGS} $< -o $@

clean:
	rm -f ${EXE}
