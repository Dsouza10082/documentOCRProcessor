package model

import (
	"bufio"
	_ "embed"
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

//go:embed extract_text_from_image.model.py
var extract_image_py string

//go:embed extract_text_from_pdf.model.py
var extract_pdf_py string


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
	tmpFile, err := os.CreateTemp("", "pdf_extractor_*.py")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %v", err)
	}
	
	if _, err := tmpFile.WriteString(extract_pdf_py); err != nil {
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
	tmpFile, err := os.CreateTemp("", "image_extractor_*.py")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %v", err)
	}
	
	if _, err := tmpFile.WriteString(extract_image_py); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("error writing script: %v", err)
	}
	
	tmpFile.Close()
	return tmpFile.Name(), nil
}