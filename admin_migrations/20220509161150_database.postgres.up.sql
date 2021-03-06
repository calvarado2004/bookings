CREATE OR REPLACE PROCEDURE database_bookings()

LANGUAGE plpgsql AS

$$ BEGIN
   IF EXISTS (SELECT FROM pg_extension WHERE extname = 'dblink') THEN
      RAISE NOTICE 'Extension dblink already exists';
   ELSE
     CREATE EXTENSION dblink;
   END IF;

   IF EXISTS (SELECT FROM pg_database WHERE datname = 'bookings') THEN
      RAISE NOTICE 'Database already exists';
   ELSE
      PERFORM dblink_exec('user={{ DB_USER }} password={{ DB_PASSWORD }}', 'CREATE DATABASE bookings');
   END IF;
END $$;

CALL database_bookings();