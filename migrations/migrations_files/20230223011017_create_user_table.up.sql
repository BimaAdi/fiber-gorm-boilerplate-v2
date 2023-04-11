CREATE TABLE IF NOT EXISTS public."user" (
	id uuid NOT NULL,
	email varchar NOT NULL,
	username varchar NOT NULL,
	"password" varchar NOT NULL,
	is_active bool NULL DEFAULT true,
	is_superuser bool NULL DEFAULT false,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	CONSTRAINT user_pkey PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_user_email ON public."user" USING btree (email);
CREATE INDEX IF NOT EXISTS idx_user_id ON public."user" USING btree (id);
CREATE UNIQUE INDEX idx_user_username ON public."user" USING btree (username);
