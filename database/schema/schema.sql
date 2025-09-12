--------------------- Creacion de tablas ---------------------

CREATE TABLE IF NOT EXISTS Users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(20) NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    --password?
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
--ContentTypes: Book, Movie, TV Show, Song, Game
--              Book saga, Movie saga, Related tv shows, Album, Game saga

CREATE TABLE IF NOT EXISTS ConsumedWorks (
    user_id INT, 
    work_id INT,
    CONSTRAINT pk_consumedworks PRIMARY KEY (user_id,work_id)
)

CREATE TABLE IF NOT EXISTS LikedWorks (
    user_id INT, 
    work_id INT,
    CONSTRAINT pk_likedworks PRIMARY KEY (user_id,work_id)
)

CREATE TABLE IF NOT EXISTS Review (
    id SERIAL PRIMARY KEY, --puede haber más de una review de la misma obra por el mismo usuario
    user_id INT NOT NULL,
    work_id INT NOT NULL,
    score INT CHECK (score >= 1 AND score <= 10) NOT NULL,
    review TEXT,
    when_watched TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    liked BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS ReviewLike (
    review_id INT, 
    user_id INT,
    when_liked TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_reviewlike PRIMARY KEY (review_id,user_id)
);

CREATE TABLE IF NOT EXISTS ReviewComment (
    id SERIAL PRIMARY KEY,
    review_id INT,
    user_id INT,
    comment VARCHAR(255),
    when_commented TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

--falta queries.sql a partir de aca
CREATE TABLE IF NOT EXISTS UserFollows (
    id_follower INT, 
    id_followed INT,
    when_followed TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_userfollows PRIMARY KEY (id_follower,id_followed)
);

CREATE TABLE IF NOT EXISTS UserFavorites (
    id_user INT,
    id_work INT,
    CONSTRAINT pk_userfavorites PRIMARY KEY (id_user,id_work)
);

--------------------- Asignacion de foreign keys ---------------------

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

ALTER TABLE ConsumedWorks ADD CONSTRAINT ConsumedWorks_Users
    FOREIGN KEY (user_id)
    REFERENCES Users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE ConsumedWorks ADD CONSTRAINT ConsumedWorks_Works
    FOREIGN KEY (work_id)
    REFERENCES Works(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE LikedWorks ADD CONSTRAINT LikedWorks_Users
    FOREIGN KEY (user_id)
    REFERENCES Users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE LikedWorks ADD CONSTRAINT LikedWorks_Users
    FOREIGN KEY (work_id)
    REFERENCES Works(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE ReviewLike ADD CONSTRAINT ReviewLike_Review
    FOREIGN KEY (review_id)
    REFERENCES Review(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE ReviewLike ADD CONSTRAINT ReviewLike_User
    FOREIGN KEY (user_id)
    REFERENCES Users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE UserFollows ADD CONSTRAINT UserFollows_UsersA
    FOREIGN KEY (id_followed)
    REFERENCES Users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE UserFollows ADD CONSTRAINT UserFollows_UsersB
    FOREIGN KEY (id_follower)--no se si esto va separado
    REFERENCES Users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE UserFavorites ADD CONSTRAINT UserFavorites_Users
    FOREIGN KEY (id_user)
    REFERENCES Users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE UserFavorites ADD CONSTRAINT UserFavorites_Works
    FOREIGN KEY (id_work)
    REFERENCES Works(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

--------------------- Funciones + Triggers ---------------------

--Control para mantener los usernames unicos (tal vez al pedo?)
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


-- Carga de la tabla ConsumedWorks al momento de cargar una review
-- Controla que la obra que se hace la review es unidad
CREATE OR REPLACE FUNCTION FN_TRIU_REVIEW()
RETURNS TRIGGER AS $$
    BEGIN
        IF EXISTS (SELECT 1 FROM Works w WHERE w.id = NEW.work_id AND w.unit) THEN --aca solo entra en inserts, porque no se puede actualizar workid
            RAISE EXCEPTION 'Solo se puede hacer review de las obras unitarias.';
        ELSE 
            IF NEW.liked AND NOT EXISTS (SELECT 1 FROM LikedWorks l WHERE l.user_id = NEW.user_id AND l.work_id = NEW.work_id) THEN 
                INSERT INTO LikedWorks VALUES (NEW.user_id,NEW.work_id);
            END IF;
            IF NOT EXISTS (SELECT 1 FROM Review r WHERE r.user_id = NEW.user_id AND r.work_id = NEW.work_id) THEN
                INSERT INTO ConsumedWorks VALUES (NEW.user_id,NEW.work_id);
            END IF;
        END IF;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER TRIU_REVIEW
    BEFORE INSERT OR UPDATE ON Review 
    FOR EACH ROW 
        EXECUTE FUNCTION FN_TRI_REVIEW();

--Eliminar de ConsumedWork debe eliminar de Liked y las review, junto con los comentarios y likes de review
--PREGUNTAR: tal vez poner un booleano "active" y desactivarlo a modo de eliminar
--!!ASEGURARSE PONER UN CARTEL DE "¿ESTAS SEGURO?"
CREATE OR REPLACE FUNCTION FN_TRD_CONSUMEDWORKS()
RETURNS TRIGGER AS $$
    BEGIN
        DELETE FROM ReviewLike lr WHERE OLD.work_id = lr.work_id;
        DELETE FROM ReviewComment rc WHERE OLD.work_id = lr.work_id;
        DELETE FROM Review r WHERE r.work_id = OLD.work_id;
        DELETE FROM LikedWorks lw WHERE lw.work_id = OLD.work_id;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER TRD_CONSUMEDWORKS
    BEFORE DELETE ON ConsumedWorks
    FOR EACH ROW 
        EXECUTE FUNCTION FN_TRD_CONSUMEDWORKS();

--s