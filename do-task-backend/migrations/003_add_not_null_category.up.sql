-- Add NOT NULL constraint to category_id in tasks table
ALTER TABLE tasks
ALTER COLUMN category_id SET NOT NULL;
