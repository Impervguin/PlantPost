import { PostFieldParser } from './parser.js';

document.addEventListener('DOMContentLoaded', () => {
    const createForm = document.getElementById('create-post-form') as HTMLFormElement;
    console.log(createForm);
    if (!createForm) return;

    const createButton = document.getElementById('create-button') as HTMLButtonElement;
    console.log(createButton);
    if (!createButton) return;

    createButton.addEventListener('click', () => {
        const fields = PostFieldParser.parseForm(createForm);
        console.log(fields);
        let postForm = new FormData();
        fields.forEach(field => field.toFormData(postForm));
        console.log(postForm);

        fetch('/api/post/create', {
            method: 'POST',
            body: postForm
        }).then(response => {
            if (response.ok) {
                window.location.href = '/view/posts';
            } else {
                console.error(response);
                throw new Error('Failed to create post');
            }
        });
    });
});