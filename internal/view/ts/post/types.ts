
export const StringType = 'string';
export const StringArrayType = 'string[]';
export const ImageType = 'image';
export const ImageArrayType = 'image[]';

export abstract class PostField {
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
        switch (this.valueType) {
            case StringType:
                this.value = formData.get(this.id);
                if (!this.value) return false;
                this.value = this.value.toString();
                break;
            case StringArrayType:
                this.value = formData.getAll(this.id).map(v => v.toString());
                if (this.value.length === 0) return false;
                break;
            case ImageType:
                this.value = formData.get(this.id);
                if (!this.value) return false;
                break;
            case ImageArrayType:
                this.value = formData.getAll(this.id);
                if (this.value.length === 0) return false;
                break;
            default:
                throw new Error(`invalid value type: ${this.valueType}`);
        }
        return true;
    }

    toFormData(formData: FormData) {
        switch (this.valueType) {
            case StringArrayType:
                for (let i = 0; i < this.value.length; i++) {
                    formData.append(`${this.name}`, this.value[i]);
                }
                break;
            case ImageArrayType:
                for (let i = 0; i < this.value.length; i++) {
                    formData.append(`${this.name}`, this.value[i]);
                }
                break;
            default:
                formData.append(this.name, this.value);
        }
    }
}
