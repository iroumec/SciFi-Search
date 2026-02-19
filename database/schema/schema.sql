-- ------------------------------------------------------------------------------------------------
-- TABLE CREATION
-- ------------------------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS users(
    user_id SERIAL,
    name VARCHAR(32) NOT NULL,
    surname VARCHAR(32) NOT NULL,
    auth_id TEXT UNIQUE NOT NULL,
    avatar_url TEXT UNIQUE,
    CONSTRAINT pk_user PRIMARY KEY (user_id)
);

-- ------------------------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS user_preferences (
    user_id INT,
    preference TEXT,
    CONSTRAINT pk_user_preferences PRIMARY KEY (user_id,preference)
);

-- ------------------------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS historic_searches (
    historic_search_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    search_string TEXT NOT NULL,
    search_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ------------------------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    user_id INT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    main_area TEXT NOT NULL,
    secondary_area TEXT,
    link TEXT,
    description TEXT,
    based_on TEXT,
    grantor TEXT,
    currency TEXT NOT NULL,
    amount TEXT NOT NULL,
    deadline TEXT NOT NULL
);

-- ------------------------------------------------------------------------------------------------
-- FOREIGN KEY ASIGNATION
-- ------------------------------------------------------------------------------------------------

ALTER TABLE user_preferences 
ADD CONSTRAINT fk_user_preferences_users 
FOREIGN KEY (user_id) 
REFERENCES users(user_id)
    ON UPDATE CASCADE
    ON DELETE CASCADE 
;

-- ------------------------------------------------------------------------------------------------

ALTER TABLE historic_searches
ADD CONSTRAINT fk_historic_searches_users
FOREIGN KEY (user_id)
REFERENCES users(user_id)
    ON UPDATE CASCADE 
    ON DELETE CASCADE
;

-- ------------------------------------------------------------------------------------------------

ALTER TABLE documents
ADD CONSTRAINT fk_documents_users
FOREIGN KEY (user_id)
REFERENCES users(user_id)
    ON UPDATE CASCADE 
    ON DELETE SET NULL
;

-- ------------------------------------------------------------------------------------------------