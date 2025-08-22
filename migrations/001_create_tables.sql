-- Create cats table
CREATE TABLE IF NOT EXISTS cats (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    years_experience INTEGER NOT NULL CHECK (years_experience >= 0),
    breed VARCHAR(255) NOT NULL,
    salary DECIMAL(10,2) NOT NULL CHECK (salary >= 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Create missions table
CREATE TABLE IF NOT EXISTS missions (
    id SERIAL PRIMARY KEY,
    cat_id INTEGER REFERENCES cats(id) ON DELETE SET NULL,
    complete BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Create targets table
CREATE TABLE IF NOT EXISTS targets (
    id SERIAL PRIMARY KEY,
    mission_id INTEGER NOT NULL REFERENCES missions(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    notes TEXT,
    complete BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_cats_deleted_at ON cats(deleted_at);
CREATE INDEX IF NOT EXISTS idx_missions_cat_id ON missions(cat_id);
CREATE INDEX IF NOT EXISTS idx_missions_deleted_at ON missions(deleted_at);
CREATE INDEX IF NOT EXISTS idx_targets_mission_id ON targets(mission_id);
CREATE INDEX IF NOT EXISTS idx_targets_deleted_at ON targets(deleted_at);