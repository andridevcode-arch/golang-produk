-- 000005_create_access.up.sql
-- usage_quotas, project_entitlements (section 3.16-3.17)

CREATE TABLE usage_quotas (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID           NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    type         VARCHAR(30)    NOT NULL DEFAULT 'free_projects',
    total_quota  INTEGER        NOT NULL DEFAULT 0 CHECK (total_quota >= 0),
    used_quota   INTEGER        NOT NULL DEFAULT 0 CHECK (used_quota >= 0),
    reset_period VARCHAR(20)    NOT NULL DEFAULT 'one_time'
                 CHECK (reset_period IN ('one_time', 'monthly', 'yearly')),
    reset_at     TIMESTAMPTZ,
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT now(),

    CONSTRAINT uq_usage_quotas_user_type UNIQUE (user_id, type)
);

CREATE TABLE project_entitlements (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID           NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    product_id  UUID           REFERENCES products (id) ON DELETE CASCADE,
    template_id UUID           REFERENCES templates (id) ON DELETE CASCADE,
    access_type VARCHAR(20)    NOT NULL DEFAULT 'free'
                CHECK (access_type IN ('free', 'premium', 'purchased')),
    expires_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),

    -- entitlement harus mengacu pada target yang valid
    CONSTRAINT ck_entitlements_has_target CHECK (product_id IS NOT NULL OR template_id IS NOT NULL)
);

CREATE INDEX idx_entitlements_user_id ON project_entitlements (user_id);
CREATE INDEX idx_entitlements_product_id ON project_entitlements (product_id);
CREATE INDEX idx_entitlements_template_id ON project_entitlements (template_id);
