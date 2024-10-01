-- +goose Up
CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR DEFAULT NULL,
    stripe_id VARCHAR DEFAULT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.usergroup (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.relation_users_usergroups (
  user_id UUID NOT NULL,
  usergroup_id UUID NOT NULL,
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE,
  FOREIGN KEY (usergroup_id) REFERENCES public.usergroup(id) ON DELETE CASCADE
);

CREATE SCHEMA IF NOT EXISTS userfiles;

CREATE TABLE IF NOT EXISTS userfiles.acl {
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  usergroup_id UUID NOT NULL,
  owner_id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (usergroup_id) REFERENCES public.usergroup(id) ON DELETE CASCADE,
  FOREIGN KEY (owner_id) REFERENCES public.users(id) ON DELETE CASCADE
  };
}

CREATE TABLE IF NOT EXISTS public.userfiles (
    url VARCHAR,
    doctype VARCHAR,
    lang VARCHAR,
    name VARCHAR,
    source VARCHAR,
    hash VARCHAR,
    mdata VARCHAR,
    stage VARCHAR,
    summary VARCHAR,
    short_summary VARCHAR,
    usergroup_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    -- when a usergroup is deleted all userfiles associated with it should
    FOREIGN KEY (usergroup_id) REFERENCES public.usergroup(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS public.users;
