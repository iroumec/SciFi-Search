--------------------- Creacion de tablas ---------------------

CREATE TABLE IF NOT EXISTS usuarios (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(50) UNIQUE CONSTRAINT uq_email NOT NULL, -- Alternative key, Se nombran para poder usarlas en el manejo de errores.
    contraseña TEXT NOT NULL, -- TEXT debido a la encriptación.
    creado_en TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS perfiles (
    id_usuario INT PRIMARY KEY
    --insignias, foto de perfil
);

ALTER TABLE perfiles ADD CONSTRAINT fk_perfiles_usuarios
    FOREIGN KEY id_usuario
    REFERENCES usuarios(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

CREATE TABLE IF NOT EXISTS pertenece_facultad (
    id_usuario INT,
    id 
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
    CONSTRAINT pk_noticias_likes PRIMARY KEY (id_noticia, id_usuario)
);

ALTER TABLE likes_noticia ADD CONSTRAINT fk_likes_noticia_usuarios
    FOREIGN KEY (id_usuario)
    REFERENCES usuarios(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

CREATE TABLE IF NOT EXISTS comentarios_noticia (
    id_noticia INT,
    id_usuario INT,
    -- Con las tres siendo primary key, entonces un usuario puede realizar más de un comentario en una publicación.
    -- Si se queire evitar el spam, podría evitarse eliminando el atributo debajo.
    id_comentario SERIAL,
    comentario TEXT NOT NULL,
    fecha_publicacion TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_comentarios_noticias PRIMARY KEY (id_noticia, id_usuario, id_comentario)
);

ALTER TABLE comentarios_noticia ADD CONSTRAINT fk_comentarios_noticia_usuarios
    FOREIGN KEY (id_usuario)
    REFERENCES usuarios(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE comentarios_noticia ADD CONSTRAINT fk_noticias_comments_noticias
    FOREIGN KEY (id_noticia)
    REFERENCES noticias(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

CREATE TABLE IF NOT EXISTS likes_comentarios (
    id_noticia INT,
    id_usuario INT,
    id_comentario INT,
    usuario_like_id INT,
    liked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_likes_comentarios PRIMARY KEY (id_noticia, id_usuario, id_comentario, usuario_like_id)
);

ALTER TABLE comments_likes ADD CONSTRAINT fk_likes_comentarios_comentarios
    FOREIGN KEY (id_noticia, id_usuario, id_comentario)
    REFERENCES comentarios_notica(id_noticia, id_usuario, id_comentario)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

CREATE TABLE IF NOT EXISTS facultades (
    id SERIAL PRIMARY KEY, 
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS perfil_facultad

CREATE TABLE IF NOT EXISTS puntajes (
    id_facultad1 INT,
    id_facultad2 INT,
    id_partido INT,
    puntos1 INT NOT NULL,
    puntos2 INT NOT NULL,
    puntosS1 INT DEFAULT NULL,
    puntosS2 INT DEFAULT NULL,
    CONSTRAINT pk_puntajes PRIMARY KEY (id_facultad1,id_facultad2,id_partido) 
);

CREATE TABLE IF NOT EXISTS partidos (
    id SERIAL PRIMARY KEY,
    id_deporte INT NOT NULL,
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
    FOREIGN KEY id_deporte
    REFERENCES deportes(id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

ALTER TABLE puntajes ADD CONSTRAINT fk_puntajes_facultad1
    FOREIGN KEY id_facultad1
    REFERENCES facultades(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE puntajes ADD CONSTRAINT fk_puntajes_facultad2
    FOREIGN KEY id_facultad2
    REFERENCES facultades(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE puntajes ADD CONSTRAINT fk_puntajes_partido
    FOREIGN KEY id_partido
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
            IF NEW.liked AND NOT EXISTS (SELECT 1 FROM liked_works l WHERE l.id_usuario = NEW.id_usuario AND l.work_id = NEW.work_id) THEN 
                INSERT INTO liked_works VALUES (NEW.id_usuario,NEW.work_id);
            END IF;
            IF NOT EXISTS (SELECT 1 FROM review r WHERE r.id_usuario = NEW.id_usuario AND r.work_id = NEW.work_id) THEN
                INSERT INTO consumed_works VALUES (NEW.id_usuario,NEW.work_id);
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
            (SELECT work_id FROM user_favourites u WHERE u.id_usuario = NEW.id_usuario))) THEN 
                RAISE EXCEPTION 'Solo se puede marcar como favorita una obra por tipo de contenido.';
        END IF;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER TRI_USER_FAVOURITES 
    BEFORE INSERT ON user_favourites
    FOR EACH ROW 
        EXECUTE FUNCTION FN_TRI_USER_FAVOURITES();

        */