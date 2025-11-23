# 🔍 Diagnostic: Problème Frontend Inpainting

## ✅ Ce qui fonctionne

### Backend
- ✅ Le backend `cuda12-diffusers-inpainting` est installé et monté
- ✅ Le fichier de configuration `/models/dreamshaper-8-inpainting.yaml` existe
- ✅ La configuration YAML est correcte:
  ```yaml
  backend: cuda12-diffusers-inpainting
  name: dreamshaper-8-inpainting
  diffusers:
    pipeline_type: StableDiffusionInpaintPipeline
    enable_parameters: "prompt,negative_prompt,...,image,mask_image,..."
  parameters:
    model: Lykon/dreamshaper-8-inpainting
  ```

### LocalAI
- ✅ LocalAI charge correctement le modèle au démarrage
- ✅ L'endpoint `/v1/images/edits` est disponible

## ❌ Le problème identifié

### Logs révélateurs

```
5:34PM DBG context local model name not found, setting to the first model first model name=sd-3.5-large-ggml
5:34PM DBG overriding empty model name in request body with value found earlier in middleware chain context localModelName=sd-3.5-large-ggml
```

### Analyse

Le **frontend envoie un nom de modèle vide** (ou ne l'envoie pas du tout) dans la requête POST à `/v1/images/edits`.

Quand LocalAI reçoit une requête sans nom de modèle:
1. Il cherche le modèle dans le contexte local → **non trouvé**
2. Il utilise le **premier modèle disponible** par défaut → `sd-3.5-large-ggml`
3. Il essaie de charger le backend associé → `stablediffusion-ggml`
4. ❌ **Erreur**: Ce backend n'est pas installé

## 🔧 Solution

### Vérifier le code frontend

Le frontend doit envoyer le paramètre `model` dans la requête. Exemple de requête correcte:

```javascript
const formData = new FormData();
formData.append('model', 'dreamshaper-8-inpainting');  // ← CRITIQUE
formData.append('image', imageFile);
formData.append('mask', maskFile);
formData.append('prompt', 'your prompt here');

const response = await fetch('http://localhost:8080/v1/images/edits', {
  method: 'POST',
  body: formData
});
```

### Points à vérifier dans le code frontend

1. **Le paramètre `model` est-il envoyé ?**
   ```javascript
   // ❌ INCORRECT - pas de model
   formData.append('prompt', prompt);
   
   // ✅ CORRECT
   formData.append('model', 'dreamshaper-8-inpainting');
   formData.append('prompt', prompt);
   ```

2. **Le nom du modèle est-il correct ?**
   ```javascript
   // ❌ INCORRECT - faute de frappe
   formData.append('model', 'dreamshaper-8-inpating');
   
   // ✅ CORRECT
   formData.append('model', 'dreamshaper-8-inpainting');
   ```

3. **Le paramètre est-il dans le bon format ?**
   - Pour `multipart/form-data`: utiliser `formData.append('model', ...)`
   - Pour `application/json`: inclure `"model": "dreamshaper-8-inpainting"` dans le JSON

### Test rapide avec curl

Pour vérifier que le backend fonctionne:

```bash
curl -X POST http://localhost:8080/v1/images/edits \
  -F "model=dreamshaper-8-inpainting" \
  -F "image=@original.png" \
  -F "mask=@mask.png" \
  -F "prompt=a beautiful sunset" \
  -F "num_inference_steps=25"
```

Si cette commande fonctionne, le problème est **définitivement dans le frontend**.

## 📋 Checklist de débogage frontend

- [ ] Ouvrir les DevTools du navigateur (F12)
- [ ] Aller dans l'onglet "Network"
- [ ] Déclencher une requête d'inpainting
- [ ] Inspecter la requête POST vers `/v1/images/edits`
- [ ] Vérifier le payload:
  - [ ] Le paramètre `model` est présent ?
  - [ ] Sa valeur est `dreamshaper-8-inpainting` ?
  - [ ] Les paramètres `image` et `mask` sont présents ?

## 🎯 Prochaines étapes

1. **Localiser le code frontend** qui fait l'appel à `/v1/images/edits`
2. **Vérifier** que le paramètre `model` est bien envoyé
3. **Corriger** si nécessaire
4. **Tester** avec les DevTools pour confirmer

---

**Où se trouve le code frontend de votre application ?**
