CREATE DATABASE go_dfw_test;

CREATE TABLE IF NOT EXISTS cars (
    id uuid PRIMARY KEY,
    make character(128),
    model character(128),
    color character(128),
    year integer
);
