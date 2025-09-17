#!/usr/bin/env python3

import sys
import json
import os
from PIL import Image
import pytesseract

def extract_image_text(image_path):
    try:
        if not os.path.exists(image_path):
            return {
                "success": False,
                "text": "",
                "error": f"Arquivo não encontrado: {image_path}",
                "pages": 0,
                "filename": os.path.basename(image_path)
            }
        
        valid_extensions = ['.png', '.jpg', '.jpeg', '.tiff', '.bmp', '.gif']
        file_extension = os.path.splitext(image_path)[1].lower()
        
        if file_extension not in valid_extensions:
            return {
                "success": False,
                "text": "",
                "error": f"Formato de arquivo não suportado: {file_extension}. Suportados: {', '.join(valid_extensions)}",
                "pages": 0,
                "filename": os.path.basename(image_path)
            }
        
        with Image.open(image_path) as image:
            if image.mode != 'RGB':
                image = image.convert('RGB')
            
            custom_config = r'--oem 3 --psm 6 -l por'
            text = pytesseract.image_to_string(image, config=custom_config)
            
            if not text.strip():
                custom_config = r'--oem 3 --psm 6 -l eng'
                text = pytesseract.image_to_string(image, config=custom_config)
            
            return {
                "success": True,
                "text": text.strip(),
                "error": "",
                "pages": 1,
                "filename": os.path.basename(image_path)
            }
            
    except FileNotFoundError:
        return {
            "success": False,
            "text": "",
            "error": f"File not found: {image_path}",
            "pages": 0,
            "filename": os.path.basename(image_path) if image_path else ""
        }
    except Exception as e:
        error_msg = str(e)
        if "tesseract is not installed" in error_msg.lower():
            error_msg = "Tesseract OCR is not installed. Install with: sudo apt-get install tesseract-ocr (Linux) or brew install tesseract (Mac)"
        
        return {
            "success": False,
            "text": "",
            "error": f"Error processing image: {error_msg}",
            "pages": 0,
            "filename": os.path.basename(image_path) if image_path else ""
        }

def main():
    if len(sys.argv) != 2:
        print(json.dumps({
            "success": False,
            "text": "",
            "error": "Usage: python image_extractor.py <image_file>",
            "pages": 0,
            "filename": ""
        }, ensure_ascii=False, indent=2))
        sys.exit(1)

    image_path = sys.argv[1]
    result = extract_image_text(image_path)
    
    print("JSON_RESULT_START")
    print(json.dumps(result, ensure_ascii=False, indent=2))
    print("JSON_RESULT_END")

if __name__ == "__main__":    
    main()