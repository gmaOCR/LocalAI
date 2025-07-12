# Solutions aux erreurs de compilation LocalAI

## Erreur 1: Protobuf non trouvé

### Description de l'erreur
```
CMake Error at tools/grpc-server/CMakeLists.txt:21 (find_package):
  Could not find a package configuration file provided by "Protobuf" with any
  of the following names:

    ProtobufConfig.cmake
    protobuf-config.cmake
```

### Analyse du problème
1. L'erreur se produit dans le fichier `llama-avx/llama.cpp/tools/grpc-server/CMakeLists.txt`
2. Ce fichier contient des appels à `find_package(Protobuf CONFIG REQUIRED)` à la ligne 21
3. Le script `prepare.sh` écrase le bon fichier `CMakeLists.txt` avec un fichier incorrect
4. Le fichier `backend/cpp/llama/CMakeLists.txt` est copié par `prepare.sh` mais contient du code incorrect

### Solution
Le problème vient du fait que `backend/cpp/llama/CMakeLists.txt` contient un code incorrect qui est copié vers `llama.cpp/tools/grpc-server/CMakeLists.txt`.

**Action corrective :** Corriger le fichier `backend/cpp/llama/CMakeLists.txt` pour qu'il utilise les cibles gRPC et Protobuf fournies par FetchContent au lieu d'appeler `find_package`.

### Fichiers modifiés
- `backend/cpp/llama/CMakeLists.txt` : Corrigé pour utiliser les cibles gRPC/Protobuf
- `backend/cpp/llama/llama.cpp/CMakeLists.txt` : Ajouté FetchContent pour gRPC
- `backend/cpp/llama/prepare.sh` : Simplifié pour éviter les conflits

### Statut
✅ **RÉSOLU** - L'erreur Protobuf a été corrigée

## Erreur 2: Répertoires manquants

### Description de l'erreur
```
CMake Error at tools/CMakeLists.txt:39 (add_subdirectory):
  add_subdirectory given source "infill" which is not an existing directory.

CMake Error at tools/CMakeLists.txt:40 (add_subdirectory):
  add_subdirectory given source "embedding" which is not an existing directory.
```

### Analyse du problème
Le fichier `tools/CMakeLists.txt` fait référence à des répertoires qui n'existent pas dans cette version de llama.cpp.

### Statut
🔄 En cours de résolution

## Erreur 3: Condition CMake invalide

### Description de l'erreur
```
CMake Error at CMakeLists.txt:267 (if):
  if given arguments:
    "GGML_CUDA_DMMV_X" ">" "64"
  Unknown arguments specified
```

### Analyse du problème
Il y a une condition `if` avec une syntaxe incorrecte dans le `CMakeLists.txt` principal.

### Statut
🔄 En cours de résolution

## Résolution des erreurs de compilation 

### ÉTAPE 7 : Correction du problème de duplication add_subdirectory(grpc-server)

**Problème identifié :**
```
CMake Error at tools/CMakeLists.txt:45 (add_subdirectory):
  The binary directory
    /fork/LocalAI/backend/cpp/llama-avx/llama.cpp/build/tools/grpc-server
  is already used to build a source directory.
```

**Cause :** Le répertoire `grpc-server` était ajouté deux fois dans le CMakeLists.txt des tools :
- Ligne 30 : `add_subdirectory(grpc-server)` 
- Ligne 45 : `add_subdirectory(grpc-server)` (dupliqué)

**Solution :** Suppression de la ligne dupliquée dans `/fork/LocalAI/backend/cpp/llama-avx/llama.cpp/tools/CMakeLists.txt`

### ÉTAPE 8 : Configuration de gRPC via pkg-config dans tools/CMakeLists.txt

**Problème :** Le CMakeLists.txt des tools utilisait `find_package(gRPC REQUIRED)` qui échouait.

**Solution :** Remplacement par la configuration pkg-config dans `/fork/LocalAI/backend/cpp/llama-avx/llama.cpp/tools/CMakeLists.txt` :

```cmake
# Find protobuf and gRPC for grpc-server using pkg-config
find_package(PkgConfig REQUIRED)
pkg_check_modules(GRPC REQUIRED grpc++)
pkg_check_modules(PROTOBUF REQUIRED protobuf)

# Create imported targets for gRPC and protobuf
if(NOT TARGET gRPC::grpc++)
    add_library(gRPC::grpc++ INTERFACE IMPORTED)
    target_include_directories(gRPC::grpc++ INTERFACE ${GRPC_INCLUDE_DIRS})
    target_link_libraries(gRPC::grpc++ INTERFACE ${GRPC_LIBRARIES})
    target_compile_options(gRPC::grpc++ INTERFACE ${GRPC_CFLAGS_OTHER})
endif()

if(NOT TARGET protobuf::libprotobuf)
    add_library(protobuf::libprotobuf INTERFACE IMPORTED)
    target_include_directories(protobuf::libprotobuf INTERFACE ${PROTOBUF_INCLUDE_DIRS})
    target_link_libraries(protobuf::libprotobuf INTERFACE ${PROTOBUF_LIBRARIES})
    target_compile_options(protobuf::libprotobuf INTERFACE ${PROTOBUF_CFLAGS_OTHER})
endif()
```

