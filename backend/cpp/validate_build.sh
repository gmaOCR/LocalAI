#!/bin/bash

# Script de validation du build LocalAI backend C++
# Usage: ./validate_build.sh [--clean]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOCALAI_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Parse arguments
CLEAN_BUILD=false
if [[ "$1" == "--clean" ]]; then
    CLEAN_BUILD=true
fi

log_info "Validation du build LocalAI backend C++"
log_info "Répertoire racine: $LOCALAI_ROOT"

# Change to LocalAI root
cd "$LOCALAI_ROOT"

# 1. Vérifier la structure des fichiers
log_info "1. Vérification de la structure des fichiers..."

required_files=(
    "backend/cpp/llama-avx/prepare.sh"
    "backend/cpp/llama-avx/grpc-server-CMakeLists.txt"
    "backend/cpp/llama-avx/grpc-server.cpp"
    "backend/cpp/llama-avx/mtmd-stub.h"
    "backend/cpp/llama-avx/mtmd-stub.cpp"
    "backend/cpp/llama-avx/Makefile"
    "backend/cpp/BUILD_FIXES.md"
    "backend/cpp/llama-avx/README.md"
)

missing_files=()
for file in "${required_files[@]}"; do
    if [[ ! -f "$file" ]]; then
        missing_files+=("$file")
    fi
done

if [[ ${#missing_files[@]} -ne 0 ]]; then
    log_error "Fichiers manquants:"
    for file in "${missing_files[@]}"; do
        echo "  - $file"
    done
    exit 1
else
    log_success "Tous les fichiers requis sont présents"
fi

# 2. Vérifier les permissions des scripts
log_info "2. Vérification des permissions des scripts..."

scripts=(
    "backend/cpp/llama-avx/prepare.sh"
)

for script in "${scripts[@]}"; do
    if [[ ! -x "$script" ]]; then
        log_warning "Correction des permissions pour $script"
        chmod +x "$script"
    fi
done

log_success "Permissions des scripts OK"

# 3. Vérifier les dépendances système
log_info "3. Vérification des dépendances système..."

dependencies=(
    "cmake:cmake"
    "make:make"
    "gcc:gcc"
    "g++:g++"
    "pkg-config:pkg-config"
    "protoc:protobuf-compiler"
    "grpc_cpp_plugin:grpc"
)

missing_deps=()
for dep in "${dependencies[@]}"; do
    cmd="${dep%%:*}"
    package="${dep##*:}"
    if ! command -v "$cmd" &> /dev/null; then
        missing_deps+=("$package")
    fi
done

if [[ ${#missing_deps[@]} -ne 0 ]]; then
    log_warning "Dépendances manquantes détectées:"
    for dep in "${missing_deps[@]}"; do
        echo "  - $dep"
    done
    log_info "Pour installer sur Ubuntu/Debian:"
    echo "  sudo apt update"
    echo "  sudo apt install -y cmake make gcc g++ pkg-config protobuf-compiler libgrpc++-dev"
else
    log_success "Toutes les dépendances système sont présentes"
fi

# 4. Test de build si demandé
if [[ "$CLEAN_BUILD" == true ]]; then
    log_info "4. Test de build complet (nettoyage + compilation)..."
    
    # Nettoyer
    log_info "Nettoyage..."
    if ! make clean 2>/dev/null; then
        log_warning "Erreur lors du nettoyage (peut être ignorée)"
    fi
    
    # Build
    log_info "Compilation avec GO_TAGS=\"tts\"..."
    if make build GO_TAGS="tts"; then
        log_success "Build complet réussi!"
    else
        log_error "Échec du build complet"
        exit 1
    fi
else
    log_info "4. Validation de la configuration (pas de build complet)"
    
    # Vérifier juste la configuration llama-avx
    cd "backend/cpp/llama-avx"
    
    if [[ ! -f "llama.cpp/CMakeLists.txt" ]]; then
        log_info "Initialisation de llama.cpp..."
        if ! bash prepare.sh; then
            log_error "Échec de l'initialisation"
            exit 1
        fi
    fi
    
    log_info "Configuration CMake..."
    cd llama.cpp
    mkdir -p build
    cd build
    
    if cmake .. -DBUILD_SHARED_LIBS=OFF -DGGML_NATIVE=OFF -DGGML_AVX=on -DGGML_AVX2=off -DGGML_AVX512=off -DGGML_FMA=off -DGGML_F16C=off -DLLAMA_CURL=OFF &>/dev/null; then
        log_success "Configuration CMake réussie"
    else
        log_error "Échec de la configuration CMake"
        exit 1
    fi
    
    cd "$LOCALAI_ROOT"
fi

# 5. Vérifier les exécutables générés (si build complet)
if [[ "$CLEAN_BUILD" == true ]]; then
    log_info "5. Vérification des exécutables générés..."
    
    expected_executables=(
        "backend-assets/grpc/llama-cpp-avx"
    )
    
    missing_executables=()
    for exe in "${expected_executables[@]}"; do
        if [[ ! -f "$exe" ]]; then
            missing_executables+=("$exe")
        elif [[ ! -x "$exe" ]]; then
            log_warning "$exe n'est pas exécutable"
        else
            log_success "$exe: $(file "$exe" | cut -d: -f2-)"
        fi
    done
    
    if [[ ${#missing_executables[@]} -ne 0 ]]; then
        log_error "Exécutables manquants:"
        for exe in "${missing_executables[@]}"; do
            echo "  - $exe"
        done
        exit 1
    else
        log_success "Tous les exécutables attendus sont présents"
    fi
    
    # Test rapide des exécutables
    log_info "Test rapide des exécutables..."
    for exe in "${expected_executables[@]}"; do
        if timeout 5 "$exe" --help &>/dev/null; then
            log_success "$exe: répond correctement"
        else
            log_warning "$exe: ne répond pas à --help (peut être normal)"
        fi
    done
fi

# 6. Résumé final
log_info "6. Résumé de la validation"

echo
log_success "✅ Validation terminée avec succès!"
echo
echo "📋 État du projet:"
echo "  - Structure des fichiers: ✅ OK"
echo "  - Permissions: ✅ OK"
echo "  - Dépendances système: $(if [[ ${#missing_deps[@]} -eq 0 ]]; then echo "✅ OK"; else echo "⚠️  Manquantes"; fi)"
if [[ "$CLEAN_BUILD" == true ]]; then
    echo "  - Build complet: ✅ OK"
    echo "  - Exécutables: ✅ OK"
else
    echo "  - Configuration CMake: ✅ OK"
fi

echo
echo "🚀 Pour utiliser le projet:"
echo "  make build GO_TAGS=\"tts\"              # Build complet"
echo "  ./validate_build.sh --clean            # Validation avec build"
echo
echo "📖 Documentation:"
echo "  backend/cpp/BUILD_FIXES.md             # Guide de dépannage"
echo "  backend/cpp/llama-avx/README.md        # Documentation backend"

exit 0
