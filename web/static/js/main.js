document.addEventListener('DOMContentLoaded', function () {

    // ==========================================
    // 1. Social Share Functionality
    // ==========================================
    function initShareButtons() {
        const shareBtns = document.querySelectorAll('.share-btn');
        shareBtns.forEach(btn => {
            // Check if listener already attached (custom attribute)
            if (btn.dataset.listenerAttached === 'true') return;

            btn.addEventListener('click', function (e) {
                e.preventDefault(); // Prevent default if it's an anchor

                const platform = btn.dataset.platform;
                // Use current page URL or data-url if provided
                const url = btn.dataset.url ? encodeURIComponent(btn.dataset.url) : encodeURIComponent(window.location.href);
                // Use current page title or data-title
                const title = btn.dataset.title ? encodeURIComponent(btn.dataset.title) : encodeURIComponent(document.title);

                let shareUrl = '';

                switch (platform) {
                    case 'twitter':
                        shareUrl = `https://twitter.com/intent/tweet?url=${url}&text=${title}`;
                        break;
                    case 'facebook':
                        shareUrl = `https://www.facebook.com/sharer/sharer.php?u=${url}`;
                        break;
                    case 'linkedin':
                        shareUrl = `https://www.linkedin.com/sharing/share-offsite/?url=${url}`;
                        break;
                    case 'whatsapp':
                        shareUrl = `https://wa.me/?text=${title}%20${url}`;
                        break;
                }

                if (shareUrl) {
                    window.open(shareUrl, '_blank', 'width=600,height=400,status=0,toolbar=0');
                }
            });

            btn.dataset.listenerAttached = 'true';
        });
    }

    // ==========================================
    // 2. Copy Link Functionality
    // ==========================================
    function initCopyLinkButtons() {
        const copyLinkBtns = document.querySelectorAll('.copy-link-btn');
        copyLinkBtns.forEach(btn => {
            if (btn.dataset.listenerAttached === 'true') return;

            btn.addEventListener('click', function (e) {
                e.preventDefault();
                const url = window.location.href;

                navigator.clipboard.writeText(url).then(() => {
                    const originalIcon = btn.innerHTML;
                    btn.innerHTML = '<i class="ph ph-check text-green-500 text-lg"></i>';
                    setTimeout(() => {
                        btn.innerHTML = originalIcon;
                    }, 2000);
                }).catch(err => {
                    console.error('Failed to copy: ', err);
                });
            });

            btn.dataset.listenerAttached = 'true';
        });
    }

    // ==========================================
    // 3. Live News Updates (Polling)
    // ==========================================
    function initLiveUpdates() {
        const container = document.getElementById('live-updates-container');
        const timer = document.getElementById('live-update-timer');

        if (!container || !timer) return; // Only runs if element exists (homepage)

        // Poll every 45 seconds
        const pollInterval = 45000;

        setInterval(() => {
            fetchLiveUpdates();
        }, pollInterval);

        function fetchLiveUpdates() {
            timer.innerText = 'Updating...';
            timer.classList.add('text-red-200');

            fetch('/api/live-updates')
                .then(response => response.json())
                .then(data => {
                    if (data && data.length > 0) {
                        updateLiveFeed(data);
                    }
                })
                .catch(err => console.error('Live update failed:', err))
                .finally(() => {
                    timer.innerText = 'Live';
                    timer.classList.remove('text-red-200');
                    // Add blink effect to indicate update
                    timer.classList.add('animate-pulse');
                    setTimeout(() => timer.classList.remove('animate-pulse'), 2000);
                });
        }

        function updateLiveFeed(items) {
            // Clear current items or diff them. For simplicity, replace.
            // But we want to animate.
            // Let's build the HTML string
            let html = '';
            items.forEach(item => {
                const time = new Date(item.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
                html += `
                <div class="p-4 hover:bg-gray-50 transition border-l-2 border-transparent hover:border-red-500 animate-fade-in-down">
                    <div class="text-[10px] text-gray-400 font-bold mb-1 flex justify-between">
                        <span class="uppercase text-red-600">${item.source}</span>
                        <span>${time}</span>
                    </div>
                    <a href="${item.url}" class="text-sm font-medium text-gray-900 leading-snug hover:text-red-600 transition block">
                        ${item.headline}
                    </a>
                </div>`;
            });

            container.innerHTML = html;
        }
    }

    // ==========================================
    // 4. Newsletter Subscription
    // ==========================================
    function initNewsletter() {
        const form = document.getElementById('newsletter-form');
        const msgContainer = document.getElementById('newsletter-message');

        if (!form) return;

        form.addEventListener('submit', function (e) {
            e.preventDefault();

            const emailInput = form.querySelector('input[type="email"]');
            const button = form.querySelector('button');
            const originalBtnContent = button.innerHTML;
            const email = emailInput.value;

            // Loading state
            button.disabled = true;
            button.innerHTML = '<i class="ph ph-spinner animate-spin"></i>';

            fetch('/api/newsletter/subscribe', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email: email })
            })
                .then(response => {
                    if (response.ok) {
                        // Success
                        msgContainer.innerHTML = '<span class="text-green-500 flex items-center gap-1"><i class="ph ph-check-circle"></i> Subscribed successfully!</span>';
                        form.reset();
                    } else {
                        // Error
                        msgContainer.innerHTML = '<span class="text-red-500 flex items-center gap-1"><i class="ph ph-warning"></i> Subscription failed. Try again.</span>';
                    }
                })
                .catch(err => {
                    console.error('Newsletter error:', err);
                    msgContainer.innerHTML = '<span class="text-red-500 flex items-center gap-1"><i class="ph ph-warning"></i> Error. Please try again.</span>';
                })
                .finally(() => {
                    button.disabled = false;
                    button.innerHTML = originalBtnContent;
                    setTimeout(() => {
                        msgContainer.innerHTML = '<i class="ph ph-shield-check"></i> We respect your privacy';
                    }, 5000);
                });
        });
    }

    // Initialize all
    initShareButtons();
    initCopyLinkButtons();
    initLiveUpdates();
    initNewsletter();

});