### ÉTAPE 9 : Amélioration du script prepare.sh

**Problèmes :** 
- Fichiers manquants : `server.cpp`, `chat.h`, `chat.cpp`, `common.h`, etc.
- Chemins incorrects (examples/server au lieu de tools/server)

**Solutions :** Amélioration du script `/fork/LocalAI/backend/cpp/llama-avx/prepare.sh` :

1. **Correction des chemins :**
   - `llama.cpp/tools/server/` au lieu de `llama.cpp/examples/server/`

2. **Ajout de copies automatiques pour tous les fichiers nécessaires :**
   ```bash
   # Copy necessary files for the grpc-server
   cp -r grpc-server.cpp llama.cpp/tools/grpc-server/ || echo "grpc-server.cpp already exists"
   if [ -f llama.cpp/common/json.hpp ] && [ ! -f llama.cpp/tools/grpc-server/json.hpp ]; then
       cp -rfv llama.cpp/common/json.hpp llama.cpp/tools/grpc-server/
   fi
   if [ -f llama.cpp/tools/server/utils.hpp ] && [ ! -f llama.cpp/tools/grpc-server/utils.hpp ]; then
       cp -rfv llama.cpp/tools/server/utils.hpp llama.cpp/tools/grpc-server/
   fi
   if [ -f llama.cpp/tools/server/server.cpp ] && [ ! -f llama.cpp/tools/grpc-server/server.cpp ]; then
       cp -rfv llama.cpp/tools/server/server.cpp llama.cpp/tools/grpc-server/
   fi
   if [ -f llama.cpp/common/common.h ] && [ ! -f llama.cpp/tools/grpc-server/common.h ]; then
       cp -rfv llama.cpp/common/common.h llama.cpp/tools/grpc-server/
   fi
   if [ -f llama.cpp/common/chat.h ] && [ ! -f llama.cpp/tools/grpc-server/chat.h ]; then
       cp -rfv llama.cpp/common/chat.h llama.cpp/tools/grpc-server/
   fi
   if [ -f llama.cpp/common/chat.cpp ] && [ ! -f llama.cpp/tools/grpc-server/chat.cpp ]; then
       cp -rfv llama.cpp/common/chat.cpp llama.cpp/tools/grpc-server/
   fi
   ```

3. **Génération automatique des fichiers web assets manquants :**
   ```bash
   # Create minimal web assets files (not needed for gRPC)
   for asset in index.html.hpp completion.js.hpp loading.html.hpp deps_daisyui.min.css.hpp deps_markdown-it.js.hpp deps_tailwindcss.js.hpp deps_vue.esm-browser.js.hpp; do
       if [ ! -f llama.cpp/tools/grpc-server/$asset ]; then
           asset_var=$(echo $asset | sed 's/\.hpp$//' | sed 's/\./_/g' | sed 's/-/_/g')
           echo "// Minimal $asset for grpc-server" > llama.cpp/tools/grpc-server/$asset
           echo "const char ${asset_var}[] = \"\";" >> llama.cpp/tools/grpc-server/$asset
           echo "const size_t ${asset_var}_len = 0;" >> llama.cpp/tools/grpc-server/$asset
       fi
   done
   ```

### ÉTAPE 10 : Optimisation du CMakeLists.txt grpc-server avec includes appropriés

**Problème :** Au lieu de copier tous les fichiers, il est plus élégant d'utiliser les bibliothèques et includes existants.

**Solution :** Configuration du CMakeLists.txt grpc-server avec includes appropriés :

```cmake
cmake_minimum_required(VERSION 3.13)
project(grpc-server)

# Find required packages
find_package(Threads REQUIRED)

# Create the grpc-server executable
add_executable(grpc-server
    grpc-server.cpp
)

# Include directories pour accéder aux headers de common et tools/server
target_include_directories(grpc-server PRIVATE 
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/../../common
    ${CMAKE_CURRENT_SOURCE_DIR}/../../tools/server
    ${CMAKE_CURRENT_SOURCE_DIR}/../../src
)

# Link libraries - use the grpc and protobuf that come with the main build
target_link_libraries(grpc-server PRIVATE 
    gRPC::grpc++ 
    protobuf::libprotobuf 
    Threads::Threads
    llama
    common
)
```

### ÉTAT ACTUEL

✅ **Problèmes résolus :**
- Configuration gRPC/protobuf via pkg-config
- Duplication add_subdirectory dans tools/CMakeLists.txt
- Copies automatiques des fichiers nécessaires via prepare.sh
- Génération des fichiers web assets manquants
- Configuration des includes et linking appropriés

⚠️ **Prochaines étapes nécessaires :**
- Vérifier que tous les headers restants (comme log.h) sont accessibles via les includes
- Tester la compilation complète `make build GO_TAGS="tts"`
- Documenter les dernières erreurs eventuelles
