document.addEventListener('DOMContentLoaded', function () {
    // Slug Generation
    const titleInput = document.getElementById('title');
    const slugInput = document.getElementById('slug');

    if (titleInput && slugInput) {
        titleInput.addEventListener('input', function () {
            if (!slugInput.dataset.manuallyChanged) {
                const slug = titleInput.value
                    .toLowerCase()
                    .replace(/[^\w ]+/g, '')
                    .replace(/ +/g, '-');
                slugInput.value = slug;
            }
        });

        slugInput.addEventListener('input', function () {
            slugInput.dataset.manuallyChanged = 'true';
        });
    }

    // Image Upload Preview
    const imageInput = document.getElementById('image-upload');
    const imagePreview = document.getElementById('image-preview');
    const imageUrlInput = document.getElementById('image_url');

    if (imageInput) {
        imageInput.addEventListener('change', function (e) {
            const file = e.target.files[0];
            if (file) {
                // Show local preview
                const reader = new FileReader();
                reader.onload = function (e) {
                    imagePreview.src = e.target.result;
                    imagePreview.style.display = 'block';
                }
                reader.readAsDataURL(file);

                // Upload to server (which uploads to Cloudflare)
                const formData = new FormData();
                formData.append('file', file);

                const csrfToken = document.querySelector('input[name="csrf_token"]')?.value;

                fetch('/admin/upload', {
                    method: 'POST',
                    headers: {
                        'X-CSRF-Token': csrfToken
                    },
                    body: formData
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.url) {
                            imageUrlInput.value = data.url;
                            // Update preview with remote URL to confirm
                            // imagePreview.src = data.url; 
                        }
                    })
                    .catch(error => {
                        console.error('Error uploading image:', error);
                        alert('Failed to upload image');
                    });
            }
        });
    }
});
