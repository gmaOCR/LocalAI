#!/usr/bin/env python3
"""
Backend TTS personnalisé pour LocalAI
Exemple d'implémentation d'un backend Text-to-Speech utilisant une librairie TTS
"""
import os
import sys
import argparse
from concurrent import futures
import grpc
import backend_pb2
import backend_pb2_grpc

# Importations pour TTS - remplacer par votre librarie TTS préférée
# Exemple avec pyttsx3 (TTS hors ligne) ou une autre librairie
try:
    import pyttsx3
    TTS_AVAILABLE = True
except ImportError:
    print("pyttsx3 non disponible, utilisation d'un mock", file=sys.stderr)
    TTS_AVAILABLE = False

# Optionnel: vous pouvez aussi utiliser gTTS, espeak, festival, etc.
try:
    from gtts import gTTS
    import pygame
    GTTS_AVAILABLE = True
except ImportError:
    GTTS_AVAILABLE = False

_ONE_DAY_IN_SECONDS = 60 * 60 * 24
MAX_WORKERS = int(os.environ.get('PYTHON_GRPC_MAX_WORKERS', '1'))

class BackendServicer(backend_pb2_grpc.BackendServicer):
    """
    Classe principale du backend TTS
    """
    
    def __init__(self):
        self.model_loaded = False
        self.tts_engine = None
        self.voice_name = "default"
        self.language = "fr"  # langue par défaut
        
    def Health(self, request, context):
        """Vérification de l'état du backend"""
        return backend_pb2.Reply(message=bytes("OK", 'utf-8'))
    
    def LoadModel(self, request, context):
        """Chargement du modèle TTS"""
        try:
            print(f"Chargement du modèle TTS: {request.Model}", file=sys.stderr)
            
            if TTS_AVAILABLE:
                # Initialisation avec pyttsx3 (TTS hors ligne)
                self.tts_engine = pyttsx3.init()
                
                # Configuration des paramètres
                voices = self.tts_engine.getProperty('voices')
                if voices:
                    # Choix de la voix (vous pouvez personnaliser cette logique)
                    for voice in voices:
                        if 'french' in voice.name.lower() or 'fr' in voice.id.lower():
                            self.tts_engine.setProperty('voice', voice.id)
                            break
                
                # Configuration du débit et du volume
                self.tts_engine.setProperty('rate', 150)  # Vitesse de parole
                self.tts_engine.setProperty('volume', 0.9)  # Volume
                
                print("Modèle TTS chargé avec succès", file=sys.stderr)
                self.model_loaded = True
                
            else:
                print("Moteur TTS mock initialisé", file=sys.stderr)
                self.model_loaded = True
                
        except Exception as err:
            print(f"Erreur lors du chargement du modèle: {err}", file=sys.stderr)
            return backend_pb2.Result(success=False, message=f"Erreur de chargement: {err}")
            
        return backend_pb2.Result(success=True, message="Modèle chargé avec succès")
    
    def TTS(self, request, context):
        """Génération audio à partir du texte"""
        if not self.model_loaded:
            return backend_pb2.Result(success=False, message="Modèle non chargé")
            
        try:
            print(f"Génération TTS pour: '{request.text}'", file=sys.stderr)
            print(f"Destination: {request.dst}", file=sys.stderr)
            print(f"Voix: {request.voice}", file=sys.stderr)
            print(f"Langue: {request.language}", file=sys.stderr)
            
            # Configuration de la voix si spécifiée
            if request.voice and request.voice != "":
                self.voice_name = request.voice
                
            # Configuration de la langue si spécifiée
            if request.language and request.language != "":
                self.language = request.language
            
            if TTS_AVAILABLE and self.tts_engine:
                # Génération avec pyttsx3
                self.tts_engine.save_to_file(request.text, request.dst)
                self.tts_engine.runAndWait()
                
            elif GTTS_AVAILABLE:
                # Alternative avec gTTS (nécessite une connexion Internet)
                tts = gTTS(text=request.text, lang=self.language[:2])  # 'fr', 'en', etc.
                tts.save(request.dst)
                
            else:
                # Mock: création d'un fichier audio vide ou utilisation d'espeak
                print("Utilisation du mock TTS", file=sys.stderr)
                self._generate_mock_audio(request.text, request.dst)
            
            print(f"Audio généré et sauvé dans: {request.dst}", file=sys.stderr)
            
        except Exception as err:
            print(f"Erreur lors de la génération TTS: {err}", file=sys.stderr)
            return backend_pb2.Result(success=False, message=f"Erreur TTS: {err}")
            
        return backend_pb2.Result(success=True, message="Audio TTS généré avec succès")
    
    def _generate_mock_audio(self, text, dst):
        """Génération d'audio mock ou avec espeak"""
        try:
            # Option 1: Utiliser espeak si disponible sur le système
            os.system(f'espeak "{text}" -w "{dst}" 2>/dev/null')
            
            # Option 2: Si espeak n'est pas disponible, créer un fichier WAV simple
            if not os.path.exists(dst):
                self._create_simple_wav(dst)
                
        except Exception as e:
            print(f"Erreur génération mock: {e}", file=sys.stderr)
            self._create_simple_wav(dst)
    
    def _create_simple_wav(self, filepath):
        """Création d'un fichier WAV simple (silence) pour les tests"""
        # En-tête WAV minimal pour un fichier de 1 seconde de silence
        import struct
        
        sample_rate = 22050
        duration = 1.0  # secondes
        num_samples = int(sample_rate * duration)
        
        with open(filepath, 'wb') as f:
            # En-tête WAV
            f.write(b'RIFF')
            f.write(struct.pack('<I', 36 + num_samples * 2))
            f.write(b'WAVE')
            f.write(b'fmt ')
            f.write(struct.pack('<I', 16))
            f.write(struct.pack('<H', 1))  # PCM
            f.write(struct.pack('<H', 1))  # Mono
            f.write(struct.pack('<I', sample_rate))
            f.write(struct.pack('<I', sample_rate * 2))
            f.write(struct.pack('<H', 2))
            f.write(struct.pack('<H', 16))
            f.write(b'data')
            f.write(struct.pack('<I', num_samples * 2))
            
            # Données audio (silence)
            for _ in range(num_samples):
                f.write(struct.pack('<h', 0))
    
    def Status(self, request, context):
        """Statut du backend"""
        return backend_pb2.StatusResponse(
            state=backend_pb2.StatusResponse.READY if self.model_loaded else backend_pb2.StatusResponse.UNINITIALIZED
        )

def serve(address):
    """Démarrage du serveur gRPC"""
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=MAX_WORKERS))
    backend_pb2_grpc.add_BackendServicer_to_server(BackendServicer(), server)
    server.add_insecure_port(address)
    print(f"Serveur TTS démarré sur {address}", file=sys.stderr)
    server.start()
    print("Serveur TTS prêt", file=sys.stderr)
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        print("Arrêt du serveur TTS", file=sys.stderr)
        server.stop(0)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Backend TTS pour LocalAI')
    parser.add_argument(
        '--addr', 
        default='localhost:50051',
        help='Adresse du serveur gRPC (défaut: localhost:50051)'
    )
    args = parser.parse_args()
    
    serve(args.addr)
