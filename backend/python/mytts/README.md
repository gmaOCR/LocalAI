# Backend TTS Personnalisé pour LocalAI

Ce backend fournit une implémentation Text-to-Speech personnalisée pour LocalAI.

## Fonctionnalités

- Support de plusieurs moteurs TTS (pyttsx3, gTTS, espeak)
- Configuration des voix et langues
- Compatible avec l'API LocalAI TTS
- Tests unitaires inclus

## Installation

```bash
make mytts
```

Ou manuellement :

```bash
bash install.sh
```

## Utilisation

### Démarrage du backend

```bash
make run
```

Ou manuellement :

```bash
bash run.sh
```

Le backend démarrera sur `localhost:50051` par défaut.

### Configuration avec LocalAI

1. Créer un fichier de configuration YAML dans le dossier `models/` :

```yaml
name: mytts-model
backend: mytts
parameters:
  model: mytts-custom
tts:
  voice: "default"
  language: "fr"
```

2. Démarrer LocalAI avec le backend :

```bash
./local-ai --models-path ./models
```

### Test via API

```bash
curl http://localhost:8080/tts -H "Content-Type: application/json" -d '{
  "model": "mytts-model",
  "input": "Bonjour, ceci est un test de synthèse vocale"
}' --output test.wav
```

## Tests

```bash
make test
```

## Personnalisation

Vous pouvez personnaliser le backend en modifiant `backend.py` :

- Changer le moteur TTS (ligne 15-30)
- Ajouter de nouvelles voix
- Modifier les paramètres audio
- Ajouter le support de nouveaux formats

## Dépendances

- grpcio et grpcio-tools (requis)
- pyttsx3 (TTS hors ligne, recommandé)
- gtts (Google TTS, nécessite Internet)
- espeak (fallback système)

## Moteurs TTS supportés

1. **pyttsx3** : TTS hors ligne, multi-plateforme
2. **gTTS** : Google Text-to-Speech (nécessite Internet)
3. **espeak** : TTS système (fallback)
4. **Mock** : Génération de fichiers audio silencieux pour les tests

## Troubleshooting

### Problème avec pyttsx3
```bash
# Ubuntu/Debian
sudo apt-get install espeak espeak-data libespeak1 libespeak-dev

# macOS
brew install espeak

# Windows: Installé par défaut avec pyttsx3
```

### Problème avec les permissions
```bash
chmod +x *.sh
```
