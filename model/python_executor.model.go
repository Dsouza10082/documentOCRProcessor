package model

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type PythonExecutor struct {
	pythonPath string
	textCache  map[string]string
	mutex      sync.RWMutex
}

type PDFTextResult struct {
	Success  bool   `json:"success"`
	Text     string `json:"text"`
	Error    string `json:"error"`
	Pages    int    `json:"pages"`
	Filename string `json:"filename"`
}

type ImageTextResult struct {
	Success  bool   `json:"success"`
	Text     string `json:"text"`
	Error    string `json:"error"`
	Pages    int    `json:"pages"`
	Filename string `json:"filename"`
}

func NewPythonExecutor() *PythonExecutor {
	executor := &PythonExecutor{
		textCache: make(map[string]string),
	}
	executor.findPythonPath()
	return executor
}

func (pe *PythonExecutor) findPythonPath() {
	pythonCommands := []string{"python3", "python"}
	
	for _, cmd := range pythonCommands {
		if path, err := exec.LookPath(cmd); err == nil {
			pe.pythonPath = path
			fmt.Printf("Python found in: %s\n", path)
			return
		}
	}
	
	log.Fatal("❌ Python not found in system. Please ensure Python is installed and in the PATH.")
}

func (pe *PythonExecutor) CheckPythonDependenciesForPDF() error {
	dependencies := []string{"PyPDF2", "pdfplumber"}
	
	fmt.Printf("🔍 Checking Python dependencies...\n")
	
	for _, dep := range dependencies {
		cmd := exec.Command(pe.pythonPath, "-c", fmt.Sprintf("import %s", strings.ToLower(dep)))
		if err := cmd.Run(); err != nil {
			fmt.Printf("Dependency %s not found. Trying to install...\n", dep)
			
			installCmd := exec.Command(pe.pythonPath, "-m", "pip", "install", dep)
			if err := installCmd.Run(); err != nil {
				return fmt.Errorf("error installing %s: %v\nTry manually: pip install %s", dep, err, dep)
			}
			fmt.Printf("✅ %s installed successfully\n", dep)
		} else {
			fmt.Printf("✅ %s already installed\n", dep)
		}
	}
	
	return nil
}

func (pe *PythonExecutor) CheckPythonDependenciesForImage() error {
	dependencies := []string{"pytesseract"}
	
	fmt.Printf("🔍 Checking Python dependencies...\n")
	
	for _, dep := range dependencies {
		cmd := exec.Command(pe.pythonPath, "-c", fmt.Sprintf("import %s", strings.ToLower(dep)))
		if err := cmd.Run(); err != nil {
			fmt.Printf("Dependency %s not found. Trying to install...\n", dep)
			
			installCmd := exec.Command(pe.pythonPath, "-m", "pip", "install", dep)
			if err := installCmd.Run(); err != nil {
				return fmt.Errorf("error installing %s: %v\nTry manually: pip install %s", dep, err, dep)
			}
			fmt.Printf("✅ %s installed successfully\n", dep)
		} else {
			fmt.Printf("✅ %s already installed\n", dep)
		}
	}
	
	return nil
}

func (pe *PythonExecutor) createPDFExtractor() (string, error) {
	script := `#!/usr/bin/env python3

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
`

	tmpFile, err := os.CreateTemp("", "pdf_extractor_*.py")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %v", err)
	}
	
	if _, err := tmpFile.WriteString(script); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("error writing script: %v", err)
	}
	
	tmpFile.Close()
	return tmpFile.Name(), nil
}

func (pe *PythonExecutor) ExtractPDFText(pdfPath string) (*PDFTextResult, error) {

	pe.mutex.RLock()
	if cachedText, exists := pe.textCache[pdfPath]; exists {
		pe.mutex.RUnlock()
		return &PDFTextResult{
			Success:  true,
			Text:     cachedText,
			Filename: filepath.Base(pdfPath),
		}, nil
	}
	pe.mutex.RUnlock()
	
	scriptPath, err := pe.createPDFExtractor()
	if err != nil {
		return nil, err
	}
	defer os.Remove(scriptPath)

	cmd := exec.Command(pe.pythonPath, scriptPath, pdfPath)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &PDFTextResult{
			Success: false,
			Error:   fmt.Sprintf("Error executing script: %v", err),
		}, nil
	}

	outputStr := string(output)
	jsonStart := strings.Index(outputStr, "JSON_RESULT_START")
	jsonEnd := strings.Index(outputStr, "JSON_RESULT_END")
	
	if jsonStart == -1 || jsonEnd == -1 {
		return &PDFTextResult{
			Success: false,
			Error:   "Invalid output format from Python script",
		}, nil
	}
	
	jsonStr := outputStr[jsonStart+len("JSON_RESULT_START"):jsonEnd]
	jsonStr = strings.TrimSpace(jsonStr)
	
	var result PDFTextResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return &PDFTextResult{
			Success: false,
			Error:   fmt.Sprintf("Error decoding JSON: %v", err),
		}, nil
	}
	
	if result.Success {
		pe.mutex.Lock()
		pe.textCache[pdfPath] = result.Text
		pe.mutex.Unlock()
	}
	
	return &result, nil
}

