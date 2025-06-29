// Inpainting functionality
class InpaintingManager {
    constructor() {
        this.originalImage = null;
        this.maskImage = null;
        this.canvas = null;
        this.ctx = null;
        this.isDrawing = false;
        this.brushSize = 20;
        
        this.initializeElements();
        this.setupEventListeners();
    }
    
    initializeElements() {
        this.canvas = document.getElementById('mask-canvas');
        this.ctx = this.canvas?.getContext('2d');
        this.brushSizeSlider = document.getElementById('brush-size');
        this.brushSizeValue = document.getElementById('brush-size-value');
    }
    
    setupEventListeners() {
        // File upload handlers
        document.getElementById('original-image')?.addEventListener('change', (e) => this.handleOriginalImageUpload(e));
        document.getElementById('mask-image')?.addEventListener('change', (e) => this.handleMaskImageUpload(e));
        document.getElementById('draw-image')?.addEventListener('change', (e) => this.handleDrawImageUpload(e));
        
        // Canvas drawing handlers
        if (this.canvas) {
            this.canvas.addEventListener('mousedown', (e) => this.startDrawing(e));
            this.canvas.addEventListener('mousemove', (e) => this.draw(e));
            this.canvas.addEventListener('mouseup', () => this.stopDrawing());
            this.canvas.addEventListener('mouseout', () => this.stopDrawing());
            
            // Touch handlers for mobile
            this.canvas.addEventListener('touchstart', (e) => this.startDrawing(e));
            this.canvas.addEventListener('touchmove', (e) => this.draw(e));
            this.canvas.addEventListener('touchend', () => this.stopDrawing());
        }
        
        // Brush size handler
        this.brushSizeSlider?.addEventListener('input', (e) => {
            this.brushSize = parseInt(e.target.value);
            this.brushSizeValue.textContent = `${this.brushSize}px`;
        });
        
        // Clear canvas handler
        document.getElementById('clear-canvas')?.addEventListener('click', () => this.clearCanvas());
        
        // Generate button handler
        document.getElementById('generate-btn')?.addEventListener('click', () => this.generateInpainting());
    }
    
    handleOriginalImageUpload(event) {
        const file = event.target.files[0];
        if (!file) return;
        
        this.loadImage(file, (img) => {
            this.originalImage = img;
            this.showImagePreview(img, 'original-preview', 'original-img', 'original-placeholder');
        });
    }
    
    handleMaskImageUpload(event) {
        const file = event.target.files[0];
        if (!file) return;
        
        this.loadImage(file, (img) => {
            this.maskImage = img;
            this.showImagePreview(img, 'mask-preview', 'mask-img', 'mask-placeholder');
        });
    }
    
    handleDrawImageUpload(event) {
        const file = event.target.files[0];
        if (!file) return;
        
        this.loadImage(file, (img) => {
            this.originalImage = img;
            this.showImagePreview(img, 'draw-preview', 'draw-img', 'draw-placeholder');
            this.setupCanvas(img);
        });
    }
    
    loadImage(file, callback) {
        const reader = new FileReader();
        reader.onload = (e) => {
            const img = new Image();
            img.onload = () => callback(img);
            img.src = e.target.result;
        };
        reader.readAsDataURL(file);
    }
    
    showImagePreview(img, previewId, imgId, placeholderId) {
        const preview = document.getElementById(previewId);
        const imgElement = document.getElementById(imgId);
        const placeholder = document.getElementById(placeholderId);
        
        if (preview && imgElement && placeholder) {
            imgElement.src = img.src;
            preview.classList.remove('hidden');
            placeholder.classList.add('hidden');
        }
    }
    
    setupCanvas(img) {
        if (!this.canvas || !this.ctx) return;
        
        // Set canvas size to match the EXACT image dimensions
        // No resizing to avoid dimension mismatches with the original image
        const width = img.naturalWidth || img.width;
        const height = img.naturalHeight || img.height;
        
        console.log('Setting up canvas with exact image dimensions:', width, 'x', height);
        
        this.canvas.width = width;
        this.canvas.height = height;
        this.canvas.style.width = Math.min(width, 512) + 'px';
        this.canvas.style.height = Math.min(height, 512) + 'px';
        
        // Clear canvas with black background (black = no inpainting, white = inpaint)
        this.ctx.fillStyle = 'black';
        this.ctx.fillRect(0, 0, width, height);
    }
    
