-- PostgreSQL
-- Os IDS são sempre gerados pelo backend com o tipo UUID v7

CREATE TABLE IF NOT EXISTS agents (
	id UUID PRIMARY KEY,
	member_id UUID NOT NULL,
	name VARCHAR(255) NOT NULL,
	avatar_url TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workspaces (
	id UUID PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_agents_member_id ON agents(member_id);
CREATE INDEX IF NOT EXISTS idx_workspaces_agent_id ON workspaces(agent_id);

CREATE TABLE IF NOT EXISTS forms (
	id UUID PRIMARY KEY,
	workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
	agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
	name VARCHAR(255) NOT NULL,
	description VARCHAR(255) NOT NULL DEFAULT '',
	is_public BOOLEAN NOT NULL DEFAULT FALSE,
	current_version_id UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS form_versions (
	id UUID PRIMARY KEY,
	form_id UUID NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
	version_number INTEGER NOT NULL CHECK (version_number > 0),
	props JSONB,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CONSTRAINT uq_form_versions_form_version UNIQUE (form_id, version_number),
	CONSTRAINT uq_form_versions_form_id_id UNIQUE (form_id, id)
);

DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'fk_forms_current_version'
	) THEN
		ALTER TABLE forms
			ADD CONSTRAINT fk_forms_current_version
			FOREIGN KEY (id, current_version_id)
			REFERENCES form_versions(form_id, id)
			DEFERRABLE INITIALLY DEFERRED;
	END IF;
END
$$;

CREATE INDEX IF NOT EXISTS idx_forms_workspace_id ON forms(workspace_id);
CREATE INDEX IF NOT EXISTS idx_forms_agent_id ON forms(agent_id);
CREATE INDEX IF NOT EXISTS idx_form_versions_form_id ON form_versions(form_id);