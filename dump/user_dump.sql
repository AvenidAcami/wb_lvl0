--
-- PostgreSQL database cluster dump
--

SET default_transaction_read_only = off;

SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

--
-- Roles
--

CREATE ROLE wb_lvl0;
ALTER ROLE wb_lvl0 WITH NOSUPERUSER INHERIT NOCREATEROLE NOCREATEDB LOGIN NOREPLICATION NOBYPASSRLS PASSWORD 'SCRAM-SHA-256$4096:DmDC6PAP1kgc2DBSNrHkDA==$HwYaV2YB8u56u/RO44xIY9cYa7/Tdy0mrlncQtWbDZU=:5yBdfHggIWHqeogKy1/bxdsBM1vadT+uRI/cXQEU/cc=';

--
-- User Configurations
--


GRANT CONNECT ON DATABASE wb_lvl0 TO wb_lvl0;
GRANT USAGE ON SCHEMA public TO wb_lvl0;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO wb_lvl0;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO wb_lvl0;






--
-- PostgreSQL database cluster dump complete
--

