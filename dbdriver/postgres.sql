
DROP TABLE if exists public.device cascade;
create table public.device (
    id            serial      not null
        constraint device_pk
            primary key,
    name          varchar(32) not null,
    lastipaddress varchar(15) not null
);

DROP TABLE if exists public.monitor cascade;
create table public.monitor (
    id           serial      not null
        constraint monitor_pk
            primary key,
    serialnumber varchar(32) not null,
    resolution   real not null
);

DROP TABLE if exists public.device_monitor cascade;
create table public.device_monitor (
    device_id  integer not null
        constraint device_monitor_device_id_fkey
            references public.device
            on update cascade on delete cascade,
    monitor_id integer not null
        constraint device_monitor_monitor_id_fkey
            references public.monitor
            on update cascade,
    constraint device_monitor_pkey
        primary key (device_id, monitor_id)
);
