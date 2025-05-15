import { PlantField } from './types.js';
import {
    PlantNameField,
    PlantLatinNameField,
    PlantDescriptionField,
    PlantHeightField,
    PlantDiameterField,
    PlantLightRelationField,
    PlantSoilTypeField,
    PlantSoilAcidityField,
    PlantSoilMoistureField,
    PlantCategoryField,
    PlantFloweringPeriodField,
    PlantWinterHardinessField,
    PlantMainPhotoField
} from './field.js';

export class PlantFieldParser {
    private static fieldMap: Record<string, new () => PlantField> = {
        'name': PlantNameField,
        'latin_name': PlantLatinNameField,
        'description': PlantDescriptionField,
        'category': PlantCategoryField,
        'file': PlantMainPhotoField,
    };

    private static specificationFieldMap: Record<string, new () => PlantField> = {
        'height': PlantHeightField,
        'diameter': PlantDiameterField,
        'flowering_period': PlantFloweringPeriodField,
        'winter_hardiness': PlantWinterHardinessField,
        'light_relation': PlantLightRelationField,
        'soil_type': PlantSoilTypeField,
        'soil_acidity': PlantSoilAcidityField,
        'soil_moisture': PlantSoilMoistureField,
    };

    static parseForm(form: HTMLFormElement): PlantField[] {
        const formData = new FormData(form);
        const activeFields: PlantField[] = [];

        Object.keys(this.fieldMap).forEach(filterType => {
            const FieldClass = this.fieldMap[filterType];
            const field = new FieldClass();
            if (field.parse(formData)) {
                activeFields.push(field);
            }
        });

        const cat = new PlantCategoryField();
        let catField: PlantField;
        catField = new PlantCategoryField();
        activeFields.forEach(obj => {
            if (obj.ID() === cat.ID()) {
                catField = obj;
                }
            }
        );

        if (!catField.Value()) return activeFields;

        let specData = new FormData();

        const category = catField.Value();
        formData.forEach((value, key) => {
            if (key.startsWith(category)) {
                specData.append(key.replace(`${category}.`, ''), value);
            }
        });

        console.log(specData);

        Object.keys(this.specificationFieldMap).forEach(fieldType => {
            const FieldClass = this.specificationFieldMap[fieldType];
            const field = new FieldClass();
            if (field.parse(specData)) {
                activeFields.push(field);
            }
        });
                

        return activeFields;
    }

    static parseMultipleForms(forms: HTMLFormElement[]): PlantField[] {
        const activeFields: PlantField[] = [];
        forms.forEach(form => {
            const fields = this.parseForm(form);
            activeFields.push(...fields);
        });
        return activeFields;
    }
}