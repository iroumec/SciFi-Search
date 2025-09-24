-- name: CrearUsuario :one
INSERT INTO usuarios (dni, nombre, email, contrasena) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: ObtenerUsuarioPorID :one
SELECT * FROM usuarios WHERE id = $1;

-- name: ObtenerUsuarioPorDNI :one
SELECT * FROM usuarios WHERE dni = $1;

-- name: ListarUsuarios :many
SELECT * FROM usuarios ORDER BY id;

-- name: ActualizarUsuario :exec
UPDATE usuarios SET nombre = $2, email = $3, contrasena = $4 WHERE id = $1;

-- name: ActualizarDNI :exec
UPDATE usuarios SET dni = $2 WHERE id = $1;

-- name: EliminarUsuario :exec
DELETE FROM usuarios WHERE id = $1;

-- name: CrearPerfil :one
INSERT INTO perfiles (id_usuario, foto) VALUES ($1, $2) RETURNING *;

-- name: CrearFacultad :one
INSERT INTO facultades (nombre) VALUES ($1) RETURNING *;

-- name: ObtenerFacultad :one
SELECT * FROM facultades WHERE id = $1;

-- name: ObtenerNombreFacultad :one
SELECT nombre FROM facultades WHERE id = $1;

-- name: ObtenerFacultades :many
SELECT * FROM facultades ORDER BY nombre;

-- name: ActualizarFacultad :exec
UPDATE facultades SET nombre = $2 WHERE id = $1;

-- name: EliminarFacultad :exec
DELETE FROM facultades WHERE id = $1;

-- name: AsignarFacultad :exec
INSERT INTO pertenece (id_usuario,id_facultad) VALUES ($1,$2);

-- name: EliminarAsginacionFacultad :exec
DELETE FROM pertenece WHERE id_usuario = $1 AND id_facultad = $2;

-- name: CrearDeporte :one
INSERT INTO deportes (nombre,masculino,foto) VALUES ($1,$2,$3) RETURNING *;

-- name: ObtenerDeporte :one
SELECT * FROM deportes WHERE id = $1;

-- name: ObtenerDeportes :many
SELECT * FROM deportes ORDER BY nombre;

-- name: ActualizarDeporte :exec
UPDATE deportes SET nombre = $2 , masculino = $3 foto = $4 WHERE id = $1;

-- name: EliminarDeporte :exec
DELETE FROM deportes WHERE id = $1;

-- name: CrearPartidoZonas :one
INSERT INTO partidos (id_deporte,tipo,zona,id_facultad1,id_facultad2,inicio,lugar) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING *;

-- name: CrearPartidoSemis :one
INSERT INTO partidos (id_deporte,tipo,zona,id_facultad1,id_facultad2,inicio,lugar) VALUES ($1,"Semifinales",null,$2,$3,$4,$5) RETURNING *;

-- name: CrearPartidoTercer :one
INSERT INTO partidos (id_deporte,tipo,zona,id_facultad1,id_facultad2,inicio,lugar) VALUES ($1,"Tercer",null,$2,$3,$4,$5) RETURNING *;

-- name: CrearPartidoFinal :one
INSERT INTO partidos (id_deporte,tipo,zona,id_facultad1,id_facultad2,inicio,lugar) VALUES ($1,"Final",null,$2,$3,$4,$5) RETURNING *;

-- name: ObtenerPartidoPorID :one
SELECT * FROM partidos WHERE id = $1;

-- name: ListarPartidosPorFacultad :many
SELECT * FROM partidos WHERE id_facultad1 = $1 OR id_facultad2 = $1;

-- name: ListarPartidosPorDeporte :many
SELECT * FROM partidos WHERE id_deporte = $1;

-- name: ListarPartidosPorFacultadYDeporte :many
SELECT * FROM partidos WHERE (id_facultad1 = $1 OR id_facultad2 = $1) AND id_deporte = $2;

-- name: MANUAL_EditarPartido :exec
UPDATE partidos SET id_deporte = $2, tipo = $3, zona = $4, id_facultad1 = $5, id_facultad2 = $6, inicio = $7, lugar = $8, cancha = $9 WHERE id = $1;

