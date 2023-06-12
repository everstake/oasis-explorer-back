CREATE TABLE IF NOT EXISTS tasks (
   tsk_id             serial constraint tasks_pk primary key,
   tsk_active         boolean default TRUE,
   tsk_title          VARCHAR(255) NOT NULL,
   tsk_start_height   int        NOT NULL,
   tsk_current_height int      NOT NULL,
   tsk_end_height     int      NOT NULL,
   tsk_batch          int      NOT NULL
);

create table if not exists blocks
(
    id           serial not null constraint blocks_pk primary key,
    total_count  int8 default 0 not null,
    last_lvl     int8 default 0 not null
);

create table if not exists day_blocks
(
    id              serial not null constraint day_blocks_pk primary key,
    day_total_count int8 default 0 not null,
    day             timestamp default date_trunc('day', now()) not null
);

create table if not exists validators
(
    id              serial not null constraint validators_pk primary key,
    address         varchar(46) not null,
    total_blk_count int8 default 0 not null,
    total_sig_count int8 default 0 not null,
    last_blk_time   timestamp default now() not null,
    last_sig_time   timestamp default now() not null
);

create table if not exists validator_day_stats
(
    id              serial not null constraint validator_day_stats_pk primary key,
    val_id          int4 references validators(id) not null,
    day_blk_count   int8 default 0 not null,
    day_sig_count   int8 default 0 not null,
    day             timestamp default date_trunc('day', now()) not null
);

insert into blocks (total_count, last_lvl) values (0, 0);
insert into day_blocks (day_total_count) values (0);