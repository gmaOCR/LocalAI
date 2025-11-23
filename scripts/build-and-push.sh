#!/usr/bin/env bash
set -euo pipefail

# Build and optionally push the inpainting image for this fork/branch
# Usage:
#   GIT_REPO=... GIT_REF=... GHCR_USER=... GHCR_TOKEN=... ./scripts/build-and-push.sh

IMAGE=${IMAGE:-ghcr.io/${GHCR_USER:-$(git config user.name)} /localai:inpainting-${GIT_REF:-latest}}
GIT_REPO=${GIT_REPO:-https://github.com/gmaOCR/LocalAI.git}
GIT_REF=${GIT_REF:-fix/inpainting-backend-response-and-cuda}

echo "Building image from ${GIT_REPO}@${GIT_REF}"
docker build -f docker/Dockerfile.inpainting -t "${IMAGE}" --build-arg GIT_REPO="${GIT_REPO}" --build-arg GIT_REF="${GIT_REF}" .

if [ -n "${GHCR_USER:-}" ] && [ -n "${GHCR_TOKEN:-}" ]; then
  echo "Logging into ghcr.io as ${GHCR_USER}"
  echo "${GHCR_TOKEN}" | docker login ghcr.io -u "${GHCR_USER}" --password-stdin
  docker push "${IMAGE}"
  echo "Pushed ${IMAGE}"
else
  echo "GHCR_USER/GHCR_TOKEN not set — image is only built locally: ${IMAGE}"
fi
