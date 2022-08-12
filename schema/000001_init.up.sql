DROP EXTENSION "uuid-ossp";
CREATE EXTENSION "uuid-ossp";
CREATE TABLE images
(
    id   uuid DEFAULT uuid_generate_v4(),
    name varchar(255) not null,
    raw bytea not null,
    created_at int not null,
    modified_at int not null
);