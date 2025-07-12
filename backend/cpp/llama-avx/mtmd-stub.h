#pragma once

#include <cstdint>
#include <memory>
#include <string>
#include <vector>

// Forward declarations for llama.cpp types
typedef struct llama_context llama_context;
typedef struct llama_model llama_model;
typedef int32_t llama_pos;

// MTMD stub types and constants
typedef struct mtmd_context mtmd_context;
typedef struct mtmd_input_chunk {
    char dummy[64]; // Make it a complete type
} mtmd_input_chunk;

// MTMD constants
#define MTMD_INPUT_CHUNK_TYPE_TEXT "text"
#define MTMD_INPUT_CHUNK_TYPE_IMAGE "image"
#define MTMD_INPUT_CHUNK_TYPE_AUDIO "audio"

// MTMD input text structure
typedef struct {
    const char* text;
    size_t n_text;
} mtmd_input_text;

// MTMD context parameters
typedef struct {
    bool use_gpu;
    int n_gpu_layers;
    int n_ctx;
    int n_batch;
    int n_ubatch;
    int n_threads;
    int n_threads_batch;
    float rope_freq_base;
    float rope_freq_scale;
    bool flash_attn;
    bool use_mmap;
    bool use_mlock;
    bool print_timings;
    int verbosity;
} mtmd_context_params;

// MTMD namespace equivalents
namespace mtmd {
    using input_chunk_ptr = std::shared_ptr<mtmd_input_chunk>;
    using input_chunks_ptr = std::unique_ptr<mtmd_input_chunk*>;
    
    struct input_chunks {
        input_chunks_ptr ptr;
        size_t chunk_count;
        input_chunks(mtmd_input_chunk** chunks) : ptr(chunks), chunk_count(0) {}
        size_t size() const { return chunk_count; }
        mtmd_input_chunk* operator[](size_t index) { return ptr ? ptr.get()[index] : nullptr; }
    };
    
    struct bitmap {
        void* ptr;
        std::string id;
        bitmap(void* p) : ptr(p) {}
        const uint8_t* data() { return static_cast<const uint8_t*>(ptr); }
        size_t n_bytes() { return 0; }
        void set_id(const char* new_id) { id = new_id; }
    };
    
    struct bitmaps {
        std::vector<bitmap> entries;
        std::vector<void*> c_ptr_storage;
        void** c_ptr() { 
            c_ptr_storage.clear();
            for (auto& entry : entries) {
                c_ptr_storage.push_back(entry.ptr);
            }
            return c_ptr_storage.data();
        }
        size_t size() const { return entries.size(); }
    };
}

// Stub function declarations
#ifdef __cplusplus
extern "C" {
#endif

// Context management stubs
mtmd_context_params mtmd_context_params_default();
mtmd_context* mtmd_init_from_file(const char* mmproj_path, llama_model* model, mtmd_context_params params);
void mtmd_free(mtmd_context* ctx);

// Capability stubs
bool mtmd_support_vision(mtmd_context* ctx);
bool mtmd_support_audio(mtmd_context* ctx);

// Input processing stubs
mtmd_input_chunk** mtmd_input_chunks_init();
int32_t mtmd_tokenize(mtmd_context* ctx, mtmd_input_chunk** chunks, mtmd_input_text* inp_txt, void** bitmaps, size_t n_bitmaps, int32_t* tokens, int32_t n_max);

// Chunk management stubs
const char* mtmd_input_chunk_get_type(const mtmd_input_chunk* chunk);
int32_t mtmd_input_chunk_get_n_pos(const mtmd_input_chunk* chunk);
const char* mtmd_input_chunk_get_id(const mtmd_input_chunk* chunk);
int32_t* mtmd_input_chunk_get_tokens_text(const mtmd_input_chunk* chunk, int32_t* n_tokens);
mtmd_input_chunk* mtmd_input_chunk_copy(const mtmd_input_chunk* chunk);

// Helper function stubs
int32_t mtmd_helper_eval_chunk_single(mtmd_context* mctx, llama_context* ctx, mtmd_input_chunk* chunk, int32_t n_past, int32_t seq_id, int32_t n_batch, bool logits_last, llama_pos* new_n_past);
void* mtmd_helper_bitmap_init_from_buf(mtmd_context* ctx, const void* data, size_t size);

// Default marker stub
const char* mtmd_default_marker();

#ifdef __cplusplus
}
#endif
