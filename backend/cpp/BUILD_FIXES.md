# Corrections de compilation pour LocalAI Backend C++

Ce document décrit les corrections apportées pour résoudre les problèmes de compilation du backend C++ de LocalAI, en particulier pour la génération des serveurs gRPC llama.cpp.

## 🎯 Objectif

Permettre la compilation réussie de LocalAI avec `make build GO_TAGS="tts"` en résolvant les conflits de définition et les dépendances manquantes dans le processus de build des backends llama.cpp.

## 🔍 Problèmes identifiés

### 1. Conflits de définition d'enums
- **Problème** : Redéfinition de l'enum `ggml_log_level` dans `mtmd-stub.h`
- **Symptôme** : Erreurs de compilation `conflicting declaration`
- **Solution** : Suppression des redéfinitions et utilisation des définitions existantes de llama.cpp

### 2. Dépendances manquantes
- **Problème** : Fonctions non résolues lors du linkage (`common_tokenize`, `common_log_*`, `common_json_parse`)
- **Symptôme** : Erreurs `undefined reference`
- **Solution** : Ajout des fichiers sources manquants au processus de build

### 3. Headers manquants
- **Problème** : Fichiers d'en-tête introuvables (`minja/chat-template.hpp`, `json-partial.h`)
- **Symptôme** : Erreurs `No such file or directory`
- **Solution** : Copie des dépendances vendor et ajout des chemins d'inclusion

## 📁 Structure des corrections

```
backend/cpp/
├── llama-avx/                    # Backend llama.cpp avec support AVX
│   ├── prepare.sh                # Script de préparation corrigé
│   ├── grpc-server-CMakeLists.txt # Configuration CMake pour grpc-server
│   ├── grpc-server.cpp           # Code source du serveur gRPC
│   ├── mtmd-stub.h               # Stubs pour multimodal (corrigé)
│   ├── mtmd-stub.cpp             # Implémentation des stubs (corrigé)
│   └── Makefile                  # Configuration de build
└── BUILD_FIXES.md                # Cette documentation
```

## 🔧 Détails des corrections

### 1. Correction du fichier `mtmd-stub.h`

**Avant :**
```cpp
typedef enum {
    GGML_LOG_LEVEL_DEBUG = 0,
    GGML_LOG_LEVEL_INFO = 1,
    GGML_LOG_LEVEL_WARN = 2,
    GGML_LOG_LEVEL_ERROR = 3
} ggml_log_level;

extern "C" {
    common_log* common_log_main();
    void common_log_add(common_log* log, ggml_log_level level, const char* format, ...);
    extern int common_log_verbosity_thold;
}
```

**Après :**
```cpp
// Les définitions ggml_log_level et common_log_* sont déjà fournies par llama.cpp
// Suppression des redéfinitions conflictuelles
```

### 2. Ajout des dépendances dans `grpc-server-CMakeLists.txt`

**Sources ajoutées :**
```cmake
set(TARGET_SRCS
    grpc-server.cpp
    backend.pb.cc
    backend.grpc.pb.cc
    utf8_range_shim.cpp
    mtmd-stub.cpp
    ../../common/common.cpp          # Fonctions utilitaires communes
    ../../common/log.cpp             # Système de logging
    ../../common/arg.cpp             # Parsing d'arguments
    ../../common/chat.cpp            # Gestion des chats
    ../../common/chat-parser.cpp     # Parseur de format chat
    ../../common/sampling.cpp        # Algorithmes de sampling
    ../../common/speculative.cpp     # Décodage spéculatif
    ../../common/json-schema-to-grammar.cpp  # Conversion JSON schema
    ../../common/regex-partial.cpp   # Expressions régulières partielles
    ../../common/json-partial.cpp    # Parsing JSON partiel
    ${CMAKE_CURRENT_BINARY_DIR}/../../common/build-info.cpp  # Infos de build
)
```

### 3. Mise à jour du script `prepare.sh`

**Nouvelles fonctionnalités :**
- Copie automatique des dépendances vendor (minja)
- Correction des chemins d'inclusion dans les fichiers copiés
- Désactivation du code multimodal pour éviter les dépendances complexes
- Génération automatique des fichiers protobuf
- Gestion des assets web

**Exemple de section clé :**
```bash
# Copy vendor dependencies
if [ -d llama.cpp/vendor/minja ]; then
    cp -rfv llama.cpp/vendor/minja llama.cpp/tools/grpc-server/
fi

# Fix grpc-server.cpp by disabling multimodal code
sed -i 's/const bool has_mtmd = ctx_server\.mctx != nullptr;/const bool has_mtmd = false; \/\/ Disabled multimodal support/' llama.cpp/tools/grpc-server/grpc-server.cpp
```

## 🚀 Utilisation

### Build standard
```bash
cd /fork/LocalAI
make build GO_TAGS="tts"
```

### Build avec nettoyage
```bash
cd /fork/LocalAI
make clean
make build GO_TAGS="tts"
```

### Build debug d'un backend spécifique
```bash
cd /fork/LocalAI/backend/cpp/llama-avx
make purge          # Nettoie les artifacts de build
bash prepare.sh     # Prépare l'environnement
make grpc-server    # Compile le serveur gRPC
```

## 🧪 Tests de validation

### 1. Vérification de la compilation
```bash
cd /fork/LocalAI
make build GO_TAGS="tts" 2>&1 | grep -E "(error|Error|erreur|Erreur)" || echo "Build successful"
```

### 2. Vérification des exécutables générés
```bash
ls -la /fork/LocalAI/backend-assets/grpc/llama-cpp-*
file /fork/LocalAI/backend-assets/grpc/llama-cpp-avx
```

### 3. Test de fonctionnalité
```bash
/fork/LocalAI/backend-assets/grpc/llama-cpp-avx --help
```

## 📋 Checklist de maintenance

Lors de futures modifications du code llama.cpp :

- [ ] Vérifier que les nouvelles dépendances sont ajoutées au `CMakeLists.txt`
- [ ] Mettre à jour le script `prepare.sh` si de nouveaux fichiers vendor sont requis
- [ ] Tester la compilation sur les variantes AVX et AVX2
- [ ] Vérifier que les stubs MTMD restent compatibles
- [ ] Documenter tout nouveau problème de compilation dans ce fichier

## 🔄 Versions testées

- **llama.cpp** : commit `de569441470332ff922c23fb0413cc957be75b25`
- **Système** : Ubuntu/Debian avec gcc 13.3.0
- **CMake** : version 3.28.3
- **gRPC** : version 1.51.1
- **protobuf** : version 24.4.0

## 🐛 Dépannage

### Erreur "conflicting declaration"
- Vérifier que `mtmd-stub.h` ne redéfinit pas d'enums existants
- S'assurer que les includes sont dans le bon ordre

### Erreur "undefined reference"
- Vérifier que tous les fichiers source nécessaires sont listés dans `CMakeLists.txt`
- Contrôler que le script `prepare.sh` copie bien tous les fichiers requis

### Erreur "No such file or directory"
- Vérifier que les dépendances vendor sont copiées
- Contrôler les chemins d'inclusion dans `CMakeLists.txt`

## 👥 Contributeurs

- Corrections initiales : Agent de développement IA
- Date : 12 juillet 2025
- Contexte : Résolution des problèmes de compilation TTS dans LocalAI

---

> **Note** : Ce document doit être mis à jour à chaque modification significative du processus de build.
