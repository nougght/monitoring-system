SELECT remove_compression_policy('metric_samples', if_exists => true)
DROP TABLE IF EXISTS metric_samples CASCADE;


DROP TABLE IF EXISTS metric_series CASCADE;
DROP TABLE IF EXISTS metric_kinds CASCADE;