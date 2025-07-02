#!/bin/bash
set -e

# Get the directory of the script
backend_dir=$(dirname "$0")

# Source the common backend library
if [ -f "$backend_dir/../common/libbackend.sh" ]; then
    source "$backend_dir/../common/libbackend.sh"
else
    echo "Error: libbackend.sh not found at $backend_dir/../common/libbackend.sh"
    exit 1
fi

# Initialize the backend
init

(
    cd "$backend_dir" || exit
    # Ensure the virtual environment is set up and activated
    ensureVenv
    # Run the tests
    runUnittests "$@"
)
