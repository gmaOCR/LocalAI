set(TARGET grpc-server)
if(DEFINED ENV{VCPKG_INSTALLED})
  set(VCPKG_INSTALLED $ENV{VCPKG_INSTALLED})
  set(CMAKE_PREFIX_PATH "${VCPKG_INSTALLED};${CMAKE_PREFIX_PATH}")
  set(CMAKE_MODULE_PATH "${VCPKG_INSTALLED}/share;${CMAKE_MODULE_PATH}")
  include_directories("${VCPKG_INSTALLED}/include")
  link_directories("${VCPKG_INSTALLED}/lib")
endif()
set(CMAKE_CXX_STANDARD 17)
cmake_minimum_required(VERSION 3.15)
set(TARGET grpc-server)
set(_PROTOBUF_LIBPROTOBUF libprotobuf)
set(_REFLECTION grpc++_reflection)

if (${CMAKE_SYSTEM_NAME} MATCHES "Darwin")
    # Set correct Homebrew install folder for Apple Silicon and Intel Macs
    if (CMAKE_HOST_SYSTEM_PROCESSOR MATCHES "arm64")
        set(HOMEBREW_DEFAULT_PREFIX "/opt/homebrew")
    else()
        set(HOMEBREW_DEFAULT_PREFIX "/usr/local")
    endif()

    link_directories("${HOMEBREW_DEFAULT_PREFIX}/lib")
    include_directories("${HOMEBREW_DEFAULT_PREFIX}/include")
endif()

find_package(absl CONFIG REQUIRED)
set(Protobuf_USE_STATIC_LIBS ON)
set(Protobuf_USE_STATIC ON)
set(Protobuf_ROOT "${VCPKG_INSTALLED}")
set(Protobuf_DIR "${VCPKG_INSTALLED}/share/protobuf")
set(gRPC_ROOT "${VCPKG_INSTALLED}")
set(gRPC_DIR "${VCPKG_INSTALLED}/share/grpc")
find_package(Protobuf CONFIG REQUIRED)
find_package(gRPC CONFIG REQUIRED)

find_program(_PROTOBUF_PROTOC protoc)
set(_GRPC_GRPCPP grpc++)
find_program(_GRPC_CPP_PLUGIN_EXECUTABLE grpc_cpp_plugin)

include_directories(${CMAKE_CURRENT_BINARY_DIR})
include_directories(${Protobuf_INCLUDE_DIRS})
include_directories("${VCPKG_INSTALLED}/include")

