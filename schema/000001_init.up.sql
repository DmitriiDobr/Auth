CREATE TABLE user_auth(
      id serial not null unique,
      username varchar(255)  not null unique,
      password varchar(255) not null
)
