![tera](assets/tera.png)

# Tera

A lightweight, language-agnostic tool for live browser reloading. Tera watches your output files and instantly reflects changes in the browser, making development faster and more efficient.

## Features

- **Live Reloading**: Instantly see your changes in the browser without manual refreshing
- **Universal File Support**: 
  - HTML, CSS, and JavaScript
  - Images, PDFs, and other static assets
  - Any file type your browser can render
- **Partial Re-rendering**: Smart updates that only refresh changed components
- **Language Agnostic**: Works with any programming language or build pipeline
- **Zero Configuration**: Just point it at your files and start developing

## Installation

```bash
go install github.com/kjabin/tera
```

## Quick Start

1. Navigate to your project directory:
```bash
cd project-dir
```

2. Start Tera:
```bash
tera index.html 
```

3. Open your browser at `http://localhost:8080`

## Usage

### Basic Command
```bash
tera [FILENAME] 
```

### Options
```bash
    --port      Specify the port number (default: 3000)
    --watch-dir Specify the root directory to be watched (default: ".")
    --exts      Filter file types to be watched (default: all)
    --help      Show help information
```

## How It Works

1. Tera monitors specified directory for file changes
2. When changes are detected, it notifies connected browsers
3. The browser updates only the necessary components

## Use Cases

- **Static Site Development**: Watch your HTML, CSS, and JavaScript files
- **Asset Preview**: Live preview of images, PDFs, or other media files
