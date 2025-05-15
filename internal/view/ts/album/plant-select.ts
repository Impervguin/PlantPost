interface SelectableItem {
    id: string;
    name: string;
    latin_name: string;
    selected?: boolean;
}

interface SelectedItem {
    id: string;
    name: string;
    latin_name: string;
}

const currentScript = document.currentScript;

export class MultiSelectPlantSearch {
    private inputElement: HTMLInputElement;
    private resultsContainer: HTMLDivElement;
    private selectedItemsContainer: HTMLDivElement;
    private debounceTimeout: number | null = null;
    private selectedItems: SelectableItem[] = [];
    private allItems: SelectableItem[] = [];

    constructor(inputElement: HTMLInputElement, searchContainer: HTMLDivElement, selectedItemsContainer: HTMLDivElement, selectedItemIds: string[]) {
        this.inputElement = inputElement;
        this.resultsContainer = searchContainer;
        this.selectedItemsContainer = selectedItemsContainer;

        if (!this.inputElement || !this.resultsContainer || !this.selectedItemsContainer) {
            console.error('Required elements not found');
            return;
        }

        this.setupEventListeners();
        this.loadInitialData(selectedItemIds);
    }

    private setupEventListeners(): void {
        this.inputElement.addEventListener('input', (e) => this.handleSearchInput(e));
        this.inputElement.addEventListener('focus', () => this.showResults());
        document.addEventListener('click', (e) => this.handleDocumentClick(e));
    }

    private async loadInitialData(initialSelectedItemIds: string[]): Promise<void> {
        // Load initial data if needed
        this.allItems = await this.plantSearchAPI('');
        for (const item of this.allItems) {
            if (initialSelectedItemIds.some(selected => selected === item.id)) {
                item.selected = true;
                this.selectedItems.push(item);
            }
        }
        this.updateSelectedDisplay();
    }

    private handleSearchInput(event: Event): void {
        const input = event.target as HTMLInputElement;
        const query = input.value.trim();

        if (this.debounceTimeout) {
            clearTimeout(this.debounceTimeout);
        }

        this.debounceTimeout = window.setTimeout(() => {
            if (query.length > 0) {
                this.performSearch(query);
            } else {
                this.showAllItems();
            }
        }, 300);
    }

    private async performSearch(query: string): Promise<void> {
        try {
            const results = await this.plantSearchAPI(query);
            this.displayResults(results);
            this.showResults();
        } catch (error) {
            console.error('Search failed:', error);
            this.hideResults();
        }
    }

    private async plantSearchAPI(query: string): Promise<SelectableItem[]> {
        let url: string;
        if (query === '') {
            url = `/api/search/plants`;
        } else {
            url = `/api/search/plants?name=${query}`;
        }
        const response = await fetch(url);
        const data = (await response.json())['plants'];
        for (const item of data) {
            item.selected = this.selectedItems.some(selected => selected.id === item.id);
        }
        return data;
    }

    private showAllItems(): void {
        this.displayResults(this.allItems.map(item => ({
            ...item,
            selected: this.selectedItems.some(selected => selected.id === item.id)
        })));
        this.showResults();
    }

    private displayResults(results: SelectableItem[]): void {
        this.resultsContainer.innerHTML = '';

        if (results.length === 0) {
            const noResults = document.createElement('div');
            noResults.className = 'px-4 py-2 text-sm text-gray-700';
            noResults.textContent = 'No results found';
            this.resultsContainer.appendChild(noResults);
            return;
        }

        results.forEach(item => {
            const resultItem = document.createElement('div');
            resultItem.className = 'flex items-center px-4 py-2 hover:bg-gray-50';
            
            const checkbox = document.createElement('input');
            checkbox.type = 'checkbox';
            checkbox.className = 'h-4 w-4 rounded-sm border-gray-300 mr-3 checked:accent-emerald-500 hover:border-emerald-500 hover:outline-offset-2 hover:outline-emerald-500';
            checkbox.checked = item.selected || false;
            checkbox.addEventListener('change', () => this.toggleItemSelection(item));

            const textContainer = document.createElement('div');
            textContainer.id = item.id;
            textContainer.className = 'flex-1 min-w-0';
            
            const title = document.createElement('div');
            title.className = 'text-sm font-medium text-gray-900 truncate';
            title.textContent = item.name;
            
            const latin_name = document.createElement('div');
            latin_name.className = 'text-xs text-gray-500 truncate';
            latin_name.textContent = item.latin_name;
            
            textContainer.appendChild(title);
            textContainer.appendChild(latin_name);

            resultItem.appendChild(checkbox);
            resultItem.appendChild(textContainer);
            this.resultsContainer.appendChild(resultItem);
        });
    }

    private toggleItemSelection(item: SelectableItem): void {
        const existingIndex = this.selectedItems.findIndex(i => i.id === item.id);
        
        if (existingIndex >= 0) {
            // Remove if already selected
            this.selectedItems.splice(existingIndex, 1);
        } else {
            // Add if not selected
            this.selectedItems.push({...item, selected: true});
        }

        this.updateSelectedDisplay();
        this.updateResultsDisplay();
    }

    private updateSelectedDisplay(): void {
        this.selectedItemsContainer.innerHTML = '';

        this.selectedItems.forEach(item => {
            const pill = document.createElement('div');
            pill.className = 'flex items-center bg-emerald-100 text-emerald-800 text-sm px-3 py-1 rounded-full';
            
            const text = document.createElement('span');
            text.className = 'mr-2';
            text.textContent = item.name;
            
            pill.appendChild(text);
            this.selectedItemsContainer.appendChild(pill);
        });
    }

    private updateResultsDisplay(): void {
        const checkboxes = this.resultsContainer.querySelectorAll<HTMLInputElement>('input[type="checkbox"]');
        checkboxes.forEach(checkbox => {
            const itemId = checkbox.closest('div')?.id;
            if (itemId) {
                checkbox.checked = this.selectedItems.some(item => item.id === itemId);
            }
        });
    }

    private showResults(): void {
        if (this.resultsContainer.children.length > 0) {
            this.resultsContainer.classList.remove('hidden');
        }
    }

    private hideResults(): void {
        this.resultsContainer.classList.add('hidden');
    }

    private handleDocumentClick(event: MouseEvent): void {
        const target = event.target as HTMLElement;
        if (!this.resultsContainer.contains(target)) {
            this.hideResults();
        }
    }

    public GetSelectedItems(): SelectedItem[] {
        return this.selectedItems;
    }
}