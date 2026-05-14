CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE universities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    country VARCHAR(100) NOT NULL,
    acceptance_rate DECIMAL(5,2) DEFAULT 0,
    type VARCHAR(50) DEFAULT '',
    category VARCHAR(50) DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE grants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    amount DECIMAL(10,2) DEFAULT 0,
    deadline VARCHAR(50) DEFAULT '',
    country VARCHAR(100) DEFAULT '',
    min_gpa DECIMAL(3,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE saved_universities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    university_id UUID REFERENCES universities(id) ON DELETE CASCADE,
    saved_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE saved_grants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    grant_id UUID REFERENCES grants(id) ON DELETE CASCADE,
    saved_at TIMESTAMP DEFAULT NOW()
);