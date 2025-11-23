#!/usr/bin/env bash
set -euo pipefail

# Build and push LocalAI with inpainting support to private registry
# Usage:
#   # Push to private registry (mercure.gregorymariani.com)
#   ./scripts/build-and-push.sh
#
#   # Push to GitHub Container Registry
#   REGISTRY=ghcr.io REGISTRY_USER=username REGISTRY_TOKEN=xxx ./scripts/build-and-push.sh
#
#   # Custom build
#   IMAGE_TAG=v1.0.0 BUILD_TYPE=cublas ./scripts/build-and-push.sh

# Configuration
REGISTRY=${REGISTRY:-mercure.gregorymariani.com}
REGISTRY_USER=${REGISTRY_USER:-}
REGISTRY_TOKEN=${REGISTRY_TOKEN:-}
GIT_REPO=${GIT_REPO:-https://github.com/gmaOCR/LocalAI.git}
GIT_REF=${GIT_REF:-local/inpainting-image}
IMAGE_NAME=${IMAGE_NAME:-localai}
IMAGE_TAG=${IMAGE_TAG:-inpainting-latest}
BUILD_TYPE=${BUILD_TYPE:-}  # empty, cublas, hipblas, etc.

# Construct full image name
if [ -n "${BUILD_TYPE}" ]; then
  FULL_IMAGE="${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}-${BUILD_TYPE}"
else
  FULL_IMAGE="${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
fi

echo "========================================="
echo "LocalAI Inpainting Build & Push"
echo "========================================="
echo "Registry:   ${REGISTRY}"
echo "Image:      ${FULL_IMAGE}"
echo "Git Repo:   ${GIT_REPO}"
echo "Git Ref:    ${GIT_REF}"
echo "Build Type: ${BUILD_TYPE:-default}"
echo "========================================="

# Build the image
echo ""
echo "Building image..."
if [ -n "${BUILD_TYPE}" ]; then
  docker build \
    -f docker/Dockerfile.inpainting \
    -t "${FULL_IMAGE}" \
    --build-arg GIT_REPO="${GIT_REPO}" \
    --build-arg GIT_REF="${GIT_REF}" \
    --build-arg BUILD_TYPE="${BUILD_TYPE}" \
    .
else
  docker build \
    -f docker/Dockerfile.inpainting \
    -t "${FULL_IMAGE}" \
    --build-arg GIT_REPO="${GIT_REPO}" \
    --build-arg GIT_REF="${GIT_REF}" \
    .
fi

echo ""
echo "✓ Build completed: ${FULL_IMAGE}"

# Push to registry if credentials are provided
if [ -n "${REGISTRY_USER}" ] && [ -n "${REGISTRY_TOKEN}" ]; then
  echo ""
  echo "Logging into ${REGISTRY} as ${REGISTRY_USER}..."
  echo "${REGISTRY_TOKEN}" | docker login "${REGISTRY}" -u "${REGISTRY_USER}" --password-stdin
  
  echo ""
  echo "Pushing ${FULL_IMAGE}..."
  docker push "${FULL_IMAGE}"
  
  echo ""
  echo "========================================="
  echo "✓ Successfully pushed: ${FULL_IMAGE}"
  echo "========================================="
else
  echo ""
  echo "========================================="
  echo "⚠ Registry credentials not provided"
  echo "Image built locally only: ${FULL_IMAGE}"
  echo ""
  echo "To push to registry, set:"
  echo "  REGISTRY_USER=<username>"
  echo "  REGISTRY_TOKEN=<token>"
  echo "========================================="
fi
