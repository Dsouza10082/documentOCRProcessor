#!/usr/bin/env python3

import sys
import json
import os
from pathlib import Path

def extract_with_pypdf2(pdf_path):
    try:
        import PyPDF2
        
        text = ""
        page_count = 0
        
        with open(pdf_path, 'rb') as file:
            pdf_reader = PyPDF2.PdfReader(file)
            page_count = len(pdf_reader.pages)
            
            for page_num in range(page_count):
                page = pdf_reader.pages[page_num]
                text += f"\n--- Página {page_num + 1} ---\n"
                text += page.extract_text()
        
        return True, text, page_count, None
    except Exception as e:
        return False, "", 0, str(e)

def extract_with_pdfplumber(pdf_path):
    """Extract text using pdfplumber (more precise)"""
    try:
        import pdfplumber
        
        text = ""
        page_count = 0
        
        with pdfplumber.open(pdf_path) as pdf:
            page_count = len(pdf.pages)
            
            for page_num, page in enumerate(pdf.pages):
                text += f"\n--- Page {page_num + 1} ---\n"
                page_text = page.extract_text()
                if page_text:
                    text += page_text
                else:
                    text += "[Page without extractable text]"
        
        return True, text, page_count, None
    except Exception as e:
        return False, "", 0, str(e)

def extract_pdf_text(pdf_path):
    """Main function for extracting text"""
    
    if not os.path.exists(pdf_path):
        return {
            "success": False,
            "text": "",
            "error": f"Arquivo não encontrado: {pdf_path}",
            "pages": 0,
            "filename": os.path.basename(pdf_path)
        }
    
    if not pdf_path.lower().endswith('.pdf'):
        return {
            "success": False,
            "text": "",
            "error": "File must have .pdf extension",
            "pages": 0,
            "filename": os.path.basename(pdf_path)
        }
    
    print(f"📄 Processando PDF: {os.path.basename(pdf_path)}")
    
    # Try with pdfplumber (more precise)
    success, text, pages, error = extract_with_pdfplumber(pdf_path)
    
    if not success:
        print(f"⚠️  pdfplumber failed: {error}")
        print("🔄 Trying with PyPDF2...")
        
        # Fallback to PyPDF2
        success, text, pages, error = extract_with_pypdf2(pdf_path)
    
    result = {
        "success": success,
        "text": text.strip(),
        "error": error if not success else "",
        "pages": pages,
        "filename": os.path.basename(pdf_path)
    }
    
    return result

def main():
    if len(sys.argv) != 2:
        print(json.dumps({
            "success": False,
            "text": "",
            "error": "Usage: python pdf_extractor.py <pdf_file>",
            "pages": 0,
            "filename": ""
        }))
        sys.exit(1)
    
    pdf_path = sys.argv[1]
    result = extract_pdf_text(pdf_path)
    
    print("JSON_RESULT_START")
    print(json.dumps(result, ensure_ascii=False, indent=2))
    print("JSON_RESULT_END")
    
    if result["success"]:
        print(f"\n✅ Extraction completed!")
        print(f"📊 Total pages: {result['pages']}")
        print(f"📝 Extracted characters: {len(result['text'])}")
    else:
        print(f"\n❌ Error in extraction: {result['error']}")

if __name__ == "__main__":
    main()