import { PlantFilterParser } from './parser.js';

document.addEventListener('DOMContentLoaded', () => {
    const forms = document.querySelectorAll<HTMLFormElement>('form.search-filters');
    const searchForms = Array.from(forms);
    
    const searchButton = document.getElementById('search-button');
    if (searchButton) {
        searchButton.addEventListener('click', () => {
            const filters = PlantFilterParser.parseMultipleForms(searchForms);
            console.log('Filters to send:', filters);

            let url = '/view/plants';
            if (filters.length > 0) {
                url += '?' + filters.map(f => f.toQueryString()).join('&');
            }
            window.location.href = url;
        });
    }
});