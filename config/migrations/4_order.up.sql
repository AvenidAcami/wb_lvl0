WITH delivery_ins AS (
INSERT INTO public.delivery_params (name, phone, zip, city, address, region, email)
VALUES (
    'Test Testov',
    '+9720000000',
    '2639809',
    'Kiryat Mozkin',
    'Ploshad Mira 15',
    'Kraiot',
    'test@gmail.com'
    )
    RETURNING id
    ),
    payment_ins AS (
INSERT INTO public.payments (
    transaction, request_id, currency, provider, amount, payment_dt,
    bank, delivery_cost, goods_total, custom_fee
)
VALUES (
    'a7686192-cbad-43bc-9116-f92554fb51c6',
    '',
    'USD',
    'wbpay',
    1817,
    1637907727,
    'alpha',
    1500,
    317,
    0
    )
    RETURNING id
    )
INSERT INTO public.orders (
    order_uid,
    track_number,
    entry,
    delivery_params_id,
    payment_id,
    locale,
    internal_signature,
    customer_id,
    delivery_service,
    shardkey,
    sm_id,
    date_created,
    oof_shard
)
SELECT
    'a7686192-cbad-43bc-9116-f92554fb51c6',
    'WBILMTESTTRACK',
    'WBIL',
    delivery_ins.id,
    payment_ins.id,
    'en',
    '',
    'test',
    'meest',
    '9',
    99,
    '2021-11-26T06:22:19Z',
    '1'
FROM delivery_ins, payment_ins;

-- Теперь вставляем items
INSERT INTO public.ordered_items (
    chrt_id,
    track_number,
    price,
    rid,
    name,
    sale,
    size,
    total_price,
    nm_id,
    brand,
    status,
    order_uid
) VALUES (
             9934930,
             'WBILMTESTTRACK',
             453,
             '8d5365b2-84e0-47aa-a5ab-3fdd10af0152',
             'Mascaras',
             30,
             '0',
             317,
             2389212,
             'Vivienne Sabo',
             202,
             'a7686192-cbad-43bc-9116-f92554fb51c6'
         );