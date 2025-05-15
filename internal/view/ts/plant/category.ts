document.addEventListener('DOMContentLoaded', () => {
    const categorySelect = document.getElementById('category') as HTMLSelectElement;

    if (!categorySelect) return;

    const categoryOptions = categorySelect.querySelectorAll('option');

    categorySelect.addEventListener('change', () => {
        const category = categorySelect.value;
        console.log(category);
        console.log(categoryOptions);

        for (let i = 0; i < categoryOptions.length; i++) {
            const specficationSection = document.getElementById(`${categoryOptions.item(i).value}-spec`) as HTMLDivElement;

            if (!specficationSection) continue;
            const option = categoryOptions[i];

            if (option.value === category) {
                specficationSection.classList.remove('hidden');
            } else {
                specficationSection.classList.add('hidden');
            }
        }
    });
});