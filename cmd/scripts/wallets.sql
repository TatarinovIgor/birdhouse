create table wallets
(
    id          integer not null
        constraint wallets_pk
            primary key
        constraint wallets_pk_2
            unique,
    merchant_id varchar(512),
    external_id varchar(512),
    key         varchar(512),
    value       varchar(2048),
    updated_at  timestamp,
    deleted_at  timestamp,
    created_at  timestamp
);
alter table wallets
    owner to postgres;

