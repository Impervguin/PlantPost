CREATE OR REPLACE FUNCTION validate_plant_specification()
RETURNS TRIGGER AS $$
DECLARE
    category_attrs JSONB;
    attr_name TEXT;
    attr_def JSONB;
    attr_value JSONB;
    attr_type TEXT;
    min_val NUMERIC;
    max_val NUMERIC;
    options TEXT[];
    option TEXT;
    valid_option BOOLEAN;
    num_value NUMERIC;
    spec_attr TEXT;
    allowed_attrs TEXT[];
    string_value TEXT;
BEGIN
    SELECT attributes INTO category_attrs
    FROM plant_category
    WHERE name = NEW.category;
    
    IF category_attrs IS NULL THEN
        RAISE EXCEPTION 'Plant category "%" does not exist', NEW.category;
    END IF;

    SELECT array_agg(key) INTO allowed_attrs
    FROM jsonb_object_keys(category_attrs) AS key;

    FOR spec_attr IN SELECT jsonb_object_keys(NEW.specification)
    LOOP
        IF NOT spec_attr = ANY(allowed_attrs) THEN
            RAISE EXCEPTION 'Specification contains extra attribute "%" not defined in category "%"', 
                           spec_attr, NEW.category;
        END IF;
    END LOOP;
    
    FOR attr_name IN SELECT jsonb_object_keys(category_attrs)
    LOOP
        attr_def := category_attrs->attr_name;
        attr_value := NEW.specification->attr_name;
        
        IF attr_value IS NULL THEN
            RAISE EXCEPTION 'Missing required attribute "%" in specification', attr_name;
        END IF;
        
        attr_type := attr_def->>'type';
        
        -- Validate based on type
        CASE attr_type
            WHEN 'float', 'number' THEN
                -- Check if value is numeric
                IF jsonb_typeof(attr_value) != 'number' THEN
                    RAISE EXCEPTION 'Attribute "%" must be a number, got %', 
                                   attr_name, jsonb_typeof(attr_value);
                END IF;
                
                BEGIN
                    num_value := (attr_value::TEXT)::NUMERIC;
                EXCEPTION WHEN OTHERS THEN
                    RAISE EXCEPTION 'Attribute "%" must be a valid number', attr_name;
                END;
                
                IF attr_type = 'number' AND num_value % 1 != 0 THEN
                    RAISE EXCEPTION 'Attribute "%" must be an integer, got %', 
                                   attr_name, attr_value;
                END IF;
                
                IF attr_def ? 'min' THEN
                  min_val := (attr_def->>'min')::NUMERIC;
                  IF num_value < min_val THEN
                      RAISE EXCEPTION 'Attribute "%" must be at least %, got %', 
                                     attr_name, min_val, attr_value;
                  END IF;
                END IF;
                
                IF attr_def ? 'max' THEN
                    max_val := (attr_def->>'max')::NUMERIC;
                    IF num_value > max_val THEN
                        RAISE EXCEPTION 'Attribute "%" must be at most %, got %', 
                                       attr_name, max_val, attr_value;
                    END IF;
                END IF;
                
            WHEN 'string' THEN
                IF jsonb_typeof(attr_value) != 'string' THEN
                    RAISE EXCEPTION 'Attribute "%" must be a string, got %', 
                                   attr_name, jsonb_typeof(attr_value);
                END IF;
                
                string_value := attr_value #>> '{}';
                
                IF attr_def ? 'options' THEN
                  options := ARRAY(SELECT jsonb_array_elements_text(attr_def->'options'));
                  valid_option := FALSE;

                  FOREACH option IN ARRAY options LOOP
                      IF option = string_value THEN
                          valid_option := TRUE;
                          EXIT;
                      END IF;
                  END LOOP;

                  IF NOT valid_option THEN
                      RAISE EXCEPTION 'Attribute "%" value "%" is not in allowed options: %', 
                                     attr_name, string_value, options;
                  END IF;
                END IF;
                
            ELSE
                RAISE EXCEPTION 'Unknown attribute type "%" for attribute "%"', 
                               attr_type, attr_name;
        END CASE;
    END LOOP;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER plant_specification_validation
BEFORE INSERT OR UPDATE ON plant
FOR EACH ROW EXECUTE FUNCTION validate_plant_specification();