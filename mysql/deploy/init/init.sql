create database test;
use test;
create table user
(
    id int auto_increment primary key,
    username varchar(64) not null
);
insert into user values(1, "小明");
insert into user values(2, "小张");
insert into user values(3, "小李");