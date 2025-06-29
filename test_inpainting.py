#!/usr/bin/env python3
"""
Test script to validate inpainting mask transmission from frontend to backend.
This script creates test images and sends them to the LocalAI inpainting endpoint
to verify that masks are transmitted correctly.
"""

import base64
import json
import requests
from PIL import Image, ImageDraw
import io
import numpy as np

def create_test_image(size=(512, 512)):
    """Create a simple test image with colored regions."""
    image = Image.new('RGB', size, color='lightblue')
    draw = ImageDraw.Draw(image)
    
    # Draw some shapes to make inpainting visible
    draw.rectangle([100, 100, 200, 200], fill='red', outline='darkred', width=3)
    draw.ellipse([300, 150, 450, 300], fill='green', outline='darkgreen', width=3)
    draw.polygon([(200, 300), (300, 250), (400, 350), (250, 400)], fill='yellow', outline='orange', width=3)
    
    return image

def create_test_mask(size=(512, 512)):
    """Create a test mask with white areas indicating inpainting regions."""
    # Create black background
    mask = Image.new('RGB', size, color='black')
    draw = ImageDraw.Draw(mask)
    
    # Draw white areas where inpainting should happen
    # Circle in the center
    draw.ellipse([200, 200, 300, 300], fill='white')
    # Rectangle on the right
    draw.rectangle([350, 100, 450, 200], fill='white')
    
    return mask

def image_to_base64(image):
    """Convert PIL image to base64 string."""
    buffer = io.BytesIO()
    image.save(buffer, format='PNG')
    return base64.b64encode(buffer.getvalue()).decode('utf-8')

def test_inpainting():
    """Test the inpainting endpoint with our test images."""
    print("Creating test images...")
    
    # Create test image and mask
    test_image = create_test_image()
    test_mask = create_test_mask()
    
    # Save test images for reference
    test_image.save('/tmp/test_image.png')
    test_mask.save('/tmp/test_mask.png')
    print("Test images saved to /tmp/test_image.png and /tmp/test_mask.png")
    
    # Convert to base64
    image_b64 = image_to_base64(test_image)
    mask_b64 = image_to_base64(test_mask)
    
    # Create inpainting data structure (matching frontend format)
    inpainting_data = {
        "image": image_b64,
        "mask_image": mask_b64
    }
    
    # Create request body (matching frontend format)
    request_body = {
        "model": "stablediffusion",  # You may need to adjust this
        "prompt": "a beautiful landscape with mountains",
        "num_inference_steps": 20,
        "guidance_scale": 7.5,
        "strength": 0.8,
        "size": "512x512",
        "n": 1,
        "file": base64.b64encode(json.dumps(inpainting_data).encode()).decode()
    }
    
    print("Sending inpainting request...")
    print(f"Image size: {test_image.size}")
    print(f"Mask size: {test_mask.size}")
    
    # Check mask values
    mask_array = np.array(test_mask)
    unique_values = np.unique(mask_array.reshape(-1, mask_array.shape[-1]), axis=0)
    print(f"Mask unique RGB values: {unique_values}")
    
    try:
        response = requests.post(
            'http://localhost:8080/v1/images/generations',
            headers={'Content-Type': 'application/json'},
            json=request_body,
            timeout=60
        )
        
        print(f"Response status: {response.status_code}")
        
        if response.status_code == 200:
            result = response.json()
            print("✅ Inpainting request successful!")
            print(f"Result: {result}")
        else:
            print(f"❌ Inpainting request failed: {response.status_code}")
            try:
                error_data = response.json()
                print(f"Error details: {error_data}")
            except:
                print(f"Error response: {response.text}")
    
    except requests.exceptions.RequestException as e:
        print(f"❌ Request failed: {e}")

if __name__ == "__main__":
    test_inpainting()