-- name: MANUAL_EliminarPartido :exec
DELETE FROM partidos WHERE id = $1;

-- name: IniciarPartido :exec
INSERT INTO puntajes (id_partido,puntos1,puntos2) VALUES ($1,0,0);

-- name: IniciarPartidoComplejo :exec
INSERT INTO puntajes (id_partido,puntos1,puntos2,puntosS1,puntosS2) VALUES ($1,0,0,0,0);

-- name: SumarPuntos :exec
CALL PR_SUMARPUNTOS($1,$2,$3);

-- name: MANUAL_ModificarPuntaje :exec
UPDATE puntajes SET puntos1 = $2, puntos2 =$3, puntosS1 = $4, puntosS2 = $5 WHERE id_partido = $1;

-- name: MANUAL_EliminarPuntaje :exec
DELETE FROM puntajes WHERE id_partido = $1;

-- name: CargarPuntajeSimpleFacultad :exec
INSERT INTO puntajes_simples (id_disciplina,id_facultad,puntos) VALUES ($1,$2,$3);

-- name: ObtenerPuntajeSimple :one
SELECT * FROM puntajes_simples WHERE id_simple = $1;

-- name: ObtenerXPuestoDisciplina :one
SELECT id_facultad FROM puntajes_simples WHERE id_disciplina = $1 ORDER BY puntos LIMIT $2 OFFSET (SELECT $2 - 1);

-- name: ListarPuntajesSimplesDisciplina :many
SELECT * FROM puntajes_simples WHERE id_disciplina = $1 ORDER BY puntos;

-- name: MANUAL_ModificarPuntajesSimples :exec
UPDATE puntajes_simples SET id_disciplina = $2 , id_facultad = $3 , puntos = $4 WHERE id_simple = $1;

-- name: MANUAL_EliminarPuntajeSimple :exec
DELETE FROM puntajes_simples WHERE id_simple = $1;

-- name: CargarParticipanteDeporte :exec
INSERT INTO participa (id_participante,id_deporte,id_facultad) VALUES ($1,$2,$3) 

-- name: MANUAL_EliminarParticipante :exec
DELETE FROM participa WHERE id_participante = $1 , id_deporte = $2 , id_facultad = $3;

-- name: CrearNoticia :one
INSERT INTO noticias (titulo, contenido, tiempo_lectura_estimado) VALUES ($1, $2, $3) RETURNING *;

-- name: ObtenerNoticia :one
SELECT * FROM noticias WHERE id = $1;

-- name: ListarNoticias :many
SELECT * FROM noticias ORDER BY publicada_en LIMIT 5 OFFSET $1;

-- name: EliminarNoticia :exec
DELETE FROM noticias WHERE id = $1;

-- name: LikearNoticia :one
INSERT INTO likes_noticia (id_noticia, id_usuario) VALUES ($1, $2) RETURNING *;

-- name: ObtenerLikesNoticia :one
SELECT COUNT(*) FROM likes_noticia WHERE id_noticia = $1;

-- name: DeslikearNoticia :exec
DELETE FROM likes_noticia WHERE id_noticia = $1 AND id_usuario = $2;

-- name: AgregarComentarioNoticia :one
INSERT INTO comentarios_noticia (id_noticia, id_usuario, comentario) VALUES ($1, $2, $3) RETURNING *;

-- name: ObtenerComentarioNoticia :one
SELECT * FROM comentarios_noticia WHERE id_comentario = $1;

-- name: ListarComantariosNoticia :many
SELECT * FROM comentarios_noticia WHERE id_noticia = $1 ORDER BY publicado_en DESC LIMIT 10 OFFSET $1;

-- name: ObtenerCantidadComentariosNoticia :one
SELECT COUNT(*) FROM comentarios_noticia WHERE id_noticia = $1;

-- name: EliminarComentarioNoticia :exec
DELETE FROM comentarios_noticia WHERE id_comentario = $1;

