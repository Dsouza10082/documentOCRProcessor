package model

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileProcessor struct {
	ToProcessDir string
	SuccessDir   string
	ErrorDir     string
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		ToProcessDir: "./learning/to-process",
		SuccessDir:   "./learning/success",
		ErrorDir:     "./learning/error",
	}
}

func (fp *FileProcessor) ReadDirectory() ([]string, error) {
	var filePaths []string
	err := filepath.WalkDir(fp.ToProcessDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			filePaths = append(filePaths, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %w", fp.ToProcessDir, err)
	}

	return filePaths, nil
}

func (fp *FileProcessor) ProcessFile(filePath string) bool {
	fmt.Printf("Processing file: %s\n", filePath)
	fileName := filepath.Base(filePath)
	if strings.Contains(strings.ToLower(fileName), "error") {
		fmt.Printf("❌ Error processing: %s\n", filePath)
		return false
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("❌ Error getting file information: %s\n", filePath)
		return false
	}
	if fileInfo.Size() == 0 {
		fmt.Printf("❌ Empty file: %s\n", filePath)
		return false
	}
	executor := NewPythonExecutor()

	var result *PDFTextResult
	if strings.HasSuffix(strings.ToLower(fileName), ".pdf") {

		if err := executor.CheckPythonDependenciesForPDF(); err != nil {
			fmt.Printf("❌ Error checking Python dependencies: %v", err)
			return false
		}

		pdfPath := filePath
		fmt.Printf("Extracting text from: %s\n", pdfPath)
		result, err := executor.ExtractPDFText(pdfPath)
		if err != nil {
			fmt.Printf("❌ Erro: %v", err)
			return false
		}

		if result.Success {

			fmt.Printf("\n✅ Success!\n")
			fmt.Printf("📊 Pages: %d\n", result.Pages)
			fmt.Printf("📝 Characters: %d\n", len(result.Text))
			fmt.Printf("💾 Text saved in memory\n")

			
			preview := result.Text
			prompt := NewPromptOCRInstance()
			promptText := prompt.GetPrompt(preview)

			fmt.Printf("💬 Prompt: %s\n", promptText)

			fmt.Printf("✅ Processed successfully: %s\n", fileName)
			
		} else {
			fmt.Printf("\n❌ Error extracting text: %s\n", result.Error)
			return false
		}
	} else if strings.HasSuffix(strings.ToLower(fileName), ".jpg") || strings.HasSuffix(strings.ToLower(fileName), ".jpeg") || strings.HasSuffix(strings.ToLower(fileName), ".png") {

		if err := executor.CheckPythonDependenciesForImage(); err != nil {
			fmt.Printf("❌ Error checking Python dependencies: %v", err)
			return false
		}

		imagePath := filePath
		fmt.Printf("Extracting text from: %s\n", imagePath)
		result, err := executor.ExtractImageText(imagePath)
		if err != nil {
			fmt.Printf("❌ Erro: %v", err)
			return false
		}

		if result.Success {

			fmt.Printf("\n✅ Success!\n")
			fmt.Printf("📝 Characters: %d\n", len(result.Text))

			preview := result.Text
			prompt := NewPromptOCRInstance()
			promptText := prompt.GetPrompt(preview)

			fmt.Printf("💬 Prompt: %s\n", promptText)

			fmt.Printf("✅ Processed successfully: %s\n", fileName)

		} else {
			fmt.Printf("\n❌ Error extracting text: %s\n", result.Error)
			return false
		}
	} else {
		fmt.Printf("❌ Error extracting text: %s\n", result.Error)
		return false
	}

	return true
}

func (fp *FileProcessor) processFileAndReturnRawMessage(filePath string) (string, bool) {

	fmt.Printf("Processing file: %s\n", filePath)

	fileName := filepath.Base(filePath)

	if strings.Contains(strings.ToLower(fileName), "error") {
		fmt.Printf("❌ Error processing: %s\n", filePath)
		return "", false
	}

	fileInfo, err := os.Stat(filePath)

	if err != nil {
		fmt.Printf("❌ Error getting file information: %s\n", filePath)
		return "", false
	}

	if fileInfo.Size() == 0 {
		fmt.Printf("❌ Empty file: %s\n", filePath)
		return "", false
	}

	executor := NewPythonExecutor()

	var result *PDFTextResult
	if strings.HasSuffix(strings.ToLower(fileName), ".pdf") {

		if err := executor.CheckPythonDependenciesForPDF(); err != nil {
			fmt.Printf("❌ Error checking Python dependencies: %v", err)
			return "", false
		}

		pdfPath := filePath
		fmt.Printf("Extracting text from: %s\n", pdfPath)
		result, err := executor.ExtractPDFText(pdfPath)
		if err != nil {
			fmt.Printf("❌ Erro: %v", err)
			return "", false
		}

		if result.Success {

			fmt.Printf("\n✅ Success!\n")
			fmt.Printf("📊 Pages: %d\n", result.Pages)
			fmt.Printf("📝 Characters: %d\n", len(result.Text))
			
			preview := result.Text
			prompt := NewPromptOCRInstance()
			promptText := prompt.GetPrompt(preview)

			fmt.Printf("💬 Prompt: %s\n", promptText)

			fmt.Printf("✅ Processed successfully: %s\n", fileName)

			return promptText, true
			
		} else {
			fmt.Printf("\n❌ Error extracting text: %s\n", result.Error)
			return "", false
		}
		
	} else if strings.HasSuffix(strings.ToLower(fileName), ".jpg") || strings.HasSuffix(strings.ToLower(fileName), ".jpeg") || strings.HasSuffix(strings.ToLower(fileName), ".png") {

		if err := executor.CheckPythonDependenciesForImage(); err != nil {
			fmt.Printf("❌ Error checking Python dependencies: %v", err)
			return "", false
		}

		imagePath := filePath
		fmt.Printf("Extracting text from: %s\n", imagePath)
		result, err := executor.ExtractImageText(imagePath)
		if err != nil {
			fmt.Printf("❌ Erro: %v", err)
			return "", false
		}

		if result.Success {

			fmt.Printf("\n✅ Success!\n")
			fmt.Printf("📝 Characters: %d\n", len(result.Text))

			preview := result.Text
			prompt := NewPromptOCRInstance()
			promptText := prompt.GetPrompt(preview)

			fmt.Printf("💬 Prompt: %s\n", promptText)

			fmt.Printf("✅ Processed successfully: %s\n", fileName)

			return promptText, true

		} else {
			fmt.Printf("\n❌ Error extracting text: %s\n", result.Error)
			return "", false
		}
	} else {
		fmt.Printf("❌ Error extracting text: %s\n", result.Error)
		return "", false
	}

}

func (fp *FileProcessor) moveFile(sourcePath, destinationDir string) error {
	fileName := filepath.Base(sourcePath)
	destinationPath := filepath.Join(destinationDir, fileName)
	if _, err := os.Stat(destinationPath); err == nil {
		timestamp := time.Now().Format("20060102_150405")
		ext := filepath.Ext(fileName)
		nameWithoutExt := strings.TrimSuffix(fileName, ext)
		fileName = fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext)
		destinationPath = filepath.Join(destinationDir, fileName)
	}
	err := os.Rename(sourcePath, destinationPath)
	if err != nil {
		return fmt.Errorf("error moving file from %s to %s: %w", sourcePath, destinationPath, err)
	}
	fmt.Printf("📁 File moved to: %s\n", destinationPath)
	return nil
}

func (fp *FileProcessor) ProcessFiles(filePaths []string) {
	fmt.Printf("\n🚀 Starting processing of %d file(s)...\n\n", len(filePaths))
	successCount := 0
	errorCount := 0
	for _, filePath := range filePaths {
		success := fp.ProcessFile(filePath)
		var destinationDir string
		if success {
			destinationDir = fp.SuccessDir
			successCount++
		} else {
			destinationDir = fp.ErrorDir
			errorCount++
		}
		if err := fp.moveFile(filePath, destinationDir); err != nil {
			log.Printf("Error moving file: %v", err)
		}
		fmt.Println()
	}

	fmt.Printf("📊 Summary of processing:\n")
	fmt.Printf("✅ Successes: %d\n", successCount)
	fmt.Printf("❌ Errors: %d\n", errorCount)
	fmt.Printf("📝 Total: %d\n", len(filePaths))
}



