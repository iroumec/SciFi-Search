-- Creacion de tablas

CREATE TABLE IF NOT EXISTS Users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(20) NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Works (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content_type_id INT NOT NULL,
    unit BOOLEAN DEFAULT FALSE,
    saga_id INT
);

CREATE TABLE IF NOT EXISTS ContentTypes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

--CREATE TABLE AssociatedWorks (
--    id SERIAL PRIMARY KEY,
--    work_id INT REFERENCES Works(id)
--);

CREATE TABLE IF NOT EXISTS Review (
    id SERIAL PRIMARY KEY,
    user_id INT,
    work_id INT,
    score INT CHECK (score >= 1 AND score <= 10),
    review TEXT
);


-- Asignacion de foreign keys

ALTER TABLE Works ADD CONSTRAINT Works_ContentType
    FOREIGN KEY (content_type_id)
    REFERENCES ContentTypes(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE Works ADD CONSTRAINT Works_Works --saga
    FOREIGN KEY (saga_id)
    REFERENCES Works(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE Review ADD CONSTRAINT Review_Users
    FOREIGN KEY (user_id)
    REFERENCES Users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE Review ADD CONSTRAINT Review_Work
    FOREIGN KEY (work_id)
    REFERENCES Works(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Funciones + Triggers

CREATE OR REPLACE FUNCTION FN_TRIU_USERNAME()
RETURNS TRIGGER AS $$
    BEGIN

        IF EXISTS (SELECT 1 FROM Users WHERE username = NEW.username) THEN
            RAISE EXCEPTION 'Ya existe otro usuario con ese username.';
        END IF;

        RETURN NEW;

    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER TRIU_USERNAME
    BEFORE INSERT OR UPDATE ON Users
    FOR EACH ROW
        EXECUTE FUNCTION FN_TRIU_USERNAME();