"""
Backend TTS personnalisé pour LocalAI utilisant pyttsx3, gTTS ou espeak.
Support de plusieurs moteurs TTS avec configuration flexible des voix et 
langues.
"""

import argparse
import signal
import sys
import os

# Add the common directory to the path
sys.path.append(os.path.join(os.path.dirname(__file__), "..", "common"))

from grpc_server import serve  # noqa: E402
from backend import BackendServicer  # noqa: E402


def main():
    parser = argparse.ArgumentParser(description="Run the MyTTS backend")
    parser.add_argument(
        "--addr", 
        default="localhost:50051", 
        help="Address to bind the server to"
    )
    args = parser.parse_args()

    # Setup signal handlers
    def signal_handler(sig, frame):
        print("Received termination signal, shutting down gracefully...")
        sys.exit(0)

    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)

    print(f"[MyTTS] Starting server on {args.addr}")
    serve(args.addr, BackendServicer)


if __name__ == "__main__":
    main()
