import faker as f
import uuid

class Plant:
    def __init__(self, id, name, latin_name, description, category, main_photo_id, specification):
        if category not in ['coniferous', 'deciduous']:
            raise ValueError('Category must be either coniferous or deciduous')
        if not isinstance(specification, Specification):
            raise ValueError('Specification must be an instance of Specification')
        if not isinstance(main_photo_id, uuid.UUID):
            raise ValueError('Main photo id must be an instance of uuid.UUID')
        if not isinstance(id, uuid.UUID):
            raise ValueError('Id must be an instance of uuid.UUID')
        self.id = id
        self.name = name
        self.latin_name = latin_name
        self.description = description
        self.category = category
        self.main_photo_id = main_photo_id
        self.specification = specification
    
    def Fake(fakr=None):
        if fakr is None:
            fakr = f.Faker()
        id = uuid.uuid4()
        name = fakr.word() + " " + fakr.word()
        latin_name = name + "us"
        description = fakr.paragraph()
        category = fakr.random_element(['coniferous', 'deciduous'])
        main_photo_id = uuid.uuid4()
        if category == 'coniferous':
            specification = Coniferous.Fake(fakr)
        elif category == 'deciduous':
            specification = Deciduous.Fake(fakr)   
        return Plant(
            id,
            name,
            latin_name,
            description,
            category,
            main_photo_id,
            specification
        )

class Specification:
    def ToDict(self):
        return NotImplementedError
    
    @staticmethod
    def Fake(fakr=None):
        return NotImplementedError


class Coniferous(Specification):
    def __init__(self, height_m, diameter_m, soil_acidity, soil_moisture, light_relation, soil_type, winter_hardiness):
        self.height_m = height_m
        self.diameter_m = diameter_m
        self.soil_acidity = soil_acidity
        self.soil_moisture = soil_moisture
        self.light_relation = light_relation
        self.soil_type = soil_type
        self.winter_hardiness = winter_hardiness
    

    
    @staticmethod
    def Fake(fakr=None):
        if fakr is None:
            fakr = f.Faker()
        return Coniferous(
            HeigthM.Fake(fakr),
            DiameterM.Fake(fakr),
            SoilAcidity.Fake(fakr),
            SoilMoisture.Fake(fakr),
            LightRelation.Fake(fakr),
            SoilType.Fake(fakr),
            WinterHardiness.Fake(fakr)
        )
    
    def ToDict(self):
        return {
            'height_m': self.height_m,
            'diameter_m': self.diameter_m,
            'soil_acidity': self.soil_acidity,
            'soil_moisture': self.soil_moisture,
            'light_relation': self.light_relation,
            'soil_type': self.soil_type,
            'winter_hardiness': self.winter_hardiness
        }


class Deciduous(Specification):
    def __init__(self, height_m, diameter_m, soil_acidity, soil_moisture, light_relation, soil_type, winter_hardiness, flowering_period):
        self.height_m = height_m
        self.diameter_m = diameter_m
        self.soil_acidity = soil_acidity
        self.soil_moisture = soil_moisture
        self.light_relation = light_relation
        self.soil_type = soil_type
        self.winter_hardiness = winter_hardiness    
        self.flowering_period = flowering_period
    
    def Fake(fakr=None):
        if fakr is None:
            fakr = f.Faker()
        return Deciduous(
            HeigthM.Fake(fakr),
            DiameterM.Fake(fakr),
            SoilAcidity.Fake(fakr),
            SoilMoisture.Fake(fakr),
            LightRelation.Fake(fakr),
            SoilType.Fake(fakr),
            WinterHardiness.Fake(fakr),
            FloweringPeriod.Fake(fakr)
        )
    
    def ToDict(self):
        return {
            'height_m': self.height_m,
            'diameter_m': self.diameter_m,
            'soil_acidity': self.soil_acidity,
            'soil_moisture': self.soil_moisture,
            'light_relation': self.light_relation,
            'soil_type': self.soil_type,
            'winter_hardiness': self.winter_hardiness,
            'flowering_period': self.flowering_period
        }



class FloweringPeriod:
    @staticmethod
    def Fake(fakr=None):
        if fakr is None:
            fakr = f.Faker()
        return fakr.random_element([
            'spring',
            'summer',
            'autumn',
            'winter',
            'january',
            'february',
            'march',
            'april',
            'may',
            'june',
            'july',
            'august',
            'september',
            'october',
            'november',
            'december'
        ])

class HeigthM:
    @staticmethod
    def Fake(fakr=None):
        if fakr is None:
            fakr = f.Faker()
        return fakr.pyfloat(min_value=0, max_value=100)

class DiameterM:
    @staticmethod
    def Fake(fakr=None):
        if fakr is None:
            fakr = f.Faker()
        return fakr.pyfloat(min_value=0, max_value=100)

class SoilAcidity:
    @staticmethod
    def Fake(fakr=None):
        if fakr is None:
            fakr = f.Faker()
        return fakr.pyint(min_value=0, max_value=100)

class SoilMoisture:
    @staticmethod
    def Fake(fakr=None):
        if fakr is None:
            fakr = f.Faker()
        return fakr.random_element([
            "dry",
            "low",
            "medium",
            "high"
        ])

class LightRelation:
    @staticmethod
    def Fake(fakr=None):
        if fakr is None:
            fakr = f.Faker()
        return fakr.random_element([
            "light",
            "halfshadow",
            "shadow"
        ])

class SoilType:
    @staticmethod
    def Fake(fakr=None):
        if fakr is None:
            fakr = f.Faker()
        return fakr.random_element([
            "light",
            "medium",
            "heavy"
        ])

class WinterHardiness:
    @staticmethod
    def Fake(fakr=None):
        if fakr is None:
            fakr = f.Faker()
        return fakr.pyint(min_value=1, max_value=11)
    
