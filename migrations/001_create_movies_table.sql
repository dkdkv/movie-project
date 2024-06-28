-- migrations/001_create_movies_table.sql
CREATE TABLE IF NOT EXISTS movies (
                                      id SERIAL PRIMARY KEY,
                                      title VARCHAR(255) NOT NULL,
                                      director VARCHAR(255) NOT NULL,
                                      release_date DATE,
                                      genre VARCHAR(100),
                                      rating DECIMAL(3,1),
                                      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                      updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                      deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_movies_title ON movies(title);
CREATE INDEX idx_movies_director ON movies(director);
CREATE INDEX idx_movies_release_date ON movies(release_date);