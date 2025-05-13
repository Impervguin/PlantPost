import { MultiSelectPlantSearch } from './plant-select.js';

document.addEventListener('DOMContentLoaded', () => {
    const searchInput = document.getElementById('search-input') as HTMLInputElement;
    const searchResults = document.getElementById('search-results') as HTMLDivElement;
    const selectedItems = document.getElementById('selected-items') as HTMLDivElement;

    if (!searchInput || !searchResults || !selectedItems) {
        console.error('Required elements not found');
        return;
    }

    let initPlantIds = searchInput.dataset.initSelectedIds;
    if (!initPlantIds) {
        console.error('initSelectedIds not found');
        return;
    }

    let plantIds = initPlantIds.split(',');
    plantIds = plantIds.map(id => id.trim());

    const search = new MultiSelectPlantSearch(searchInput, searchResults, selectedItems, plantIds);

    const updateButton = document.getElementById('update-button') as HTMLButtonElement;
    const nameInput = document.getElementById('name') as HTMLInputElement;
    const descriptionInput = document.getElementById('description') as HTMLTextAreaElement;
    
    updateButton.addEventListener('click', () => {
        const selectedItems = search.GetSelectedItems();
        const name = nameInput.value;
        const description = descriptionInput.value;

        let albumId = window.location.pathname.split('/')[3];

        let plantIdsToRemove: string[] = [];
        let plantIdsToAdd: string[] = [];
        for (const item of selectedItems) {
            if (plantIds.indexOf(item.id) == -1) {
                plantIdsToAdd.push(item.id);
            }
        }
        for (const item of plantIds) {
            if (selectedItems.findIndex(i => i.id === item) == -1) {
                plantIdsToRemove.push(item);
            }
        }

        // update name
        fetch(`/api/album/name/${albumId}`, {
            method: 'PUT',
            body: JSON.stringify({name: name}),
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => {
            if (!response.ok) {
                console.error(response);
            }
        });

        // update description
        fetch(`/api/album/description/${albumId}`, {
            method: 'PUT',
            body: JSON.stringify({description: description}),
            headers: {
                'Content-Type': 'application/json'
            }
        }).then(response => {
            if (!response.ok) {
                console.error(response);
            }
        });

        // remove plants
        for (const id of plantIdsToRemove) {
            fetch(`/api/album/remove/${albumId}`, {
                method: 'DELETE',
                body: JSON.stringify({plant_id: id}),
                headers: {
                    'Content-Type': 'application/json'
                }
            }).then(response => {
                if (!response.ok) {
                    console.error(response);
                }
            });
        }

        // add plants
        for (const id of plantIdsToAdd) {
            fetch(`/api/album/add/${albumId}`, {
                method: 'POST',
                body: JSON.stringify({plant_id: id}),
                headers: {
                    'Content-Type': 'application/json'
                }
            }).then(response => {
                if (!response.ok) {
                    console.error(response);
                }
            });
        }
        window.location.href = '/view/album/' + albumId;
    });
});