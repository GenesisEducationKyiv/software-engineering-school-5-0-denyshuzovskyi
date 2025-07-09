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
    location_name VARCHAR(60)         NOT NULL,
    frequency     frequency           NOT NULL,
    status        subscription_status NOT NULL,
    created_at    TIMESTAMP           NOT NULL,
    updated_at    TIMESTAMP           NOT NULL,
    UNIQUE (subscriber_id, location_name)
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