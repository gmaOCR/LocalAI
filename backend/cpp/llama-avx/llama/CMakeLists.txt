cmake_minimum_required(VERSION 3.13)
project(llama-backend)

# Configuration options
option(LLAMA_BUILD_GRPC "Build gRPC tools" ON)
option(LLAMA_BUILD_COMMON "Build common library" ON)
option(LLAMA_BUILD_EXAMPLES "Build examples" ON)

# Set build type to Release by default
if(NOT CMAKE_BUILD_TYPE)
    set(CMAKE_BUILD_TYPE Release)
endif()

# Add the main llama.cpp project
add_subdirectory(llama.cpp)

# Configure the main target to export the llama library
set_target_properties(llama PROPERTIES
    OUTPUT_NAME "llama"
    LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin
    ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin
    RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin
)