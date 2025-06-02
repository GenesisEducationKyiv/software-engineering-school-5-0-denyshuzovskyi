CREATE TABLE location
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(60) NOT NULL UNIQUE
);

CREATE TABLE subscriber
(
    id         SERIAL PRIMARY KEY,
    email      VARCHAR(60) NOT NULL UNIQUE,
    created_at TIMESTAMP   NOT NULL
);

CREATE TABLE subscription
(
    id            SERIAL PRIMARY KEY,
    subscriber_id INT                 NOT NULL
        REFERENCES subscriber (id) ON DELETE CASCADE,
    location_id   INT                 NOT NULL
        REFERENCES location (id) ON DELETE CASCADE,
    frequency     frequency           NOT NULL,
    status        subscription_status NOT NULL,
    created_at    TIMESTAMP           NOT NULL,
    updated_at    TIMESTAMP           NOT NULL,
    UNIQUE (subscriber_id, location_id)
);


CREATE TABLE token
(
    token           CHAR(36) PRIMARY KEY,
    subscription_id INT        NOT NULL
        REFERENCES subscription (id) ON DELETE CASCADE,
    type            token_type NOT NULL,
    created_at      TIMESTAMP  NOT NULL,
    expires_at      TIMESTAMP  NOT NULL,
    used_at         TIMESTAMP
);

CREATE TABLE weather
(
    location_id  INT           NOT NULL REFERENCES location (id),
    last_updated TIMESTAMP     NOT NULL,
    fetched_at   TIMESTAMP     NOT NULL,
    temperature  NUMERIC(5, 2) NOT NULL,
    humidity     NUMERIC(5, 2) NOT NULL CHECK (humidity BETWEEN 0 AND 100),
    description  VARCHAR(400)  NOT NULL,
    PRIMARY KEY (location_id, last_updated)
);