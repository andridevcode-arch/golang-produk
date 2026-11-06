-- 000004_create_billing.up.sql
-- plans, subscriptions, orders, order_items, payments (section 3.11-3.15)

CREATE TABLE plans (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(50)    NOT NULL,
    type          VARCHAR(20)    NOT NULL CHECK (type IN ('free', 'premium')),
    price         DECIMAL(12,2)  NOT NULL DEFAULT 0 CHECK (price >= 0),
    duration_days INTEGER,
    features      JSONB          NOT NULL DEFAULT '[]'::jsonb,
    is_active     BOOLEAN        NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ    NOT NULL DEFAULT now(),

    CONSTRAINT ck_plans_duration_positive CHECK (duration_days IS NULL OR duration_days > 0)
);

CREATE TABLE subscriptions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID           NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    plan_id    UUID           NOT NULL REFERENCES plans (id),
    status     VARCHAR(20)    NOT NULL DEFAULT 'pending'
               CHECK (status IN ('active', 'cancelled', 'expired', 'pending')),
    start_date TIMESTAMPTZ    NOT NULL DEFAULT now(),
    end_date   TIMESTAMPTZ,
    auto_renew BOOLEAN        NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions (user_id);
CREATE INDEX idx_subscriptions_plan_id ON subscriptions (plan_id);

-- cegah langganan aktif ganda per pengguna
CREATE UNIQUE INDEX uq_subscriptions_single_active
    ON subscriptions (user_id) WHERE status = 'active';

CREATE TABLE orders (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID           NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    total_amount DECIMAL(12,2)  NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
    status       VARCHAR(20)    NOT NULL DEFAULT 'pending'
                 CHECK (status IN ('pending', 'paid', 'cancelled', 'refunded')),
    currency     CHAR(3)        NOT NULL DEFAULT 'IDR',
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE INDEX idx_orders_user_id ON orders (user_id);
CREATE INDEX idx_orders_status ON orders (status);

CREATE TABLE order_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id    UUID           NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    product_id  UUID           REFERENCES products (id) ON DELETE SET NULL,
    template_id UUID           REFERENCES templates (id) ON DELETE SET NULL,
    description VARCHAR(255)   NOT NULL,
    unit_price  DECIMAL(12,2)  NOT NULL DEFAULT 0 CHECK (unit_price >= 0),
    quantity    INTEGER        NOT NULL DEFAULT 1 CHECK (quantity > 0),
    subtotal    DECIMAL(12,2)  NOT NULL DEFAULT 0 CHECK (subtotal >= 0),
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE INDEX idx_order_items_order_id ON order_items (order_id);
CREATE INDEX idx_order_items_product_id ON order_items (product_id);
CREATE INDEX idx_order_items_template_id ON order_items (template_id);

CREATE TABLE payments (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id       UUID           NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    amount         DECIMAL(12,2)  NOT NULL CHECK (amount >= 0),
    method         VARCHAR(50)    NOT NULL,
    status         VARCHAR(20)    NOT NULL DEFAULT 'pending'
                   CHECK (status IN ('pending', 'paid', 'failed', 'refunded')),
    transaction_id VARCHAR(255),
    payment_url    TEXT,
    paid_at        TIMESTAMPTZ,
    created_at     TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE INDEX idx_payments_order_id ON payments (order_id);

-- id transaksi provider unik jika ada (untuk idempoten webhook)
CREATE UNIQUE INDEX uq_payments_transaction_id
    ON payments (transaction_id) WHERE transaction_id IS NOT NULL;
