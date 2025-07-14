#!/bin/bash
# Script de compilation du binaire principal LocalAI avec support CUDA/cuBLAS
set -e

export BUILD_TYPE=cublas
export CGO_LDFLAGS="-lcublas -lcudart -L/usr/local/cuda/lib64/ -L/usr/local/cuda/lib64/stubs/ -lcuda"

echo "[build-cuda.sh] Compilation du binaire local-ai avec support CUDA/cuBLAS..."
go build -o local-ai ./
echo "[build-cuda.sh] Terminé. Le binaire local-ai est prêt."
