-- #############
-- enqueue phase
BEGIN;
  INSERT INTO que_jobs (
    priority
    , job_class
    , job_id
  ) values (
    1
    ,'helloJob'
    , 463
  );
COMMIT;


-- #############
-- worker phase 1
-- lock job
SELECT pg_try_advisory_lock(463);

  -- worker phase 2
  -- execute work func
  BEGIN;
    -- success
    INSERT INTO test (
      name
      ,count
    ) values (
      'test'+now
      , 1
    );
  COMMIT;

-- worker phase 3
-- delete job
BEGIN;
  DELETE FROM que_jobs WHERE job_id = 463;
COMMIT;

-- worker phase 4
-- release lock
SELECT pg_advisory_unlock(463);
