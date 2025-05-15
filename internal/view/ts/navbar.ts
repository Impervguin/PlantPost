document.addEventListener('DOMContentLoaded', () => {
    const profileButton = document.getElementById('user-menu-button') as HTMLButtonElement;
    const profileMenu = document.getElementById('user-menu') as HTMLDivElement;

    if (!profileButton || !profileMenu) return;

    profileButton.addEventListener('click', () => {
        const isExpanded = profileButton.getAttribute('aria-expanded') === 'true';
        profileButton.setAttribute('aria-expanded', String(!isExpanded));
        if (!isExpanded) {
            profileMenu.classList.remove('hidden');
            profileMenu.classList.add('transition', 'ease-out', 'duration-100');
            profileMenu.classList.add('transform', 'opacity-0', 'scale-95');
            setTimeout(() => {
                profileMenu.classList.remove('transform', 'opacity-0', 'scale-95');
                profileMenu.classList.add('transform', 'opacity-100', 'scale-100');
            }, 20);
        } else {
            profileMenu.classList.add('transition', 'ease-in', 'duration-75');
            profileMenu.classList.remove('transform', 'opacity-100', 'scale-100');
            profileMenu.classList.add('transform', 'opacity-0', 'scale-95');
            setTimeout(() => {
                profileMenu.classList.add('hidden');
            }, 75);
        }
    });
    
    // Close when clicking outside
    document.addEventListener('click', (e: Event) => {
        if (!profileMenu.classList.contains('hidden') && 
            !profileButton.contains(e.target as Node) && 
            !profileMenu.contains(e.target as Node)) {
            profileMenu.classList.add('transition', 'ease-in', 'duration-75');
            profileMenu.classList.remove('transform', 'opacity-100', 'scale-100');
            profileMenu.classList.add('transform', 'opacity-0', 'scale-95');
            
            setTimeout(() => {
                profileMenu.classList.add('hidden');
                profileButton.setAttribute('aria-expanded', 'false');
            }, 75);
        }
    });
});