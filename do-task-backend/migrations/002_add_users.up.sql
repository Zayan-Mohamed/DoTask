-- Add users table and update tasks table to include user_id

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add user_id column to tasks table
ALTER TABLE tasks ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;

-- Add index for user_id
CREATE INDEX idx_tasks_user_id ON tasks(user_id);

-- Update categories table to include user_id
ALTER TABLE categories ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;
CREATE INDEX idx_categories_user_id ON categories(user_id);

-- Insert a default user for existing data (with bcrypt hash for "password123")
INSERT INTO users (id, name, email, password) VALUES 
('00000000-0000-0000-0000-000000000001', 'John Doe', 'john@example.com', '$2b$10$fqiEgGOhpIVFkbdL0vAdN.Q0HCirE0fh.OJy26ryze5Zo7AH4EFbG');

-- Update existing tasks and categories to belong to the default user
UPDATE tasks SET user_id = '00000000-0000-0000-0000-000000000001' WHERE user_id IS NULL;
UPDATE categories SET user_id = '00000000-0000-0000-0000-000000000001' WHERE user_id IS NULL;
