# 🔧 Résolution du problème Backend Inpainting

## ❌ Problème identifié

```
5:34PM ERR Backend not found: stablediffusion-ggml
5:34PM ERR Failed to load model sd-3.5-large-ggml with backend stablediffusion-ggml
```

## 🔍 Analyse

### 1. **Architecture modulaire de LocalAI**

LocalAI v3.2+ utilise une architecture où les backends sont **séparés du binaire principal** et doivent être installés/chargés dynamiquement.

### 2. **Logique de sélection de modèle**

Voici comment LocalAI sélectionne le modèle:

```
1. Le frontend envoie: model="dreamshaper-8-inpainting"
2. Le middleware cherche une config pour ce modèle
3. Si non trouvé → utilise le premier modèle disponible (sd-3.5-large-ggml)
4. Charge le backend associé (stablediffusion-ggml)
5. ❌ ERREUR: Backend non installé dans l'image
```

**Log explicatif:**
```
5:34PM DBG context local model name not found, setting to the first model
5:34PM DBG overriding empty model name in request body with value found earlier
```

## ✅ Solutions

### Solution 1: Installer le backend stablediffusion-ggml

Le backend doit être installé dans l'image Docker. Modifions le Dockerfile:

```dockerfile
# Après la copie du binaire, installer les backends nécessaires
RUN /usr/local/bin/local-ai backends install stablediffusion-ggml
```

### Solution 2: Utiliser un backend Python (diffusers)

Alternative avec le backend Python qui supporte l'inpainting:

```dockerfile
RUN /usr/local/bin/local-ai backends install diffusers
```

### Solution 3: Créer une configuration de modèle

Créer un fichier de configuration pour votre modèle:

**`/models/dreamshaper-8-inpainting.yaml`**
```yaml
name: dreamshaper-8-inpainting
backend: diffusers  # ou stablediffusion-ggml si installé
parameters:
  model: runwayml/stable-diffusion-inpainting
  # Ou pour ggml:
  # model: /models/dreamshaper-8-inpainting.gguf
```

## 🚀 Solution recommandée

### Étape 1: Modifier le Dockerfile pour inclure les backends

```dockerfile
# Dans docker/Dockerfile.inpainting, après la copie du binaire:

# Install required backends for inpainting
RUN /usr/local/bin/local-ai backends install stablediffusion-ggml && \
    /usr/local/bin/local-ai backends install diffusers
```

### Étape 2: Rebuild et push

```bash
./scripts/build-and-push.sh
```

### Étape 3: Configuration côté serveur

Créer les configurations de modèles dans `/models/`:

**`/models/dreamshaper-8-inpainting.yaml`**
```yaml
name: dreamshaper-8-inpainting
backend: stablediffusion-ggml
parameters:
  model: /models/dreamshaper-8-inpainting.gguf
  step: 25
  cfg_scale: 7.0
```

## 🔄 Workaround immédiat (sans rebuild)

Si vous ne voulez pas rebuilder l'image maintenant:

### Option A: Installer le backend au runtime

```bash
docker exec -it votre-conteneur /usr/local/bin/local-ai backends install stablediffusion-ggml
```

### Option B: Utiliser le modèle existant

Modifiez votre frontend pour utiliser le modèle déjà chargé:

```javascript
// Au lieu de:
model: "dreamshaper-8-inpainting"

// Utilisez:
model: "sd-3.5-large-ggml"
```

### Option C: Monter les backends depuis l'hôte

```bash
docker run -d -p 8080:8080 \
  -v /path/to/models:/models \
  -v /path/to/backends:/backends \  # ← Montez les backends
  mercure.gregorymariani.com/localai:inpainting-latest
```

## 📋 Vérification des backends disponibles

Pour voir quels backends sont installés:

```bash
# Dans le conteneur
docker exec votre-conteneur ls -la /backends

# Ou via l'API
curl http://localhost:8080/backends/available
```

## 🎯 Prochaines étapes recommandées

1. **Court terme**: Utiliser le modèle `sd-3.5-large-ggml` qui est déjà chargé
2. **Moyen terme**: Créer une nouvelle image avec les backends pré-installés
3. **Long terme**: Configurer un système de gestion de modèles automatique

---

**Voulez-vous que je modifie le Dockerfile pour inclure les backends nécessaires et rebuilder l'image?**
