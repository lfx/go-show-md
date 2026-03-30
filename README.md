# Go Markdown Viewer

A single-service Go web application that renders markdown files from multiple directories with Mermaid diagram support, code highlighting, drag-and-drop file handling, and live reload on file changes.

## Features

- ✅ **Multi-directory support**: Watch multiple directories for markdown files
- ✅ **Mermaid diagrams**: Full support for Mermaid diagram rendering
- ✅ **Code highlighting**: Syntax highlighting for various programming languages
- ✅ **Smart file selection**: Automatically ignore noisy directories like `.git`, `node_modules`, `venv`, and more.
- ✅ **Custom ignore patterns**: Configure global and per-project ignore lists.
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

The application creates a `config.json` file to store watched directories and ignore patterns:

```json
{
  "watched_directories": [
    "/path/to/your/docs"
  ],
  "ignored_patterns": [
    "temp_*",
    "backup_*.md"
  ],
  "port": 8080,
  "host": "127.0.0.1"
}
```

### Ignoring Files & Directories

The application uses a "smart selection" mechanism to filter markdown files:

1. **Default Ignore List**: Common noisy directories like `.git`, `.venv`, `venv`, `node_modules`, `vendor`, and `.env` are ignored by default.
2. **Global Configuration**: You can create a file at `~/.go-show-md-ignore` with custom patterns (one per line).
3. **Per-Project Configuration**: You can place a `.go-show-md-ignore` file inside any watched directory to add project-specific ignore patterns.
4. **App Configuration**: Use the `ignored_patterns` array in `config.json` to specify patterns globally across all watched directories.

Patterns are matched using standard shell-style wildcards (e.g., `test_*.md`, `build/`).

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

## Releases

Automated builds and releases are managed via [GoReleaser](https://goreleaser.com/) and GitHub Actions.

- **GitHub Releases:** Every push to the `main` branch automatically triggers a new release package. A version tag is auto-generated, and macOS native binaries (both Intel `amd64` and Apple Silicon `arm64`) are published to the Releases page as `.tar.gz` archives.
- **Archive Contents:** The downloaded archive includes the compiled executable alongside the required `static/` and `templates/` folders so the application works seamlessly right out of the box.

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
