create table profiles (
  id int primary key,
  name varchar(25) not null,
  username varchar(25) not null unique,
  profile_picture_uri varchar(250)
);

create index username_idx_profiles on profiles (username);
