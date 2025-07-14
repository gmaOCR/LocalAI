#!/bin/bash
# Script pour restaurer et recompiler le backend Piper gRPC pour LocalAI
# Usage : ./scripts/setup-backend-piper-grpc.sh [-j N]

set -euo pipefail

JOBS=16
if [[ $# -ge 2 && $1 == "-j" ]]; then
  JOBS=$2
fi

# 1. Initialiser les sous-modules (llama.cpp, etc.)
echo "[INFO] Initialisation des sous-modules git..."
git submodule update --init --recursive

# 2. Nettoyer les builds précédents
echo "[INFO] Nettoyage des builds précédents..."
make clean || true
find . -type d -name build -exec rm -rf {} +

# 3. Recréer le dossier backend-assets si besoin
mkdir -p backend-assets/grpc

# 4. Recompiler le backend Piper gRPC (et tous les backends nécessaires)
echo "[INFO] Compilation de LocalAI et des backends (Piper gRPC inclus)..."
make -j $JOBS build

# 5. Vérifier la présence du binaire Piper gRPC
echo "[INFO] Vérification du binaire Piper gRPC..."
if [ -f backend-assets/grpc/piper ] || [ -f backend-assets/grpc/piper-bin ]; then
  echo "[OK] Backend Piper gRPC compilé avec succès."
else
  echo "[ERREUR] Le binaire Piper gRPC est manquant !"
  exit 1
fi

echo "[INFO] LocalAI et le backend Piper gRPC sont prêts."
