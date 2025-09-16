ALTER TABLE plant_category ADD COLUMN photo_id uuid;
ALTER TABLE plant_category ADD CONSTRAINT plant_category_photo_id_fkey FOREIGN KEY (photo_id) REFERENCES file(id);