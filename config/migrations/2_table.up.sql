CREATE SEQUENCE public.delivery_params_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.delivery_params_id_seq OWNER TO postgres;

CREATE SEQUENCE public.payments_payment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.payments_payment_id_seq OWNER TO postgres;

CREATE TABLE public.delivery_params (
                                        id bigint NOT NULL DEFAULT nextval('public.delivery_params_id_seq'::regclass),
                                        name text NOT NULL,
                                        phone text NOT NULL,
                                        zip text NOT NULL,
                                        city text NOT NULL,
                                        address text NOT NULL,
                                        region text NOT NULL,
                                        email text NOT NULL,
                                        CONSTRAINT delivery_params_pkey PRIMARY KEY (id)
);
ALTER TABLE public.delivery_params OWNER TO postgres;

CREATE TABLE public.payments (
                                 id bigint NOT NULL DEFAULT nextval('public.payments_payment_id_seq'::regclass),
                                 transaction uuid NOT NULL,
                                 request_id text,
                                 currency text NOT NULL,
                                 provider text NOT NULL,
                                 amount integer NOT NULL,
                                 payment_dt bigint NOT NULL,
                                 bank text NOT NULL,
                                 delivery_cost integer NOT NULL,
                                 goods_total integer NOT NULL,
                                 custom_fee integer NOT NULL,
                                 CONSTRAINT payments_pkey PRIMARY KEY (id)
);
ALTER TABLE public.payments OWNER TO postgres;

CREATE TABLE public.orders (
                               order_uid uuid NOT NULL,
                               track_number text NOT NULL,
                               entry text NOT NULL,
                               delivery_params_id bigint NOT NULL,
                               payment_id bigint NOT NULL,
                               locale text NOT NULL,
                               internal_signature text,
                               customer_id text NOT NULL,
                               delivery_service text NOT NULL,
                               shardkey text NOT NULL,
                               sm_id bigint NOT NULL,
                               date_created timestamp with time zone,
                               oof_shard text NOT NULL,
                               CONSTRAINT orders_pkey PRIMARY KEY (order_uid),
                               CONSTRAINT fk_delivery_params FOREIGN KEY (delivery_params_id) REFERENCES public.delivery_params(id),
                               CONSTRAINT fk_payment FOREIGN KEY (payment_id) REFERENCES public.payments(id)
);
ALTER TABLE public.orders OWNER TO postgres;

CREATE TABLE public.ordered_items (
                                      chrt_id integer NOT NULL,
                                      track_number text NOT NULL,
                                      price integer NOT NULL,
                                      rid uuid NOT NULL,
                                      name text NOT NULL,
                                      sale integer NOT NULL,
                                      size text NOT NULL,
                                      total_price integer NOT NULL,
                                      nm_id bigint NOT NULL,
                                      brand text NOT NULL,
                                      status smallint NOT NULL,
                                      order_uid uuid NOT NULL,
                                      CONSTRAINT fk_order_uid FOREIGN KEY (order_uid) REFERENCES public.orders(order_uid) NOT VALID
);
ALTER TABLE public.ordered_items OWNER TO postgres;

ALTER SEQUENCE public.delivery_params_id_seq OWNED BY public.delivery_params.id;
ALTER SEQUENCE public.payments_payment_id_seq OWNED BY public.payments.id;

GRANT ALL ON SEQUENCE public.delivery_params_id_seq TO wb_lvl0;
GRANT ALL ON SEQUENCE public.payments_payment_id_seq TO wb_lvl0;
