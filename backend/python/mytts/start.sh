#!/bin/bash
# Script pour démarrer le backend MyTTS manuellement

echo "=== Démarrage du backend MyTTS ==="

# Aller dans le répertoire du backend
cd "$(dirname "$0")"

# Vérifier les dépendances
echo "Vérification des dépendances..."
python3 -c "import pyttsx3, gtts, grpc" 2>/dev/null
if [ $? -ne 0 ]; then
    echo "⚠️  Certaines dépendances manquent. Installation..."
    pip install -r requirements.txt
fi

# Démarrer le backend
echo "🚀 Démarrage du backend MyTTS sur localhost:50051"
echo "Pour arrêter: Ctrl+C"
echo

python3 main.py --addr localhost:50051
