# ExpectoLang

## Building
- `git clone --recurse-submodules https://github.com/Noituri/ExpectoLang.git`
- `cd ExpectoLang/llvm/bindings/go`
- `./build.sh -DCMAKE_BUILD_TYPE=Debug -DLLVM_TARGETS_TO_BUILD=host -DBUILD_SHARED_LIBS=ON`
- `cd ../../..`
- `go build .`
