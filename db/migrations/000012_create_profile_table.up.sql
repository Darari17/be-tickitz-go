CREATE TABLE
  public.profile (
    user_id uuid NOT NULL,
    firstname character varying(100) NULL,
    lastname character varying(100) NULL,
    phone_number character varying(20) NULL,
    avatar text NULL,
    point integer NULL DEFAULT 0,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NULL
  );

ALTER TABLE
  public.profile
ADD
  CONSTRAINT profile_pkey PRIMARY KEY (user_id)