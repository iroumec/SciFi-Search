--------------------- Creacion de tablas ---------------------

CREATE TABLE IF NOT EXISTS usuarios (
    id SERIAL PRIMARY KEY,
    DNI VARCHAR(9) UNIQUE CONSTRAINT uq_dni NOT NULL, -- Alternative key
    nombre VARCHAR(50) NOT NULL,
    email VARCHAR(50) UNIQUE CONSTRAINT uq_email NOT NULL, -- Se nombran para poder usarlas en el manejo de errores.
    contraseña TEXT NOT NULL, -- TEXT debido a la encriptación.
    creado_en TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS perfiles (
    id_usuario INT PRIMARY KEY,
    foto TEXT
    --insignias
);

CREATE TABLE IF NOT EXISTS pertenece (
    id_usuario INT,
    id_facultad INT,
    CONSTRAINT pk_pertenece PRIMARY KEY (id_usuario, id_facultad)
);

CREATE TABLE IF NOT EXISTS noticias (
    id SERIAL PRIMARY KEY,
    titulo TEXT NOT NULL,
    contenido TEXT NOT NULL,
    visualizaciones INT DEFAULT 0,
    tiempo_lectura_estimado INT NOT NULL,
    publicada_en TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS likes_noticia (
    id_noticia INT,
    id_usuario INT,
    likeado_en TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_noticias_likes PRIMARY KEY (id_noticia, id_usuario)
);

CREATE TABLE IF NOT EXISTS comentarios_noticia (
    id_noticia INT,
    id_usuario INT,
    id_comentario SERIAL,
    comentario TEXT NOT NULL,
    publicado_en TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_comentarios_noticias PRIMARY KEY (id_noticia, id_usuario, id_comentario)
);

CREATE TABLE IF NOT EXISTS likes_comentario (
    id_noticia INT,
    id_usuario INT,
    id_comentario INT,
    id_usuario_like INT,
    liked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_likes_comentarios PRIMARY KEY (id_noticia, id_usuario, id_comentario, id_usuario_like)
);

CREATE TABLE IF NOT EXISTS facultades (
    id SERIAL PRIMARY KEY, 
    nombre VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS perfiles_facultad (--altas posibilidades de volar
    id_facultad INT PRIMARY KEY
    --insignias y otras cosas
);

CREATE TABLE IF NOT EXISTS puntajes (
    id_partido INT PRIMARY KEY,
    puntos1 INT NOT NULL,
    puntos2 INT NOT NULL,
    puntosS1 INT DEFAULT NULL,
    puntosS2 INT DEFAULT NULL
    --puntaje_techo: cuando puntosS1 llega a puntaje_techo, suma puntos en puntos1 (puede cambiar en base al deporte).
);

CREATE TABLE IF NOT EXISTS partidos (--maybe fusionar partidos y puntajes?
    id SERIAL PRIMARY KEY,
    id_deporte INT,
    tipo VARCHAR(20),
    zona CHAR DEFAULT 'A',
    id_facultad1 INT,
    id_facultad2 INT,
    inicio TIMESTAMP,
    lugar VARCHAR(255),
    cancha TEXT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS partidos_historicos (--hay que hacer tabla historica de "participo"
    id_partido INT PRIMARY KEY,
    id_deporte INT NOT NULL,
    tipo VARCHAR(20),
    zona CHAR DEFAULT 'A',
    id_facultad1 INT NOT NULL,
    id_facultad2 INT NOT NULL,
    inicio TIMESTAMP NOT NULL,
    fin TIMESTAMP NOT NULL,
    lugar VARCHAR(255),
    cancha TEXT DEFAULT NULL,
    puntos1 INT NOT NULL,
    puntos2 INT NOT NULL,
    puntosS1 INT DEFAULT NULL,
    puntosS2 INT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS deportes (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255),
    foto TEXT
);

CREATE TABLE IF NOT EXISTS puntajes_simples ( -- para culturales, ajedrez y cross?
    id_disciplina INT,
    id_facultad INT,
    puntos INT NOT NULL,
    nombre VARCHAR(255),
    CONSTRAINT pk_puntajes_simples PRIMARY KEY (id_disciplina,id_facultad)
);

CREATE TABLE IF NOT EXISTS participa (
    id_participante INT,
    id_partido INT,
    id_facultad INT,
    CONSTRAINT pk_participa PRIMARY KEY (id_participante,id_partido,id_facultad)
);

--------------------- Foreign Keys ---------------------

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

ALTER TABLE likes_noticia ADD CONSTRAINT fk_likes_noticia_usuarios
    FOREIGN KEY (id_usuario)
    REFERENCES usuarios(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE likes_noticia ADD CONSTRAINT fk_likes_noticia_noticia
    FOREIGN KEY (id_noticia)
    REFERENCES noticias(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE perfiles ADD CONSTRAINT fk_perfiles_usuarios
    FOREIGN KEY (id_usuario)
    REFERENCES usuarios(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE puntajes_simples ADD CONSTRAINT fk_puntajes_simples_deportes 
    FOREIGN KEY (id_disciplina)
    REFERENCES deportes(id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

ALTER TABLE puntajes_simples ADD CONSTRAINT fk_puntajes_simples_facultad 
    FOREIGN KEY (id_facultad)
    REFERENCES facultades(id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

ALTER TABLE participa ADD CONSTRAINT fk_participa_usuario
    FOREIGN KEY (id_participante)
    REFERENCES usuarios(id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

ALTER TABLE participa ADD CONSTRAINT fk_participa_partido
    FOREIGN KEY (id_partido)
    REFERENCES partidos(id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

ALTER TABLE participa ADD CONSTRAINT fk_participa_facultad
    FOREIGN KEY (id_facultad)
    REFERENCES facultades(id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

ALTER TABLE partidos ADD CONSTRAINT fk_partidos_deportes
    FOREIGN KEY (id_deporte)
    REFERENCES deportes(id)
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

ALTER TABLE partidos ADD CONSTRAINT fk_partidos_facultad1
    FOREIGN KEY (id_facultad1)
    REFERENCES facultades(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE partidos ADD CONSTRAINT fk_partidos_facultad2
    FOREIGN KEY (id_facultad2)
    REFERENCES facultades(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE puntajes ADD CONSTRAINT fk_puntajes_partido
    FOREIGN KEY (id_partido)
    REFERENCES partidos(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE perfiles_facultad ADD CONSTRAINT fk_perfiles_facultad_facultades
    FOREIGN KEY (id_facultad)
    REFERENCES facultades(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE likes_comentario ADD CONSTRAINT fk_likes_comentario_comentarios
    FOREIGN KEY (id_noticia, id_usuario, id_comentario)
    REFERENCES comentarios_noticia(id_noticia, id_usuario, id_comentario)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE pertenece ADD CONSTRAINT fk_pertenece_usuario
    FOREIGN KEY (id_usuario)
    REFERENCES usuarios(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

ALTER TABLE pertenece ADD CONSTRAINT fk_pertenece_facultad
    FOREIGN KEY (id_facultad)
    REFERENCES facultades(id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;


--------------------- Funciones + Triggers ---------------------
-- Creo que podrían estar en otro archivo, porque el sqlc no los usa y solo los usa docker.
-- Tampoco creo que usa los alter table.

CREATE OR REPLACE FUNCTION FN_TRI_PARTICIPA()
RETURNS TRIGGER AS $$
    BEGIN
        IF(new.id_facultad NOT IN (SELECT id_facultad FROM pertenece p WHERE p.id_usuario = NEW.id_usuario))THEN
            RAISE EXCEPTION 'Una persona puede participar en un partido por una unica facultad.'
        END IF;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER TRI_PARTICIPA 
BEFORE INSERT ON participa 
FOR EACH ROW EXECUTE FUNCTION FN_TRI_PARTICIPA;

--------------------- Procedimientos ---------------------

CREATE OR REPLACE PROCEDURE PR_SUMARPUNTOS(id_partido INT, nro_equipo INT, cant_puntos INT)
LANGUAGE 'plpgsql' AS $$
BEGIN
    IF(nro_equipo == 1)THEN
        UPDATE puntajes this SET this.puntos1 = this.puntos1+cant_puntos WHERE this.id_partido = id_partido;
    ELSE
        UPDATE puntajes this SET this.puntos2 = this.puntos2+cant_puntos WHERE this.id_partido = id_partido;
    END IF;
END;
$$;

CREATE OR REPLACE PROCEDURE PR_FINALIZARPARTIDO(id_partido INT)
LANGUAGE 'plpgsql' AS $$
BEGIN
    INSERT INTO partidos_historicos (id_partido,id_deporte,tipo,zona,id_facultad1,id_facultad2,inicio,fin,lugar,cancha,puntos1,puntos2,puntosS1,puntosS2)
    (SELECT pa.id,pa.id_deporte,pa.tipo,pa.zona,pa.id_facultad1,pa.id_facultad2,pa.inicio,NOW() as fin,pa.lugar,pa.cancha,pu.puntos1,pu.puntos2,pu.puntosS1,pu.puntosS2
    FROM partidos pa 
    LEFT JOIN puntajes pu ON (pa.id = pu.id_partido) 
    WHERE id = id_partido);
    DELETE FROM partidos WHERE id = id_partido;
    DELETE FROM puntajes WHERE id = id_partido;
END;
$$;
