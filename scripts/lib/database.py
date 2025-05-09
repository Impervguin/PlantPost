import psycopg2 as pg
from lib.plant import Plant
import json

class PlantDatabase:
    def __init__(self, host, port, database, user, password):
        self.conn = pg.connect(host=host, port=port, database=database, user=user, password=password)
        self.conn.autocommit = True
    
    def Insert(self, plant : Plant):
        raise NotImplementedError
    
    def InsertMany(self, plants : [Plant]):
        raise NotImplementedError
    
    def AnalyzeQuery(self, query):
        raise NotImplementedError
    
    @staticmethod
    def ID() -> str:
        raise NotImplementedError
    
    @staticmethod
    def SelectPreambula():
        raise NotImplementedError

class JsonDatabase(PlantDatabase):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
    
    def Insert(self, plant : Plant):
        spec = plant.specification.ToDict()
        cur = self.conn.cursor()
        cur.execute("INSERT INTO file (id, name, url) VALUES (%s, %s, %s)", (plant.main_photo_id.hex, plant.main_photo_id.hex+'.jpg', plant.main_photo_id.hex+'.jpg'))
        cur.execute(
            "INSERT INTO plant (id, name, latin_name, description, category, main_photo_id, specification) VALUES (%s,%s, %s, %s, %s, %s, %s)",
            (plant.id.hex, plant.name, plant.latin_name, plant.description, plant.category, plant.main_photo_id.hex, json.dumps(spec))
        )

        cur.close()
    
    def InsertMany(self, plants : [Plant]):
        cur = self.conn.cursor()
        photo_query = "INSERT INTO file (id, name, url) VALUES "
        photo_vals = []
        for plant in plants[:-1]:
            photo_query += "(%s, %s, %s), "
            photo_vals.extend((plant.main_photo_id.hex, plant.main_photo_id.hex+'.jpg', plant.main_photo_id.hex+'.jpg'))
        plant = plants[-1]
        photo_query += "(%s, %s, %s);"
        photo_vals.extend((plant.main_photo_id.hex, plant.main_photo_id.hex+'.jpg', plant.main_photo_id.hex+'.jpg'))
        cur.execute(photo_query, photo_vals)

        query = "INSERT INTO plant (id, name, latin_name, description, category, main_photo_id, specification) VALUES "
        vals = []
        for plant in plants[:-1]:
            spec = plant.specification.ToDict()
            query += "(%s, %s, %s, %s, %s, %s, %s), "
            vals.extend((plant.id.hex, plant.name, plant.latin_name, plant.description, plant.category, plant.main_photo_id.hex, json.dumps(spec)))
        plant = plants[-1]
        spec = plant.specification.ToDict()
        query += "(%s, %s, %s, %s, %s, %s, %s);"
        vals.extend((plant.id.hex, plant.name, plant.latin_name, plant.description, plant.category, plant.main_photo_id.hex, json.dumps(spec)))
        cur.execute(query, vals)
        cur.close()
    
    @staticmethod
    def SelectPreambula():
        return """with p as (SELECT 
        p.id,
        p.name AS plant_name,
        p.latin_name,
        p.description,
        p.category,
        f.url AS main_photo_url,
        (p.specification->>'height_m')::float AS height_m,
        (p.specification->>'diameter_m')::float AS diameter_m,
        (p.specification->>'soil_acidity')::numeric AS soil_acidity,
        p.specification->>'soil_moisture' AS soil_moisture,
        p.specification->>'light_relation' AS light_relation,
        p.specification->>'soil_type' AS soil_type,
        (p.specification->>'winter_hardiness')::numeric AS winter_hardiness,
        p.specification->>'flowering_period' AS flowering_period
        FROM plant p
        JOIN file f ON p.main_photo_id = f.id)
        """

    def AnalyzeQuery(self, query):
        cur = self.conn.cursor()
        cur.execute("EXPLAIN ANALYZE " + self.SelectPreambula() + " " + query)
        analyzeRows = []
        for row in cur.fetchall():
            analyzeRows.append(row[0])
        return analyzeRows

    @staticmethod
    def ID():
        return "json"
    
