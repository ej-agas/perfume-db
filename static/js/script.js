document.addEventListener('DOMContentLoaded', function() {
    htmx.logAll();
});

document.addEventListener('alpine:init', () => {
    Alpine.store('modal', {
        isOpen: false,
        init() {
            // Add event listener for Escape key
            document.addEventListener('keydown', (e) => {
                if (e.key === 'Escape' && this.isOpen) {
                    this.close();
                }
            });
        },
        show() {
            this.isOpen = true;
            document.body.style.overflow = 'hidden';
            // Focus trap for better accessibility
            setTimeout(() => {
                const focusable = document.querySelector('[x-data] [x-trap]')?.querySelector('button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])');
                if (focusable) focusable.focus();
            }, 100);
        },
        close() {
            this.isOpen = false;
            document.body.style.overflow = '';
            const content = document.getElementById('modal-content');
            if (content) content.innerHTML = '';
        },
        // Getter for x-show
        open() {
            return this.isOpen;
        }
    });
});