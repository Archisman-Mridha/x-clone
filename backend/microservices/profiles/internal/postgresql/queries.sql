-- name: CreateProfile :exec
insert
into
  profiles(
		id,
		name,
		username
	)
values
	(
		$1,
		$2,
		$3
	);

-- name: GetProfilePreviews :many
select
  id,
  name,
  username
from
  profiles
where
  id = any($1::int[]);
