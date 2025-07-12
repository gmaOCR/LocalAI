#!/bin/bash

set -e

# First, ensure we have the complete llama.cpp structure
if [ ! -f llama.cpp/CMakeLists.txt ]; then
    echo "Copying complete llama.cpp structure from llama/ directory..."
    if [ -d llama/llama.cpp ]; then
        # Copy the entire structure
        cp -r llama/llama.cpp/* llama.cpp/
        echo "Copied llama.cpp structure from llama/"
    else
        echo "Error: llama/llama.cpp directory not found"
        exit 1
    fi
fi

# Create the grpc-server directory if it doesn't exist
mkdir -p llama.cpp/tools/grpc-server

# Copy necessary files for the grpc-server
cp grpc-server.cpp llama.cpp/tools/grpc-server/
cp grpc-server-CMakeLists.txt llama.cpp/tools/grpc-server/CMakeLists.txt

# Copy required headers and libraries
cp -rfv llama/llama.cpp/vendor/nlohmann/json.hpp llama.cpp/tools/grpc-server/
cp -rfv llama/llama.cpp/tools/server/utils.hpp llama.cpp/tools/grpc-server/
cp -rfv llama/llama.cpp/vendor/cpp-httplib/httplib.h llama.cpp/tools/grpc-server/

# Copy server.cpp for inclusion
if [ -f llama/llama.cpp/tools/server/server.cpp ]; then
    cp -rfv llama/llama.cpp/tools/server/server.cpp llama.cpp/tools/grpc-server/
fi

# Copy MTMD stub files
cp -rfv mtmd-stub.h llama.cpp/tools/grpc-server/
cp -rfv mtmd-stub.cpp llama.cpp/tools/grpc-server/

# Fix missing header includes in copied files
echo "Fixing missing header includes..."
for file in llama.cpp/tools/grpc-server/utils.hpp llama.cpp/tools/grpc-server/server.cpp; do
    if [ -f "$file" ]; then
        # Remove mtmd.h and mtmd-helper.h includes
        sed -i '/^#include "mtmd\.h"$/d' "$file"
        sed -i '/^#include "mtmd-helper\.h"$/d' "$file"
        # Add mtmd-stub.h include
        sed -i '/^#include <cinttypes>$/a #include "mtmd-stub.h"' "$file"
        # For server.cpp, also add mtmd-stub.h after speculative.h if not already added
        if [[ "$file" == *"server.cpp" ]]; then
            sed -i '/^#include "speculative\.h"$/a #include "mtmd-stub.h"' "$file"
        fi
        
        # Fix utils.hpp specific issues
        if [[ "$file" == *"utils.hpp" ]]; then
            # Fix include paths
            sed -i 's|#include <cpp-httplib/httplib.h>|#include "httplib.h"|g' "$file"
            sed -i 's|#include <nlohmann/json.hpp>|#include "json.hpp"|g' "$file"
            # Fix string comparisons for mtmd types
            sed -i 's/strcmp(type, MTMD_INPUT_CHUNK_TYPE_IMAGE) == 0/strcmp(type, MTMD_INPUT_CHUNK_TYPE_IMAGE) == 0/g' "$file"
            sed -i 's/strcmp(type, MTMD_INPUT_CHUNK_TYPE_AUDIO) == 0/strcmp(type, MTMD_INPUT_CHUNK_TYPE_AUDIO) == 0/g' "$file"
            sed -i 's/strcmp(type, MTMD_INPUT_CHUNK_TYPE_TEXT) == 0/strcmp(type, MTMD_INPUT_CHUNK_TYPE_TEXT) == 0/g' "$file"
            # Fix the comparison in utils.hpp
            sed -i 's/mtmd_input_chunk_get_type(chunk.get()) == MTMD_INPUT_CHUNK_TYPE_IMAGE/strcmp(mtmd_input_chunk_get_type(chunk.get()), MTMD_INPUT_CHUNK_TYPE_IMAGE) == 0/g' "$file"
            # Fix n_tokens type from size_t to int32_t
            sed -i 's/size_t n_tokens;/int32_t n_tokens;/g' "$file"
            sed -i 's/for (size_t i = 0; i < n_tokens; ++i)/for (int32_t i = 0; i < n_tokens; ++i)/g' "$file"
        fi
        
        echo "Fixed includes in $file"
    fi
done

# Copy vendor dependencies
if [ -d llama.cpp/vendor/minja ]; then
    cp -rfv llama.cpp/vendor/minja llama.cpp/tools/grpc-server/
fi

# Fix grpc-server.cpp by disabling multimodal code
if [ -f llama.cpp/tools/grpc-server/grpc-server.cpp ]; then
    echo "Disabling multimodal code in grpc-server.cpp..."
    # Force has_mtmd to false to skip multimodal code paths
    sed -i 's/const bool has_mtmd = ctx_server\.mctx != nullptr;/const bool has_mtmd = false; \/\/ Disabled multimodal support/' llama.cpp/tools/grpc-server/grpc-server.cpp
    
    # Create a Python script to fix the multimodal code blocks
    cat > fix_grpc_server.py << 'PYEOF'
import re

# Read the file
with open('llama.cpp/tools/grpc-server/grpc-server.cpp', 'r') as f:
    content = f.read()

# Replace the first mtmd block (around line 466)
pattern1 = r'(if \(has_mtmd\) \{[^}]*?mtmd_input_text inp_txt = \{[^}]*?\};[^}]*?mtmd::input_chunks chunks\(mtmd_input_chunks_init\(\)\);[^}]*?auto bitmaps_c_ptr = bitmaps\.c_ptr\(\);[^}]*?int32_t tokenized = mtmd_tokenize\([^;]*?bitmaps_c_ptr\.size\(\)\);[^}]*?if \(tokenized != 0\) \{[^}]*?\}[^}]*?server_tokens tmp\(chunks, true\);[^}]*?inputs\.push_back\(std::move\(tmp\)\);[^}]*?\})'
replacement1 = r'''if (has_mtmd) {
                // Multimodal code disabled - using non-multimodal fallback
                auto tokenized_prompts = tokenize_input_prompts(ctx_server.vocab, prompt, true, true);
                for (auto & p : tokenized_prompts) {
                    auto tmp = server_tokens(p, ctx_server.mctx != nullptr);
                    inputs.push_back(std::move(tmp));
                }
            }'''

# Replace the second mtmd block (around line 659)
pattern2 = r'(if \(has_mtmd\) \{[^}]*?mtmd_input_text inp_txt = \{[^}]*?\};[^}]*?mtmd::input_chunks chunks\(mtmd_input_chunks_init\(\)\);[^}]*?auto bitmaps_c_ptr = bitmaps\.c_ptr\(\);[^}]*?int32_t tokenized = mtmd_tokenize\([^;]*?bitmaps_c_ptr\.size\(\)\);[^}]*?if \(tokenized != 0\) \{[^}]*?\}[^}]*?server_tokens tmp\(chunks, true\);[^}]*?inputs\.push_back\(std::move\(tmp\)\);[^}]*?\})'
replacement2 = r'''if (has_mtmd) {
                // Multimodal code disabled - using non-multimodal fallback
                auto tokenized_prompts = tokenize_input_prompts(ctx_server.vocab, prompt, true, true);
                for (auto & p : tokenized_prompts) {
                    auto tmp = server_tokens(p, ctx_server.mctx != nullptr);
                    inputs.push_back(std::move(tmp));
                }
            }'''

# Apply replacements
content = re.sub(pattern1, replacement1, content, flags=re.DOTALL)
content = re.sub(pattern2, replacement2, content, flags=re.DOTALL)

# Write back
with open('llama.cpp/tools/grpc-server/grpc-server.cpp', 'w') as f:
    f.write(content)

print("Fixed multimodal code blocks in grpc-server.cpp")
PYEOF

    # Run the Python script
    python3 fix_grpc_server.py
    rm fix_grpc_server.py
    
    echo "Disabled multimodal code in grpc-server.cpp"
fi

# Remove main function from server.cpp to avoid conflict with grpc-server.cpp
if [ -f llama.cpp/tools/grpc-server/server.cpp ]; then
    echo "Removing main function from server.cpp..."
    # Find the main function and remove it including its closing brace
    sed -i '/^int main(int argc, char \*\* argv) {$/,/^}$/d' llama.cpp/tools/grpc-server/server.cpp
    echo "Removed main function from server.cpp"
fi

# Copy web assets if they exist
if [ -f llama.cpp/tools/server/public/index.html.gz ]; then
    cp -rfv llama.cpp/tools/server/public/index.html.gz llama.cpp/tools/grpc-server/
fi
if [ -f llama.cpp/tools/server/public/loading.html ]; then
    cp -rfv llama.cpp/tools/server/public/loading.html llama.cpp/tools/grpc-server/
fi

# Create the UTF8 range shim to fix linking issues
cat > llama.cpp/tools/grpc-server/utf8_range_shim.cpp << 'EOF'
// UTF8 Range compatibility shim
// This provides missing symbols for utf8_range compatibility

#include <stddef.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// Stub implementation for missing utf8_range_IsValid function
bool utf8_range_IsValid(const char* data, size_t length) {
    // Simple validation - assume valid for now
    // In a production environment, this should be a proper UTF-8 validation
    (void)data;   // Suppress unused parameter warning
    (void)length; // Suppress unused parameter warning
    return true;
}

#ifdef __cplusplus
}
#endif
EOF

# Generate protobuf files
echo "Generating protobuf files"
cd llama.cpp/tools/grpc-server

# Create backend.proto if it doesn't exist
if [ ! -f backend.proto ]; then
    cp ../../../../../../backend/backend.proto .
fi

# Generate protobuf files
if [ ! -f backend.pb.h ] || [ ! -f backend.grpc.pb.h ]; then
    protoc --cpp_out=. --grpc_out=. --plugin=protoc-gen-grpc=$(which grpc_cpp_plugin) backend.proto
fi

cd ../../..

# Add grpc-server to the tools CMakeLists.txt
set +e
if [ -f llama.cpp/tools/CMakeLists.txt ]; then
    if grep -q "add_subdirectory(grpc-server)" llama.cpp/tools/CMakeLists.txt; then
        echo "grpc-server already added"
    else
        echo "add_subdirectory(grpc-server)" >> llama.cpp/tools/CMakeLists.txt
        echo "Added grpc-server to tools CMakeLists.txt"
    fi
else
    echo "Warning: llama.cpp/tools/CMakeLists.txt not found, creating minimal version..."
    cat > llama.cpp/tools/CMakeLists.txt << 'EOF'
# dependencies

find_package(Threads REQUIRED)

# flags

llama_add_compile_flags()

# tools

add_subdirectory(grpc-server)
EOF
    echo "Created minimal tools CMakeLists.txt"
fi
set -e

echo "prepare.sh completed successfully"