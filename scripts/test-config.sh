#!/usr/bin/env bash
# Script de test pour vérifier la configuration avant le build
set -euo pipefail

echo "========================================="
echo "Test de configuration LocalAI Inpainting"
echo "========================================="
echo ""

# Vérifier que nous sommes sur la bonne branche
CURRENT_BRANCH=$(git branch --show-current)
echo "✓ Branche actuelle: ${CURRENT_BRANCH}"

if [ "${CURRENT_BRANCH}" != "local/inpainting-image" ]; then
  echo "⚠ ATTENTION: Vous n'êtes pas sur la branche 'local/inpainting-image'"
  echo "  Branche actuelle: ${CURRENT_BRANCH}"
fi

# Vérifier que les fichiers d'inpainting existent
echo ""
echo "Vérification des fichiers d'inpainting..."
FILES_TO_CHECK=(
  "core/http/endpoints/openai/inpainting.go"
  "core/http/endpoints/openai/inpainting_test.go"
  "docker/Dockerfile.inpainting"
  "scripts/build-and-push.sh"
  "scripts/registry-config.sh"
)

ALL_OK=true
for file in "${FILES_TO_CHECK[@]}"; do
  if [ -f "$file" ]; then
    echo "  ✓ $file"
  else
    echo "  ✗ MANQUANT: $file"
    ALL_OK=false
  fi
done

# Vérifier que l'endpoint est enregistré
echo ""
echo "Vérification de l'endpoint /v1/images/inpainting..."
if grep -q "/v1/images/inpainting" core/http/routes/openai.go; then
  echo "  ✓ Endpoint enregistré dans les routes"
else
  echo "  ✗ Endpoint NON trouvé dans les routes"
  ALL_OK=false
fi

# Vérifier la configuration du registry
echo ""
echo "Configuration du registry:"
source scripts/registry-config.sh 2>/dev/null || true
echo "  Registry: ${REGISTRY:-non défini}"
echo "  Image:    ${IMAGE_NAME:-non défini}:${IMAGE_TAG:-non défini}"

if [ -n "${REGISTRY_USER:-}" ]; then
  echo "  User:     ${REGISTRY_USER} (défini)"
else
  echo "  User:     (non défini - build local uniquement)"
fi

if [ -n "${REGISTRY_TOKEN:-}" ]; then
  echo "  Token:    ******* (défini)"
else
  echo "  Token:    (non défini - build local uniquement)"
fi

# Vérifier Docker
echo ""
echo "Vérification de Docker..."
if command -v docker &> /dev/null; then
  DOCKER_VERSION=$(docker --version)
  echo "  ✓ Docker installé: ${DOCKER_VERSION}"
  
  if docker info &> /dev/null; then
    echo "  ✓ Docker daemon accessible"
  else
    echo "  ✗ Docker daemon non accessible"
    ALL_OK=false
  fi
else
  echo "  ✗ Docker non installé"
  ALL_OK=false
fi

# Résumé
echo ""
echo "========================================="
if [ "$ALL_OK" = true ]; then
  echo "✓ Tous les tests sont OK!"
  echo ""
  echo "Pour builder et pousser l'image:"
  echo "  1. source scripts/registry-config.sh"
  echo "  2. export REGISTRY_USER='votre-username'"
  echo "  3. export REGISTRY_TOKEN='votre-token'"
  echo "  4. ./scripts/build-and-push.sh"
  echo ""
  echo "Ou utilisez le fichier .env.registry:"
  echo "  1. cp .env.registry.example .env.registry"
  echo "  2. # Éditez .env.registry avec vos credentials"
  echo "  3. source .env.registry"
  echo "  4. source scripts/registry-config.sh"
  echo "  5. ./scripts/build-and-push.sh"
else
  echo "✗ Certains tests ont échoué"
  echo "Veuillez corriger les erreurs ci-dessus"
  exit 1
fi
echo "========================================="
