CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE admission_analyses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    university_id UUID NOT NULL,
    admission_chance_percentage DECIMAL(5,2) DEFAULT 0,
    weak_areas TEXT[] DEFAULT '{}',
    improvement_suggestions TEXT[] DEFAULT '{}',
    gap_analysis TEXT DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE roadmaps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE roadmap_steps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    roadmap_id UUID REFERENCES roadmaps(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    month_number INTEGER DEFAULT 1,
    category VARCHAR(100) DEFAULT ''
);

CREATE TABLE essay_reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    essay_text TEXT NOT NULL,
    university_name VARCHAR(255) DEFAULT '',
    program_name VARCHAR(255) DEFAULT '',
    grammar_score DECIMAL(4,2) DEFAULT 0,
    coherence_score DECIMAL(4,2) DEFAULT 0,
    uniqueness_score DECIMAL(4,2) DEFAULT 0,
    relevance_score DECIMAL(4,2) DEFAULT 0,
    improvement_suggestions TEXT[] DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW()
);