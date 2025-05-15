import { PostFieldParser } from './parser.js';

document.addEventListener('DOMContentLoaded', () => {
    const updateForm = document.getElementById('update-post-form') as HTMLFormElement;
    console.log(updateForm);
    if (!updateForm) return;

    const updateButton = document.getElementById('update-button') as HTMLButtonElement;
    console.log(updateButton);
    if (!updateButton) return;

    updateButton.addEventListener('click', () => {
        const fields = PostFieldParser.parseForm(updateForm);
        console.log(fields);
        let postForm = new FormData();
        fields.forEach(field => field.toFormData(postForm));
        console.log(postForm);

        fetch('/api/post/text/' + window.location.pathname.split('/')[3], {
            method: 'PUT',
            body: postForm
        }).then(response => {
            if (response.ok) {
                window.location.href = '/view/posts';
            } else {
                console.error(response);
                throw new Error('Failed to update post');
            }
        });
    });
});