CREATE USER quetest;
CREATE DATABASE quetest OWNER quetest;


CREATE TABLE que_jobs (
  priority    smallint    NOT NULL DEFAULT 100,
  run_at      timestamptz NOT NULL DEFAULT now(),
  job_id      bigserial   NOT NULL,
  job_class   text        NOT NULL,
  args        json        NOT NULL DEFAULT '[]'::json,
  error_count integer     NOT NULL DEFAULT 0,
  last_error  text,
  queue       text        NOT NULL DEFAULT '',

  CONSTRAINT que_jobs_pkey PRIMARY KEY (queue, priority, run_at, job_id)
);

COMMENT ON TABLE que_jobs IS '3';

CREATE TABLE item (
  id          serial    PRIMARY KEY
  ,name       text      NOT NULL
  ,created_at TIMESTAMP NOT NULL default CURRENT_TIMESTAMP
  ,updated_at TIMESTAMP NOT NULL default CURRENT_TIMESTAMP
);

CREATE TABLE item_attribute (
  id         serial     PRIMARY KEY
  ,item_id    integer    NOT NULL
  ,price      integer    NOT NULL
  ,category   text       NOT NULL
  ,created_at TIMESTAMP NOT NULL default CURRENT_TIMESTAMP
  ,updated_at TIMESTAMP NOT NULL default CURRENT_TIMESTAMP
);

BEGIN;

INSERT INTO item (id, name) VALUES(1, 'item1');
INSERT INTO item_attribute (id, item_id, price, category) VALUES(1, 1, 650, 'food');

INSERT INTO item (id, name) VALUES(2, 'item2');
INSERT INTO item_attribute (id, item_id, price, category) VALUES(2, 2, 300, 'snack');

INSERT INTO item (id, name) VALUES(3, 'item3');
INSERT INTO item_attribute (id, item_id, price, category) VALUES(3, 3, 100, 'drink');

COMMIT;
