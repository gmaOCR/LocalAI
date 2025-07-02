#!/bin/bash
# Script de test rapide pour le backend MyTTS

echo "=== Test du backend MyTTS ==="
echo

# Vérifier que LocalAI est en cours d'exécution
if ! curl -s http://localhost:8080/v1/models >/dev/null 2>&1; then
    echo "❌ LocalAI ne semble pas être en cours d'exécution sur le port 8080"
    echo "Démarrez LocalAI avec: make run"
    exit 1
fi

echo "✓ LocalAI est en cours d'exécution"

# Test de base du TTS
echo "🎤 Test de génération TTS..."
curl -s http://localhost:8080/tts \
    -H "Content-Type: application/json" \
    -d '{
        "backend": "mytts",
        "model": "mytts-francais", 
        "input": "Bonjour, ceci est un test de synthèse vocale avec le backend MyTTS pour LocalAI."
    }' \
    -o test_output.wav

if [ $? -eq 0 ] && [ -f test_output.wav ] && [ -s test_output.wav ]; then
    echo "✓ Fichier audio généré avec succès: test_output.wav"
    
    # Essayer de jouer le fichier si possible
    if command -v aplay >/dev/null 2>&1; then
        echo "🔊 Lecture du fichier audio..."
        aplay test_output.wav 2>/dev/null
    elif command -v afplay >/dev/null 2>&1; then
        echo "🔊 Lecture du fichier audio..."
        afplay test_output.wav
    else
        echo "ℹ️  Pour écouter le fichier: aplay test_output.wav (ou votre lecteur audio)"
    fi
else
    echo "❌ Échec de la génération audio"
    echo "Vérifiez que le backend mytts est correctement configuré"
fi

echo
echo "=== Test terminé ==="
