#!/bin/bash

APP_NAME="MD Viewer"
BINARY_NAME="go-show-md"
BUNDLE_NAME="${APP_NAME}.app"

echo "Building Go binary..."
go build -o "$BINARY_NAME"

if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

echo "Creating .app bundle structure..."
rm -rf "$BUNDLE_NAME"
mkdir -p "$BUNDLE_NAME/Contents/MacOS"
mkdir -p "$BUNDLE_NAME/Contents/Resources"

echo "Copying binary and assets..."
cp "$BINARY_NAME" "$BUNDLE_NAME/Contents/MacOS/"
cp -r templates "$BUNDLE_NAME/Contents/MacOS/"
cp -r static "$BUNDLE_NAME/Contents/MacOS/"

# Optional: Copy existing config if it exists
if [ -f config.json ]; then
    cp config.json "$BUNDLE_NAME/Contents/MacOS/"
fi

echo "Creating Info.plist..."
cat > "$BUNDLE_NAME/Contents/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>$BINARY_NAME</string>
    <key>CFBundleIdentifier</key>
    <string>com.liudas.go-show-md</string>
    <key>CFBundleName</key>
    <string>$APP_NAME</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0</string>
    <key>LSUIElement</key>
    <true/>
</dict>
</plist>
EOF

chmod +x "$BUNDLE_NAME/Contents/MacOS/$BINARY_NAME"

echo "Done! You can now drag '$BUNDLE_NAME' to your Applications folder."
