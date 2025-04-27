CREATE TABLE rate_limits (
    client_id TEXT PRIMARY KEY,
    capacity INTEGER NOT NULL ,
    rate_per_sec INTEGER NOT NULL
);



