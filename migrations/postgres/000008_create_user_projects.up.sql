-- 000008_create_user_projects.up.sql

CREATE TABLE IF NOT EXISTS user_projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    template_id UUID NOT NULL REFERENCES templates(id) ON DELETE RESTRICT,
    template_version_id UUID NOT NULL REFERENCES template_versions(id) ON DELETE RESTRICT,
    title VARCHAR(150) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'published', 'archived')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_user_projects_user_created_at
    ON user_projects(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_projects_user_status
    ON user_projects(user_id, status);
CREATE INDEX IF NOT EXISTS idx_user_projects_template_id
    ON user_projects(template_id);
CREATE INDEX IF NOT EXISTS idx_user_projects_template_version_id
    ON user_projects(template_version_id);
