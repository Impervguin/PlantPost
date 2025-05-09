export interface PlantFilterJSON {
    type: string;
    params: Record<string, any>;
}

export abstract class PlantFilter {
    abstract type: string;
    abstract params: Record<string, any>;
    
    abstract parse(formData: FormData): boolean;

    abstract toQueryString(): string;
    
    toJSON(): PlantFilterJSON {
        return {
            type: this.type,
            params: this.params
        };
    }

    protected getStringValue(formData: FormData, field: string): string | null {
        const value = formData.get(field);
        return value ? value.toString() : null;
    }
    
    protected getNumberValue(formData: FormData, field: string): number | null {
        const value = formData.get(field);
        return value ? parseFloat(value.toString()) : null;
    }
    
    protected getStringArray(formData: FormData, field: string): string[] {
        return formData.getAll(field).map(v => v.toString());
    }
}