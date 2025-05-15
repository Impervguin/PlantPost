export abstract class PostFilter {
    abstract type: string;
    abstract params: Record<string, any>;
    
    abstract parse(formData: FormData): boolean;

    abstract toQueryString(): string;

    protected getStringValue(formData: FormData, field: string): string | null {
        const value = formData.get(field);
        return value ? value.toString() : null;
    }

    protected getStringArray(formData: FormData, field: string): string[] {
        return formData.getAll(field).map(v => v.toString());
    }
}