\set previous_version 'v1.14.0-ee'
\set next_version 'v1.15.0-ee'
SELECT openreplay_version()                       AS current_version,
       openreplay_version() = :'previous_version' AS valid_previous,
       openreplay_version() = :'next_version'     AS is_next
\gset

\if :valid_previous
\echo valid previous DB version :'previous_version', starting DB upgrade to :'next_version'
BEGIN;
SELECT format($fn_def$
CREATE OR REPLACE FUNCTION openreplay_version()
    RETURNS text AS
$$
SELECT '%1$s'
$$ LANGUAGE sql IMMUTABLE;
$fn_def$, :'next_version')
\gexec

--
ALTER TABLE IF EXISTS events_common.requests
    ADD COLUMN transfer_size bigint NULL;

ALTER TABLE IF EXISTS public.sessions
    ADD COLUMN IF NOT EXISTS timezone text NULL;

ALTER TABLE IF EXISTS public.projects
    ADD COLUMN IF NOT EXISTS platform public.platform NOT NULL DEFAULT 'web';

COMMIT;

\elif :is_next
\echo new version detected :'next_version', nothing to do
\else
\warn skipping DB upgrade of :'next_version', expected previous version :'previous_version', found :'current_version'
\endif