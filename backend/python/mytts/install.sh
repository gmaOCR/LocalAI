#!/bin/bash
set -e

backend_dir=$(dirname $0)
if [ -d $backend_dir/common ]; then
    source $backend_dir/common/libbackend.sh
else
    source $backend_dir/../common/libbackend.sh
fi

installRequirements

echo "Installation terminée!"
echo ""
echo "Usage:"
echo "  python main.py --addr localhost:50051"
echo ""
echo "Configuration d'exemple pour LocalAI:"
echo "backend: mytts"
echo "name: mon-tts"
echo "parameters:"
echo "  model: mytts-model"
echo "tts:"
echo "  voice: fr"