    startDrawing(event) {
        if (!this.canvas || !this.ctx) return;
        
        this.isDrawing = true;
        const rect = this.canvas.getBoundingClientRect();
        const x = (event.clientX || event.touches[0].clientX) - rect.left;
        const y = (event.clientY || event.touches[0].clientY) - rect.top;
        
        this.ctx.beginPath();
        this.ctx.moveTo(x, y);
        
        event.preventDefault();
    }
    
    draw(event) {
        if (!this.isDrawing || !this.canvas || !this.ctx) return;
        
        const rect = this.canvas.getBoundingClientRect();
        const x = (event.clientX || event.touches[0].clientX) - rect.left;
        const y = (event.clientY || event.touches[0].clientY) - rect.top;
        
        this.ctx.lineWidth = this.brushSize;
        this.ctx.lineCap = 'round';
        this.ctx.strokeStyle = 'white';
        this.ctx.lineTo(x, y);
        this.ctx.stroke();
        this.ctx.beginPath();
        this.ctx.moveTo(x, y);
        
        event.preventDefault();
    }
    
    stopDrawing() {
        if (!this.isDrawing) return;
        this.isDrawing = false;
        this.ctx?.beginPath();
    }
    
    clearCanvas() {
        if (!this.canvas || !this.ctx) return;
        this.ctx.fillStyle = 'black';
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
    }
    
    async generateInpainting() {
        const model = document.getElementById('image-model')?.value;
        const positivePrompt = document.getElementById('positive-prompt')?.value;
        const negativePrompt = document.getElementById('negative-prompt')?.value || '';
        const steps = parseInt(document.getElementById('steps')?.value) || 20;
        const guidanceScale = parseFloat(document.getElementById('guidance-scale')?.value) || 7.5;
        const strength = parseFloat(document.getElementById('strength')?.value) || 0.8;
        
        if (!model || !positivePrompt) {
            alert('Please select a model and enter a prompt');
            return;
        }
        
        if (!this.originalImage) {
            alert('Please upload an original image');
            return;
        }
        
        // Get mask image
        let maskImageData;
        
        // Check which mode is active by looking for active tab button
        const uploadTab = document.querySelector('button[x-text*="Upload"]');
        const isUploadMode = uploadTab && uploadTab.classList.contains('bg-purple-600');
        
        let maskWidth, maskHeight;
        
        if (isUploadMode) {
            if (!this.maskImage) {
                alert('Please upload a mask image');
                return;
            }
            
            // Check dimensions match
            const origWidth = this.originalImage.naturalWidth || this.originalImage.width;
            const origHeight = this.originalImage.naturalHeight || this.originalImage.height;
            maskWidth = this.maskImage.naturalWidth || this.maskImage.width;
            maskHeight = this.maskImage.naturalHeight || this.maskImage.height;
            
            console.log('Image dimensions check:');
            console.log('- Original image:', origWidth, 'x', origHeight);
            console.log('- Mask image:', maskWidth, 'x', maskHeight);
            
            if (origWidth !== maskWidth || origHeight !== maskHeight) {
                alert(`Dimension mismatch! Original image is ${origWidth}x${origHeight} but mask is ${maskWidth}x${maskHeight}. Please ensure both images have the same dimensions.`);
                return;
            }
            
            maskImageData = this.imageToBase64(this.maskImage);
        } else {
            if (!this.canvas) {
                alert('Please draw a mask on the canvas');
                return;
            }
            
            // Check dimensions match for canvas mode
            const origWidth = this.originalImage.naturalWidth || this.originalImage.width;
            const origHeight = this.originalImage.naturalHeight || this.originalImage.height;
            maskWidth = this.canvas.width;
            maskHeight = this.canvas.height;
            
            console.log('Canvas dimensions check:');
            console.log('- Original image:', origWidth, 'x', origHeight);
            console.log('- Canvas mask:', maskWidth, 'x', maskHeight);
            
            if (origWidth !== maskWidth || origHeight !== maskHeight) {
                alert(`Dimension mismatch! Original image is ${origWidth}x${origHeight} but canvas is ${maskWidth}x${maskHeight}. Please reload and ensure the canvas matches the image size.`);
                return;
            }
            
            // Get mask data and add debug info
            const maskDataURL = this.canvas.toDataURL('image/png');
            maskImageData = maskDataURL.split(',')[1];
            
            // Debug: Check if canvas has any white pixels
            const imageData = this.ctx.getImageData(0, 0, this.canvas.width, this.canvas.height);
            const data = imageData.data;
            let hasWhitePixels = false;
            let whitePixelCount = 0;
            
            for (let i = 0; i < data.length; i += 4) {
                // Check if pixel is white (or close to white)
                if (data[i] > 200 && data[i + 1] > 200 && data[i + 2] > 200) {
                    hasWhitePixels = true;
                    whitePixelCount++;
                }
            }
            
            console.log('Canvas mask debug:');
            console.log('- Canvas size:', this.canvas.width, 'x', this.canvas.height);
            console.log('- Has white pixels:', hasWhitePixels);
            console.log('- White pixel count:', whitePixelCount);
            console.log('- Total pixels:', (data.length / 4));
            console.log('- Mask data URL length:', maskDataURL.length);
            
            if (!hasWhitePixels) {
                if (!confirm('Warning: Your mask appears to be empty (all black). This means no inpainting will occur. Do you want to continue anyway?')) {
                    return;
                }
            }
        }
        
        // Show loading state
        this.setLoadingState(true);
        
        try {
            // Prepare the image data
            const originalImageData = this.imageToBase64(this.originalImage);
            
            // Debug: Log final data sizes
            console.log('Final data preparation:');
            console.log('- Original image base64 length:', originalImageData.length);
            console.log('- Mask image base64 length:', maskImageData.length);
            
            // Create the inpainting data
            const inpaintingData = {
                image: originalImageData,
                mask_image: maskImageData
            };
            
            console.log('Inpainting data object created, JSON size:', JSON.stringify(inpaintingData).length);
            
            // Create the request
            const requestBody = {
                model: model,
                prompt: negativePrompt ? `${positivePrompt}|${negativePrompt}` : positivePrompt,
                num_inference_steps: steps,
                guidance_scale: guidanceScale,
                strength: strength,
                size: "512x512",
                n: 1,
                // Pass the inpainting data as a base64 encoded JSON string
                file: btoa(JSON.stringify(inpaintingData))
            };
            
            console.log('Sending inpainting request with file size:', requestBody.file.length);
            
            const response = await fetch('v1/images/generations', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestBody),
            });
            
