create table autotrade.platform
(
    id          int auto_increment,
    name        varchar(15)                                                     not null,
    name_cn     varchar(10)                                                     null comment 'Chinese name',
    create_time timestamp default CURRENT_TIMESTAMP                             not null,
    update_time timestamp default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP not null,
    constraint platform_pk
        primary key (id),
    constraint platform_name_uindex
        unique (name)
);

create table autotrade.quantity
(
    id          int auto_increment,
    name        varchar(15)                                                           not null,
    platform    varchar(15)                                                           not null,
    free        decimal(18, 10) default 0.0000000000                                  not null,
    locked      decimal(18, 10) default 0.0000000000                                  not null,
    create_time timestamp       default CURRENT_TIMESTAMP                             not null,
    update_time timestamp       default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP not null,
    constraint quantity_pk
        primary key (id),
    constraint quantity_symbol_platform_uindex
        unique (name, platform)
);

create table autotrade.account
(
    id          int auto_increment,
    platform    varchar(15)                                                           not null,
    usdt        decimal(18, 10) default 0.0000000000                                  not null comment 'usdt  count',
    usdt_locked decimal(18, 10) default 0.0000000000                                  not null,
    btc         decimal(18, 10) default 0.0000000000                                  not null comment 'btc count',
    btc_locked  decimal(18, 10) default 0.0000000000                                  not null,
    eth         decimal(18, 10) default 0.0000000000                                  not null comment 'eth count',
    eth_locked  decimal(18, 10) default 0.0000000000                                  not null,
    bnb         decimal(18, 10) default 0.0000000000                                  not null comment 'bnb count',
    bnb_locked  decimal(18, 10) default 0.0000000000                                  not null,
    eos         decimal(18, 10) default 0.0000000000                                  not null comment 'eos count',
    eos_locked  decimal(18, 10) default 0.0000000000                                  not null,
    xrp         decimal(18, 10) default 0.0000000000                                  not null comment 'xrp count',
    xrp_locked  decimal(18, 10) default 0.0000000000                                  not null,
    create_time timestamp       default CURRENT_TIMESTAMP                             not null,
    update_time timestamp       default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP not null,
    constraint account_pk
        primary key (id),
    constraint account_platform_uindex
        unique (platform)
)
    comment 'account status';

create table autotrade.strategy_low_buy_high_sell
(
    id                  int auto_increment,
    symbol              varchar(15)                                                           not null,
    coin_name           varchar(11)                                                           not null,
    platform            varchar(15)                                                           not null,
    quantity            decimal(18, 10) default 0                                             not null comment 'total count',
    spend               decimal(18, 10) default 0                                             not null comment 'total spend, But not the actual total consumption',
    position_average    decimal(18, 10) default 0                                             not null comment 'Average purchase price',
    last_spend          decimal(18, 10) default 0                                             not null,
    spend_coefficient   double(4, 3)    default 2                                             not null,
    target_profit_point double(4, 3)    default 1.05                                          not null comment 'Sell when the current price/average position price is higher than this value',
    target_sell_price   decimal(18, 10) default 0                                             not null,
    target_buy_point    double(4, 3)    default 0.92                                          not null comment 'Buy when the current price/average price is lower than this value',
    target_buy_price    decimal(18, 10) default 0                                             not null,
    month_average       decimal(18, 10) default 0                                             not null comment 'Average market price for a month (30 days)',
    status              int(1)          default 0                                             not null,
    actual_cost         decimal(18, 10) default 0                                             not null comment 'Actual total cost, which can be used to calculate the return',
    create_time         datetime        default CURRENT_TIMESTAMP                             not null,
    update_time         datetime        default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP not null,
    constraint strategy_low_buy_high_sell_pk
        primary key (id),
    constraint strategy_low_buy_high_sell_symbol_platform_uindex
        unique (symbol, platform)
);

