--------------------- Creacion de tablas ---------------------

CREATE TABLE IF NOT EXISTS usuarios (
    id SERIAL PRIMARY KEY,
    email VARCHAR(50) UNIQUE CONSTRAINT uq_email NOT NULL, -- Alternative key, Se nombran para poder usarlas en el manejo de errores.
    contraseña TEXT NOT NULL, -- TEXT debido a la encriptación.
    creado_en TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS perfiles (
    id_usuario INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
);

CREATE TABLE IF NOT EXISTS noticias (
    id SERIAL PRIMARY KEY,
    titulo TEXT NOT NULL,
    contenido TEXT NOT NULL,
    publicada_en TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    tiempo_lectura_estimado TIMESTAMP,
    visualizaciones INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS likes_noticia (
    id_noticia INT,
    id_usuario INT,
    likeado_en TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_news_likes PRIMARY KEY (id_noticia, id_usuario)
);

ALTER TABLE likes_noticia ADD CONSTRAINT fk_likes_noticia_usuarios
    FOREIGN KEY (id_usuario)
    REFERENCES usuarios(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

CREATE TABLE IF NOT EXISTS comentarios_noticia (
    new_id INT,
    user_id INT,
    -- Con las tres siendo primary key, entonces un usuario puede realizar más de un comentario en una publicación.
    -- Si se queire evitar el spam, podría evitarse eliminando el atributo debajo.
    comment_id SERIAL,
    comment TEXT NOT NULL,
    published_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    tiempo_estimado_lectura TIMESTAMP,
    visualizaciones INT DEFAULT 0,
    CONSTRAINT pk_news_comments PRIMARY KEY (new_id, user_id, comment_id)
);

ALTER TABLE news_comments ADD CONSTRAINT fk_news_comments_users
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE news_comments ADD CONSTRAINT fk_news_comments_news
    FOREIGN KEY (new_id)
    REFERENCES news(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

CREATE TABLE IF NOT EXISTS comments_likes (
    new_id INT,
    user_id INT,
    comment_id INT,
    user_like_id INT,
    liked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_comments_likes PRIMARY KEY (new_id, user_id, comment_id, user_like_id)
);

ALTER TABLE comments_likes ADD CONSTRAINT fk_comments_likes_comments
    FOREIGN KEY (new_id, user_id, comment_id)
    REFERENCES news_comments(new_id, user_id, comment_id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

CREATE TABLE IF NOT EXISTS facultades (
    id SERIAL PRIMARY KEY, 
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS puntajes (
    facultad1_id INT,
    facultad2_id INT,
    partido_id INT,
    puntos1 INT NOT NULL,
    puntos2 INT NOT NULL,
    puntosS1 INT DEFAULT NULL,
    puntosS2 INT DEFAULT NULL,
    CONSTRAINT pk_puntaje PRIMARY KEY (facultad1_id,facultad2_id,partido_id) 
);

CREATE TABLE IF NOT EXISTS partidos (
    id SERIAL PRIMARY KEY,
    deporte_id INT NOT NULL,
    incio TIMESTAMP,
    fin TIMESTAMP DEFAULT NULL,
    lugar VARCHAR(255),
    cancha TEXT DEFAULT NULL,
    zona CHAR DEFAULT NULL,
    tipo VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS deportes (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255)
);

ALTER TABLE partidos ADD CONSTRAINT fk_partidos_deportes
    FOREIGN KEY deporte_id
    REFERENCES deportes(id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

ALTER TABLE puntajes ADD CONSTRAINT fk_puntajes_facultad1
    FOREIGN KEY facultad1_id
    REFERENCES facultades(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE puntajes ADD CONSTRAINT fk_puntajes_facultad2
    FOREIGN KEY facultad2_id
    REFERENCES facultades(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE puntajes ADD CONSTRAINT fk_puntajes_partido
    FOREIGN KEY partido_id
    REFERENCES partidos(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

--------------------- Funciones + Triggers ---------------------

/*
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

        */