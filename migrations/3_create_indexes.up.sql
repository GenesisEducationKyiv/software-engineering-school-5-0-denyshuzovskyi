CREATE INDEX idx_token_subscription_id_type ON token (subscription_id, type);
CREATE INDEX idx_weather_location_id_last_updated ON weather (location_id, last_updated DESC);