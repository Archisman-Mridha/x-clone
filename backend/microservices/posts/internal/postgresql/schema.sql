create table posts (
	id serial primary key,
	owner_id int not null,
	description varchar(255) not null,
	created_at timestamp not null default now()
);

create index owner_id_idx_posts on posts (owner_id, created_at desc);
