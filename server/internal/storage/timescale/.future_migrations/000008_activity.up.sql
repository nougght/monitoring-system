
-- -- интервалы сессий пользователя
-- CREATE TABLE user_sessions (
--   agent_id  UUID NOT NULL,
--   logon_at  TIMESTAMPTZ NOT NULL,
--   logoff_at TIMESTAMPTZ,
--   os_user   TEXT NOT NULL,
--   remote    BOOLEAN,
--   meta      JSONB
-- );

-- SELECT create_hypertable('user_sessions', 'logon_at',
--                          chunk_time_interval => INTERVAL '7 day');

-- -- интервалы активности
-- CREATE TABLE activity_intervals (
--   agent_id   UUID NOT NULL,
--   started_at TIMESTAMPTZ NOT NULL,
--   ended_at   TIMESTAMPTZ,
--   kind       SMALLINT NOT NULL,  
--   app_id     INT,                
--   category   SMALLINT,            
--   title      TEXT,               
--   meta       JSONB
-- );

-- SELECT create_hypertable('activity_intervals', 'started_at',
--                          chunk_time_interval => INTERVAL '1 day');
             
-- -- события по активности
-- CREATE TABLE activity_events (
--   time     TIMESTAMPTZ NOT NULL,
--   agent_id UUID NOT NULL,
--   kind     SMALLINT NOT NULL,     
--   app_id   INT,
--   meta     JSONB
-- );

-- SELECT create_hypertable('activity_events', 'started_at',
--                          chunk_time_interval => INTERVAL '1 day');
-- -- доступность агента
-- CREATE TABLE agent_availability (
--   agent_id   UUID NOT NULL,
--   started_at TIMESTAMPTZ NOT NULL,
--   ended_at   TIMESTAMPTZ,
--   state      SMALLINT NOT NULL  
-- );


-- SELECT create_hypertable('agent_availability', 'started_at',
--                          chunk_time_interval => INTERVAL '7 day');