CREATE TYPE status AS ENUM ('REGISTERED', 'NEW', 'INVALID', 'PROCESSING', 'PROCESSED');
CREATE TABLE orders (
    order_id TEXT NOT NULL UNIQUE,
    user_id SERIAL PRIMARY KEY REFERENCES users (id),
    accrual INTEGER,
    order_status status NOT NULL,
    uploaded DATE NOT NULL DEFAULT CURRENT_DATE
);
