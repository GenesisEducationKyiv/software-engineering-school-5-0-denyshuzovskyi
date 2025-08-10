CREATE TYPE frequency AS ENUM ('hourly', 'daily');
CREATE TYPE subscription_status AS ENUM ('pending', 'confirmed');
CREATE TYPE token_type AS ENUM ('confirmation', 'unsubscribe');