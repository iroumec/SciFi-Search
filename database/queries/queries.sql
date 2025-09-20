-- name: ObtenerUsuarioPorID :one
SELECT * FROM usuarios WHERE id = $1;

-- name: ObtenerUsuarioPorDNI :one
SELECT * FROM usuarios WHERE dni = $1;

-- name: ListarUsuarios :many
SELECT * FROM usuarios ORDER BY id;

-- name: CrearUsuario :one
INSERT INTO usuarios (dni, nombre, email, contraseña) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: ActualizarUsuario :exec
UPDATE usuarios SET nombre = $2, email = $3 WHERE id = $1;

-- name: ActualizarDNI :exec
UPDATE usuarios SET dni = $2 WHERE id = $1;

-- name: CrearPerfil :one
INSERT INTO perfiles (id_usuario, image) VALUES ($1, $2) RETURNING *;

-- name: EliminarUsuario :exec
DELETE FROM usuarios WHERE id = $1;

-- name: CrearNoticia :one
INSERT INTO noticias (titulo, contenido, tiempo_lectura_estimado) VALUES ($1, $2, $3) RETURNING *;

-- name: ListarNoticias :many
SELECT * FROM noticias ORDER BY publicada_en LIMIT 5 OFFSET $1;

-- name: LikearNoticia :one
INSERT INTO likes_noticia (id_noticia, id_usuario) VALUES ($1, $2) RETURNING *;

-- name: DeslikearNoticia :exec
DELETE FROM likes_noticia WHERE id_noticia = $1 AND id_usuario = $2;

-- name: ObtenerLikesNoticia :one
SELECT COUNT(*) FROM likes_noticia WHERE id_noticia = $1;

-- name: AgregarComentario :one
INSERT INTO comentarios_noticia (id_noticia, id_usuario, comentario) VALUES ($1, $2, $3) RETURNING *;

-- name: EliminarComentario :exec
DELETE FROM comentarios_noticia WHERE id_noticia = $1 AND id_usuario = $2 AND id_comentario = $3;

-- name: ObtenerComentariosNoticia :one
SELECT COUNT(*) FROM comentarios_noticia WHERE id_noticia = $1;

-- name: LikearComentario :one
INSERT INTO likes_comentario (id_noticia, id_usuario, id_comentario, id_usuario_like) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: DeslikearComentario :exec
DELETE FROM likes_comentario WHERE id_noticia = $1 AND id_usuario = $2 AND id_comentario = $3 AND id_usuario_like = $4;

-- name: ObtenerLikesComentario :one
SELECT COUNT(*) FROM likes_comentario WHERE id_noticia = $1 AND id_comentario = $2;

-- name: ListarComentarios :many
SELECT * FROM comentarios_noticia ORDER BY publicado_en LIMIT 10 OFFSET $1;

-- name: CrearDeporte :one
INSERT INTO deportes (nombre,foto) VALUES ($1,$2) RETURNING *;

-- name: ObtenerDeporte :one
SELECT * FROM deportes WHERE id = $1;

-- name: ObtenerDeportes :many
SELECT * FROM deportes 

-- name: CrearFacultad :one
INSERT INTO facultades (nombre) VALUES ($1) RETURNING *;


