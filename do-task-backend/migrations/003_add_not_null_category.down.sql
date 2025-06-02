-- Remove NOT NULL constraint from category_id in tasks table
ALTER TABLE tasks
ALTER COLUMN category_id DROP NOT NULL;
