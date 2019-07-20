use autotrade;
create table account
(
    id   varchar(10)                          not null,
    usdt decimal(18, 10) default 0.0000000000 not null comment 'usdt count',
    btc  decimal(18, 10) default 0.0000000000 not null comment 'btc count',
    eth  decimal(18, 10) default 0.0000000000 not null comment 'eth count',
    bnb  decimal(18, 10) default 0.0000000000 not null comment 'bnb count',
    eos  decimal(18, 10) default 0.0000000000 not null comment 'eos count',
    xrp  decimal(18, 10) default 0.0000000000 not null comment 'xrp count',
    constraint account_pk
        primary key (id)
)
    comment 'account status';

create table strategy_low_buy_high_sell
(
    id                  int auto_increment,
    symol               varchar(15)                  not null,
    platform            varchar(15)                  not null,
#     quantity            decimal(18, 10) default 0    not null comment 'total count',
    spend               decimal(18, 10) default 0    not null comment 'total spend',
    position_average    decimal(18, 10) default 0    not null comment 'Average purchase price',
    last_spend          decimal(18, 10) default 0    not null,
    target_profit_point double(4, 3)    default 1.05 not null comment 'Sell when the current price/average position price is higher than this value',
    target_sell_price   decimal(18, 10) default 0    null,
    target_buy_point    double(4, 3)    default 0.92 not null comment 'Buy when the current price/average price is lower than this value',
    target_buy_price    decimal(18, 10) default 0    not null,
    month_average       decimal(18, 10) default 0    not null comment 'Average market price for a month (30 days)',
    constraint strategy_low_buy_high_sell_pk
        primary key (id)
);

