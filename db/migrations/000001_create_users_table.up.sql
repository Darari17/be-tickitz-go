CREATE TABLE
  public.users (
    id uuid NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    role character varying(50) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NULL
  );

ALTER TABLE
  public.users
ADD
  CONSTRAINT users_pkey PRIMARY KEY (id)