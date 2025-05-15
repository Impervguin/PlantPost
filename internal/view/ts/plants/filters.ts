import { PlantFilter } from './types.js';


export abstract class MinMaxFilter extends PlantFilter {
    abstract name: string;
    params: { min: number; max: number } = { min: 0, max: 0 };
    toQueryString(): string {
        return `${this.type}=${this.params.min}-${this.params.max}`;
    }

    parse(formData: FormData): boolean {
        const min = this.getNumberValue(formData, `${this.name}-min`);
        const max = this.getNumberValue(formData, `${this.name}-max`);
        if (min === null || max === null) return false;
        this.params.min = min;
        this.params.max = max;
        return true;
    }
}

export abstract class StringFilter extends PlantFilter {
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

export abstract class OptionArrayFilter extends PlantFilter {
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

export class PlantNameFilter extends StringFilter {
    type: string = 'name';
    name = 'name';
}

export class PlantLatinNameFilter extends StringFilter {
    type: string = 'latin_name'; // var for query string
    name = 'latin-name'; // var for form data
}

export class PlantHeightFilter extends MinMaxFilter {
    type: string = 'height';
    name = 'height';
}

export class PlantDiameterFilter extends MinMaxFilter {
    type: string = 'diameter';
    name = 'diameter';
}

export class PlantLightRelationFilter extends OptionArrayFilter { // var for query string
    type: string = 'light_relation'; // var for form data
    name = 'light-relation';
}

export class PlantSoilTypeFilter extends OptionArrayFilter {
    type: string = 'soil_type';
    name = 'soil-type';
}

export class PlantSoilAcidityFilter extends MinMaxFilter {
    type: string = 'soil_acidity';
    name = 'soil-acidity';
}

export class PlantSoilMoistureFilter extends OptionArrayFilter {
    type: string = 'soil_moisture';
    name = 'soil-moisture';
}

export class PlantCategoryFilter extends OptionArrayFilter {
    type: string = 'category';
    name = 'category';
}

export class PlantFloweringPeriodFilter extends OptionArrayFilter {
    type: string = 'flowering_period';
    name = 'flowering-period';
}

export class PlantWinterHardinessFilter extends MinMaxFilter {
    type: string = 'winter_hardiness';
    name = 'winter-hardiness';
}


