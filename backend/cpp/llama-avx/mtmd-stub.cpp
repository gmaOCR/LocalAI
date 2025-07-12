#include "mtmd-stub.h"
#include <cstring>
#include <iostream>
#include "ggml.h"

// Context management stubs
mtmd_context_params mtmd_context_params_default() {
    mtmd_context_params params;
    memset(&params, 0, sizeof(params));
    params.use_gpu = false;
    params.n_gpu_layers = 0;
    params.n_ctx = 2048;
    params.n_batch = 512;
    params.n_ubatch = 512;
    params.n_threads = 1;
    params.n_threads_batch = 1;
    params.rope_freq_base = 10000.0f;
    params.rope_freq_scale = 1.0f;
    params.flash_attn = false;
    params.use_mmap = true;
    params.use_mlock = false;
    params.print_timings = false;
    params.verbosity = 0;
    return params;
}

mtmd_context* mtmd_init_from_file(const char* mmproj_path, llama_model* model, mtmd_context_params params) {
    (void)mmproj_path;
    (void)model;
    (void)params;
    // Return nullptr to indicate no multimodal support
    return nullptr;
}

void mtmd_free(mtmd_context* ctx) {
    (void)ctx;
    // No-op for stub
}

// Capability stubs
bool mtmd_support_vision(mtmd_context* ctx) {
    (void)ctx;
    return false;
}

bool mtmd_support_audio(mtmd_context* ctx) {
    (void)ctx;
    return false;
}

// Input processing stubs
mtmd_input_chunk** mtmd_input_chunks_init() {
    return nullptr;
}

int32_t mtmd_tokenize(mtmd_context* ctx, mtmd_input_chunk** chunks, mtmd_input_text* inp_txt, void** bitmaps, size_t n_bitmaps, int32_t* tokens, int32_t n_max) {
    (void)ctx;
    (void)chunks;
    (void)inp_txt;
    (void)bitmaps;
    (void)n_bitmaps;
    (void)tokens;
    (void)n_max;
    return 0;
}

// Chunk management stubs
const char* mtmd_input_chunk_get_type(const mtmd_input_chunk* chunk) {
    (void)chunk;
    return "text";
}

int32_t mtmd_input_chunk_get_n_pos(const mtmd_input_chunk* chunk) {
    (void)chunk;
    return 0;
}

const char* mtmd_input_chunk_get_id(const mtmd_input_chunk* chunk) {
    (void)chunk;
    return "";
}

int32_t* mtmd_input_chunk_get_tokens_text(const mtmd_input_chunk* chunk, int32_t* n_tokens) {
    (void)chunk;
    if (n_tokens) *n_tokens = 0;
    return nullptr;
}

mtmd_input_chunk* mtmd_input_chunk_copy(const mtmd_input_chunk* chunk) {
    (void)chunk;
    return nullptr;
}

// Helper function stubs
int32_t mtmd_helper_eval_chunk_single(mtmd_context* mctx, llama_context* ctx, mtmd_input_chunk* chunk, int32_t n_past, int32_t seq_id, int32_t n_batch, bool logits_last, llama_pos* new_n_past) {
    (void)mctx;
    (void)ctx;
    (void)chunk;
    (void)n_past;
    (void)seq_id;
    (void)n_batch;
    (void)logits_last;
    if (new_n_past) *new_n_past = n_past;
    return 0;
}

void* mtmd_helper_bitmap_init_from_buf(mtmd_context* ctx, const void* data, size_t size) {
    (void)ctx;
    (void)data;
    (void)size;
    return nullptr;
}

// Default marker stub
const char* mtmd_default_marker() {
    return "";
}
