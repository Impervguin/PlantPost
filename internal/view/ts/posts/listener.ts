import { PostFilterParser } from './parser.js';

document.addEventListener('DOMContentLoaded', () => {
    const searchForm = document.getElementById('search-filters') as HTMLFormElement;
    console.log(searchForm);
    if (!searchForm) return;

    const searchButton = document.getElementById('search-button') as HTMLButtonElement;
    console.log(searchButton);
    if (!searchButton) return;

    searchButton.addEventListener('click', () => {
        const filters = PostFilterParser.parseForm(searchForm);

        let url = '/view/posts';
        if (filters.length > 0) {
            url += '?' + filters.map(f => f.toQueryString()).join('&');
        }
        window.location.href = url;
    });
});