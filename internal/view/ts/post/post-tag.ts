class TagManager {
    private container: HTMLElement;
  
    constructor(containerId: string) {
      this.container = document.getElementById(containerId) as HTMLElement;
      if (!this.container) {
        throw new Error(`Container element with ID ${containerId} not found`);
      }
  
      this.init();
    }
  
    private init(): void {
      // Add initial tag create button
      const button = this.container.querySelector('#create-tag-button') as HTMLButtonElement;
      if (!button) {
        this.addTagButton();
      } else {
        button.addEventListener('click', () => this.replaceWithTagInput(button));
      }

      // Add delete tag buttons
      const deleteButtons = this.container.querySelectorAll('.delete-tag-button') as NodeListOf<HTMLButtonElement>;
      if (deleteButtons.length > 0) {
        deleteButtons.forEach(button => button.addEventListener('click', () => this.removeTagField(button)));
      }
    }
  
    private createTagInput(value: string = ''): HTMLDivElement {
      const div = document.createElement('div');
      div.className = 'relative transition-all duration-200 ease-out opacity-0 animate-in fade-in relative group h-[42px]';
      
      div.innerHTML = `
        <input 
          type="text" 
          name="tags[]" 
          value="${value}"
          class="block h-full w-full rounded-md border-gray-300 shadow-sm focus:border-emerald-500 focus:ring-emerald-500 sm:text-sm pr-8 transition-all duration-200"
        >
        <button 
          type="button"
          class="absolute right-1 top-1/2 -translate-y-1/2 text-gray-400 hover:text-emerald-600 transition-colors duration-200"
        >
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
            <path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" />
          </svg>
        </button>
      `;
  
      const button = div.querySelector('button') as HTMLButtonElement;
      button.addEventListener('click', () => this.removeTagField(button));
  
      setTimeout(() => {
        div.classList.remove('opacity-0');
      }, 10);
  
      return div;
    }
  
    private addTagButton(): void {
      const buttonDiv = document.createElement('div');
      buttonDiv.className = 'transition-all duration-200 ease-out opacity-0 animate-in fade-in h-[42px]';
      
      buttonDiv.innerHTML = `
        <button 
          type="button"
          id="create-tag-button"
          class="w-full h-full flex items-center justify-center rounded-md border-2 border-dashed border-gray-300 hover:border-emerald-500 text-gray-400 hover:text-emerald-600 p-2 transition-all duration-200 hover:scale-[1.02]"
        >
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5 mr-1">
            <path d="M10.75 4.75a.75.75 0 00-1.5 0v4.5h-4.5a.75.75 0 000 1.5h4.5v4.5a.75.75 0 001.5 0v-4.5h4.5a.75.75 0 000-1.5h-4.5v-4.5z" />
          </svg>
          Add Tag
        </button>
      `;
  
      const button = buttonDiv.querySelector('button') as HTMLButtonElement;
      button.addEventListener('click', () => this.replaceWithTagInput(button));
  
      this.container.appendChild(buttonDiv);
      
      setTimeout(() => {
        buttonDiv.classList.remove('opacity-0');
      }, 10);
    }
  
    private replaceWithTagInput(button: HTMLButtonElement): void {
      const buttonDiv = button.parentElement as HTMLDivElement;
      buttonDiv.classList.add('opacity-0', 'scale-95');
      
      setTimeout(() => {
        const inputDiv = this.createTagInput();
        this.container.insertBefore(inputDiv, buttonDiv);
        buttonDiv.remove();
        this.addTagButton();
        (inputDiv.querySelector('input') as HTMLInputElement).focus();
      }, 200); // Match with transition duration
    }
  
    private removeTagField(button: HTMLButtonElement): void {
      const inputDiv = button.closest('.relative') as HTMLDivElement;
      inputDiv.classList.add('opacity-0', 'scale-95');
      
      setTimeout(() => {
        inputDiv.remove();
        
        if (this.container.children.length === 0) {
          this.addTagButton();
        }
      }, 200); // Match with transition duration
    }
  }
  
  document.addEventListener('DOMContentLoaded', () => {
    new TagManager('tags-input-container');
  });