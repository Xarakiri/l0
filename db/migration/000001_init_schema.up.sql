CREATE TABLE delivery
(
    id      varchar(128) PRIMARY KEY,
    name    varchar(128),
    phone   varchar(16),
    zip     varchar(128),
    city    varchar(128),
    address varchar(256),
    region  varchar(256),
    email   varchar(128)
);

CREATE TABLE payment
(
    transaction   varchar(128) PRIMARY KEY,
    request_id    varchar(128),
    currency      varchar(128),
    provider      varchar(128),
    amount        int,
    payment_dt    bigint,
    bank          varchar(128),
    delivery_cost int,
    goods_total   int,
    custom_fee    int
);

CREATE TABLE item
(
    rid          varchar(256) PRIMARY KEY,
    chrt_id      bigint,
    track_number varchar(256),
    price        int,
    name         varchar(128),
    sale         int,
    size         varchar(256),
    total_price  int,
    nm_id        bigint,
    brand        varchar(256),
    status       int
);

CREATE TABLE "order"
(
    order_uid          varchar(256) PRIMARY KEY,
    track_number       varchar(128),
    entry              varchar(128),
    delivery_id        varchar(128),
    payment_id         varchar(128),
    locale             varchar(128),
    internal_signature varchar(128),
    customer_id        varchar(128),
    delivery_service   varchar(128),
    shardkey           varchar(128),
    sm_id              int,
    date_created       timestamp,
    oof_shard          varchar(128),

    FOREIGN KEY (delivery_id) REFERENCES delivery (id),
    FOREIGN KEY (payment_id) REFERENCES payment (transaction)
);

CREATE TABLE order_item
(
    order_id varchar(256) REFERENCES "order" (order_uid) ON UPDATE CASCADE ON DELETE CASCADE,
    item_id  varchar(256) REFERENCES item (rid) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX ON "order" ("order_uid");

-- STORED PROCEDURES
-- SET DELIVERY
CREATE OR REPLACE PROCEDURE set_delivery(
    _id varchar,
    _name varchar,
    _phone varchar,
    _zip varchar,
    _city varchar,
    _address varchar,
    _region varchar,
    _email varchar
)
    language plpgsql
AS
$$
BEGIN
    INSERT INTO delivery VALUES (_id, _name, _phone, _zip, _city, _address, _region, _email);
END
$$;


-- SET PAYMENT
CREATE OR REPLACE PROCEDURE set_payment(
    _transaction varchar,
    _request_id varchar,
    _currency varchar,
    _provider varchar,
    _amount int,
    _payment_dt bigint,
    _bank varchar,
    _delivery_cost int,
    _goods_total int,
    _custom_fee int
)
    language plpgsql
AS
$$
BEGIN
    INSERT INTO payment
    VALUES (_transaction, _request_id, _currency, _provider, _amount, _payment_dt, _bank, _delivery_cost, _goods_total,
            _custom_fee);
END;
$$;

-- SET ITEM
CREATE OR REPLACE PROCEDURE set_item(
    _rid varchar,
    _chrt_id bigint,
    _track_number varchar,
    _price int,
    _name varchar,
    _sale int,
    _size varchar,
    _total_price int,
    _nm_id bigint,
    _brand varchar,
    _status int
)
    language plpgsql
AS
$$
BEGIN
    INSERT INTO item
    VALUES (_rid, _chrt_id, _track_number, _price, _name, _sale, _size, _total_price, _nm_id, _brand, _status);
END;
$$;

-- SET ORDER
CREATE OR REPLACE PROCEDURE set_order(
    _order_uid varchar,
    _track_number varchar,
    _entry varchar,
    _delivery_id varchar,
    _payment_id varchar,
    _locale varchar,
    _internal_signature varchar,
    _customer_id varchar,
    _delivery_service varchar,
    _shardkey varchar,
    _sm_id int,
    _date_created timestamp,
    _oof_shard varchar
)
    language plpgsql
AS
$$
BEGIN
    INSERT INTO "order"
    VALUES (_order_uid, _track_number, _entry, _delivery_id, _payment_id, _locale,
            _internal_signature, _customer_id, _delivery_service, _shardkey, _sm_id,
            _date_created, _oof_shard);
END ;
$$;

-- SET ORDER_ITEM
CREATE OR REPLACE PROCEDURE set_order_item(
    _order_id varchar,
    _item_id varchar
)
    language plpgsql
AS
$$
BEGIN
    INSERT INTO order_item VALUES (_order_id, _item_id);
END;
$$;
