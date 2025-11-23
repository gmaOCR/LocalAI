#!/usr/bin/env bash
# Exemple complet de build et push vers mercure.gregorymariani.com
# Ce script montre toutes les étapes

set -euo pipefail

echo "========================================="
echo "Exemple: Build & Push LocalAI Inpainting"
echo "========================================="
echo ""

# Étape 1: Vérifier la configuration
echo "Étape 1/4: Vérification de la configuration..."
./scripts/test-config.sh

# Étape 2: Charger la configuration du registry
echo ""
echo "Étape 2/4: Chargement de la configuration..."
source scripts/registry-config.sh

# Étape 3: Demander les credentials (ou les charger depuis .env.registry)
echo ""
echo "Étape 3/4: Configuration des credentials..."
if [ -f ".env.registry" ]; then
  echo "Chargement des credentials depuis .env.registry..."
  source .env.registry
else
  echo "Fichier .env.registry non trouvé."
  echo ""
  read -p "Voulez-vous entrer vos credentials maintenant? (y/n) " -n 1 -r
  echo
  if [[ $REPLY =~ ^[Yy]$ ]]; then
    read -p "Registry username: " REGISTRY_USER
    read -sp "Registry token/password: " REGISTRY_TOKEN
    echo
    export REGISTRY_USER
    export REGISTRY_TOKEN
  else
    echo "Build local uniquement (pas de push vers le registry)"
  fi
fi

# Étape 4: Build et push
echo ""
echo "Étape 4/4: Build et push de l'image..."
echo ""
./scripts/build-and-push.sh

echo ""
echo "========================================="
echo "✓ Terminé!"
echo "========================================="

if [ -n "${REGISTRY_USER:-}" ]; then
  echo ""
  echo "Image disponible sur:"
  echo "  mercure.gregorymariani.com/localai:inpainting-latest"
  echo ""
  echo "Pour l'utiliser:"
  echo "  docker pull mercure.gregorymariani.com/localai:inpainting-latest"
  echo "  docker run -d -p 8080:3007 mercure.gregorymariani.com/localai:inpainting-latest"
fi
