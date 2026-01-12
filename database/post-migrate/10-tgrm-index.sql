-- ------------------------------------------------------------------------------------------------
-- Trigramas
-- ------------------------------------------------------------------------------------------------

-- Verificación de que la extensión esté instalada
-- y de que las tablas necesarias existan.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_schema = 'public'
            AND table_name = 'user_preferences'
    ) THEN
        CREATE INDEX IF NOT EXISTS idx_user_preferences_preference_trgm
        ON public.user_preferences
        USING gin (preference gin_trgm_ops);
    END IF;
END $$;

-- ------------------------------------------------------------------------------------------------

-- Creación de índice.
-- Esto tiene el objetivo de mejorar la eficiencia y la escalabilidad.
CREATE INDEX IF NOT EXISTS idx_user_preferences_preference_trgm 
ON user_preferences 
USING gin (preference gin_trgm_ops);

-- ------------------------------------------------------------------------------------------------

-- Definición de umbral de similitud.
ALTER SYSTEM SET pg_trgm.similarity_threshold = 0.4;
SELECT pg_reload_conf();


-- ------------------------------------------------------------------------------------------------