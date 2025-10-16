-- CREACIÓN DE TABLAS
CREATE TABLE IF NOT EXISTS users(
    user_id INT PRIMARY KEY,
    --username VARCHAR(16),
    --email TEXT NOT NULL,
    name VARCHAR(32) NOT NULL,
    middlename VARCHAR(32),
    surname VARCHAR(32) NOT NULL,
    --password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS preferences (
    preference TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS user_preferences (
    user_id INT,
    preference TEXT,
    CONSTRAINT pk_user_preferences PRIMARY KEY (user_id,preference)
);

CREATE TABLE IF NOT EXISTS historic_searches (
    historic_search_id SERIAL PRIMARY KEY,
    user_id INT,
    search_string TEXT NOT NULL
);

-- ASIGNACION DE CLAVES FORÁNEAS 
ALTER TABLE user_preferences 
ADD CONSTRAINT fk_user_preferences_users 
FOREIGN KEY (user_id) 
REFERENCES users(user_id)
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
FOREIGN KEY (user_id)
REFERENCES users(user_id)
    ON UPDATE CASCADE 
    ON DELETE CASCADE
;