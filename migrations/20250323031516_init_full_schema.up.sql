CREATE TABLE guests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    name TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    name TEXT,
    password_hash TEXT NOT NULL,
    phone TEXT,
    is_verified BOOLEAN DEFAULT false,
    is_admin BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    printful_id BIGINT,
    external_id TEXT,
    name TEXT,
    synced BOOLEAN,
    thumbnail_url TEXT
);

CREATE TABLE variants (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    printful_id BIGINT,
    product_id BIGINT REFERENCES products(id),
    name TEXT,
    sku TEXT,
    retail_price TEXT,
    thumbnail_url TEXT,
    size TEXT,
    color TEXT,
    availability_status TEXT
);

CREATE TABLE shipping_addresses (
    id SERIAL PRIMARY KEY,
    full_name TEXT NOT NULL,
    address_line1 TEXT NOT NULL,
    address_line2 TEXT,
    city TEXT NOT NULL,
    state TEXT NOT NULL,
    postal_code TEXT NOT NULL,
    country TEXT NOT NULL,
    phone TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    guest_name TEXT,
    guest_email TEXT,
    shipping_address TEXT NOT NULL,
    billing_address TEXT,
    status TEXT DEFAULT 'pending',
    total_amount NUMERIC NOT NULL,
    currency TEXT DEFAULT 'USD',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    variant_id BIGINT NOT NULL,
    product_name TEXT NOT NULL,
    size TEXT,
    color TEXT,
    quantity INTEGER NOT NULL DEFAULT 1,
    price_each NUMERIC NOT NULL,
    total_price NUMERIC
);
