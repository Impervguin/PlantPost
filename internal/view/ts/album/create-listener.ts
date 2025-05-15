import { MultiSelectPlantSearch } from './plant-select.js';

document.addEventListener('DOMContentLoaded', () => {
    console.log('loaded');
    const searchInput = document.getElementById('search-input') as HTMLInputElement;
    const searchResults = document.getElementById('search-results') as HTMLDivElement;
    const selectedItems = document.getElementById('selected-items') as HTMLDivElement;

    if (!searchInput || !searchResults || !selectedItems) {
        console.error('Required elements not found');
        return;
    }

    const search = new MultiSelectPlantSearch(searchInput, searchResults, selectedItems, []);

    const createButton = document.getElementById('create-button') as HTMLButtonElement;
    const nameInput = document.getElementById('name') as HTMLInputElement;
    const descriptionInput = document.getElementById('description') as HTMLTextAreaElement;
    
    createButton.addEventListener('click', () => {
        const selectedItems = search.GetSelectedItems();
        const name = nameInput.value;
        const description = descriptionInput.value;
        let postForm = new FormData();
        postForm.append('name', name);
        postForm.append('description', description);
        for (const item of selectedItems) {
            postForm.append('plant_ids', item.id);
        }
        if (selectedItems.length === 0) {
            postForm.append('plant_ids', '');
        }
        fetch('/api/album/create', {
            method: 'POST',
            body: postForm
        }).then(response => {
            if (response.ok) {
                window.location.href = '/view/albums/';
            } else {
                console.error('Failed to create album:', response);
            }
        });
    });
});