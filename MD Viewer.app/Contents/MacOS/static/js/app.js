// Add Directory functionality
document.addEventListener('DOMContentLoaded', function() {
    const addDirectoryBtn = document.getElementById('addDirectoryBtn');
    const directoryInput = document.getElementById('directoryInput');
    const directoryMessage = document.getElementById('directoryMessage');

    if (addDirectoryBtn) {
        addDirectoryBtn.addEventListener('click', function() {
            const directory = directoryInput.value.trim();
            if (!directory) {
                showMessage(directoryMessage, 'Please enter a directory path', 'error');
                return;
            }

            fetch('/api/add-directory', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ directory: directory })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    showMessage(directoryMessage, data.message, 'success');
                    directoryInput.value = '';
                    setTimeout(() => {
                        location.reload();
                    }, 1000);
                } else {
                    showMessage(directoryMessage, data.message || 'Failed to add directory', 'error');
                }
            })
            .catch(error => {
                showMessage(directoryMessage, 'Error: ' + error.message, 'error');
            });
        });
    }

    // File upload functionality
    const uploadArea = document.getElementById('uploadArea');
    const fileInput = document.getElementById('fileInput');
    const uploadMessage = document.getElementById('uploadMessage');

    if (uploadArea && fileInput) {
        uploadArea.addEventListener('click', function() {
            fileInput.click();
        });

        uploadArea.addEventListener('dragover', function(e) {
            e.preventDefault();
            uploadArea.classList.add('dragover');
        });

        uploadArea.addEventListener('dragleave', function(e) {
            e.preventDefault();
            uploadArea.classList.remove('dragover');
        });

        uploadArea.addEventListener('drop', function(e) {
            e.preventDefault();
            uploadArea.classList.remove('dragover');

            const files = e.dataTransfer.files;
            if (files.length > 0) {
                handleFileUpload(files[0]);
            }
        });

        fileInput.addEventListener('change', function(e) {
            if (e.target.files.length > 0) {
                handleFileUpload(e.target.files[0]);
            }
        });
    }

    function handleFileUpload(file) {
        if (!file.name.endsWith('.md') && !file.name.endsWith('.markdown')) {
            showMessage(uploadMessage, 'Please upload a markdown file (.md or .markdown)', 'error');
            return;
        }

        const formData = new FormData();
        formData.append('file', file);

        fetch('/api/upload', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showMessage(uploadMessage, data.message, 'success');
                setTimeout(() => {
                    location.reload();
                }, 1000);
            } else {
                showMessage(uploadMessage, data.message || 'Failed to upload file', 'error');
            }
        })
        .catch(error => {
            showMessage(uploadMessage, 'Error: ' + error.message, 'error');
        });
    }

    function showMessage(element, message, type) {
        element.textContent = message;
        element.className = 'message ' + type;
        element.style.display = 'block';
    }
});
