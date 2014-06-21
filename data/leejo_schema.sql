--
-- This is the schema used in leejo backend
--

create table leejo_user (
  user_id serial,
  username varchar(255) not null default '',
  gender char(1) not null default 'M',
  primary key (user_id)
);

create table leejo_skill (
  user_id integer not null default 0,
  skill_name varchar(255) not null default '',
  primary key (user_id, skill_name)
);

create table leejo_interest (
  user_id integer default 0,
  interest_name varchar(255) not null default '',
  primary key (user_id, interest_name)
);
