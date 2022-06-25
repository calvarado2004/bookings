CREATE OR REPLACE PROCEDURE database_bookings_drop()

LANGUAGE plpgsql AS

$$ BEGIN
   IF EXISTS (SELECT FROM pg_extension WHERE extname = 'dblink') THEN
      RAISE NOTICE 'Extension dblink already exists';
   ELSE
     CREATE EXTENSION dblink;
   END IF;

   IF EXISTS (SELECT FROM pg_database WHERE datname = 'bookings') THEN
      PERFORM dblink_exec('host=postgres-headless-svc user=postgres dbname=postgres password=WbrTpN3g7q', 'DROP DATABASE bookings');
      DROP EXTENSION dblink;
   END IF;
END $$;

CALL database_bookings_drop();

