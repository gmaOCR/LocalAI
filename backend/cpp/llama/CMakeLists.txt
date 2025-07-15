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

# Proto file
get_filename_component(hw_proto "../../../../../../backend/backend.proto" ABSOLUTE)
get_filename_component(hw_proto_path "${hw_proto}" PATH)

# Generated sources
set(hw_proto_srcs "${CMAKE_CURRENT_BINARY_DIR}/backend.pb.cc")
set(hw_proto_hdrs "${CMAKE_CURRENT_BINARY_DIR}/backend.pb.h")
set(hw_grpc_srcs "${CMAKE_CURRENT_BINARY_DIR}/backend.grpc.pb.cc")
set(hw_grpc_hdrs "${CMAKE_CURRENT_BINARY_DIR}/backend.grpc.pb.h")

add_custom_command(
      OUTPUT "${hw_proto_srcs}" "${hw_proto_hdrs}" "${hw_grpc_srcs}" "${hw_grpc_hdrs}"
      COMMAND ${_PROTOBUF_PROTOC}
      ARGS --grpc_out "${CMAKE_CURRENT_BINARY_DIR}"
        --cpp_out "${CMAKE_CURRENT_BINARY_DIR}"
        -I "${hw_proto_path}"
        --plugin=protoc-gen-grpc="${_GRPC_CPP_PLUGIN_EXECUTABLE}"
        "${hw_proto}"
      DEPENDS "${hw_proto}")

# hw_grpc_proto
add_library(hw_grpc_proto
  ${hw_grpc_srcs}
  ${hw_grpc_hdrs}
  ${hw_proto_srcs}
  ${hw_proto_hdrs} )

add_executable(${TARGET} grpc-server.cpp utils.hpp json.hpp httplib.h)

target_include_directories(${TARGET} PRIVATE ../llava)
target_include_directories(${TARGET} PRIVATE ${CMAKE_SOURCE_DIR})

target_link_libraries(${TARGET} PRIVATE common llama mtmd ${CMAKE_THREAD_LIBS_INIT} absl::flags hw_grpc_proto
  absl::flags_parse
  gRPC::${_REFLECTION}
  gRPC::${_GRPC_GRPCPP}
  protobuf::${_PROTOBUF_LIBPROTOBUF})
target_compile_features(${TARGET} PRIVATE cxx_std_11)
if(TARGET BUILD_INFO)
  add_dependencies(${TARGET} BUILD_INFO)
endif()
