CREATE TABLE
  public.seats (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
    seat_code character varying(10) NOT NULL
  );

ALTER TABLE
  public.seats
ADD
  CONSTRAINT seats_pkey PRIMARY KEY (id)