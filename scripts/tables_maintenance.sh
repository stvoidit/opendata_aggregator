select_tables="
SELECT
    CONCAT(schemaname,'.',quote_ident(tablename))
FROM
    pg_catalog.pg_tables
WHERE
    schemaname NOT IN ('pg_catalog', 'information_schema')
ORDER BY
    schemaname
    , tablename
"
tables=$(psql -U opendata_app -d opendata_aggregator -p 7777 -qtc "$select_tables")
echo $tables | xargs -n 1 | xargs -I {} /bin/bash -c 'psql -U opendata_app -d opendata_aggregator -p 7777 -c "VACUUM (VERBOSE, ANALYZE, SKIP_LOCKED) {}"'
