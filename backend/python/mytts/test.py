#!/usr/bin/env python3
"""
Tests unitaires pour le backend TTS personnalisé
"""
import unittest
import grpc
import backend_pb2
import backend_pb2_grpc
import subprocess
import time
import os
import signal
import tempfile

class TestTTSBackend(unittest.TestCase):
    
    @classmethod
    def setUpClass(cls):
        """Démarrage du serveur pour les tests"""
        cls.server_process = subprocess.Popen(
            ['python3', 'backend.py', '--addr', 'localhost:50051'],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        # Attendre que le serveur démarre
        time.sleep(3)
    
    @classmethod 
    def tearDownClass(cls):
        """Arrêt du serveur après les tests"""
        if cls.server_process:
            cls.server_process.terminate()
            cls.server_process.wait()
    
    def test_health_check(self):
        """Test de vérification de l'état du backend"""
        with grpc.insecure_channel("localhost:50051") as channel:
            stub = backend_pb2_grpc.BackendStub(channel)
            response = stub.Health(backend_pb2.HealthMessage())
            self.assertEqual(response.message, b"OK")
    
    def test_load_model(self):
        """Test de chargement du modèle"""
        with grpc.insecure_channel("localhost:50051") as channel:
            stub = backend_pb2_grpc.BackendStub(channel)
            response = stub.LoadModel(backend_pb2.ModelOptions(Model="mytts-model"))
            self.assertTrue(response.success)
            self.assertIn("succès", response.message)
    
    def test_tts_generation(self):
        """Test de génération TTS"""
        with grpc.insecure_channel("localhost:50051") as channel:
            stub = backend_pb2_grpc.BackendStub(channel)
            
            # Chargement du modèle d'abord
            load_response = stub.LoadModel(backend_pb2.ModelOptions(Model="mytts-model"))
            self.assertTrue(load_response.success)
            
            # Génération TTS
            with tempfile.NamedTemporaryFile(suffix=".wav", delete=False) as tmp_file:
                tts_request = backend_pb2.TTSRequest(
                    text="Bonjour, ceci est un test de génération TTS",
                    model="mytts-model",
                    dst=tmp_file.name,
                    voice="default",
                    language="fr"
                )
                
                tts_response = stub.TTS(tts_request)
                self.assertTrue(tts_response.success)
                
                # Vérifier que le fichier a été créé
                self.assertTrue(os.path.exists(tmp_file.name))
                self.assertGreater(os.path.getsize(tmp_file.name), 0)
                
                # Nettoyage
                os.unlink(tmp_file.name)
    
    def test_status(self):
        """Test du statut du backend"""
        with grpc.insecure_channel("localhost:50051") as channel:
            stub = backend_pb2_grpc.BackendStub(channel)
            
            # Chargement du modèle d'abord
            stub.LoadModel(backend_pb2.ModelOptions(Model="mytts-model"))
            
            # Vérification du statut
            response = stub.Status(backend_pb2.HealthMessage())
            self.assertEqual(response.state, backend_pb2.StatusResponse.READY)

if __name__ == '__main__':
    unittest.main()
