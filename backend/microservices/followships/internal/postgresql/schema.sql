create table followships (
  id int primary key,
  follower_id int not null,
  followee_id int not null,
  unique (follower_id, followee_id)
);

create index follower_id_idx_followships on followships (follower_id, followee_id);

create index followee_id_idx_followships on followships (followee_id, follower_id);