message(STATUS "Using protobuf version ${Protobuf_VERSION} | Protobuf_INCLUDE_DIRS: ${Protobuf_INCLUDE_DIRS} | CMAKE_CURRENT_BINARY_DIR: ${CMAKE_CURRENT_BINARY_DIR}")

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
  protobuf::${_PROTOBUF_LIBPROTOBUF}
  "${VCPKG_INSTALLED}/lib/libprotobuf.a"
  "${VCPKG_INSTALLED}/lib/libabsl_flags.a"
  "${VCPKG_INSTALLED}/lib/libabsl_flags_parse.a"
  "${VCPKG_INSTALLED}/lib/libabsl_base.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log.a"
  "${VCPKG_INSTALLED}/lib/libabsl_strings.a"
  "${VCPKG_INSTALLED}/lib/libabsl_synchronization.a"
  "${VCPKG_INSTALLED}/lib/libabsl_time.a"
  "${VCPKG_INSTALLED}/lib/libabsl_cord.a"
  "${VCPKG_INSTALLED}/lib/libabsl_hash.a"
  "${VCPKG_INSTALLED}/lib/libabsl_raw_hash_set.a"
  "${VCPKG_INSTALLED}/lib/libabsl_raw_logging_internal.a"
  "${VCPKG_INSTALLED}/lib/libabsl_throw_delegate.a"
  "${VCPKG_INSTALLED}/lib/libabsl_int128.a"
  "${VCPKG_INSTALLED}/lib/libabsl_malloc_internal.a"
  "${VCPKG_INSTALLED}/lib/libabsl_stacktrace.a"
  "${VCPKG_INSTALLED}/lib/libabsl_symbolize.a"
  "${VCPKG_INSTALLED}/lib/libabsl_demangle_internal.a"
  "${VCPKG_INSTALLED}/lib/libabsl_str_format_internal.a"
  "${VCPKG_INSTALLED}/lib/libabsl_strings_internal.a"
  "${VCPKG_INSTALLED}/lib/libabsl_time_zone.a"
  "${VCPKG_INSTALLED}/lib/libabsl_bad_any_cast_impl.a"
  "${VCPKG_INSTALLED}/lib/libabsl_bad_optional_access.a"
  "${VCPKG_INSTALLED}/lib/libabsl_bad_variant_access.a"
  "${VCPKG_INSTALLED}/lib/libabsl_base.a"
  "${VCPKG_INSTALLED}/lib/libabsl_city.a"
  "${VCPKG_INSTALLED}/lib/libabsl_civil_time.a"
  "${VCPKG_INSTALLED}/lib/libabsl_cleanup.a"
  "${VCPKG_INSTALLED}/lib/libabsl_compare.a"
  "${VCPKG_INSTALLED}/lib/libabsl_compressed_tuple.a"
  "${VCPKG_INSTALLED}/lib/libabsl_container_common.a"
  "${VCPKG_INSTALLED}/lib/libabsl_container_memory.a"
  "${VCPKG_INSTALLED}/lib/libabsl_debugging.a"
  "${VCPKG_INSTALLED}/lib/libabsl_debugging_internal.a"
  "${VCPKG_INSTALLED}/lib/libabsl_dynamic_annotations.a"
  "${VCPKG_INSTALLED}/lib/libabsl_endian.a"
  "${VCPKG_INSTALLED}/lib/libabsl_errno_saver.a"
  "${VCPKG_INSTALLED}/lib/libabsl_examine_stack.a"
  "${VCPKG_INSTALLED}/lib/libabsl_failure_signal_handler.a"
  "${VCPKG_INSTALLED}/lib/libabsl_fast_type_id.a"
  "${VCPKG_INSTALLED}/lib/libabsl_fixed_array.a"
  "${VCPKG_INSTALLED}/lib/libabsl_flat_hash_map.a"
  "${VCPKG_INSTALLED}/lib/libabsl_flat_hash_set.a"
  "${VCPKG_INSTALLED}/lib/libabsl_function_ref.a"
  "${VCPKG_INSTALLED}/lib/libabsl_graphcycles_internal.a"
  "${VCPKG_INSTALLED}/lib/libabsl_has_ostream_operator.a"
  "${VCPKG_INSTALLED}/lib/libabsl_hashtable_debug.a"
  "${VCPKG_INSTALLED}/lib/libabsl_hashtable_debug_hooks.a"
  "${VCPKG_INSTALLED}/lib/libabsl_hashtablez_sampler.a"
  "${VCPKG_INSTALLED}/lib/libabsl_if_constexpr.a"
  "${VCPKG_INSTALLED}/lib/libabsl_inlined_vector.a"
  "${VCPKG_INSTALLED}/lib/libabsl_inlined_vector_internal.a"
  "${VCPKG_INSTALLED}/lib/libabsl_kernel_timeout_internal.a"
  "${VCPKG_INSTALLED}/lib/libabsl_layout.a"
  "${VCPKG_INSTALLED}/lib/libabsl_leak_check.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_entry.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_flags.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_globals.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_initialize.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_append_truncated.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_check_impl.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_check_op.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_conditions.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_config.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_flags.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_fnmatch.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_format.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_globals.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_log_impl.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_log_sink_set.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_message.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_nullguard.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_nullstream.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_proto.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_strip.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_structured.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_structured_proto.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_internal_voidify.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_severity.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_sink.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_sink_registry.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_streamer.a"
  "${VCPKG_INSTALLED}/lib/libabsl_log_structured.a"
  "${VCPKG_INSTALLED}/lib/libabsl_low_level_hash.a"
  "${VCPKG_INSTALLED}/lib/libabsl_malloc_internal.a"
  "${VCPKG_INSTALLED}/lib/libabsl_memory.a"
  "${VCPKG_INSTALLED}/lib/libabsl_meta.a"
  "${VCPKG_INSTALLED}/lib/libabsl_no_destructor.a"
  "${VCPKG_INSTALLED}/lib/libabsl_node_hash_map.a"
  "${VCPKG_INSTALLED}/lib/libabsl_node_hash_set.a"
  "${VCPKG_INSTALLED}/lib/libabsl_node_slot_policy.a"
  "${VCPKG_INSTALLED}/lib/libabsl_non_temporal_arm_intrinsics.a"
  "${VCPKG_INSTALLED}/lib/libabsl_non_temporal_memcpy.a"
  "${VCPKG_INSTALLED}/lib/libabsl_nullability.a"
  "${VCPKG_INSTALLED}/lib/libabsl_numeric.a"
  "${VCPKG_INSTALLED}/lib/libabsl_numeric_representation.a"
  "${VCPKG_INSTALLED}/lib/libabsl_optional.a"
  "${VCPKG_INSTALLED}/lib/libabsl_overload.a"
  "${VCPKG_INSTALLED}/lib/libabsl_periodic_sampler.a"
  "${VCPKG_INSTALLED}/lib/libabsl_poison.a"
  "${VCPKG_INSTALLED}/lib/libabsl_prefetch.a"
  "${VCPKG_INSTALLED}/lib/libabsl_pretty_function.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_bit_gen_ref.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_distributions.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_distribution_caller.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_distribution_test_util.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_fast_uniform_bits.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_fastmath.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_generate_real.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_iostream_state_saver.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_mock_helpers.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_nonsecure_base.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_pcg_engine.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_platform.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_pool_urbg.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_randen.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_randen_engine.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_randen_hwaes.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_randen_hwaes_impl.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_randen_slow.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_salted_seed_seq.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_seed_material.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_traits.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_uniform_helper.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_internal_wide_multiply.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_random.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_seed_gen_exception.a"
  "${VCPKG_INSTALLED}/lib/libabsl_random_seed_sequences.a"
  "${VCPKG_INSTALLED}/lib/libabsl_raw_hash_map.a"
  "${VCPKG_INSTALLED}/lib/libabsl_raw_hash_set.a"
  "${VCPKG_INSTALLED}/lib/libabsl_raw_logging_internal.a"
  "${VCPKG_INSTALLED}/lib/libabsl_sample_recorder.a"
  "${VCPKG_INSTALLED}/lib/libabsl_scoped_set_env.a"
  "${VCPKG_INSTALLED}/lib/libabsl_span.a"
  "${VCPKG_INSTALLED}/lib/libabsl_spinlock_wait.a"
  "${VCPKG_INSTALLED}/lib/libabsl_stacktrace.a"
  "${VCPKG_INSTALLED}/lib/libabsl_status.a"
  "${VCPKG_INSTALLED}/lib/libabsl_statusor.a"
  "${VCPKG_INSTALLED}/lib/libabsl_str_format.a"
  "${VCPKG_INSTALLED}/lib/libabsl_str_format_internal.a"
  "${VCPKG_INSTALLED}/lib/libabsl_strerror.a"
  "${VCPKG_INSTALLED}/lib/libabsl_string_view.a"
  "${VCPKG_INSTALLED}/lib/libabsl_strings.a"
  "${VCPKG_INSTALLED}/lib/libabsl_strings_internal.a"
  "${VCPKG_INSTALLED}/lib/libabsl_symbolize.a"
  "${VCPKG_INSTALLED}/lib/libabsl_synchronization.a"
  "${VCPKG_INSTALLED}/lib/libabsl_throw_delegate.a"
  "${VCPKG_INSTALLED}/lib/libabsl_time.a"
  "${VCPKG_INSTALLED}/lib/libabsl_time_zone.a"
  "${VCPKG_INSTALLED}/lib/libabsl_utility.a"
  "${VCPKG_INSTALLED}/lib/libabsl_variant.a"
  "${VCPKG_INSTALLED}/lib/libabsl_vlog_config_internal.a"
  "${VCPKG_INSTALLED}/lib/libabsl_vlog_is_on.a"
)
target_compile_features(${TARGET} PRIVATE cxx_std_11)
if(TARGET BUILD_INFO)
  add_dependencies(${TARGET} BUILD_INFO)
endif()
