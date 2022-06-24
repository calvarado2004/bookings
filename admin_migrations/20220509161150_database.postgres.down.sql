DO
$do$
BEGIN
   IF EXISTS (SELECT FROM pg_extension WHERE extname = 'dblink') THEN
      RAISE NOTICE 'Extension dblink already exists';
   ELSE
     CREATE EXTENSION dblink;
   END IF;

   IF EXISTS (SELECT FROM pg_database WHERE datname = 'bookings') THEN
      PERFORM dblink_exec('dbname=' || current_database()
                        , 'DROP DATABASE bookings');
      DROP EXTENSION dblink;
   END IF;
END
$do$;
