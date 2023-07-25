CREATE TABLE post (
  id serial primary key not null,
  title varchar(100) not null,
  description varchar(500) not null,
  latitude float,
  longitude float,
  likes integer not null default 0,
  dislikes integer not null default 0
);