-- name: LikearComentarioNoticia :one
INSERT INTO likes_comentario (id_comentario,id_usuario) VALUES ($1, $2) RETURNING *;

-- name: ObtenerLikesComentarioNoticia :one
SELECT COUNT(*) FROM likes_comentario WHERE id_noticia = $1 AND id_comentario = $2;

-- name: DeslikearComentarioNoticia :exec
DELETE FROM likes_comentario WHERE id_comentario = $1 AND id_usuario = $2;

-- name: FinalizarPartido :exec
CALL PR_FINALIZARPARTIDO($1,$2);

-- name: MANUAL_CrearPartidoHistorico :exec
INSERT INTO partidos_historicos VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14);

-- name: ObtenerPartidoHistorico :one
SELECT * FROM partidos_historicos WHERE id_partido = $1;

-- name: ListarPartidosHistoricos :many
SELECT * FROM partidos_historicos ORDER BY fin DESC;

-- name: ListarPartidosHistoricosPorAnio :many
SELECT * FROM partidos_historicos WHERE $1 = EXTRACT(YEAR FROM fin) ORDER BY fin DESC;

-- name: MANUAL_ActualizarPartidoHistorico :exec
UPDATE partidos_historicos SET (id_deporte,tipo,zona,id_facultad1,id_facultad2,inicio,fin,lugar,cancha,puntos1,puntos2,puntosS1,puntosS2) = ($2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) WHERE id_partido = $1;

-- name: MANUAL_EliminarPartidoHistorico :exec
DELETE FROM partidos_historicos WHERE id_partido = $1;

-- name: CerrarPuntajeSimple :exec
CALL PR_CERRAR_PUNTAJE_SIMPLE($1,$2);

-- name: CerrarDisciplina :exec
CALL PR_CERRAR_DISCIPLINA($1,$2);

-- name: MANUAL_CrearSimpleHistorico :exec
INSERT INTO simples_historicos VALUES ($1,$2,$3,$4,$5);

-- name: ObtenerSimpleHistorico :one
SELECT * FROM simples_historicos WHERE id_simple = $1;

-- name: ListarSimplesHistoricos :many
SELECT * FROM simples_historicos ORDER BY anio,id_disciplina DESC;

-- name: MANUAL_ActualizarSimpleHistorico :exec
UPDATE simples_historicos SET id_disciplina = $2, id_facultad = $3, anio = $4, puntos = $5 WHERE id_simple = $1;

-- name: MANUAL_EliminarSimpleHistorico :exec
DELETE FROM simples_historicos WHERE id_simple = $1;

-- name: FinalizarOlimpiadas :exec
CALL PR_FINALIZAR_OLIMPIADAS();

-- name: MANUAL_CrearParticipaHistorico :exec
INSERT INTO participa_historico VALUES ($1,$2,$3,$4);

-- name: ObtenerParticipaHistorico :one
SELECT * FROM participa_historico WHERE id_participante = $1 AND id_deporte = $2 AND id_facultad = $3 AND anio = $4;

-- name: ListarParticipaHistorico :many
SELECT * FROM participa_historico ORDER BY anio,id_facultad,id_deporte,id_participante DESC;

-- name: ListarParticipacionesHistoricas :many
SELECT * FROM participa_historico WHERE id_participante = $1 ORDER BY id_facultad,anio,id_deporte DESC;

-- name: CrearPublicacionPartido :exec
INSERT INTO publicaciones (id_partido,link_fotos) VALUES ($1,$2);

-- name: CrearPublicacionSimple :exec 
INSERT INTO publicaciones (id_simple,link_fotos) VALUES ($1,$2);

-- name: CrearPublicacionDisciplina :exec
INSERT INTO publicaciones (id_disciplina,link_fotos) VALUES ($1,$2);

-- name: ObtenerPublicacion :one
SELECT * FROM publicaciones WHERE id = $1;

-- name: ListarPublicaciones :many
SELECT * FROM publicaciones ORDER BY fecha DESC LIMIT $1 OFFSET $2;

-- name: ActualizarPublicacion :exec
UPDATE publicaciones 