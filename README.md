# Go OCR Text Extractor

<img width="1536" height="1024" alt="gopher_reading_python_mind" src="https://github.com/user-attachments/assets/4096c4fa-4ee2-4f4f-b1dc-ef72fdddecc9" />


A powerful and free OCR (Optical Character Recognition) solution for the Go community that bridges the gap between Go's efficiency and Python's rich AI ecosystem.

## 🎯 Why This Project Exists

The OCR market is flooded with expensive, proprietary solutions that often come with limitations:
- **Costly licensing fees** for commercial OCR APIs
- **Limited language support** in many solutions
- **Vendor lock-in** with cloud-based services
- **Complex integration** requiring specialized knowledge
- **Poor accuracy** on varied document types
- **No local processing** options for sensitive documents

This project addresses these pain points by providing a **completely free**, **locally-run** OCR solution that leverages the power of established Python AI libraries while maintaining Go's performance and simplicity.

## 🌟 The Python Advantage

Python has already solved many complex problems in the AI and machine learning space with mature, battle-tested libraries. Rather than reinventing the wheel in Go, this project creates a bridge that allows Go developers to harness these powerful Python capabilities:

- **Tesseract OCR** - Google's industry-leading OCR engine
- **PIL (Python Imaging Library)** - Robust image processing
- **PyPDF2 & pdfplumber** - Comprehensive PDF text extraction
- **Extensive language support** - Over 100 languages supported by Tesseract

## 🚀 Features

- **PDF Text Extraction**: Extract text from PDF documents using multiple extraction methods
- **Image OCR**: Convert images to text with high accuracy
- **Multilingual Support**: Supports Portuguese, English, and 100+ other languages
- **Caching System**: Built-in memory cache to avoid reprocessing the same files
- **Fallback Mechanisms**: Multiple extraction methods ensure maximum compatibility
- **Thread-Safe**: Concurrent processing with mutex protection
- **Error Handling**: Comprehensive error reporting and recovery
- **Free & Open Source**: No licensing fees or API limits

## 📋 Prerequisites

### Python Dependencies

The project automatically handles Python dependency installation, but you can install them manually:

```bash
pip install PyPDF2 pdfplumber pytesseract Pillow
```

### Tesseract OCR Installation

#### Windows
1. Download the installer from [GitHub Tesseract releases](https://github.com/UB-Mannheim/tesseract/wiki)
2. Run the installer and follow the setup wizard
3. Add Tesseract to your system PATH:
   - Default installation path: `C:\Program Files\Tesseract-OCR`
   - Add this path to your Windows PATH environment variable
4. Restart your command prompt

#### macOS
Using Homebrew:
```bash
brew install tesseract
```

Using MacPorts:
```bash
sudo port install tesseract
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt-get update
sudo apt-get install tesseract-ocr
sudo apt-get install tesseract-ocr-por  # For Portuguese language support
```

#### Linux (CentOS/RHEL/Fedora)
```bash
sudo yum install tesseract tesseract-langpack-por
```
or for newer versions:
```bash
sudo dnf install tesseract tesseract-langpack-por
```

## 🏗️ Architecture

### Core Components

#### PythonExecutor
The main orchestrator that manages Python script execution and handles:
- Python environment detection
- Dependency management
- Script execution
- Result parsing
- Caching management

#### PDF Text Extraction
- **Primary Method**: pdfplumber (more accurate)
- **Fallback Method**: PyPDF2 (broader compatibility)
- **Automatic Selection**: Chooses the best method based on document type

#### Image OCR Processing
- **Language Detection**: Attempts Portuguese first, falls back to English
- **Image Preprocessing**: Automatic RGB conversion
- **Format Support**: PNG, JPG, JPEG, TIFF, BMP, GIF

#### Prompt System
- **LangChain Integration**: Structured prompt templates for AI processing
- **Flexible Configuration**: Customizable prompt parameters
- **JSON Output**: Structured response format

### Result Structures

```go
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
```

## 🔧 Usage Example

```go
package main

import (
    "fmt"
    "log"
    "your-project/model"
)

func main() {
    // Initialize the Python executor
    executor := model.NewPythonExecutor()
    
    // Check and install PDF dependencies
    if err := executor.CheckPythonDependenciesForPDF(); err != nil {
        log.Fatal(err)
    }
    
    // Extract text from PDF
    result, err := executor.ExtractPDFText("document.pdf")
    if err != nil {
        log.Fatal(err)
    }
    
    if result.Success {
        fmt.Printf("Extracted %d characters from %d pages\n", 
                   len(result.Text), result.Pages)
        fmt.Println(result.Text)
    } else {
        fmt.Printf("Error: %s\n", result.Error)
    }
    
    // Extract text from image
    imageResult, err := executor.ExtractImageText("scanned_document.png")
    if err != nil {
        log.Fatal(err)
    }
    
    if imageResult.Success {
        fmt.Printf("OCR Result: %s\n", imageResult.Text)
    }
    
    // Use with AI prompts
    promptInstance := model.NewPromptOCRInstance()
    prompt := promptInstance.GetPrompt(result.Text)
    // Process with your preferred LLM...
}
```

## 🤝 Contributing to the Go AI Community

This project represents a commitment to advancing AI capabilities within the Go ecosystem. By providing free access to powerful OCR functionality, we aim to:

- **Lower barriers to entry** for developers interested in document processing
- **Accelerate innovation** in Go-based AI applications
- **Foster collaboration** between Go and Python communities
- **Enable experimentation** without cost constraints
- **Support education** and research initiatives

## 📚 Dependencies & Credits

This project stands on the shoulders of giants. We extend our gratitude to the following projects and their maintainers:

### Python Libraries
- **[Tesseract OCR](https://github.com/tesseract-ocr/tesseract)** - Google's open-source OCR engine
- **[pytesseract](https://github.com/madmaze/pytesseract)** - Python wrapper for Tesseract
- **[PyPDF2](https://github.com/py-pdf/PyPDF2)** - Pure Python PDF library
- **[pdfplumber](https://github.com/jsvine/pdfplumber)** - Detailed PDF text extraction
- **[Pillow (PIL)](https://python-pillow.org/)** - Python Imaging Library

### Go Libraries
- **[LangChain Go](https://github.com/tmc/langchaingo)** - Go implementation of LangChain for LLM integration

## 🔍 Error Handling

The system provides comprehensive error handling:
- **Dependency Check**: Automatic verification and installation of Python packages
- **File Validation**: Format and existence verification
- **Graceful Fallbacks**: Multiple extraction methods for maximum compatibility
- **Detailed Logging**: Clear error messages and debugging information

## 🎯 Future Enhancements

- Support for more document formats (DOCX, RTF, etc.)
- Batch processing capabilities
- Configuration file support
- Docker containerization
- REST API wrapper
- Performance benchmarking tools

## 📄 License

This project is open-source and free to use. Please refer to the LICENSE file for details.

## 🤝 Contributing

We welcome contributions from the community! Whether it's bug reports, feature requests, or code contributions, every effort helps make this tool better for everyone.

---

*Built with ❤️ for the Go community. Empowering developers to build amazing AI applications without breaking the bank.*
