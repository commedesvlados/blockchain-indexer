CREATE TABLE gas_logs
(
    id            serial       not null unique,
    block_hash    varchar(255) not null unique,
    block_number  int not null unique,
    average_gas   int not null
);

CREATE TABLE token_logs
(
    id            serial       not null unique,
    block_hash    varchar(255) not null unique,
    block_number  int not null unique,
    log_name      varchar(255) not null,
    log_index     int not null,
    address_from  varchar(255) not null unique,
    address_to    varchar(255) not null unique,
    token_value   varchar(255) not null
);

