import { PlantFilter, PlantFilterJSON } from './types.js';
import {
    PlantNameFilter,
    PlantLatinNameFilter,
    PlantHeightFilter,
    PlantLightRelationFilter,
    PlantSoilTypeFilter,
    PlantSoilAcidityFilter,
    PlantSoilMoistureFilter,
    PlantCategoryFilter,
    PlantFloweringPeriodFilter,
    PlantWinterHardinessFilter,
    PlantDiameterFilter
} from './filters.js';

export class PlantFilterParser {
    private static filterMap: Record<string, new () => PlantFilter> = {
        'name': PlantNameFilter,
        'latin_name': PlantLatinNameFilter,
        'height': PlantHeightFilter,
        'light_relation': PlantLightRelationFilter,
        'soil_type': PlantSoilTypeFilter,
        'soil_acidity': PlantSoilAcidityFilter,
        'soil_moisture': PlantSoilMoistureFilter,
        'category': PlantCategoryFilter,
        'flowering_period': PlantFloweringPeriodFilter,
        'winter_hardiness': PlantWinterHardinessFilter,
        'diameter': PlantDiameterFilter,
    };

    static parseForm(form: HTMLFormElement): PlantFilter[] {
        const formData = new FormData(form);
        const activeFilters: PlantFilter[] = [];

        Object.keys(this.filterMap).forEach(filterType => {
            const FilterClass = this.filterMap[filterType];
            const filter = new FilterClass();
            if (filter.parse(formData)) {
                activeFilters.push(filter);
            }
        });
        
        return activeFilters;
    }

    static parseMultipleForms(forms: HTMLFormElement[]): PlantFilter[] {
        const activeFilters: PlantFilter[] = [];
        forms.forEach(form => {
            const filters = this.parseForm(form);
            activeFilters.push(...filters);
        });
        return activeFilters;
    }
}