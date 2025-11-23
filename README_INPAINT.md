# LocalAI inpainting image (fork)

This folder contains helper artifacts to build a Docker image of LocalAI
that includes the inpainting changes from the branch `fix/inpainting-backend-response-and-cuda`.

Files added:
- `docker/Dockerfile.inpainting` — Dockerfile that clones this repo/branch and builds the binary.
- `docker-compose.override.yml` — convenience override to run the custom image with `docker compose`.
- `scripts/build-and-push.sh` — helper to build and push the image to GHCR.

Quickstart (local-only):

```bash
# from /fork/LocalAI
docker build -f docker/Dockerfile.inpainting -t localai:inpainting --build-arg GIT_REPO=https://github.com/gmaOCR/LocalAI.git --build-arg GIT_REF=fix/inpainting-backend-response-and-cuda .
docker compose -f docker-compose.yml -f docker-compose.override.yml up
```

To publish to GitHub Container Registry:

```bash
# set GHCR_USER and GHCR_TOKEN env vars
GHCR_USER=your-user GHCR_TOKEN=xxx ./scripts/build-and-push.sh
```

Notes:
- Adjust `./cmd/local-ai` path in the Dockerfile if your build target differs.
- If you want the official AIO stack to use this image, replace the service image reference in
  your deployment compose or use this override in your deployment.
