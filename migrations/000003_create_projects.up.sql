-- 000003_create_projects.up.sql
-- user_projects, published_pages, media_assets (section 3.8-3.10)
-- user_projects.id menjadi referensi projectId di MongoDB

CREATE TABLE user_projects (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID           NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    template_id UUID           NOT NULL REFERENCES templates (id),
    title       VARCHAR(150)   NOT NULL,
    description TEXT,
    status      VARCHAR(20)    NOT NULL DEFAULT 'draft'
                CHECK (status IN ('draft', 'published', 'archived')),
    is_active   BOOLEAN        NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE INDEX idx_user_projects_user_id ON user_projects (user_id);
CREATE INDEX idx_user_projects_template_id ON user_projects (template_id);
CREATE INDEX idx_user_projects_status ON user_projects (status);

CREATE TABLE published_pages (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id   UUID           NOT NULL REFERENCES user_projects (id) ON DELETE CASCADE,
    slug         VARCHAR(150)   NOT NULL UNIQUE,
    url          TEXT           NOT NULL,
    is_active    BOOLEAN        NOT NULL DEFAULT true,
    published_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE INDEX idx_published_pages_project_id ON published_pages (project_id);

CREATE TABLE media_assets (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID           NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    project_id UUID           REFERENCES user_projects (id) ON DELETE SET NULL,
    file_name  VARCHAR(255)   NOT NULL,
    storage_key TEXT          NOT NULL,
    file_url   TEXT,
    file_type  VARCHAR(20)    NOT NULL DEFAULT 'other'
               CHECK (file_type IN ('image', 'audio', 'video', 'other')),
    file_size  BIGINT         NOT NULL DEFAULT 0 CHECK (file_size >= 0),
    mime_type  VARCHAR(100)   NOT NULL DEFAULT 'application/octet-stream',
    created_at TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE INDEX idx_media_assets_user_id ON media_assets (user_id);
CREATE INDEX idx_media_assets_project_id ON media_assets (project_id);
