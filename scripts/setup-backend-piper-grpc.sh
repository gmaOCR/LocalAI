# Force le chemin d'installation vcpkg pour éviter toute ambiguïté
export VCPKG_INSTALLED="/home/greg/installation/vcpkg/installed/x64-linux"

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
PROTOC_BIN="$VCPKG_INSTALLED/tools/protobuf/protoc"
PROTO_FILE="$PROJECT_ROOT/backend/backend.proto"
PROTOBUF_HEADER="$VCPKG_INSTALLED/include/google/protobuf/stubs/common.h"
if [ -x "$PROTOC_BIN" ]; then
  echo "[INFO] Version de protoc utilisée :"
  "$PROTOC_BIN" --version
else
  echo "[ERREUR] Le binaire protoc n'est pas exécutable : $PROTOC_BIN"
  exit 1
fi
if [ -f "$PROTOBUF_HEADER" ]; then
  echo "[INFO] Version des headers protobuf (common.h) :"
  grep 'define PROTOBUF_VERSION ' "$PROTOBUF_HEADER" || grep 'define GOOGLE_PROTOBUF_VERSION ' "$PROTOBUF_HEADER" || echo "[WARN] Impossible de trouver la version dans common.h"
else
  echo "[ERREUR] Header protobuf introuvable : $PROTOBUF_HEADER"
  exit 1
fi

# 1. Régénérer les fichiers protobuf avec le protoc de vcpkg
if [ -x "$PROTOC_BIN" ] && [ -f "$PROTO_FILE" ]; then
  echo "[INFO] Suppression renforcée des anciens fichiers protobuf générés..."
  for d in backend/cpp/llama backend/cpp/llama-avx backend/cpp/llama-avx2 backend/cpp/llama-avx512 backend/cpp/llama-fallback backend/cpp/llama-grpc; do
    GEN_DIR="$PROJECT_ROOT/$d/llama.cpp/build"
    # Supprime tous les anciens fichiers générés (tous les .pb.* et .grpc.pb.*)
    find "$GEN_DIR" -type f \( -name 'backend.pb.*' -o -name 'backend.grpc.pb.*' \) -exec rm -f {} +
  done
  echo "[INFO] Régénération des fichiers protobuf avec $PROTOC_BIN..."
  for d in backend/cpp/llama backend/cpp/llama-avx backend/cpp/llama-avx2 backend/cpp/llama-avx512 backend/cpp/llama-fallback backend/cpp/llama-grpc; do
    GEN_DIR="$PROJECT_ROOT/$d/llama.cpp/build"
    mkdir -p "$GEN_DIR"
    GRPC_PLUGIN="$VCPKG_INSTALLED/tools/grpc/grpc_cpp_plugin"
    if [ ! -x "$GRPC_PLUGIN" ]; then
      GRPC_PLUGIN=$(which grpc_cpp_plugin || true)
    fi
    "$PROTOC_BIN" \
      -I "$PROJECT_ROOT/backend" \
      -I "$VCPKG_INSTALLED/include" \
      -I "$PROJECT_ROOT" \
      --cpp_out="$GEN_DIR" \
      --grpc_out="$GEN_DIR" \
      --plugin=protoc-gen-grpc="$GRPC_PLUGIN" \
      "$PROTO_FILE"
  done
fi
#!/bin/bash
# Script pour restaurer et recompiler le backend Piper gRPC pour LocalAI
# Usage : ./scripts/setup-backend-piper-grpc.sh [-j N]

set -euo pipefail



# --- Utilisation de vcpkg pour protobuf et abseil ---
VCPKG_ROOT="/home/greg/installation/vcpkg"
VCPKG_BIN="$VCPKG_ROOT/vcpkg"
VCPKG_INSTALLED="$VCPKG_ROOT/installed/x64-linux"

echo "[INFO] Installation de protobuf et abseil via vcpkg si nécessaire..."
$VCPKG_BIN install protobuf abseil || true

export PATH="$VCPKG_INSTALLED/bin:$PATH"
export LD_LIBRARY_PATH="$VCPKG_INSTALLED/lib:$LD_LIBRARY_PATH"
export PKG_CONFIG_PATH="$VCPKG_INSTALLED/lib/pkgconfig:${PKG_CONFIG_PATH:-}"
export CMAKE_PREFIX_PATH="$VCPKG_INSTALLED:$VCPKG_ROOT:$CMAKE_PREFIX_PATH"
export CPATH="$VCPKG_INSTALLED/include:${CPATH:-}"

# Forcer l'utilisation des includes/libs vcpkg pour protobuf et abseil
export CXXFLAGS="-I$VCPKG_INSTALLED/include ${CXXFLAGS:-}"
export LDFLAGS="-L$VCPKG_INSTALLED/lib ${LDFLAGS:-}"
export LIBRARY_PATH="$VCPKG_INSTALLED/lib:${LIBRARY_PATH:-}"
export LD_LIBRARY_PATH="$VCPKG_INSTALLED/lib:$LD_LIBRARY_PATH"
echo "[INFO] Utilisation de protobuf et abseil via vcpkg ($VCPKG_INSTALLED)"

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


# 4. Régénérer les fichiers protobuf AVANT compilation (sécurité)
echo "[INFO] Régénération forcée des fichiers protobuf avec le protoc de vcpkg avant compilation..."
for d in backend/cpp/llama backend/cpp/llama-avx backend/cpp/llama-avx2 backend/cpp/llama-avx512 backend/cpp/llama-fallback backend/cpp/llama-grpc; do
  GEN_DIR="$PROJECT_ROOT/$d/llama.cpp/build"
  mkdir -p "$GEN_DIR"
  GRPC_PLUGIN="$VCPKG_INSTALLED/tools/grpc/grpc_cpp_plugin"
  if [ ! -x "$GRPC_PLUGIN" ]; then
    GRPC_PLUGIN=$(which grpc_cpp_plugin || true)
  fi
  "$PROTOC_BIN" \
    -I "$PROJECT_ROOT/backend" \
    -I "$VCPKG_INSTALLED/include" \
    -I "$PROJECT_ROOT" \
    --cpp_out="$GEN_DIR" \
    --grpc_out="$GEN_DIR" \
    --plugin=protoc-gen-grpc="$GRPC_PLUGIN" \
    "$PROTO_FILE"
done

# 5. Recompiler le backend Piper gRPC (et tous les backends nécessaires)
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
