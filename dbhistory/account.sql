use autotrade;
create table account
(
    id  varchar(10)           not null,
    btc float(8, 8) default 0 not null comment 'btc count',
    eth float(8, 8) default 0 not null comment 'eth count',
    bnb float(8, 8) default 0 not null comment 'bnb count',
    eos float(8, 8) default 0 not null comment 'eos count',
    xrp float(8, 8) default 0 not null comment 'xrp count',
    constraint account_pk
        primary key (id)
)
    comment 'account status';
