-- name: CreatePost :one
insert into
  posts(
		owner_id,
		description
	)
values
	(
		$1,
		$2
	)
returning
  id;

-- name: GetPostsOfUser :many
select
  id,
  owner_id,
  description,
  created_at
from
  posts
where
  owner_id = $1
limit $2
offset $3;

-- name: GetPosts :many
select
  id,
  owner_id,
  description,
  created_at
from
  posts
where
  id = any ($1::int[])
order by
  created_at desc;
