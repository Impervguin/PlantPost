
export const StringType = 'string';
export const FloatType = 'float';
export const IntType = 'int';
export const ImageType = 'image';

export abstract class PlantField {
    public abstract id: string; // var for input parsing
    abstract name: string; // var for output form data
    abstract valueType: string;
    value: any;

    ID(): string {
        return this.id;
    }

    Value(): any {
        return this.value;
    }

    parse(formData: FormData): boolean {
        const value = formData.get(this.id);
        if (!value) return false;
        switch (this.valueType) {
            case StringType:
                this.value = value.toString();
                break;
            case FloatType:
                this.value = parseFloat(value.toString());
                break;
            case IntType:
                this.value = parseInt(value.toString());
                break;
            case ImageType:
                this.value = value;
                break;
            default:
                throw new Error(`invalid value type: ${this.valueType}`);
        }
        return true;
    }

    toFormData(formData: FormData) {
        formData.append(this.name, this.value);
    }
}

export abstract class PlantSpecificationField extends PlantField {
    toFormData(formData: FormData)  {
        let specJson = formData.get('specification')?.toString();
        if (!specJson) specJson = "{}";
        let spec = JSON.parse(specJson);
        console.log(spec);
        spec[this.name] = this.value;
        formData.delete('specification');
        formData.append('specification', JSON.stringify(spec));
    }
}
