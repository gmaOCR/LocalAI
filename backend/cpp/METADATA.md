# Métadonnées du projet LocalAI Backend C++

## Informations du build

- **Date de dernière modification** : 12 juillet 2025
- **Version llama.cpp** : `de569441470332ff922c23fb0413cc957be75b25`
- **Environnement testé** : Ubuntu/Debian avec gcc 13.3.0
- **Status** : ✅ Fonctionnel

## Checksums des fichiers critiques

### Scripts principaux
- `prepare.sh` : Script de préparation de l'environnement
- `grpc-server-CMakeLists.txt` : Configuration CMake
- `mtmd-stub.h` : Stubs multimodal (version corrigée)
- `mtmd-stub.cpp` : Implémentation des stubs

### Validation
- `validate_build.sh` : Script de validation du build

## Compatibilité

### Versions supportées
- CMake >= 3.10
- GCC >= 7.0 / Clang >= 8.0
- gRPC >= 1.30
- protobuf >= 3.12

### Architectures
- ✅ x86_64 avec AVX
- ✅ x86_64 avec AVX2 (variante séparée)
- ❌ ARM64 (non testé)
- ❌ x86_64 sans AVX (non supporté)

## Historique des modifications

### v1.0 (Juillet 2025)
- Résolution des conflits de définition `ggml_log_level`
- Ajout de toutes les dépendances manquantes
- Support des assets minja et des templates de chat
- Désactivation propre du code multimodal
- Documentation complète et scripts de validation

## Maintenance

### Tâches périodiques
- [ ] Vérifier les mises à jour de llama.cpp
- [ ] Tester la compatibilité avec les nouvelles versions de gRPC
- [ ] Mettre à jour la documentation

### Signalement de bugs
Les problèmes peuvent être signalés via :
1. Tests avec `validate_build.sh --clean`
2. Vérification des logs de compilation
3. Consultation de `BUILD_FIXES.md`

---

**Note** : Ce fichier est automatiquement généré et doit être mis à jour lors de chaque modification majeure.
