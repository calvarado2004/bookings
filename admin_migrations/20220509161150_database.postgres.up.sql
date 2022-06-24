CREATE TEMPORARY TABLE temp_table_name(
   column_list
);

DO
$do$
BEGIN
   IF EXISTS (SELECT FROM pg_extension WHERE extname = 'dblink') THEN
      RAISE NOTICE 'Extension dblink already exists';
   ELSE
     CREATE EXTENSION dblink;
   END IF;

   IF EXISTS (SELECT FROM pg_database WHERE datname = 'bookings') THEN
      RAISE NOTICE 'Database already exists';
   ELSE
      PERFORM dblink_exec('dbname=' || current_database()
                        , 'CREATE DATABASE bookings');
   END IF;
END
$do$;
