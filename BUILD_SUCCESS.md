# ✅ Build Réussi!

L'image LocalAI avec support inpainting a été construite avec succès!

## 📦 Image construite

```
mercure.gregorymariani.com/localai:inpainting-latest
```

## 🚀 Prochaines étapes

### 1. Tester l'image localement (optionnel)

```bash
# Lancer le conteneur
docker run -d -p 8080:8080 \
  --name localai-test \
  mercure.gregorymariani.com/localai:inpainting-latest

# Vérifier que ça fonctionne
curl http://localhost:8080/readyz

# Tester l'endpoint inpainting
curl http://localhost:8080/v1/images/inpainting

# Arrêter le test
docker stop localai-test && docker rm localai-test
```

### 2. Pousser vers le registry

```bash
# Option A: Avec variables d'environnement
REGISTRY_USER="votre-username" \
REGISTRY_TOKEN="votre-token" \
./scripts/build-and-push.sh

# Option B: Avec fichier .env.registry
cp .env.registry.example .env.registry
# Éditer .env.registry avec vos credentials
source .env.registry
source scripts/registry-config.sh
./scripts/build-and-push.sh

# Option C: Push manuel
docker login mercure.gregorymariani.com -u votre-username
docker push mercure.gregorymariani.com/localai:inpainting-latest
```

### 3. Déployer sur votre serveur

```bash
# Sur votre serveur
docker pull mercure.gregorymariani.com/localai:inpainting-latest
docker run -d -p 8080:8080 \
  --name localai-inpainting \
  -v /path/to/models:/models \
  mercure.gregorymariani.com/localai:inpainting-latest
```

## 📊 Informations de build

- **Branche**: local/inpainting-image
- **Repository**: https://github.com/gmaOCR/LocalAI.git
- **Endpoint inpainting**: ✅ Inclus (`/v1/images/inpainting`)
- **Taille de l'image**: Vérifier avec `docker images`

## 🔍 Vérification de l'endpoint inpainting

Une fois déployé, testez l'endpoint:

```bash
curl -X POST http://votre-serveur:8080/v1/images/inpainting \
  -F "model=dreamshaper-8-inpainting" \
  -F "prompt=a beautiful sunset" \
  -F "steps=25" \
  -F "image=@original.png" \
  -F "mask=@mask.png"
```

## 📝 Notes

- L'image est actuellement **locale uniquement**
- Pour la pousser vers mercure.gregorymariani.com, fournissez vos credentials
- Le build a pris environ 25 secondes (grâce au cache Docker)
- L'image finale est optimisée (multi-stage build)

---

**Prêt à pousser vers le registry!** 🚀
