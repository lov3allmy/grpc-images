CREATE EXTENSION "uuid-ossp";
CREATE TABLE images
(
    id   uuid DEFAULT uuid_generate_v4(),
    name varchar(255) not null,
    raw bytea not null,
    created_at timestamp not null,
    modified_at timestamp not null
);