------------------------------------------------| Creacion de tablas |------------------------------------------------

---------------------------
-- Usuarios y facultades --
---------------------------
CREATE TABLE IF NOT EXISTS usuarios (
    id SERIAL PRIMARY KEY,
    DNI VARCHAR(9) UNIQUE CONSTRAINT uq_dni NOT NULL, -- Alternative key
    nombre VARCHAR(50) NOT NULL,
    email VARCHAR(50) UNIQUE CONSTRAINT uq_email NOT NULL, -- Se nombran para poder usarlas en el manejo de errores.
    contrasena TEXT NOT NULL, -- TEXT debido a la encriptación.
    creado_en TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS perfiles (
    id_usuario INT PRIMARY KEY,
    foto TEXT
    --insignias
);

CREATE TABLE IF NOT EXISTS facultades (
    id SERIAL PRIMARY KEY, 
    nombre VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS perfiles_facultad (--altas posibilidades de volar, si on vuela CRUD
    id_facultad INT PRIMARY KEY
    --insignias y otras cosas
);

CREATE TABLE IF NOT EXISTS pertenece (
    id_usuario INT,
    id_facultad INT,
    CONSTRAINT pk_pertenece PRIMARY KEY (id_usuario, id_facultad)
);

-----------------------------------
-- Deportes, partidos y puntajes --
-----------------------------------
CREATE TABLE IF NOT EXISTS deportes (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255),
    masculino BOOLEAN, --(true = masculino, false = femenino, null = sin genero)
    foto TEXT --se guardaría aca?
);

-- Para todo deporte menos Culturales, Cross, Ajedrez, Egames, Tenis de mesa
CREATE TABLE IF NOT EXISTS partidos (--maybe fusionar partidos y puntajes?
    id SERIAL PRIMARY KEY,
    id_deporte INT NOT NULL,
    tipo VARCHAR(20) NOT NULL, -- zonas , semifinal , 3y4 , final
    zona CHAR NOT NULL DEFAULT 'A', -- se podría hacer A hasta R zonas, S semi, F 3y4 y final, medio feo
    id_facultad1 INT NOT NULL,
    id_facultad2 INT NOT NULL,
    inicio TIMESTAMP NOT NULL,
    lugar VARCHAR(255) NOT NULL,
    cancha TEXT
); 

--S1 y S2 (y puntaje_techo) es null para TODO menos Voley.
CREATE TABLE IF NOT EXISTS puntajes (
    id_partido INT PRIMARY KEY,
    puntos1 INT NOT NULL,
    puntos2 INT NOT NULL,
    puntosS1 INT,
    puntosS2 INT
    --puntaje_techo: cuando puntosS1 llega a puntaje_techo, suma puntos en puntos1 (puede cambiar en base al deporte).
);

--Para culturales, cross y ajedrez (e-games?)
CREATE TABLE IF NOT EXISTS puntajes_simples (
    id_simple SERIAL PRIMARY KEY,
    id_disciplina INT NOT NULL,
    id_facultad INT NOT NULL,
    puntos INT NOT NULL
);

CREATE TABLE IF NOT EXISTS participa (
    id_participante INT,
    id_deporte INT,
    id_facultad INT,
    CONSTRAINT pk_participa PRIMARY KEY (id_participante,id_deporte,id_facultad)
);

--------------
-- Noticias --
--------------
CREATE TABLE IF NOT EXISTS noticias (
    id SERIAL PRIMARY KEY,
    titulo TEXT NOT NULL,
    contenido TEXT NOT NULL,
    visualizaciones INT NOT NULL DEFAULT 0,
    tiempo_lectura_estimado INT NOT NULL,
    publicada_en TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS likes_noticia (
    id_noticia INT,
    id_usuario INT,
    likeado_en TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_noticias_likes PRIMARY KEY (id_noticia, id_usuario)
);

CREATE TABLE IF NOT EXISTS comentarios_noticia (
    id_comentario SERIAL PRIMARY KEY,
    id_noticia INT,
    id_usuario INT,
    comentario TEXT NOT NULL CHECK (BTRIM(comentario) <> ''),
    publicado_en TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS likes_comentario_noticia (
    id_comentario INT,
    id_usuario INT,
    liked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_likes_comentarios_noticia PRIMARY KEY (id_comentario, id_usuario)
);

----------------
-- Historicos --
----------------
CREATE TABLE IF NOT EXISTS partidos_historicos (
    id_partido INT PRIMARY KEY,
    id_deporte INT NOT NULL,
    tipo VARCHAR(20) NOT NULL,
    zona CHAR NOT NULL DEFAULT 'A',
    id_facultad1 INT NOT NULL,
    id_facultad2 INT NOT NULL,
    inicio TIMESTAMP NOT NULL,
    fin TIMESTAMP NOT NULL,
    lugar VARCHAR(255) NOT NULL,
    cancha TEXT,
    puntos1 INT NOT NULL,
    puntos2 INT NOT NULL,
    puntosS1 INT,
    puntosS2 INT
);

CREATE TABLE IF NOT EXISTS simples_historicos (
    id_simple INT PRIMARY KEY,
    id_disciplina INT NOT NULL,
    id_facultad INT NOT NULL,
    anio INT NOT NULL DEFAULT EXTRACT(YEAR FROM CURRENT_TIMESTAMP),
    puntos INT NOT NULL
);

CREATE TABLE IF NOT EXISTS participa_historico (
    id_participante INT,
    id_deporte INT,
    id_facultad INT,
    anio INT DEFAULT EXTRACT (YEAR FROM (CURRENT_TIMESTAMP)),
    CONSTRAINT pk_participa PRIMARY KEY (id_participante,id_deporte,id_facultad,anio)
);

-------------------
-- Publicaciones --
-------------------
CREATE TABLE IF NOT EXISTS publicaciones ( --se genera el posteo una vez finaliza el partido/disciplina
    id_publicacion SERIAL PRIMARY KEY,
    id_partido INT,
    id_simple INT,
    id_disciplina INT,
    link_fotos TEXT, --url a un drive dedicado a fotos?
    fecha TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    visualizaciones INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS comentarios_publicacion (
    id SERIAL PRIMARY KEY,
    id_publicacion INT NOT NULL,
    id_usuario INT NOT NULL,
    comentario VARCHAR(512) NOT NULL CHECK (BTRIM(comentario) <> '')
    --posibilidad de subir foto?
);

CREATE TABLE IF NOT EXISTS likes_publicacion (
    id_publicacion INT,
    id_usuario INT,
    likeado_en TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_likes_publicacion PRIMARY KEY (id_publicacion, id_usuario)
);

CREATE TABLE IF NOT EXISTS respuestas_comentario_publicacion (
    id SERIAL PRIMARY KEY,
    id_comentario INT,
    id_usuario INT,
    comentario VARCHAR(512) NOT NULL CHECK (BTRIM(comentario) <> '')
);

CREATE TABLE IF NOT EXISTS likes_comentario_publicacion (
    id_comentario INT,
    id_usuario INT,
    likeado_en TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_likes_comentario_publicacion PRIMARY KEY (id_comentario, id_usuario)
);


------------------------------------------------| Foreign Keys y Restricciones de Integridad |------------------------------------------------
---------------------------
-- Usuarios y facultades --
---------------------------
ALTER TABLE perfiles ADD CONSTRAINT fk_perfiles_usuarios
    FOREIGN KEY (id_usuario) REFERENCES usuarios(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE perfiles_facultad ADD CONSTRAINT fk_perfiles_facultad_facultades
    FOREIGN KEY (id_facultad) REFERENCES facultades(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE pertenece ADD CONSTRAINT fk_pertenece_usuario
    FOREIGN KEY (id_usuario) REFERENCES usuarios(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE pertenece ADD CONSTRAINT fk_pertenece_facultad
    FOREIGN KEY (id_facultad) REFERENCES facultades(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

-----------------------------------
-- Deportes, partidos y puntajes --
-----------------------------------
ALTER TABLE partidos ADD CONSTRAINT fk_partidos_deportes
    FOREIGN KEY (id_deporte) REFERENCES deportes(id)
    NOT DEFERRABLE  INITIALLY IMMEDIATE
;

ALTER TABLE partidos ADD CONSTRAINT fk_partidos_facultades1
    FOREIGN KEY (id_facultad1) REFERENCES facultades(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE partidos ADD CONSTRAINT fk_partidos_facultades2
    FOREIGN KEY (id_facultad2) REFERENCES facultades(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE puntajes ADD CONSTRAINT fk_puntajes_partidos
    FOREIGN KEY (id_partido) REFERENCES partidos(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE puntajes_simples ADD CONSTRAINT fk_puntajes_simples_deportes 
    FOREIGN KEY (id_disciplina) REFERENCES deportes(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE puntajes_simples ADD CONSTRAINT fk_puntajes_simples_facultades 
    FOREIGN KEY (id_facultad) REFERENCES facultades(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE participa ADD CONSTRAINT fk_participa_usuarios
    FOREIGN KEY (id_participante) REFERENCES usuarios(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE participa ADD CONSTRAINT fk_participa_deportes
    FOREIGN KEY (id_deporte) REFERENCES deportes(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE participa ADD CONSTRAINT fk_participa_facultades
    FOREIGN KEY (id_facultad) REFERENCES facultades(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

--------------
-- Noticias --
--------------
ALTER TABLE comentarios_noticia ADD CONSTRAINT fk_comentarios_noticia_usuarios
    FOREIGN KEY (id_usuario) REFERENCES usuarios(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE comentarios_noticia ADD CONSTRAINT fk_comentarios_noticia_noticias
    FOREIGN KEY (id_noticia) REFERENCES noticias(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE likes_noticia ADD CONSTRAINT fk_likes_noticia_usuarios
    FOREIGN KEY (id_usuario) REFERENCES usuarios(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE likes_noticia ADD CONSTRAINT fk_likes_noticia_noticias
    FOREIGN KEY (id_noticia) REFERENCES noticias(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE likes_comentario_noticia ADD CONSTRAINT fk_likes_comentario_noticia_comentarios
    FOREIGN KEY (id_comentario) REFERENCES comentarios_noticia(id_comentario)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE likes_comentario_noticia ADD CONSTRAINT fk_likes_comentario_noticia_usuarios
    FOREIGN KEY (id_usuario) REFERENCES usuarios(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

----------------
-- Historicos --
----------------
ALTER TABLE partidos_historicos ADD CONSTRAINT fk_partidos_historicos_deportes
    FOREIGN KEY (id_deporte) REFERENCES deportes(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE partidos_historicos ADD CONSTRAINT fk_partidos_historicos_facultades1
    FOREIGN KEY (id_facultad1) REFERENCES facultades(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE partidos_historicos ADD CONSTRAINT fk_partidos_historicos_facultades2
    FOREIGN KEY (id_facultad2) REFERENCES facultades(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE simples_historicos ADD CONSTRAINT fk_simples_historicos_deportes
    FOREIGN KEY (id_disciplina) REFERENCES deportes(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE simples_historicos ADD CONSTRAINT fk_simples_historicos_facultades
    FOREIGN KEY (id_facultad) REFERENCES facultades(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE participa_historico ADD CONSTRAINT fk_participa_historico_usuarios
    FOREIGN KEY (id_participante) REFERENCES usuarios(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE participa_historico ADD CONSTRAINT fk_participa_historico_deportes
    FOREIGN KEY (id_deporte) REFERENCES deportes(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE participa_historico ADD CONSTRAINT fk_participa_historico_facultades
    FOREIGN KEY (id_facultad) REFERENCES facultades(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

-------------------
-- Publicaciones --
-------------------
ALTER TABLE publicaciones ADD CONSTRAINT check_asociacion_historico
    CHECK(( 
        (id_partido IS NULL AND id_simple IS NOT NULL) --publicacion particular de facultad y disciplina
        OR 
        (id_partido IS NOT NULL AND id_simple IS NULL) --publicacion de partido fac vs fac
        ) OR 
        id_disciplina IS NOT NULL --publicacion de disciplina/deporte general
    )
;

ALTER TABLE publicaciones ADD CONSTRAINT fk_publicaciones_partidos_historicos
    FOREIGN KEY (id_partido) REFERENCES partidos_historicos(id_partido)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE publicaciones ADD CONSTRAINT fk_publicaciones_simples_historicos
    FOREIGN KEY (id_simple) REFERENCES simples_historicos(id_simple)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE comentarios_publicacion ADD CONSTRAINT fk_comentarios_publicacion_publicaciones
    FOREIGN KEY (id_publicacion) REFERENCES publicaciones(id_publicacion)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE comentarios_publicacion ADD CONSTRAINT fk_comentarios_publicacion_usuarios
    FOREIGN KEY (id_usuario) REFERENCES usuarios(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE likes_publicacion ADD CONSTRAINT fk_likes_publicacion_usuarios
    FOREIGN KEY (id_usuario) REFERENCES usuarios(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE likes_publicacion ADD CONSTRAINT fk_likes_publicacion_publicaciones
    FOREIGN KEY (id_publicacion) REFERENCES publicaciones(id_publicacion)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE respuestas_comentario_publicacion ADD CONSTRAINT fk_respuestas_comentario_comentarios_publicacion
    FOREIGN KEY (id_comentario) REFERENCES comentarios_publicacion(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE respuestas_comentario_publicacion ADD CONSTRAINT fk_respuestas_comentario_usuarios
    FOREIGN KEY (id_usuario) REFERENCES usuarios(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE likes_comentario_publicacion ADD CONSTRAINT fk_likes_comentario_publicacion_usuarios
    FOREIGN KEY (id_usuario) REFERENCES usuarios(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;

ALTER TABLE likes_comentario_publicacion ADD CONSTRAINT fk_likes_comentario_publicacion_comentarios_publicacion
    FOREIGN KEY (id_comentario) REFERENCES comentarios_publicacion(id)
    NOT DEFERRABLE INITIALLY IMMEDIATE
;


------------------------------------------------| Funciones + Triggers |------------------------------------------------
-- Creo que podrían estar en otro archivo, porque el sqlc no los usa y solo los usa docker.
-- Tampoco creo que usa los alter table.

CREATE OR REPLACE FUNCTION FN_TRI_PARTICIPA()
RETURNS TRIGGER AS $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pertenece p WHERE p.id_usuario = NEW.id_participante AND p.id_facultad = NEW.id_facultad)THEN
            RAISE EXCEPTION 'Una persona debe participar en un partido por una facultad de la que pertenezca.';
        END IF;
        IF EXISTS (SELECT 1 FROM participa p WHERE NEW.id_deporte = p.id_deporte AND NEW.id_participante = p.id_participante)THEN 
            RAISE EXCEPTION 'Una persona solo puede participar en un deporte por una facultad.';
        END IF;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER TRI_PARTICIPA 
BEFORE INSERT ON participa 
FOR EACH ROW EXECUTE FUNCTION FN_TRI_PARTICIPA;

------------------------------------------------| Procedimientos |------------------------------------------------

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

CREATE OR REPLACE PROCEDURE PR_FINALIZARPARTIDO(id_partido INT,link TEXT)
LANGUAGE 'plpgsql' AS $$
BEGIN
    INSERT INTO partidos_historicos 
        (id_partido,
        id_deporte,
        tipo,
        zona,
        id_facultad1,
        id_facultad2,
        inicio,
        fin,
        lugar,
        cancha,
        puntos1,
        puntos2,
        puntosS1,
        puntosS2
        )
    (SELECT 
        pa.id,
        pa.id_deporte,
        pa.tipo,
        pa.zona,
        pa.id_facultad1,
        pa.id_facultad2,
        pa.inicio,
        NOW() as fin,
        pa.lugar,
        pa.cancha,
        pu.puntos1,
        pu.puntos2,
        pu.puntosS1,
        pu.puntosS2
    FROM partidos pa LEFT JOIN puntajes pu ON (pa.id = pu.id_partido) 
    WHERE id = id_partido);

    DELETE FROM partidos WHERE id = id_partido;
    DELETE FROM puntajes WHERE id = id_partido;

    --Publicacion automática al finalizar el partido
    INSERT INTO publicaciones (id_partido,link_fotos) VALUES (id_partido,link);
END;
$$;

--Solo se debe ejecutar cuando se confirman los resultados, cierra UN puntaje de UNA disciplina de UNA facultad 
CREATE OR REPLACE PROCEDURE PR_CERRAR_PUNTAJE_SIMPLE(id_simple INT, link TEXT)
LANGUAGE 'plpgsql' AS $$
BEGIN
    INSERT INTO simples_historicos (id_simple,id_disciplina,id_facultad,puntos) 
        (SELECT this.id_simple,
                this.id_disciplina,
                this.id_facultad,
                this.puntos
                FROM puntajes_simples this WHERE this.id_simple = id_simple);
    DELETE FROM puntajes_simples this WHERE this.id_simple = id_simple;
    
    --Publicacion automática al finalizar
    INSERT INTO publicaciones (id_simple,link_fotos) VALUES (id_simple,link);
END;
$$;

--Solo se debe ejecutar para disciplinas culturales, ajedrez o cross (e-games?)
CREATE OR REPLACE PROCEDURE PR_CERRAR_DISCIPLINA(id_disciplina INT, link TEXT)
LANGUAGE 'plpgsql' AS $$
BEGIN
    INSERT INTO simples_historicos (id_simple,id_disciplina,id_facultad,puntos) 
        (SELECT this.id_simple,
                this.id_disciplina,
                this.id_facultad,
                this.puntos
                FROM puntajes_simples this WHERE this.id_disciplina = id_disciplina);
    DELETE FROM puntajes_simples this WHERE this.id_disciplina = id_disciplina;
    
    --Publicacion automática al finalizar
    INSERT INTO publicaciones (id_disciplina,link_fotos) VALUES (id_disciplina,link);
END;
$$;

--Verifica que no haya partidos pendientes, si no es el caso: pasa las participaciones a historicos. 
CREATE OR REPLACE PROCEDURE PR_FINALIZAR_OLIMPIADAS()
LANGUAGE 'plpgsql' AS $$
BEGIN
    IF (EXISTS (SELECT 1 FROM partidos)) OR (EXISTS (SELECT 1 FROM puntajes_simples)) THEN 
        RAISE EXCEPTION 'Error. Quedan partidos sin finalizar.';
    ELSE
        INSERT INTO participa_historico (id_participante,id_deporte,id_facultad) (
            SELECT * FROM participa
        );
        DELETE FROM participa;
    END IF;
END;
$$;