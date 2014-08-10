--
-- This is the schema used in leejo backend
--

create table leejo_user (
  user_id serial,
  username varchar(255) not null default '',
  gender char(1) not null default 'M',
  primary key (user_id)
);

create table leejo_user_skill (
  user_skill_id serial,
  user_id integer not null default 0,
  skill_name varchar(255) not null default '',
  primary key (user_skill_id)
);

create table leejo_user_interest (
  user_interest_id serial,
  user_id integer default 0,
  interest_name varchar(255) not null default '',
  primary key (user_interest_id)
);

create table leejo_api_client (
  id varchar(255) default '',
  secret varchar(255) default '',
  redirect_uri varchar(255) default '',
  primary key (id)
);

create table leejo_api_authdata (
  id serial,
  code varchar(255) default '' unique,
  user_id integer default 0,
  client_id varchar(255) default '',
  scope varchar(255) default '',
  state varchar(255) default '',
  redirect_uri varchar(255) default '',
  expired_timestamp integer,
  created_timestamp integer,
  primary key (id)
);

create table leejo_api_access (
  id serial,
  access_token varchar(255) default '' unique,
  refresh_token varchar(255) default '',
  client_id varchar(255) default '',
  user_id integer default 0,
  scope varchar(255) default '',
  expired_timestamp integer,
  created_timestamp integer,
  primary key (id)
);
