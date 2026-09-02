-- виды метрик, описаны в proto и согласованы с агентом
CREATE TABLE
    metric_kinds (
        kind SMALLINT PRIMARY KEY,
        key TEXT UNIQUE NOT NULL, 
        unit TEXT NOT NULL,     -- единица измерения
        agg SMALLINT NOT NULL,  -- тип агрегации 
        label_name TEXT,        -- название метки
        description TEXT,
        deprecated BOOLEAN NOT NULL DEFAULT false  -- старые метрики не удаляются
    );

-- уникальный id для метрик по агенту, виду и метке метрики
CREATE TABLE metric_series (
  id       INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  agent_id UUID NOT NULL REFERENCES agents(id),
  kind     SMALLINT NOT NULL REFERENCES metric_kinds(kind),  -- вид метрики
  label    TEXT NOT NULL DEFAULT '',   -- метка для метрик одного вида(например названия дисков)
  UNIQUE (agent_id, kind, label)
);


-- основные числовые метрики 
-- без primary key
CREATE TABLE metric_samples (
  time      TIMESTAMPTZ NOT NULL,
  series_id INT NOT NULL,    -- без foreign key
  value     DOUBLE PRECISION NOT NULL
);
-- разделение на чанки по 6 часов
SELECT create_hypertable('metric_samples', 'time',
                         chunk_time_interval => INTERVAL '6 hour');

CREATE INDEX ON metric_samples (series_id, time DESC);

ALTER TABLE metric_samples SET (
  timescaledb.compress,
  timescaledb.compress_segmentby = 'series_id',
  timescaledb.compress_orderby   = 'time DESC'
);
-- сжатие на x2 интервала чанка
SELECT add_compression_policy('metric_samples', INTERVAL '12 hours');
