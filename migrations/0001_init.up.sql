CREATE TABLE subscriptions (
    id uuid PRIMARY KEY,
    service_name text NOT NULL,
    price integer NOT NULL CHECK (price >= 0),
    user_id uuid NOT NULL,
    start_date date NOT NULL,
    end_date date NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
