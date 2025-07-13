#!/bin/bash

set -e

# Copy necessary files for the grpc-server
cp -r grpc-server.cpp llama.cpp/tools/grpc-server/
cp -rfv llama.cpp/vendor/nlohmann/json.hpp llama.cpp/tools/grpc-server/
cp -rfv llama.cpp/tools/server/utils.hpp llama.cpp/tools/grpc-server/
cp -rfv llama.cpp/vendor/cpp-httplib/httplib.h llama.cpp/tools/grpc-server/

# Copy the new CMakeLists.txt to the grpc-server directory
cp -f CMakeLists.txt llama.cpp/tools/grpc-server/

# Add grpc-server to the tools CMakeLists.txt
set +e
if grep -q "add_subdirectory(grpc-server)" llama.cpp/tools/CMakeLists.txt; then
    echo "grpc-server already added"
else
    echo "add_subdirectory(grpc-server)" >> llama.cpp/tools/CMakeLists.txt
fi
set -e