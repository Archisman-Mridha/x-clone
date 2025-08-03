-- name: CreateFollowship :exec
insert
into
  followships (
    follower_id,
    followee_id
  )
values
  (
    $1,
    $2
  );

-- name: DeleteFollowship :exec
delete
from
  followships
where
  follower_id = $1 and
  followee_id = $2;

-- name: GetFollowship :exec
select 1
from
  followships
where
  follower_id = $1 and
  followee_id = $2
limit 1;

-- name: GetFollowers :many
select
  follower_id
from
  followships
where
  followee_id = $1
limit $2
offset $3;

-- name: GetFollowees :many
select
  followee_id
from
  followships
where
  follower_id = $1
limit $2
offset $3;

-- name: GetFollowerAndFolloweeCounts :one
select
	(select count(*) from followships where followships.followee_id = $1) as follower_count,
	(select count(*) from followships where followships.follower_id = $1) as followee_count;
