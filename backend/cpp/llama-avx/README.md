# Backend llama.cpp avec support AVX

Ce répertoire contient la configuration pour compiler le backend llama.cpp avec les optimisations AVX.

## 📋 Aperçu

Ce backend utilise une version modifiée de llama.cpp optimisée pour les processeurs supportant les instructions AVX. Il génère un serveur gRPC (`grpc-server`) qui permet à LocalAI de communiquer avec les modèles de langage.

## 🔧 Fichiers principaux

### Scripts de build
- `prepare.sh` - Script de préparation de l'environnement de build
- `Makefile` - Configuration de compilation principale

### Sources
- `grpc-server.cpp` - Code source du serveur gRPC principal
- `mtmd-stub.h` / `mtmd-stub.cpp` - Stubs pour le support multimodal (désactivé)
- `grpc-server-CMakeLists.txt` - Configuration CMake pour le serveur gRPC

### Configuration
- `CMakeLists.txt` - Configuration CMake globale
- `patches/` - Patches spécifiques à appliquer

## 🚀 Utilisation

### Build automatique (recommandé)
Le build est automatiquement déclenché par le Makefile principal de LocalAI :
```bash
cd /fork/LocalAI
make build GO_TAGS="tts"
```

### Build manuel
Pour compiler manuellement ce backend :
```bash
cd /fork/LocalAI/backend/cpp/llama-avx

# Nettoyer l'environnement
make purge

# Préparer l'environnement
bash prepare.sh

# Compiler le serveur gRPC
make grpc-server
```

## 🔍 Structure après build

```
llama-avx/
├── llama.cpp/                    # Sources llama.cpp clonées
│   ├── build/                    # Artifacts de compilation
│   │   └── bin/grpc-server       # Exécutable final
│   └── tools/grpc-server/        # Sources préparées pour gRPC
├── grpc-server                   # Exécutable copié (résultat final)
└── [fichiers de configuration]
```

## ⚙️ Variables d'environnement

### Variables CMake utilisées
- `BUILD_SHARED_LIBS=OFF` - Build statique
- `GGML_NATIVE=OFF` - Désactive les optimisations natives
- `GGML_AVX=on` - Active le support AVX
- `GGML_AVX2=off` - Désactive AVX2
- `GGML_AVX512=off` - Désactive AVX-512
- `GGML_FMA=off` - Désactive FMA
- `GGML_F16C=off` - Désactive F16C
- `LLAMA_CURL=OFF` - Désactive libcurl

### Variables Makefile
- `LLAMA_VERSION` - Version/commit de llama.cpp à utiliser
- `CMAKE_ARGS` - Arguments supplémentaires pour CMake
- `BUILD_TYPE` - Type de build (vide pour AVX)

## 🔧 Processus de préparation

Le script `prepare.sh` effectue les étapes suivantes :

1. **Initialisation** : Copie des sources llama.cpp depuis le répertoire `llama/`
2. **Préparation des sources** : 
   - Copie des fichiers nécessaires vers `llama.cpp/tools/grpc-server/`
   - Application des corrections d'headers
   - Copie des dépendances vendor (minja)
3. **Corrections automatiques** :
   - Désactivation du code multimodal
   - Suppression de la fonction main conflictuelle
   - Correction des chemins d'inclusion
4. **Génération des assets** :
   - Création des fichiers protobuf
   - Copie des assets web
   - Génération du shim UTF8

## 🐛 Problèmes connus et solutions

### 1. Erreurs de compilation avec multimodal
**Symptôme** : Erreurs liées aux types `mtmd_*`
**Solution** : Le code multimodal est automatiquement désactivé par `prepare.sh`

### 2. Dépendances manquantes
**Symptôme** : `undefined reference` lors du linkage
**Solution** : Tous les fichiers sources nécessaires sont automatiquement inclus

### 3. Headers introuvables
**Symptôme** : `No such file or directory` pour minja ou autres
**Solution** : Les dépendances vendor sont copiées automatiquement

## 📊 Performance

Cette configuration est optimisée pour :
- Processeurs x86-64 avec support AVX
- Performance CPU sans accélération GPU
- Déploiements sur serveurs standards

**Benchmark typique** :
- Tokens/sec : Variable selon le modèle et le CPU
- Mémoire : ~2-8GB selon la taille du modèle
- CPU : Utilisation optimisée des instructions AVX

## 🔄 Maintenance

### Mise à jour de llama.cpp
1. Modifier `LLAMA_VERSION` dans le Makefile principal
2. Tester la compilation : `make purge && make grpc-server`
3. Vérifier la compatibilité des stubs MTMD
4. Mettre à jour cette documentation si nécessaire

### Ajout de nouvelles dépendances
1. Modifier `grpc-server-CMakeLists.txt`
2. Mettre à jour `prepare.sh` si des fichiers doivent être copiés
3. Tester la compilation complète

### Debug
Pour diagnostiquer les problèmes :
```bash
# Build verbose
make grpc-server VERBOSE=1

# Logs détaillés
bash -x prepare.sh

# Vérification des dépendances
ldd grpc-server
```

## 📝 Changelog

### Version actuelle (Juillet 2025)
- ✅ Résolution des conflits de définition ggml_log_level
- ✅ Ajout de toutes les dépendances manquantes  
- ✅ Support automatique des assets minja
- ✅ Désactivation propre du code multimodal
- ✅ Build reproductible et maintenable

---

> **Note** : Ce README doit être synchronisé avec `/fork/LocalAI/backend/cpp/BUILD_FIXES.md` lors de modifications importantes.
