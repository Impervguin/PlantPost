import { PostFilter } from './types.js';

export abstract class StringFilter extends PostFilter {
    abstract name: string;
    params: { [key: string]: string } = {};

    toQueryString(): string {
        if (Object.keys(this.params).length !== 1) {
            throw new Error('StringFilter must have exactly one key');
        }
        const key = Object.keys(this.params)[0];
        return `${this.type}=${this.params[key]}`;
    }

    parse(formData: FormData): boolean {
        const value = this.getStringValue(formData, this.name);
        if (!value) return false;
        this.params[this.name] = value;
        return true;
    }
}

export abstract class OptionArrayFilter extends PostFilter {
    abstract name: string;
    params: { possibleValues: string[] } = { possibleValues: [] };
    toQueryString(): string {
        return `${this.type}=${this.params.possibleValues.join(',')}`;
    }

    parse(formData: FormData): boolean {
        const possibleValues = this.getStringArray(formData, `${this.name}[]`);
        if (possibleValues.length === 0) return false;
        this.params.possibleValues = possibleValues;
        return true;
    }
}

export class PostTitleFilter extends StringFilter {
    type: string = 'title';
    name = 'title';
}

export class PostTagsFilter extends OptionArrayFilter {
    type: string = 'tags';
    name = 'tags';
}

export class PostAuthorFilter extends StringFilter {
    type: string = 'author';
    name = 'author';
}