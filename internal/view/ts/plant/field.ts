import { PlantField, PlantSpecificationField, StringType, FloatType, IntType, ImageType } from './types.js';

export class PlantNameField extends PlantField {
    id: string = 'name';
    name: string = 'name';
    valueType: string = StringType;
}

export class PlantLatinNameField extends PlantField {
    id: string = 'latin-name';
    name: string = 'latin_name';
    valueType: string = StringType;
}

export class PlantDescriptionField extends PlantField {
    id: string = 'description';
    name: string = 'description';
    valueType: string = StringType;
}

export class PlantHeightField extends PlantSpecificationField {
    id: string = 'height';
    name: string = 'height_m';
    valueType: string = FloatType;
}

export class PlantDiameterField extends PlantSpecificationField {
    id: string = 'diameter';
    name: string = 'diameter_m';
    valueType: string = FloatType;
}

export class PlantLightRelationField extends PlantSpecificationField {
    id: string = 'light-relation';
    name: string = 'light_relation';
    valueType: string = StringType;
}

export class PlantSoilTypeField extends PlantSpecificationField {
    id: string = 'soil-type';
    name: string = 'soil_type';
    valueType: string = StringType;
}

export class PlantSoilAcidityField extends PlantSpecificationField {
    id: string = 'soil-acidity';
    name: string = 'soil_acidity';
    valueType: string = IntType;
}

export class PlantSoilMoistureField extends PlantSpecificationField {
    id: string = 'soil-moisture';
    name: string = 'soil_moisture';
    valueType: string = StringType;
}

export class PlantCategoryField extends PlantField {
    id: string = 'category';
    name: string = 'category';
    valueType: string = StringType;
}

export class PlantFloweringPeriodField extends PlantSpecificationField {
    id: string = 'flowering-period';
    name: string = 'flowering_period';
    valueType: string = StringType;
}

export class PlantWinterHardinessField extends PlantSpecificationField {
    id: string = 'winter-hardiness';
    name: string = 'winter_hardiness';
    valueType: string = IntType;
}

export class PlantMainPhotoField extends PlantField {
    id: string = 'photo';
    name: string = 'file';
    valueType: string = ImageType;
}