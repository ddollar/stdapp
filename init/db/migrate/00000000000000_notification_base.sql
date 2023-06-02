CREATE
OR REPLACE FUNCTION table_changed () RETURNS TRIGGER AS $$
BEGIN
        PERFORM pg_notify(TG_ARGV[0], CASE TG_OP
        WHEN 'INSERT' THEN NEW.ctid
        WHEN 'UPDATE' THEN NEW.ctid
        ELSE OLD.ctid
        END::VARCHAR);
        RETURN NULL;
END;
$$ LANGUAGE plpgsql;