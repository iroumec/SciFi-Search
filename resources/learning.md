# Learning

En PostgreSQL **da lo mismo en rendimiento** usar `VARCHAR(255)` o `TEXT`. A diferencia de MySQL, Postgres no limita internamente `VARCHAR(n)` a `n` bytes: solo lo valida al insertar/actualizar.

## Diferencias prácticas

- `VARCHAR(255)`

  - Tiene una restricción de longitud máxima.
  - Si intentás guardar un hash mayor (no es el caso de bcrypt, siempre \~60 chars), fallará.
  - Es útil si querés **validar** que no se guarden cadenas más largas de lo esperado.

- `TEXT`

  - No tiene límite.
  - Ideal si no necesitás restricciones de longitud.
  - Más flexible para futuros cambios (si más adelante cambiás de bcrypt a Argon2, por ejemplo, que produce hashes más largos).

---

## Recomendación

- Para contraseñas **TEXT** es la opción más robusta y flexible.
- Si preferís un límite “de seguridad”, poné `VARCHAR(255)` (más que suficiente para bcrypt/argon2/scrypt).

---

¿Querés que te prepare un `ALTER TABLE` directo a `TEXT` y además te muestre cómo recrear el `sqlc` con ese cambio?
