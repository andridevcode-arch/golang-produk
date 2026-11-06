-- 000010_create_published_pages.up.sql

CREATE TABLE published_pages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL UNIQUE REFERENCES user_projects(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    slug VARCHAR(150) NOT NULL,
    custom_path VARCHAR(255) NOT NULL,
    preview_token VARCHAR(64) NOT NULL UNIQUE,
    allow_preview BOOLEAN NOT NULL DEFAULT true,
    is_active BOOLEAN NOT NULL DEFAULT false,
    view_count BIGINT NOT NULL DEFAULT 0 CHECK (view_count >= 0),
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_published_pages_user_slug UNIQUE (user_id, slug)
);

CREATE INDEX idx_published_pages_user_slug ON published_pages(user_id, slug);
