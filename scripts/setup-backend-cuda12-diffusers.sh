#!/bin/bash
# Script d'automatisation pour préparer et fiabiliser le backend cuda12-diffusers
set -e


# Détecter le chemin absolu du dossier du backend, peu importe où est lancé le script
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/../backends/cuda12-diffusers"
cd "$BACKEND_DIR"

# Exporter la variable d'environnement CUDA pour tout le setup
export LD_LIBRARY_PATH=/usr/local/cuda/lib64:$LD_LIBRARY_PATH

# 1. Corriger le shebang et les fins de ligne du run.sh
if ! head -1 run.sh | grep -q '#!/bin/bash'; then
  sed -i '1i #!/bin/bash' run.sh
fi
dos2unix run.sh
chmod +x run.sh

# 2. Créer/mettre à jour le venv Python
if [ ! -d venv ]; then
  python3.12 -m venv venv
fi
# S'assurer que python pointe sur python3.12 dans le venv
ln -sf python3.12 venv/bin/python
source venv/bin/activate
pip install --upgrade pip

# 3. Installer les dépendances critiques
pip install protobuf grpcio

# 4. Installer les requirements spécifiques si présents
for req in requirements.txt requirements-cublas12.txt requirements-cpu.txt; do
  if [ -f "$req" ]; then
    pip install -r "$req"
  fi
done

# 4b. Appliquer le patch stable_diffusion_inpaint_support.patch si présent
PATCH_PATH="$SCRIPT_DIR/../patch/stable_diffusion_inpaint_support.patch"
if [ -f "$PATCH_PATH" ]; then
  echo "[setup-backend.sh] Application du patch stable_diffusion_inpaint_support.patch"
  patch -N -p0 < "$PATCH_PATH" || echo "[setup-backend.sh] Patch déjà appliqué ou erreur mineure."
fi

# 5. Régénérer les fichiers protobuf avec la version installée
if [ -f backend.proto ]; then
  python -m pip install --upgrade pip
  python -m pip install grpcio-tools
  python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. backend.proto
fi

# 6. Exporter la bonne variable d'environnement CUDA dans run.sh si absent
if ! grep -q 'export LD_LIBRARY_PATH=/usr/local/cuda/lib64:$LD_LIBRARY_PATH' run.sh; then
  sed -i '2i export LD_LIBRARY_PATH=/usr/local/cuda/lib64:$LD_LIBRARY_PATH' run.sh
fi

echo "[setup-backend.sh] Backend cuda12-diffusers prêt à l'emploi."
