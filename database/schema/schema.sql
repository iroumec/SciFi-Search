--------------------- Creacion de tablas ---------------------

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(20) UNIQUE CONSTRAINT uq_usuario NOT NULL, -- Alternative key.
    name VARCHAR(20) NOT NULL,
    password VARCHAR(20) NOT NULL,
    email VARCHAR(50) UNIQUE CONSTRAINT uq_email NOT NULL, -- Se nombran para poder usarlas en el manejo de errores.
    password TEXT NOT NULL, -- TEXT debido a la encriptación.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS news (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    published_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    estimated_lecture_time TIMESTAMP,
    views INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS news_likes (
    new_id INT,
    user_id INT,
    liked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_news_likes PRIMARY KEY (new_id, user_id)
);

ALTER TABLE news_likes ADD CONSTRAINT fk_news_likes_users
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

CREATE TABLE IF NOT EXISTS news_comments (
    id SERIAL PRIMARY KEY,
    comment TEXT NOT NULL,
    published_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    tiempo_estimado_lectura TIMESTAMP,
    visualizaciones INT DEFAULT 0
);

ALTER TABLE news_comments ADD CONSTRAINT fk_news_comments_users
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;


CREATE TABLE IF NOT EXISTS comments_likes (
    new_id INT,
    user_id INT,
    comment_id SERIAL INT,
    liked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_news_likes PRIMARY KEY (new_id, user_id, comment_id)
);


--------------------- Asignacion de foreign keys ---------------------

ALTER TABLE works ADD CONSTRAINT fk_works_content_type
    FOREIGN KEY (content_type_id)
    REFERENCES content_types(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE works ADD CONSTRAINT fk_works_works --saga
    FOREIGN KEY (saga_id)
    REFERENCES works(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE review ADD CONSTRAINT fk_review_users
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE review ADD CONSTRAINT fk_review_works
    FOREIGN KEY (work_id)
    REFERENCES works(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE consumed_works ADD CONSTRAINT fk_consumed_works_users
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE consumed_works ADD CONSTRAINT fk_consumed_works_works
    FOREIGN KEY (work_id)
    REFERENCES works(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE liked_works ADD CONSTRAINT fk_liked_works_users
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE liked_works ADD CONSTRAINT fk_liked_works_works
    FOREIGN KEY (work_id)
    REFERENCES works(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE review_like ADD CONSTRAINT fk_review_like_review
    FOREIGN KEY (review_id)
    REFERENCES review(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE review_like ADD CONSTRAINT fk_review_like_users
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE user_follows ADD CONSTRAINT fk_user_follows_followed_user
    FOREIGN KEY (followed_id)
    REFERENCES users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE user_follows ADD CONSTRAINT fk_user_follows_follower_user
    FOREIGN KEY (follower_id) -- Separate foreign key constraint for follower_id, as both followed_id and follower_id reference Users(id) independently
    REFERENCES users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE user_favourites ADD CONSTRAINT fk_user_favourites_users
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE user_favourites ADD CONSTRAINT fk_user_favourites_works
    FOREIGN KEY (work_id)
    REFERENCES works(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

--------------------- Funciones + Triggers ---------------------
-- Creo que podrían estar en otro archivo, porque el sqlc no los usa y solo los usa docker.
-- Tampoco creo que usa los alter table.

-- Carga de la tabla ConsumedWorks al momento de cargar una review
-- Controla que la obra que se hace la review es unidad
CREATE OR REPLACE FUNCTION FN_TRIU_REVIEW()
RETURNS TRIGGER AS $$
    BEGIN
        IF EXISTS (SELECT 1 FROM works w WHERE w.id = NEW.work_id AND w.unit) THEN --aca solo entra en inserts, porque no se puede actualizar workid
            RAISE EXCEPTION 'Solo se puede hacer review de las obras unitarias.';
        ELSE 
            IF NEW.liked AND NOT EXISTS (SELECT 1 FROM liked_works l WHERE l.user_id = NEW.user_id AND l.work_id = NEW.work_id) THEN 
                INSERT INTO liked_works VALUES (NEW.user_id,NEW.work_id);
            END IF;
            IF NOT EXISTS (SELECT 1 FROM review r WHERE r.user_id = NEW.user_id AND r.work_id = NEW.work_id) THEN
                INSERT INTO consumed_works VALUES (NEW.user_id,NEW.work_id);
            END IF;
        END IF;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER TRIU_REVIEW
    BEFORE INSERT OR UPDATE ON review 
    FOR EACH ROW 
        EXECUTE FUNCTION FN_TRIU_REVIEW();

--Eliminar de ConsumedWork debe eliminar de Liked y las review, junto con los comentarios y likes de review
--PREGUNTAR: tal vez poner un booleano "active" y desactivarlo a modo de eliminar
--!!ASEGURARSE PONER UN CARTEL DE "¿ESTAS SEGURO?"
CREATE OR REPLACE FUNCTION FN_TRD_CONSUMEDWORKS()
RETURNS TRIGGER AS $$
    BEGIN
        DELETE FROM review_like lr WHERE OLD.work_id = lr.work_id;
        DELETE FROM review_comment rc WHERE OLD.work_id = lr.work_id;
        DELETE FROM review r WHERE r.work_id = OLD.work_id;
        DELETE FROM liked_works lw WHERE lw.work_id = OLD.work_id;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER TRD_CONSUMEDWORKS
    BEFORE DELETE ON consumed_works
    FOR EACH ROW 
        EXECUTE FUNCTION FN_TRD_CONSUMEDWORKS();

--Verifica que la obra a marcar es unitaria.
--Verifica que el usuario haya marcado unicamente un favorito por tipo de contenido
CREATE OR REPLACE FUNCTION FN_TRI_USER_FAVOURITES()
RETURNS TRIGGER AS $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM works WHERE id = NEW.work_id AND unit) THEN 
            RAISE EXCEPTION 'Para marcar una obra como favorita, esta debe ser unitaria.';
        END IF;
        
        IF ((SELECT content_type_id FROM works w WHERE NEW.work_id = id) --obtengo el tipo de contenido del nuevo
            IN 
            (SELECT content_type_id FROM works m WHERE m.id
            IN --obtengo lista de los tipos de contenido de los favoritos del usuario
            (SELECT work_id FROM user_favourites u WHERE u.user_id = NEW.user_id))) THEN 
                RAISE EXCEPTION 'Solo se puede marcar como favorita una obra por tipo de contenido.';
        END IF;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER TRI_USER_FAVOURITES 
    BEFORE INSERT ON user_favourites
    FOR EACH ROW 
        EXECUTE FUNCTION FN_TRI_USER_FAVOURITES();