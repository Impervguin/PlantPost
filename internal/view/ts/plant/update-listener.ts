import { PlantFieldParser } from './parser.js';

document.addEventListener('DOMContentLoaded', () => {
    const updateForm = document.getElementById('update-plant-form') as HTMLFormElement;
    console.log(updateForm);
    if (!updateForm) return;

    const updateButton = document.getElementById('update-button') as HTMLButtonElement;
    console.log(updateButton);
    if (!updateButton) return;

    // Get plant ID from URL
    // http://domain.com/view/plant/12345678-1234-1234-1234-123456789012/update
    const url = new URL(window.location.href);
    const pathParts = url.pathname.split('/');
    const plantId = pathParts[3]; 

    updateButton.addEventListener('click', () => {
        const fields = PlantFieldParser.parseForm(updateForm);
        console.log(fields);
        let postForm = new FormData();
        fields.forEach(field => field.toFormData(postForm));
        console.log(postForm);

        fetch(`/api/plant/specification/${plantId}`, {
            method: 'PUT',
            body: postForm
        }).then(response => {
            if (response.ok) {
                window.location.href = `/view/plant/${plantId}`;
            } else {
                console.error(response);
                throw new Error('Failed to update plant');
            }
        });
    });
});