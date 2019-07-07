#!/bin/bash
rm example.s
#go run . 2>&1 |cat >example.ll
llc -relocation-model=pic example.ll
clang example.s
./a.out
