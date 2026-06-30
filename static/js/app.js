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
                // Open the uploaded file right away
                window.location.href = '/view?file=' + encodeURIComponent(data.path);
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

    // Remove Directory functionality
    const removeMessage = document.getElementById('removeMessage');
    const removeButtons = document.querySelectorAll('.remove-dir-btn');

    removeButtons.forEach(btn => {
        btn.addEventListener('click', function() {
            const directory = this.getAttribute('data-dir');
            if (!directory) return;

            const li = this.closest('li.directory-item');
            
            fetch('/api/remove-directory', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ directory: directory })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    if (removeMessage) showMessage(removeMessage, data.message, 'success');
                    if (li) li.remove();
                    
                    // Optimistically remove directory group from available files section
                    const fileGroups = document.querySelectorAll('.directory-group');
                    fileGroups.forEach(group => {
                        const h3 = group.querySelector('h3.directory-name');
                        if (h3 && h3.textContent.trim() === directory) {
                            group.remove();
                        }
                    });

                    // Hide empty list message if necessary or show it if empty
                    const remainingItems = document.querySelectorAll('.directory-item');
                    if (remainingItems.length === 0) {
                        const dirList = document.querySelector('.directory-list');
                        if (dirList) {
                            dirList.insertAdjacentHTML('afterend', '<p class="empty-message">No directories being watched. Add one above!</p>');
                            dirList.remove();
                        }
                    }

                    setTimeout(() => {
                        if (removeMessage) removeMessage.style.display = 'none';
                    }, 3000);
                } else {
                    if (removeMessage) showMessage(removeMessage, data.message || 'Failed to untrack directory', 'error');
                }
            })
            .catch(error => {
                if (removeMessage) showMessage(removeMessage, 'Error: ' + error.message, 'error');
            });
        });
    });
});
