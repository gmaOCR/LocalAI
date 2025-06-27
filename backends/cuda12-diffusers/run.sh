#!/bin/bash
set -e
backend_dir=$(dirname "$0")
venv_dir="$backend_dir/venv"

# Vérifie que python3 est disponible
if ! command -v python3 &> /dev/null; then
    echo "[run.sh] Erreur : python3 n'est pas installé." >&2
    exit 1
fi

# Supprime le venv s'il est corrompu (liens cassés)
if [ -d "$venv_dir" ] && [ ! -x "$venv_dir/bin/python3" ]; then
    echo "[run.sh] Suppression du venv corrompu (liens cassés)"
    rm -rf "$venv_dir"
fi

# Crée le venv si nécessaire
if [ ! -d "$venv_dir" ]; then
    echo "[run.sh] Création du venv Python dans $venv_dir..."
    python3 -m venv "$venv_dir" || { echo "Erreur lors de la création du venv"; ls -l "$backend_dir"; exit 1; }
    echo "[run.sh] Contenu du dossier après tentative de création du venv :"
    ls -l "$venv_dir/bin/"
    if [ ! -x "$venv_dir/bin/python3" ]; then
        echo "[run.sh] Le venv ne contient pas d'interpréteur python3 exécutable ! Abandon."
        exit 1
    fi
    source "$venv_dir/bin/activate"
    if [ -f "$backend_dir/requirements.txt" ]; then
        echo "[run.sh] Installation des dépendances Python..."
        pip install --upgrade pip
        pip install -r "$backend_dir/requirements.txt" || { echo "Erreur lors de l'installation des dépendances"; exit 1; }
    fi
else
    source "$venv_dir/bin/activate"
fi

if [ -d "$backend_dir/common" ]; then
    source "$backend_dir/common/libbackend.sh"
else
    source "$backend_dir/../common/libbackend.sh"
fi

# Lancer le backend avec le Python du venv
echo "[run.sh] Lancement du backend Python..."
ls -l "$venv_dir/bin/"
exec "$venv_dir/bin/python3" "$backend_dir/backend.py" "$@"