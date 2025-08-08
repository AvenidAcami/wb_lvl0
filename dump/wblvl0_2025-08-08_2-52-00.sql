--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4
-- Dumped by pg_dump version 17.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: delivery_params; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.delivery_params (
    id bigint NOT NULL,
    name text NOT NULL,
    phone text NOT NULL,
    zip text NOT NULL,
    city text NOT NULL,
    address text NOT NULL,
    region text NOT NULL,
    email text NOT NULL
);


ALTER TABLE public.delivery_params OWNER TO postgres;

--
-- Name: delivery_params_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.delivery_params_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.delivery_params_id_seq OWNER TO postgres;

--
-- Name: delivery_params_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.delivery_params_id_seq OWNED BY public.delivery_params.id;


--
-- Name: ordered_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ordered_items (
    item_id bigint NOT NULL,
    order_id bigint NOT NULL,
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
    status smallint NOT NULL
);


ALTER TABLE public.ordered_items OWNER TO postgres;

--
-- Name: ordered_items_item_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ordered_items_item_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ordered_items_item_id_seq OWNER TO postgres;

--
-- Name: ordered_items_item_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ordered_items_item_id_seq OWNED BY public.ordered_items.item_id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

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
    date_created timestamp with time zone DEFAULT now(),
    oof_shard text NOT NULL
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- Name: payments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payments (
    id bigint NOT NULL,
    transaction uuid NOT NULL,
    request_id text,
    currency text NOT NULL,
    provider text NOT NULL,
    amount integer NOT NULL,
    payment_dt bigint NOT NULL,
    bank text NOT NULL,
    delivery_cost integer NOT NULL,
    goods_total integer NOT NULL,
    custom_fee integer NOT NULL
);


ALTER TABLE public.payments OWNER TO postgres;

--
-- Name: payments_payment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.payments_payment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payments_payment_id_seq OWNER TO postgres;

--
-- Name: payments_payment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.payments_payment_id_seq OWNED BY public.payments.id;


--
-- Name: delivery_params id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.delivery_params ALTER COLUMN id SET DEFAULT nextval('public.delivery_params_id_seq'::regclass);


--
-- Name: ordered_items item_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ordered_items ALTER COLUMN item_id SET DEFAULT nextval('public.ordered_items_item_id_seq'::regclass);


--
-- Name: payments id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments ALTER COLUMN id SET DEFAULT nextval('public.payments_payment_id_seq'::regclass);


--
-- Data for Name: delivery_params; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.delivery_params (id, name, phone, zip, city, address, region, email) FROM stdin;
\.


--
-- Data for Name: ordered_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.ordered_items (item_id, order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) FROM stdin;
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.orders (order_uid, track_number, entry, delivery_params_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) FROM stdin;
\.


--
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.payments (id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) FROM stdin;
\.


--
-- Name: delivery_params_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.delivery_params_id_seq', 1, false);


--
-- Name: ordered_items_item_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ordered_items_item_id_seq', 1, false);


--
-- Name: payments_payment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.payments_payment_id_seq', 1, false);


--
-- Name: delivery_params delivery_params_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.delivery_params
    ADD CONSTRAINT delivery_params_pkey PRIMARY KEY (id);


--
-- Name: ordered_items ordered_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ordered_items
    ADD CONSTRAINT ordered_items_pkey PRIMARY KEY (item_id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (order_uid);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- Name: orders delivery_params; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT delivery_params FOREIGN KEY (delivery_params_id) REFERENCES public.delivery_params(id);


--
-- Name: orders payment; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT payment FOREIGN KEY (payment_id) REFERENCES public.payments(id);


--
-- Name: SEQUENCE delivery_params_id_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON SEQUENCE public.delivery_params_id_seq TO wb_lvl0;


--
-- Name: SEQUENCE ordered_items_item_id_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON SEQUENCE public.ordered_items_item_id_seq TO wb_lvl0;


--
-- Name: SEQUENCE payments_payment_id_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON SEQUENCE public.payments_payment_id_seq TO wb_lvl0;


--
-- PostgreSQL database dump complete
--

