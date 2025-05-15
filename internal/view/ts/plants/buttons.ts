document.addEventListener('DOMContentLoaded', () => {
    const filterButtons = document.querySelectorAll('.filter-close-open');
    
    filterButtons.forEach(button => {
        
        button.addEventListener('click', () => {
            const filterName = button.getAttribute('aria-controls')?.replace('filter-section-', '');
            if (!filterName) return;
            
            const filterSection = document.getElementById(`filter-section-${filterName}`);
            if (!filterSection) return;
            
            const isExpanded = button.getAttribute('aria-expanded') === 'true';
            button.setAttribute('aria-expanded', String(!isExpanded));
            filterSection.classList.toggle('hidden');
            
            const icons = button.querySelectorAll('svg');
            icons.forEach(icon => icon.classList.toggle('hidden'));
        });
    });
});