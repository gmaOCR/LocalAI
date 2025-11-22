#!/usr/bin/env bash
#
# Script pour tester le fix inpainting en local
# Usage: ./test-inpainting-fix.sh
#

set -e

echo "=================================================="
echo "Test du fix inpainting - LocalAI"
echo "=================================================="
echo ""

# Vérifier qu'on est dans le bon répertoire
if [ ! -f "Makefile" ] || [ ! -d "gallery" ]; then
    echo "❌ Erreur: Ce script doit être exécuté depuis la racine du projet LocalAI"
    exit 1
fi

# Vérifier que Go est installé
if ! command -v go &> /dev/null; then
    echo "❌ Erreur: Go n'est pas installé. Installez Go 1.21+ pour continuer."
    exit 1
fi

echo "✅ Pré-requis OK"
echo ""

# Étape 1: Compiler le backend stablediffusion-ggml
echo "📦 Étape 1/3 : Compilation du backend stablediffusion-ggml..."
echo "   (Cela peut prendre plusieurs minutes la première fois)"
if make backends/stablediffusion-ggml; then
    echo "✅ Backend compilé avec succès"
else
    echo "❌ Échec de la compilation du backend"
    exit 1
fi
echo ""

# Étape 2: Préparer l'environnement de test
echo "🔧 Étape 2/3 : Préparation de l'environnement de test..."
if make prepare-test; then
    echo "✅ Environnement préparé"
else
    echo "❌ Échec de la préparation"
    exit 1
fi
echo ""

# Étape 3: Lancer le test stablediffusion
echo "🧪 Étape 3/3 : Lancement du test stablediffusion..."
echo "   (Cela peut prendre 5-10 minutes : téléchargement du modèle + génération)"
echo ""

if make test-stablediffusion; then
    echo ""
    echo "=================================================="
    echo "✅ TEST RÉUSSI !"
    echo "=================================================="
    echo ""
    echo "Le fix fonctionne correctement. Vous pouvez maintenant:"
    echo "  1. Commiter vos changements: git add gallery/index.yaml"
    echo "  2. Pousser sur la branche: git push origin fix/inpainting-single"
    echo "  3. Vérifier que le CI passe sur GitHub"
    echo ""
else
    echo ""
    echo "=================================================="
    echo "❌ TEST ÉCHOUÉ"
    echo "=================================================="
    echo ""
    echo "Le test a échoué. Consultez les logs ci-dessus pour plus de détails."
    echo "Pour plus d'informations, consultez TEST_INPAINTING_LOCAL.md"
    echo ""
    exit 1
fi
