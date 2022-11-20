SET statement_timeout = 60000; -- 60 seconds
SET lock_timeout = 30000; -- 30 seconds

--gopg:split
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--gopg:split
CREATE OR REPLACE FUNCTION uuid() RETURNS uuid AS $$
	BEGIN
		RETURN uuid_generate_v4();
	END;
$$ LANGUAGE plpgsql;

--gopg:split
CREATE OR REPLACE FUNCTION update_datetime()	
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_update = now();
    RETURN NEW;	
END;
$$ language 'plpgsql';

--gopg:split
CREATE TABLE IF NOT EXISTS public.currencies (
    uuid uuid NOT NULL DEFAULT uuid(),
    symbol text NOT NULL,
    "description" text NOT NULL,
    source text NOT NULL,
    deleted boolean NOT NULL DEFAULT false,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT currencies_pkey PRIMARY KEY (symbol)
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'calculation_types') THEN
        CREATE TYPE calculation_types AS ENUM ('mult','div');
    END IF;
END$$;


--gopg:split
CREATE TABLE IF NOT EXISTS public.currency_rates (
    uuid uuid NOT NULL DEFAULT uuid(),
    symbol_from text NOT NULL,
    symbol_to text NOT NULL,
    rate numeric NOT NULL,
    calculation_type calculation_types NOT NULL DEFAULT 'mult',
    source text NOT NULL,
    deleted boolean NOT NULL DEFAULT false,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_currency_rates_from
        FOREIGN KEY (symbol_from)
            REFERENCES currencies(symbol),
    CONSTRAINT fk_currency_rates_to
        FOREIGN KEY (symbol_to)
            REFERENCES currencies(symbol)
);

CREATE TRIGGER update_currency_rates
BEFORE UPDATE ON currency_rates
FOR EACH ROW EXECUTE PROCEDURE update_datetime();

CREATE INDEX IF NOT EXISTS currency_rates_uuid_idx ON public.currency_rates USING btree (symbol_from,symbol_to);
