#!/usr/bin/env bash
# Configuration pour le build et push vers mercure.gregorymariani.com
# Source ce fichier avant d'exécuter build-and-push.sh
#
# Usage:
#   source scripts/registry-config.sh
#   ./scripts/build-and-push.sh

# Registry configuration
export REGISTRY="mercure.gregorymariani.com"
export IMAGE_NAME="localai"
export IMAGE_TAG="inpainting-latest"

# Git configuration (déjà configuré par défaut dans build-and-push.sh)
export GIT_REPO="https://github.com/gmaOCR/LocalAI.git"
export GIT_REF="local/inpainting-image"

# Build type (laissez vide pour CPU, ou utilisez: cublas, hipblas, etc.)
export BUILD_TYPE=""

# Registry credentials (à définir avant le push)
# export REGISTRY_USER="votre-username"
# export REGISTRY_TOKEN="votre-token"

echo "Configuration chargée pour ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
echo ""
echo "Pour pousser l'image, définissez:"
echo "  export REGISTRY_USER='votre-username'"
echo "  export REGISTRY_TOKEN='votre-token'"
echo ""
echo "Puis exécutez: ./scripts/build-and-push.sh"
