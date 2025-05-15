document.addEventListener('DOMContentLoaded', () => {
    const deleteButton = document.getElementById('delete-post-button') as HTMLButtonElement;
    if (!deleteButton) return;

    deleteButton.addEventListener('click', () => {
        const dialog = document.getElementById('delete-post-dialog') as HTMLDivElement;
        dialog.style.display = 'block';
    });

    const cancelButton = document.getElementById('delete-cancel-button') as HTMLButtonElement;
    if (!cancelButton) return;

    cancelButton.addEventListener('click', () => {
        const dialog = document.getElementById('delete-post-dialog') as HTMLDivElement;
        dialog.style.display = 'none';
    });

    const confirmButton = document.getElementById('delete-confirm-button') as HTMLButtonElement;
    if (!confirmButton) return;

    confirmButton.addEventListener('click', () => {
        const dialog = document.getElementById('delete-post-dialog') as HTMLDivElement;
        dialog.style.display = 'none';
        fetch(`/api/post/delete/${window.location.pathname.split('/')[3]}`, {
            method: 'DELETE'
        }).then(response => {
            if (response.ok) {
                window.location.href = '/view/posts';
            } else {
                console.error(response);
                throw new Error('Failed to delete post');
            }
        });
    });

    
});