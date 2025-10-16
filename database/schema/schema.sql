-- CREACIÓN DE TABLAS
CREATE TABLE IF NOT EXISTS users(
    username VARCHAR(16) PRIMARY KEY,
    email TEXT NOT NULL,
    name VARCHAR(32) NOT NULL,
    middlename VARCHAR(32),
    surname VARCHAR(32) NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS preferences (
    preference TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS user_preferences (
    username VARCHAR(16),
    preference TEXT,
    CONSTRAINT pk_user_preferences PRIMARY KEY (username,preference)
);

CREATE TABLE IF NOT EXISTS historic_searches (
    historic_search_id SERIAL PRIMARY KEY,
    username VARCHAR(16),
    search_string TEXT NOT NULL
);

-- ASIGNACION DE CLAVES FORÁNEAS 
ALTER TABLE user_preferences 
ADD CONSTRAINT fk_user_preferences_users 
FOREIGN KEY (username) 
REFERENCES users(username)
    ON UPDATE CASCADE
    ON DELETE CASCADE 
;

ALTER TABLE user_preferences 
ADD CONSTRAINT fk_user_preferences_preferences
FOREIGN KEY (preference)
REFERENCES preferences(preference)
    ON UPDATE CASCADE
    ON DELETE CASCADE 
;

ALTER TABLE historic_searches
ADD CONSTRAINT fk_historic_searches_users
FOREIGN KEY (username)
REFERENCES users(username)
    ON UPDATE CASCADE 
    ON DELETE CASCADE
;