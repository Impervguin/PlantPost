INSERT INTO plant_category (name, attributes, photo_id) VALUES
    ('coniferous', 
     '{
    "height_m": {
        "type": "float",
        "min": 0.0
    },
    "diameter_m": {
        "type": "float",
        "min": 0.0
    }, 
    "soil_acidity": {
        "type": "number",
        "min": 0.0
    },
    "soil_moisture": {
        "type": "string",
        "options": [
            "dry",
            "low",
            "medium",
            "high"
        ]
    },
    "light_relation": {
        "type": "string",
        "options": [
            "light",
            "halfshadow",
            "shadow"
        ]
    }, 
    "soil_type": {
        "type": "string",
        "options": [
            "light",
            "medium",
            "heavy"
            ]
    }, 
    "winter_hardiness": {
        "type": "number",
        "min": 1.0,
        "max": 11.0
    }
}', 
     NULL),
    ('deciduous', 
     '{
    "height_m": {
        "type": "float",
        "min": 0.0
    },
    "diameter_m": {
        "type": "float",
        "min": 0.0
    }, 
    "soil_acidity": {
        "type": "number",
        "min": 0.0
    },
    "soil_moisture": {
        "type": "string",
        "options": [
            "dry",
            "low",
            "medium",
            "high"
        ]
    },
    "light_relation": {
        "type": "string",
        "options": [
            "light",
            "halfshadow",
            "shadow"
        ]
    }, 
    "soil_type": {
        "type": "string",
        "options": [
            "light",
            "medium",
            "heavy"
            ]
    }, 
    "winter_hardiness": {
        "type": "number",
        "min": 1.0,
        "max": 11.0
    },
    "flowering_period": {
        "type": "string",
        "options": [
            "spring",
            "summer",
            "autumn",
            "winter",
            "january",
            "february",
            "march",
            "april",
            "may",
            "june",
            "july",
            "august",
            "september",
            "october",
            "november", 
            "december"
        ]
    }
}', NULL);
