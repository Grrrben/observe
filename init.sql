CREATE TABLE IF NOT EXISTS project(
    id SERIAL PRIMARY KEY,
    hash VARCHAR(80) UNIQUE,
    project_name VARCHAR(255),
    description VARCHAR(255),
    address VARCHAR(255),
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS observation (
    id SERIAL PRIMARY KEY,
    project_hash VARCHAR(80) REFERENCES project(hash),
    image VARCHAR(255),
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