func (pe *PythonExecutor) ExtractImageText(imagePath string) (*ImageTextResult, error) {

	pe.mutex.RLock()
	if cachedText, exists := pe.textCache[imagePath]; exists {
		pe.mutex.RUnlock()
		return &ImageTextResult{
			Success:  true,
			Text:     cachedText,
			Filename: filepath.Base(imagePath),
		}, nil
	}
	pe.mutex.RUnlock()
	
	scriptPath, err := pe.CreateImageExtractor()
	if err != nil {
		return nil, err
	}
	defer os.Remove(scriptPath)

	cmd := exec.Command(pe.pythonPath, scriptPath, imagePath)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &ImageTextResult{
			Success: false,
			Error:   fmt.Sprintf("Error executing script: %v", err),
		}, nil
	}

	outputStr := string(output)
	jsonStart := strings.Index(outputStr, "JSON_RESULT_START")
	jsonEnd := strings.Index(outputStr, "JSON_RESULT_END")
	
	if jsonStart == -1 || jsonEnd == -1 {
		return &ImageTextResult{
			Success: false,
			Error:   "Invalid output format from Python script",
		}, nil
	}
	
	jsonStr := outputStr[jsonStart+len("JSON_RESULT_START"):jsonEnd]
	jsonStr = strings.TrimSpace(jsonStr)
	
	var result ImageTextResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return &ImageTextResult{
			Success: false,
			Error:   fmt.Sprintf("Error decoding JSON: %v", err),
		}, nil
	}
	
	if result.Success {
		pe.mutex.Lock()
		pe.textCache[imagePath] = result.Text
		pe.mutex.Unlock()
	}
	
	return &result, nil
}

func (pe *PythonExecutor) saveTextToFile(pdfPath, outputPath string) error {
	pe.mutex.RLock()
	text, exists := pe.textCache[pdfPath]
	pe.mutex.RUnlock()
	
	if !exists {
		return fmt.Errorf("text not found in cache for: %s", pdfPath)
	}
	
	return os.WriteFile(outputPath, []byte(text), 0644)
}

func (pe *PythonExecutor) readOutput(reader io.Reader, label string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "JSON_RESULT") {
			fmt.Printf("[%s] %s\n", label, line)
		}
	}
}

func (pe *PythonExecutor) showCachedTexts() {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	
	fmt.Printf("\n💾 Texts in Memory (%d files):\n", len(pe.textCache))
	fmt.Printf("═══════════════════════════════════\n")
	
	for path, text := range pe.textCache {
		filename := filepath.Base(path)
		preview := text
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		
		fmt.Printf("📄 %s\n", filename)
		fmt.Printf("   Size: %d characters\n", len(text))
		fmt.Printf("   Preview: %s\n\n", strings.ReplaceAll(preview, "\n", " "))
	}
}


func (pe *PythonExecutor) readFile(filename string) (string, error) {

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", filename)
	}
	
	if !strings.HasSuffix(strings.ToLower(filename), ".py") {
		return "", fmt.Errorf("file must have .py extension: %s", filename)
	}
	
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	
	return string(content), nil
}

func (pe *PythonExecutor) executeFile(filename string) error {

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}
	
	fmt.Printf("\n=== Executando: %s ===\n", absPath)

	cmd := exec.Command(pe.pythonPath, absPath)
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %v", err)
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	go pe.readOutput(stdout, "STDOUT")
	go pe.readOutput(stderr, "STDERR")
	
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error during execution: %v", err)
	}
	
	fmt.Printf("\n=== Execution completed ===\n")
	return nil
}

func (pe *PythonExecutor) validatePythonCode(content string) error {
	tmpFile, err := os.CreateTemp("", "validate_*.py")
	if err != nil {
		return fmt.Errorf("error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	
	if _, err := tmpFile.WriteString(content); err != nil {
		return fmt.Errorf("error writing temporary file: %v", err)
	}
	tmpFile.Close()
	
	cmd := exec.Command(pe.pythonPath, "-m", "py_compile", tmpFile.Name())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("syntax error in Python code: %v", err)
	}
	
	return nil
}

func (pe *PythonExecutor) CreateImageExtractor() (string, error) {
	script := `#!/usr/bin/env python3

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
`

	tmpFile, err := os.CreateTemp("", "image_extractor_*.py")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %v", err)
	}
	
	if _, err := tmpFile.WriteString(script); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("error writing script: %v", err)
	}
	
	tmpFile.Close()
	return tmpFile.Name(), nil
}