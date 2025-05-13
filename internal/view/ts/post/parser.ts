import { PostField } from './types.js';
import {
    PostTitleField,
    PostContentField,
    PostTagsField,
    PostPhotosField
} from './field.js';

export class PostFieldParser {
    private static fieldMap: Record<string, new () => PostField> = {
        'title': PostTitleField,
        'content': PostContentField,
        'tags': PostTagsField,
        'photos': PostPhotosField,
    };

    static parseForm(form: HTMLFormElement): PostField[] {
        const formData = new FormData(form);
        const activeFields: PostField[] = [];

        Object.keys(this.fieldMap).forEach(filterType => {
            const FieldClass = this.fieldMap[filterType];
            const field = new FieldClass();
            if (field.parse(formData)) {
                activeFields.push(field);
            }
        });     

        return activeFields;
    }

    static parseMultipleForms(forms: HTMLFormElement[]): PostField[] {
        const activeFields: PostField[] = [];
        forms.forEach(form => {
            const fields = this.parseForm(form);
            activeFields.push(...fields);
        });
        return activeFields;
    }
}