import { PostFilter } from './types.js';
import {
    PostTitleFilter,
    PostTagsFilter,
    PostAuthorFilter
} from './filters.js';

export class PostFilterParser {
    private static filterMap: Record<string, new () => PostFilter> = {
        'title': PostTitleFilter,
        'tags': PostTagsFilter,
        'author': PostAuthorFilter
    };

    static parseForm(form: HTMLFormElement): PostFilter[] {
        const formData = new FormData(form);
        const activeFilters: PostFilter[] = [];

        Object.keys(this.filterMap).forEach(filterType => {
            const FilterClass = this.filterMap[filterType];
            const filter = new FilterClass();
            if (filter.parse(formData)) {
                activeFilters.push(filter);
            }
        });
        
        return activeFilters;
    }

    static parseMultipleForms(forms: HTMLFormElement[]): PostFilter[] {
        const activeFilters: PostFilter[] = [];
        forms.forEach(form => {
            const filters = this.parseForm(form);
            activeFilters.push(...filters);
        });
        return activeFilters;
    }
}