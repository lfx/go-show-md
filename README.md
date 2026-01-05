# Go Markdown Viewer

A single-service Go web application that renders markdown files from multiple directories with Mermaid diagram support, code highlighting, drag-and-drop file handling, and live reload on file changes.

## Features

- ✅ **Multi-directory support**: Watch multiple directories for markdown files
- ✅ **Mermaid diagrams**: Full support for Mermaid diagram rendering
- ✅ **Code highlighting**: Syntax highlighting for various programming languages
- ✅ **Live reload**: Automatic page refresh when markdown files change
- ✅ **Drag & drop**: Upload markdown files via drag-and-drop
- ✅ **Dark theme**: Beautiful dark theme optimized for reading
- ✅ **No configuration pollution**: Source directories remain untouched

## Quick Start

1. **Build the application**:
   ```bash
   go build -o go-show-md
   ```

2. **Run the server**:
   ```bash
   ./go-show-md
   ```

3. **Open your browser**:
   Navigate to `http://127.0.0.1:8080`

## Usage

### Adding Directories to Watch

1. On the home page, enter a directory path in the "Add Directory to Watch" input field
2. Click "Add Directory" button
3. The directory will be added to the watch list and all markdown files will be displayed

### Uploading Files

1. Drag and drop a markdown file onto the upload area
2. Or click the upload area to select a file
3. The file will be copied to `./watched-files/` directory and displayed in the file list

### Viewing Files

1. Click on any markdown file from the list
2. The file will be rendered with:
   - Syntax highlighting for code blocks
   - Mermaid diagram rendering
   - GitHub-flavored Markdown support

### Live Reload

When you edit a markdown file in your watched directories:
- The changes are detected automatically
- If you have the file open in the viewer, the page will refresh automatically
- No manual refresh needed!

## Configuration

The application creates a `config.json` file to store watched directories:

```json
{
  "watched_directories": [
    "/path/to/your/docs"
  ],
  "port": 8080,
  "host": "127.0.0.1"
}
```

You can manually edit this file to add or remove directories.

## Project Structure

```
go-show-md/
├── main.go                 # Entry point
├── config.json            # Configuration (auto-generated)
├── internal/
│   ├── config/            # Config management
│   ├── watcher/           # File system watcher
│   ├── renderer/          # Markdown renderer
│   └── handlers/          # HTTP handlers
├── templates/             # HTML templates
├── static/
│   ├── css/              # Stylesheets
│   └── js/               # JavaScript
└── watched-files/         # Uploaded files (auto-generated)
```

## Dependencies

- [goldmark](https://github.com/yuin/goldmark) - Markdown parser
- [goldmark-highlighting](https://github.com/yuin/goldmark-highlighting) - Syntax highlighting
- [fsnotify](https://github.com/fsnotify/fsnotify) - File system watcher
- [gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket support

## Security

- Server binds to `127.0.0.1` (localhost only) for security
- Path validation prevents directory traversal attacks
- No authentication required for personal use

## License

[MIT](https://www.tldrlegal.com/license/mit-license)