            const result = await response.json();
            
            if (result.error) {
                throw new Error(result.error.message || 'Generation failed');
            }
            
            if (!result.data || !result.data[0] || !result.data[0].url) {
                throw new Error('Invalid response format');
            }
            
            this.displayResults(result.data[0].url);
            
        } catch (error) {
            console.error('Inpainting error:', error);
            alert('Error generating inpainting: ' + error.message);
        } finally {
            this.setLoadingState(false);
        }
    }
    
    imageToBase64(img) {
        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');
        
        // Use natural dimensions to preserve original image size
        const width = img.naturalWidth || img.width;
        const height = img.naturalHeight || img.height;
        
        canvas.width = width;
        canvas.height = height;
        
        console.log('Converting image to base64 with dimensions:', width, 'x', height);
        
        ctx.drawImage(img, 0, 0, width, height);
        const dataUrl = canvas.toDataURL('image/png');
        
        console.log('Generated base64 data URL length:', dataUrl.length);
        
        return dataUrl.split(',')[1];
    }
    
    setLoadingState(loading) {
        const generateBtn = document.getElementById('generate-btn');
        const generateText = document.getElementById('generate-text');
        const generateLoader = document.getElementById('generate-loader');
        
        if (generateBtn && generateText && generateLoader) {
            generateBtn.disabled = loading;
            generateText.classList.toggle('hidden', loading);
            generateLoader.classList.toggle('hidden', !loading);
        }
    }
    
    displayResults(resultUrl) {
        const resultsSection = document.getElementById('results-section');
        const originalResult = document.getElementById('original-result');
        const inpaintedResult = document.getElementById('inpainted-result');
        
        if (resultsSection && originalResult && inpaintedResult) {
            // Show original image
            originalResult.innerHTML = `<img src="${this.originalImage.src}" class="max-w-full h-auto rounded-lg">`;
            
            // Show inpainted result
            inpaintedResult.innerHTML = `<img src="${resultUrl}" class="max-w-full h-auto rounded-lg">`;
            
            // Show results section
            resultsSection.classList.remove('hidden');
            
            // Scroll to results
            resultsSection.scrollIntoView({ behavior: 'smooth' });
        }
    }
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new InpaintingManager();
});