class EavDatabase(PlantDatabase):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
    
    def Insert(self, plant : Plant):
        cur = self.conn.cursor()

        cur.execute(
            "INSERT INTO file (id, name, url) VALUES (%s, %s, %s)",
            (plant.main_photo_id.hex, plant.main_photo_id.hex+'.jpg', plant.main_photo_id.hex+'.jpg')
        )
    

        cur.execute(
            "INSERT INTO plant (id, name, latin_name, description, category_name, main_photo_id) VALUES (%s, %s, %s, %s, %s, %s);",
            (plant.id.hex, plant.name, plant.latin_name, plant.description, plant.category, plant.main_photo_id.hex)
        )

        spec = plant.specification.ToDict()

        for key, value in spec.items():
            cur.execute(
                "SELECT id, data_type FROM attribute WHERE name = %s;",
                (key,)
            )
            attr_result = cur.fetchone()
            if attr_result is None:
                raise ValueError('Attribute not found: ' + key)
            attr_id = attr_result[0]
            data_type = attr_result[1]
            if data_type == 'float':
                if not isinstance(value, float):
                    raise ValueError('Float attribute must be a float:' + key + ':' + str(value))
                cur.execute(
                    "INSERT INTO plant_attribute_value (plant_id, attribute_id, float_value, number_value, text_value) VALUES (%s, %s, %s, %s, %s);",
                    (plant.id.hex, attr_id, value, None, None)
                )
            elif data_type == 'number':
                if not isinstance(value, int):
                    raise ValueError('Number attribute must be a int:' + key + ':' + str(value))
                cur.execute(
                    "INSERT INTO plant_attribute_value (plant_id, attribute_id, float_value, number_value, text_value) VALUES (%s, %s, %s, %s, %s);",
                    (plant.id.hex, attr_id, None, value, None)
                )
            elif data_type in ['select', 'string']:
                if not isinstance(value, str):
                    raise ValueError('String attribute must be a string' + key + ':' + str(value))
                cur.execute(
                    "INSERT INTO plant_attribute_value (plant_id, attribute_id, float_value, number_value, text_value) VALUES (%s, %s, %s, %s, %s);",
                    (plant.id.hex, attr_id, None, None, value)
                )
            
            

        cur.close()
    def InsertMany(self, plants : [Plant]):
        cur = self.conn.cursor()
        photo_query = "INSERT INTO file (id, name, url) VALUES "
        photo_vals = []
        for plant in plants[:-1]:
            photo_query += "(%s, %s, %s), "
            photo_vals.extend((plant.main_photo_id.hex, plant.main_photo_id.hex+'.jpg', plant.main_photo_id.hex+'.jpg'))
        plant = plants[-1]
        photo_query += "(%s, %s, %s);"
        photo_vals.extend((plant.main_photo_id.hex, plant.main_photo_id.hex+'.jpg', plant.main_photo_id.hex+'.jpg'))
        cur.execute(photo_query, photo_vals)

        query = "INSERT INTO plant (id, name, latin_name, description, category_name, main_photo_id) VALUES "
        vals = []
        for plant in plants[:-1]:
            query += "(%s, %s, %s, %s, %s, %s), "
            vals.extend((plant.id.hex, plant.name, plant.latin_name, plant.description, plant.category, plant.main_photo_id.hex))
        plant = plants[-1]
        query += "(%s, %s, %s, %s, %s, %s);"
        vals.extend((plant.id.hex, plant.name, plant.latin_name, plant.description, plant.category, plant.main_photo_id.hex))
        cur.execute(query, vals)

        keys = set()
        for plant in plants:
            spec = plant.specification.ToDict()
            keys.update(spec.keys())

        for key in keys:
            cur.execute(
                "SELECT id, data_type FROM attribute WHERE name = %s;",
                (key,)
            )
            attr_result = cur.fetchone()
            if attr_result is None:
                raise ValueError('Attribute not found: ' + key)
            attr_id = attr_result[0]
            data_type = attr_result[1]
            query = "INSERT INTO plant_attribute_value (plant_id, attribute_id, float_value, number_value, text_value) VALUES "
            vals = []
            for plant in plants:
                plspec = plant.specification.ToDict()
                if key not in plspec:
                    continue
                value = plspec[key]
                query += "(%s, %s, %s, %s, %s), "
                if data_type == 'float':
                    if not isinstance(plspec[key], float):
                        raise ValueError('Float attribute must be a float:' + key + ':' + str(plspec[key]))
                    vals.extend((plant.id.hex, attr_id, plspec[key], None, None))
                elif data_type == 'number':
                    if not isinstance(plspec[key], int):
                        raise ValueError('Number attribute must be a int:' + key + ':' + str(plspec[key]))
                    vals.extend((plant.id.hex, attr_id, None, plspec[key], None))
                elif data_type in ['select', 'string']:
                    if not isinstance(plspec[key], str):
                        raise ValueError('String attribute must be a string' + key + ':' + str(plspec[key]))
                    vals.extend((plant.id.hex, attr_id, None, None, plspec[key]))
            query = ';'.join(query.rsplit(',', 1))
            cur.execute(query, vals)
    @staticmethod
    def SelectPreambula():
        return """with p as (SELECT 
        p.id,
        p.name AS plant_name,
        p.latin_name,
        p.description,
        p.category_name AS category,
        MAX(CASE WHEN a.name = 'height_m' THEN pav.float_value END) AS height_m,
        MAX(CASE WHEN a.name = 'diameter_m' THEN pav.float_value END) AS diameter_m,
        MAX(CASE WHEN a.name = 'soil_acidity' THEN pav.number_value END) AS soil_acidity,
        MAX(CASE WHEN a.name = 'soil_moisture' THEN pav.text_value END) AS soil_moisture,
        MAX(CASE WHEN a.name = 'light_relation' THEN pav.text_value END) AS light_relation,
        MAX(CASE WHEN a.name = 'soil_type' THEN pav.text_value END) AS soil_type,
        MAX(CASE WHEN a.name = 'winter_hardiness' THEN pav.number_value END) AS winter_hardiness,
        MAX(CASE WHEN a.name = 'flowering_period' THEN pav.text_value END) AS flowering_period
    FROM plant p
    LEFT JOIN plant_attribute_value pav ON p.id = pav.plant_id
    LEFT JOIN attribute a ON pav.attribute_id = a.id
    GROUP BY p.id, p.name, p.latin_name, p.description, p.category_name)"""

    def AnalyzeQuery(self, query):
        cur = self.conn.cursor()
        cur.execute("EXPLAIN ANALYZE " + self.SelectPreambula() + " " + query)
        analyzeRows = []
        for row in cur.fetchall():
            analyzeRows.append(row[0])
        return analyzeRows
    
    @staticmethod
    def ID():
        return "eav"
