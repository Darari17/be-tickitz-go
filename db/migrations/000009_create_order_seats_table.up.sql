CREATE TABLE
  public.order_seats (
    orders_id integer NOT NULL,
    seats_id integer NOT NULL
  );

ALTER TABLE
  public.order_seats
ADD
  CONSTRAINT order_seats_pkey PRIMARY KEY (orders_id, seats_id)