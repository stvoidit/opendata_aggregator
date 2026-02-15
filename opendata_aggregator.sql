--
-- PostgreSQL database dump
--

-- Dumped from database version 15.3
-- Dumped by pg_dump version 15.3 (Ubuntu 15.3-1.pgdg22.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: opendata_aggregator; Type: DATABASE; Schema: -; Owner: -
--

CREATE DATABASE opendata_aggregator WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'ru_RU.UTF-8';


\connect opendata_aggregator

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: accounting_statements; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA accounting_statements;


--
-- Name: egr; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA egr;


--
-- Name: service_management; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA service_management;


--
-- Name: btree_gin; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS btree_gin WITH SCHEMA egr;


--
-- Name: EXTENSION btree_gin; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION btree_gin IS 'support for indexing common datatypes in GIN';


--
-- Name: pg_stat_statements; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_stat_statements WITH SCHEMA public;


--
-- Name: EXTENSION pg_stat_statements; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pg_stat_statements IS 'track planning and execution statistics of all SQL statements executed';


--
-- Name: pg_trgm; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA egr;


--
-- Name: EXTENSION pg_trgm; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pg_trgm IS 'text similarity measurement and index searching based on trigrams';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA service_management;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: extract_expenses(jsonb); Type: FUNCTION; Schema: accounting_statements; Owner: -
--

CREATE FUNCTION accounting_statements.extract_expenses(balance_data jsonb) RETURNS bigint
    LANGUAGE sql PARALLEL SAFE
    AS $$
        SELECT 
        COALESCE((balance_data -> 'ФинРез' -> 'СебестПрод' ->> 'СумОтч')::bigint, 0) 
        + COALESCE((balance_data -> 'ФинРез' -> 'КомРасход' ->> 'СумОтч')::bigint, 0) 
        + COALESCE((balance_data -> 'ФинРез' -> 'УпрРасход' ->> 'СумОтч')::bigint, 0) 
        + COALESCE((balance_data -> 'ФинРез' -> 'ПроцУпл' ->> 'СумОтч')::bigint, 0) 
        + COALESCE((balance_data -> 'ФинРез' -> 'РасхОбДеят' ->> 'СумОтч')::bigint, 0) 
        + COALESCE((balance_data -> 'ФинРез' -> 'ПрочРасход' ->> 'СумОтч')::bigint, 0) 
        + COALESCE((balance_data -> 'ФинРез' -> 'НалПрибДох' ->> 'СумОтч')::bigint, 0);
$$;


--
-- Name: FUNCTION extract_expenses(balance_data jsonb); Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON FUNCTION accounting_statements.extract_expenses(balance_data jsonb) IS 'функция извлечения суммы расходов';


--
-- Name: extract_income(jsonb); Type: FUNCTION; Schema: accounting_statements; Owner: -
--

CREATE FUNCTION accounting_statements.extract_income(balance_data jsonb) RETURNS bigint
    LANGUAGE sql PARALLEL SAFE
    AS $$
        SELECT
        COALESCE((balance_data -> 'ФинРез' -> 'Выруч' ->> 'СумОтч')::bigint, 0) 
        + COALESCE((balance_data -> 'ФинРез' -> 'ДоходОтУчаст' ->> 'СумОтч')::bigint, 0) 
        + COALESCE((balance_data -> 'ФинРез' -> 'ПроцПолуч' ->> 'СумОтч')::bigint, 0) 
        + COALESCE((balance_data -> 'ФинРез' -> 'ПрочДоход' ->> 'СумОтч')::bigint, 0);
$$;


--
-- Name: extract_income_tax(jsonb); Type: FUNCTION; Schema: accounting_statements; Owner: -
--

CREATE FUNCTION accounting_statements.extract_income_tax(balance_data jsonb) RETURNS bigint
    LANGUAGE sql PARALLEL SAFE
    AS $$
        SELECT
        CASE
        WHEN balance_data ->> 'ВерсФорм' = '5.03' THEN COALESCE((balance_data -> 'ФинРез' -> 'НалПрибДох' ->> 'СумОтч')::bigint, 0)
        ELSE COALESCE((balance_data -> 'ФинРез' -> 'НалПриб' ->> 'СумОтч')::bigint, 0)
    END;
$$;


--
-- Name: extract_address_text(jsonb); Type: FUNCTION; Schema: egr; Owner: -
--

CREATE FUNCTION egr.extract_address_text(address jsonb) RETURNS text
    LANGUAGE sql IMMUTABLE PARALLEL SAFE
    AS $$
        SELECT 
        CASE WHEN address ? 'СвАдрЮЛФИАС' THEN 
        upper(concat_ws(' ',
        address -> 'СвАдрЮЛФИАС' -> 'НаселенПункт' ->> 'Вид',
        address -> 'СвАдрЮЛФИАС' -> 'НаселенПункт' ->> 'Наим',
        address -> 'СвАдрЮЛФИАС' -> 'МуниципРайон' ->> 'Наим',
        address -> 'СвАдрЮЛФИАС'-> 'ЭлУлДорСети' ->> 'Тип',
        address -> 'СвАдрЮЛФИАС'-> 'ЭлУлДорСети' ->> 'Наим',
        address -> 'СвАдрЮЛФИАС'-> 'Здание' ->> 'Тип',
        address -> 'СвАдрЮЛФИАС'-> 'Здание' ->> 'Номер',
        address -> 'СвАдрЮЛФИАС'-> 'ПомещЗдания' ->> 'Тип',
        address -> 'СвАдрЮЛФИАС'-> 'ПомещЗдания' ->> 'Номер'
        ))
    WHEN address ? 'АдресРФ' THEN 
        upper(concat_ws(' ',  
        address -> 'АдресРФ'->>'Индекс',
        address -> 'АдресРФ'->'Регион'->>'НаимРегион',
        address -> 'АдресРФ'->'Регион'->>'ТипРегион',
        address -> 'АдресРФ'->'Город'->>'ТипГород',
        address -> 'АдресРФ'->'Город'->>'НаимГород',  
        address -> 'АдресРФ'->'НаселПункт'->>'ТипНаселПункт',
        address -> 'АдресРФ'->'НаселПункт'->>'НаимНаселПункт',
        address -> 'АдресРФ'->'Район'->>'ТипРайон',
        address -> 'АдресРФ'->'Район'->>'НаимРайон',
        address -> 'АдресРФ'->'Улица'->>'ТипУлица',
        address -> 'АдресРФ'->'Улица'->>'НаимУлица',
        address -> 'АдресРФ'->>'Дом',
        address -> 'АдресРФ'->>'Кварт'
        ))
    ELSE NULL END
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: egr_search; Type: TABLE; Schema: egr; Owner: -
--

CREATE TABLE egr.egr_search (
    date_discharge date NOT NULL,
    ogrn text NOT NULL,
    inn text NOT NULL,
    is_legal boolean NOT NULL,
    avg_number_employees integer,
    tax_regime jsonb,
    name text,
    name_full text,
    status text,
    address jsonb,
    date_registration date,
    date_liquidation date,
    chief jsonb,
    kpp text,
    tax_authority jsonb,
    idx_name_full tsvector,
    okved jsonb,
    okpo text
);


--
-- Name: extract_egrip_search(date, text, text, jsonb); Type: FUNCTION; Schema: egr; Owner: -
--

CREATE FUNCTION egr.extract_egrip_search(date_discharge date, inn text, ogrn text, value jsonb) RETURNS SETOF egr.egr_search
    LANGUAGE sql IMMUTABLE PARALLEL SAFE
    AS $$
SELECT
    date_discharge
    , ogrn
    , inn
    , FALSE AS is_legal
    , NULL::int AS avg_number_employees
    , NULL::jsonb AS tax_regime
    , upper(concat_ws(' '::TEXT, ((value -> 'СвФЛ'::TEXT) -> 'ФИОРус'::TEXT) ->> 'Фамилия'::TEXT, ((value -> 'СвФЛ'::TEXT) -> 'ФИОРус'::TEXT) ->> 'Имя'::TEXT, ((value -> 'СвФЛ'::TEXT) -> 'ФИОРус'::TEXT) ->> 'Отчество'::TEXT)) AS name
    , upper(concat_ws(' '::TEXT, ((value -> 'СвФЛ'::TEXT) -> 'ФИОРус'::TEXT) ->> 'Фамилия'::TEXT, ((value -> 'СвФЛ'::TEXT) -> 'ФИОРус'::TEXT) ->> 'Имя'::TEXT, ((value -> 'СвФЛ'::TEXT) -> 'ФИОРус'::TEXT) ->> 'Отчество'::TEXT)) AS name_full
    , (
            (
                value -> 'СвСтатус'::TEXT
        ) -> 'СвСтатусИП'::TEXT
    ) ->> 'НаимСтатус'::TEXT AS status
    , value -> 'СвАдрМЖ'::TEXT AS address
    , COALESCE(
            (
                value -> 'СвРегИП'::TEXT
        ) ->> 'ДатаРег'::TEXT
    , (
                value -> 'СвРегИП'::TEXT
        ) ->> 'ДатаОГРНИП'::TEXT
    )::date AS date_registration
    , (
            (
                (
                    value -> 'СвПрекращ'::TEXT
            ) -> 'СвСтатус'::TEXT
        ) ->> 'ДатаПрекращ'::TEXT
    )::date AS date_liquidation
    , value -> 'СвФЛ'::TEXT AS chief
    , NULL AS kpp
    , (
        value -> 'СвУчетНО'::TEXT
    ) -> 'СвНО'::TEXT AS tax_authority
    , to_tsvector(
            upper(concat_ws(' '::TEXT, ((value -> 'СвФЛ'::TEXT) -> 'ФИОРус'::TEXT) ->> 'Фамилия'::TEXT, ((value -> 'СвФЛ'::TEXT) -> 'ФИОРус'::TEXT) ->> 'Имя'::TEXT, ((value -> 'СвФЛ'::TEXT) -> 'ФИОРус'::TEXT) ->> 'Отчество'::TEXT))
    ) AS idx_name_full
    , CASE
            WHEN (
            value -> 'СвОКВЭД'::TEXT
        ) ? 'СвОКВЭДОсн'::TEXT THEN jsonb_build_object(
            'code'
            , (
                (
                    value -> 'СвОКВЭД'::TEXT
                ) -> 'СвОКВЭДОсн'::TEXT
            ) ->> 'КодОКВЭД'::TEXT
            , 'title'
            , (
                (
                    value -> 'СвОКВЭД'::TEXT
                ) -> 'СвОКВЭДОсн'::TEXT
            ) ->> 'НаимОКВЭД'::TEXT
        )
        ELSE NULL::jsonb
    END AS okved
    , NULL AS okpo
$$;


--
-- Name: extract_egrul_search(date, text, text, jsonb); Type: FUNCTION; Schema: egr; Owner: -
--

CREATE FUNCTION egr.extract_egrul_search(date_discharge date, inn text, ogrn text, value jsonb) RETURNS SETOF egr.egr_search
    LANGUAGE sql IMMUTABLE PARALLEL SAFE
    AS $$
SELECT 
dd.date_discharge
, dd.ogrn
, dd.inn
, dd.is_legal
, sr."колраб" AS avg_number_employees
, jsonb_build_object(
        'есхн'
        , rn."есхн"
        , 'усн'
        , rn."усн"
        , 'енвд'
        , rn."енвд"
        , 'срп'
        , rn."срп"
    ) AS tax_regime
, dd.name
, dd.name_full
, dd.status
, dd.address
, dd.date_registration
, dd.date_liquidation
, dd.chief
, dd.kpp
, dd.tax_authority
, dd.idx_name_full
, dd.okved
, o.okpo

FROM (SELECT
    date_discharge
    , ogrn
    , inn
    , TRUE AS is_legal
    , upper(((value -> 'СвНаимЮЛ'::TEXT) -> 'СвНаимЮЛСокр'::TEXT) ->> 'НаимСокр'::TEXT) AS name
    , upper((value -> 'СвНаимЮЛ'::TEXT) ->> 'НаимЮЛПолн'::TEXT) AS name_full
    ,(
        (
            value -> 'СвСтатус'::TEXT
        ) -> 'СвСтатус'::TEXT
    ) ->> 'НаимСтатусЮЛ'::TEXT AS status
    , value -> 'СвАдресЮЛ'::TEXT AS address
    ,(
        (
            value -> 'СвОбрЮЛ'::TEXT
        ) ->> 'ДатаРег'::TEXT
    )::date AS date_registration
    ,(
        (
            value -> 'СвПрекрЮЛ'::TEXT
        ) ->> 'ДатаПрекрЮЛ'::TEXT
    )::date AS date_liquidation
    ,(
        value -> 'СведДолжнФЛ'::TEXT
    ) -> 0 AS chief
    , value ->> 'КПП'::TEXT AS kpp
    ,(
        value -> 'СвУчетНО'::TEXT
    ) -> 'СвНО'::TEXT AS tax_authority
    , to_tsvector(
        upper((value -> 'СвНаимЮЛ'::TEXT) ->> 'НаимЮЛПолн'::TEXT)
    ) AS idx_name_full
    , CASE
            WHEN (
            value -> 'СвОКВЭД'::TEXT
        ) ? 'СвОКВЭДОсн'::TEXT THEN jsonb_build_object(
            'code'
            , (
                (
                    value -> 'СвОКВЭД'::TEXT
                ) -> 'СвОКВЭДОсн'::TEXT
            ) ->> 'КодОКВЭД'::TEXT
            , 'title'
            , (
                (
                    value -> 'СвОКВЭД'::TEXT
                ) -> 'СвОКВЭДОсн'::TEXT
            ) ->> 'НаимОКВЭД'::TEXT
        )
        ELSE NULL::jsonb
    END AS okved
    ) AS dd

    
LEFT JOIN "сведенияосреднчислработников" sr ON
    sr."иннюл" = dd.inn
LEFT JOIN "режимналогоплательщика" rn ON
    rn."иннюл" = dd.inn
LEFT JOIN accounting_statements.okpo o ON
    o.inn = dd.inn
$$;


--
-- Name: extract_termination_status(jsonb); Type: FUNCTION; Schema: egr; Owner: -
--

CREATE FUNCTION egr.extract_termination_status(egr_data jsonb) RETURNS text
    LANGUAGE sql PARALLEL SAFE
    AS $$
    SELECT egr_data -> 'СвПрекрЮЛ' -> 'СпПрекрЮЛ' ->> 'НаимСпПрекрЮЛ';
$$;


--
-- Name: trigger_proc_egrip_search(); Type: FUNCTION; Schema: egr; Owner: -
--

CREATE FUNCTION egr.trigger_proc_egrip_search() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
INSERT INTO egr.egr_search(
    date_discharge,
    ogrn,
    inn,
    is_legal,
    avg_number_employees,
    tax_regime,
    name,
    name_full,
    status,
    address,
    date_registration,
    date_liquidation,
    chief,
    kpp,
    tax_authority,
    idx_name_full,
    okved,
    okpo
)
SELECT (egr.extract_egrip_search(NEW.датавып, NEW.инн, NEW.огрн, NEW.egrip_data)).*
ON CONFLICT (inn, ogrn) DO UPDATE
SET
    date_discharge = EXCLUDED.date_discharge
    , ogrn = EXCLUDED.ogrn
    , inn = EXCLUDED.inn
    , is_legal = EXCLUDED.is_legal
    , avg_number_employees = EXCLUDED.avg_number_employees
    , tax_regime = EXCLUDED.tax_regime
    , "name" = EXCLUDED."name"
    , name_full = EXCLUDED.name_full
    , status = EXCLUDED.status
    , address = EXCLUDED.address
    , date_registration = EXCLUDED.date_registration
    , date_liquidation = EXCLUDED.date_liquidation
    , chief = EXCLUDED.chief
    , kpp = EXCLUDED.kpp
    , tax_authority = EXCLUDED.tax_authority
    , idx_name_full = EXCLUDED.idx_name_full
    , okved = EXCLUDED.okved
    , okpo = EXCLUDED.okpo
;
RETURN NULL;
END;
$$;


--
-- Name: trigger_proc_egrul_search(); Type: FUNCTION; Schema: egr; Owner: -
--

CREATE FUNCTION egr.trigger_proc_egrul_search() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
INSERT INTO egr.egr_search(
    date_discharge,
    ogrn,
    inn,
    is_legal,
    avg_number_employees,
    tax_regime,
    name,
    name_full,
    status,
    address,
    date_registration,
    date_liquidation,
    chief,
    kpp,
    tax_authority,
    idx_name_full,
    okved,
    okpo
)
SELECT (egr.extract_egrul_search(NEW.датавып, NEW.инн, NEW.огрн, NEW.egrul_data)).*
ON CONFLICT (inn, ogrn) DO UPDATE
SET
    date_discharge = EXCLUDED.date_discharge
    , ogrn = EXCLUDED.ogrn
    , inn = EXCLUDED.inn
    , is_legal = EXCLUDED.is_legal
    , avg_number_employees = EXCLUDED.avg_number_employees
    , tax_regime = EXCLUDED.tax_regime
    , "name" = EXCLUDED."name"
    , name_full = EXCLUDED.name_full
    , status = EXCLUDED.status
    , address = EXCLUDED.address
    , date_registration = EXCLUDED.date_registration
    , date_liquidation = EXCLUDED.date_liquidation
    , chief = EXCLUDED.chief
    , kpp = EXCLUDED.kpp
    , tax_authority = EXCLUDED.tax_authority
    , idx_name_full = EXCLUDED.idx_name_full
    , okved = EXCLUDED.okved
    , okpo = EXCLUDED.okpo
;
RETURN NULL;
END;
$$;


--
-- Name: upsert_egrul(date, text, text, jsonb); Type: PROCEDURE; Schema: egr; Owner: -
--

CREATE PROCEDURE egr.upsert_egrul(IN n_date date, IN n_inn text, IN n_ogrn text, IN n_value jsonb)
    LANGUAGE plpgsql
    AS $$
BEGIN 
    IF EXISTS(SELECT 1 FROM egr.егрюл e WHERE e.инн = n_inn AND e.огрн = n_ogrn) THEN 
        UPDATE egr.егрюл 
        SET датавып = n_date, инн = n_inn, огрн = n_ogrn, egrul_data = n_value
        WHERE инн = n_inn AND огрн = n_ogrn AND датавып < n_date;
    ELSE
        INSERT INTO egr.егрюл (датавып,инн,огрн,egrul_data) VALUES (n_date, n_inn, n_ogrn, n_value);
    END IF;
END
$$;


--
-- Name: баланс; Type: TABLE; Schema: accounting_statements; Owner: -
--

CREATE TABLE accounting_statements."баланс" (
    "иннюл" text NOT NULL,
    balance_data jsonb DEFAULT '{}'::jsonb NOT NULL,
    "год" integer NOT NULL
);


--
-- Name: balance_mat_view; Type: MATERIALIZED VIEW; Schema: accounting_statements; Owner: -
--

CREATE MATERIALIZED VIEW accounting_statements.balance_mat_view AS
 SELECT (("баланс".balance_data ->> 'ОтчетГод'::text))::integer AS year,
    "баланс"."год" AS doc_date,
    "баланс"."иннюл" AS inn,
    (("баланс".balance_data -> 'СвНП'::text) ->> 'ОКПО'::text) AS okpo,
    accounting_statements.extract_income("баланс".balance_data) AS income,
    accounting_statements.extract_expenses("баланс".balance_data) AS expenses,
    accounting_statements.extract_income_tax("баланс".balance_data) AS income_tax,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN COALESCE((((("баланс".balance_data -> 'ФинРез'::text) -> 'НалПрибДох'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
            ELSE NULL::bigint
        END AS tax_usn,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'НеМатФинАкт'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
            ELSE COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'ВнеОбА'::text) -> 'НематАкт'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS intangible_assets,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'МатВнеАкт'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
            ELSE COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'ОснСр'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS basic_assets,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'НеМатФинАкт'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
            ELSE COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'ПрочВнеОбА'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS other_non_current_assets,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'НеМатФинАкт'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
            ELSE COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'ВнеОбА'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS non_current_assets,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'Запасы'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
            ELSE COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'ВнеОбА'::text) ->> 'Запасы'::text))::bigint, (0)::bigint)
        END AS stocks,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN COALESCE(((((("баланс".balance_data -> 'ОтчетИзмКап'::text) -> 'Чист'::text) -> 'Актив'::text) ->> 'На31ДекОтч'::text))::bigint, (0)::bigint)
            ELSE COALESCE((((("баланс".balance_data -> 'ОтчетИзмКап'::text) -> 'ЧистАктив'::text) ->> 'На31ДекОтч'::text))::bigint, (0)::bigint)
        END AS net_aassets,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN NULL::bigint
            ELSE COALESCE((((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'ОбА'::text) -> 'ДебЗад'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS accounts_receivable,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'ДенежнСр'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
            ELSE COALESCE(((((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'ВнеОбА'::text) -> 'ОбА'::text) -> 'ДенежнСр'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS cash_and_equivalents,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN ((COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'Запасы'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint) + COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'ДенежнСр'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)) + COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'ФинВлож'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint))
            ELSE COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) -> 'ОбА'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS current_assets,
    COALESCE((((("баланс".balance_data -> 'Баланс'::text) -> 'Актив'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint) AS total_assets,
    COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Пассив'::text) -> 'КапРез'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint) AS capital_and_reserves,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Пассив'::text) -> 'ДлгЗаемСредств'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
            ELSE COALESCE((((((("баланс".balance_data -> 'Баланс'::text) -> 'Пассив'::text) -> 'ДолгосрОбяз'::text) -> 'ЗаемСредств'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS borrowed_funds_long_term,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN COALESCE((((((("баланс".balance_data -> 'Баланс'::text) -> 'Пассив'::text) -> 'ДлгЗаемСредств'::text) -> 'КртЗаемСредств'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
            ELSE COALESCE((((((("баланс".balance_data -> 'Баланс'::text) -> 'Пассив'::text) -> 'КраткосрОбяз'::text) -> 'ЗаемСредств'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS borrowed_funds_short_term,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Пассив'::text) -> 'КредитЗадолж'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
            ELSE COALESCE((((((("баланс".balance_data -> 'Баланс'::text) -> 'Пассив'::text) -> 'КраткосрОбяз'::text) -> 'КредитЗадолж'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS accounts_payable,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN COALESCE(((((("баланс".balance_data -> 'Баланс'::text) -> 'Пассив'::text) -> 'ДрКраткосрОбяз'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
            ELSE COALESCE((((((("баланс".balance_data -> 'Баланс'::text) -> 'Пассив'::text) -> 'КраткосрОбяз'::text) -> 'ПрочОбяз'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS other_short_term_liabilities,
    COALESCE((((("баланс".balance_data -> 'Баланс'::text) -> 'Пассив'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint) AS total_liabilities,
    COALESCE((((("баланс".balance_data -> 'ФинРез'::text) -> 'Выруч'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint) AS revenue,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN NULL::bigint
            ELSE COALESCE((((("баланс".balance_data -> 'ФинРез'::text) -> 'СебестПрод'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS cost_of_sales,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN NULL::bigint
            ELSE COALESCE((((("баланс".balance_data -> 'ФинРез'::text) -> 'ВаловаяПрибыль'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS gross_profit,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN NULL::bigint
            ELSE COALESCE((((("баланс".balance_data -> 'ФинРез'::text) -> 'КомРасход'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS commercial_expenses,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN NULL::bigint
            ELSE COALESCE((((("баланс".balance_data -> 'ФинРез'::text) -> 'УпрРасход'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS management_expenses,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN NULL::bigint
            ELSE COALESCE((((("баланс".balance_data -> 'ФинРез'::text) -> 'ПрибПрод'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS profit_from_sale,
    COALESCE((((("баланс".balance_data -> 'ФинРез'::text) -> 'ПроцУпл'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint) AS interest_payable,
    COALESCE((((("баланс".balance_data -> 'ФинРез'::text) -> 'ПрочДоход'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint) AS other_income,
    COALESCE((((("баланс".balance_data -> 'ФинРез'::text) -> 'ПрочРасход'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint) AS other_expenses,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN NULL::bigint
            ELSE COALESCE((((("баланс".balance_data -> 'ФинРез'::text) -> 'ПрибУбДоНал'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS profit_before_taxation,
        CASE
            WHEN (("баланс".balance_data ->> 'ВерсФорм'::text) = '5.03'::text) THEN NULL::bigint
            ELSE COALESCE((((("баланс".balance_data -> 'ФинРез'::text) -> 'ТекНалПриб'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint)
        END AS "сurrent_income_tax",
    COALESCE((((("баланс".balance_data -> 'ФинРез'::text) -> 'ЧистПрибУб'::text) ->> 'СумОтч'::text))::bigint, (0)::bigint) AS net_profit
   FROM accounting_statements."баланс"
  WITH NO DATA;


--
-- Name: COLUMN balance_mat_view.okpo; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.okpo IS 'ОКПО';


--
-- Name: COLUMN balance_mat_view.income; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.income IS 'Доходы';


--
-- Name: COLUMN balance_mat_view.expenses; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.expenses IS 'Расходы';


--
-- Name: COLUMN balance_mat_view.income_tax; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.income_tax IS 'Налог на прибыль';


--
-- Name: COLUMN balance_mat_view.tax_usn; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.tax_usn IS 'Налог по УСН';


--
-- Name: COLUMN balance_mat_view.intangible_assets; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.intangible_assets IS 'Нематериальные активы';


--
-- Name: COLUMN balance_mat_view.basic_assets; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.basic_assets IS 'Основные средства';


--
-- Name: COLUMN balance_mat_view.other_non_current_assets; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.other_non_current_assets IS 'Прочие внеоборотные активы';


--
-- Name: COLUMN balance_mat_view.non_current_assets; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.non_current_assets IS 'Внеоборотные активы';


--
-- Name: COLUMN balance_mat_view.stocks; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.stocks IS 'Запасы';


--
-- Name: COLUMN balance_mat_view.net_aassets; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.net_aassets IS 'Чистые активы';


--
-- Name: COLUMN balance_mat_view.accounts_receivable; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.accounts_receivable IS 'Дебиторская задолженность';


--
-- Name: COLUMN balance_mat_view.cash_and_equivalents; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.cash_and_equivalents IS 'Денежные средства и денежные эквиваленты';


--
-- Name: COLUMN balance_mat_view.current_assets; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.current_assets IS 'Оборотные активы';


--
-- Name: COLUMN balance_mat_view.total_assets; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.total_assets IS 'Активы всего';


--
-- Name: COLUMN balance_mat_view.capital_and_reserves; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.capital_and_reserves IS 'Капитал и резервы';


--
-- Name: COLUMN balance_mat_view.borrowed_funds_long_term; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.borrowed_funds_long_term IS 'Заёмные средства (долгосрочные)';


--
-- Name: COLUMN balance_mat_view.borrowed_funds_short_term; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.borrowed_funds_short_term IS 'Заёмные средства (краткосрочные)';


--
-- Name: COLUMN balance_mat_view.accounts_payable; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.accounts_payable IS 'Кредиторская задолженность';


--
-- Name: COLUMN balance_mat_view.other_short_term_liabilities; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.other_short_term_liabilities IS 'Прочие краткосрочные обязательства';


--
-- Name: COLUMN balance_mat_view.total_liabilities; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.total_liabilities IS 'Пассивы всего';


--
-- Name: COLUMN balance_mat_view.revenue; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.revenue IS 'Выручка';


--
-- Name: COLUMN balance_mat_view.cost_of_sales; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.cost_of_sales IS 'Себестоимость продаж';


--
-- Name: COLUMN balance_mat_view.gross_profit; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.gross_profit IS 'Валовая прибыль';


--
-- Name: COLUMN balance_mat_view.commercial_expenses; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.commercial_expenses IS 'Коммерческие расходы';


--
-- Name: COLUMN balance_mat_view.management_expenses; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.management_expenses IS 'Управленческие расходы';


--
-- Name: COLUMN balance_mat_view.profit_from_sale; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.profit_from_sale IS 'Прибыль (убыток) от продажи';


--
-- Name: COLUMN balance_mat_view.interest_payable; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.interest_payable IS 'Проценты к уплате';


--
-- Name: COLUMN balance_mat_view.other_income; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.other_income IS 'Прочие доходы';


--
-- Name: COLUMN balance_mat_view.other_expenses; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.other_expenses IS 'Прочие расходы';


--
-- Name: COLUMN balance_mat_view.profit_before_taxation; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.profit_before_taxation IS 'Прибыль (убыток) до налогообложения';


--
-- Name: COLUMN balance_mat_view."сurrent_income_tax"; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view."сurrent_income_tax" IS 'Текущий налог на прибыль';


--
-- Name: COLUMN balance_mat_view.net_profit; Type: COMMENT; Schema: accounting_statements; Owner: -
--

COMMENT ON COLUMN accounting_statements.balance_mat_view.net_profit IS 'Чистая прибыль (убыток)';


--
-- Name: okpo; Type: MATERIALIZED VIEW; Schema: accounting_statements; Owner: -
--

CREATE MATERIALIZED VIEW accounting_statements.okpo AS
 SELECT DISTINCT ON (b."иннюл") b."иннюл" AS inn,
    ((b.balance_data -> 'СвНП'::text) ->> 'ОКПО'::text) AS okpo
   FROM accounting_statements."баланс" b
  WHERE (((b.balance_data -> 'СвНП'::text) ->> 'ОКПО'::text) <> ''::text)
  WITH NO DATA;


--
-- Name: егрип; Type: TABLE; Schema: egr; Owner: -
--

CREATE TABLE egr."егрип" (
    "датавып" date NOT NULL,
    "инн" text NOT NULL,
    "огрн" text NOT NULL,
    egrip_data jsonb DEFAULT '{}'::jsonb NOT NULL
);


--
-- Name: егрюл; Type: TABLE; Schema: egr; Owner: -
--

CREATE TABLE egr."егрюл" (
    "датавып" date NOT NULL,
    "инн" text NOT NULL,
    "огрн" text NOT NULL,
    egrul_data jsonb DEFAULT '{}'::jsonb NOT NULL
);


--
-- Name: handbook_okved; Type: MATERIALIZED VIEW; Schema: egr; Owner: -
--

CREATE MATERIALIZED VIEW egr.handbook_okved AS
 WITH aggd AS (
         SELECT eu.egrul_data AS value
           FROM egr."егрюл" eu
        UNION ALL
         SELECT ep.egrip_data AS value
           FROM egr."егрип" ep
        ), okveds_j AS (
         SELECT (aggd.value -> 'СвОКВЭД'::text) AS okveds
           FROM aggd
          WHERE (aggd.value ? 'СвОКВЭД'::text)
        ), unwrap_okved AS (
         SELECT (okveds_j.okveds -> 'СвОКВЭДОсн'::text) AS okved_j
           FROM okveds_j
        UNION ALL
         SELECT jsonb_array_elements((okveds_j.okveds -> 'СвОКВЭДДоп'::text)) AS okved_j
           FROM okveds_j
        )
 SELECT DISTINCT ON (o.code) o.code,
    o.title,
    o.vers
   FROM ( SELECT (unwrap_okved.okved_j ->> 'КодОКВЭД'::text) AS code,
            (unwrap_okved.okved_j ->> 'НаимОКВЭД'::text) AS title,
            (unwrap_okved.okved_j ->> 'ПрВерсОКВЭД'::text) AS vers
           FROM unwrap_okved
          WHERE (unwrap_okved.okved_j ? 'КодОКВЭД'::text)) o
  ORDER BY o.code
  WITH NO DATA;


--
-- Name: tax_authority; Type: MATERIALIZED VIEW; Schema: egr; Owner: -
--

CREATE MATERIALIZED VIEW egr.tax_authority AS
 SELECT egr_search.inn,
    egr_search.ogrn,
    (egr_search.tax_authority ->> 'КодНО'::text) AS code,
    (egr_search.tax_authority ->> 'НаимНО'::text) AS name
   FROM egr.egr_search
  WHERE (egr_search.tax_authority IS NOT NULL)
  WITH NO DATA;


--
-- Name: handbook_tax_authority; Type: MATERIALIZED VIEW; Schema: egr; Owner: -
--

CREATE MATERIALIZED VIEW egr.handbook_tax_authority AS
 SELECT DISTINCT tax_authority.code,
    upper(tax_authority.name) AS name
   FROM egr.tax_authority
  ORDER BY tax_authority.code
  WITH NO DATA;


--
-- Name: statuses; Type: TABLE; Schema: egr; Owner: -
--

CREATE TABLE egr.statuses (
    status_name text NOT NULL,
    status_full_name text NOT NULL
);


--
-- Name: handbook_categories_ip; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.handbook_categories_ip (
    category text NOT NULL,
    subcategory text NOT NULL
);


--
-- Name: TABLE handbook_categories_ip; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.handbook_categories_ip IS 'категории ИП';


--
-- Name: COLUMN handbook_categories_ip.category; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.handbook_categories_ip.category IS 'Категория';


--
-- Name: COLUMN handbook_categories_ip.subcategory; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.handbook_categories_ip.subcategory IS 'Подкатегория';


--
-- Name: hotels; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.hotels (
    federal_number text NOT NULL,
    type text,
    full_name text,
    short_name text,
    region text,
    inn text,
    ogrn text,
    address text,
    phone text,
    fax text,
    email text,
    site text,
    owner text DEFAULT ''::text NOT NULL
);


--
-- Name: hotels_classification; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.hotels_classification (
    federal_number text NOT NULL,
    date_issued date,
    date_end date,
    category text,
    license_number text,
    registration_number text
);


--
-- Name: hotels_rooms; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.hotels_rooms (
    federal_number text NOT NULL,
    category text,
    rooms integer,
    seats integer
);


--
-- Name: hotels_view; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.hotels_view AS
 SELECT (((row_to_json(h.*))::jsonb || jsonb_build_object('classification', kg_j.k)) || jsonb_build_object('rooms', ng_j.n)) AS hotel_data
   FROM ((public.hotels h
     LEFT JOIN LATERAL ( SELECT COALESCE(jsonb_agg(row_to_json(hc.*) ORDER BY hc.date_issued), '[]'::jsonb) AS k
           FROM public.hotels_classification hc
          WHERE (hc.federal_number = h.federal_number)) kg_j ON (true))
     LEFT JOIN LATERAL ( SELECT COALESCE(jsonb_agg(row_to_json(hr.*) ORDER BY hr.category), '[]'::jsonb) AS n
           FROM public.hotels_rooms hr
          WHERE (hr.federal_number = h.federal_number)) ng_j ON (true))
  ORDER BY h.federal_number DESC;


--
-- Name: iplegallistcomplete; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.iplegallistcomplete (
    nameofdebtor text,
    addressofdebtororganization text,
    actualaddressofdebtororganization text,
    numberofenforcementproceeding text NOT NULL,
    dateofinstitutionproceeding text,
    totalnumberofenforcementproceedings text,
    executivedocumenttype text,
    dateofexecutivedocument text,
    numberofexecutivedocument text NOT NULL,
    objectofexecutivedocuments text,
    objectofexecution text,
    datecompleteipreason text DEFAULT 0 NOT NULL,
    departmentsofbailiffs text,
    addressofdepartmentsofbailiff text,
    debtortaxpayeridentificationnumber text,
    taxpayeridentificationnumberoforganizationcollector text
);


--
-- Name: TABLE iplegallistcomplete; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.iplegallistcomplete IS 'Исполнительные производства в отношении юридических лиц';


--
-- Name: COLUMN iplegallistcomplete.nameofdebtor; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.nameofdebtor IS 'Наименование юридического лица';


--
-- Name: COLUMN iplegallistcomplete.addressofdebtororganization; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.addressofdebtororganization IS 'Адрес организации - должника';


--
-- Name: COLUMN iplegallistcomplete.actualaddressofdebtororganization; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.actualaddressofdebtororganization IS 'Фактический адрес организации должника';


--
-- Name: COLUMN iplegallistcomplete.numberofenforcementproceeding; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.numberofenforcementproceeding IS 'Номер исполнительного производства';


--
-- Name: COLUMN iplegallistcomplete.dateofinstitutionproceeding; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.dateofinstitutionproceeding IS 'Дата возбуждения исполнительного производства';


--
-- Name: COLUMN iplegallistcomplete.totalnumberofenforcementproceedings; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.totalnumberofenforcementproceedings IS 'Номер сводного производства по взыскателю или должнику';


--
-- Name: COLUMN iplegallistcomplete.executivedocumenttype; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.executivedocumenttype IS 'Тип исполнительного документа';


--
-- Name: COLUMN iplegallistcomplete.dateofexecutivedocument; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.dateofexecutivedocument IS 'Дата исполнительного документа';


--
-- Name: COLUMN iplegallistcomplete.numberofexecutivedocument; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.numberofexecutivedocument IS 'Номер исполнительного документа';


--
-- Name: COLUMN iplegallistcomplete.objectofexecutivedocuments; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.objectofexecutivedocuments IS 'Требования исполнительного документа';


--
-- Name: COLUMN iplegallistcomplete.objectofexecution; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.objectofexecution IS 'Предмет исполнения';


--
-- Name: COLUMN iplegallistcomplete.datecompleteipreason; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.datecompleteipreason IS 'Дата, причина окончания или прекращения ИП (статья, часть, пункт основания)';


--
-- Name: COLUMN iplegallistcomplete.departmentsofbailiffs; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.departmentsofbailiffs IS 'Наименование отдела';


--
-- Name: COLUMN iplegallistcomplete.addressofdepartmentsofbailiff; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.addressofdepartmentsofbailiff IS 'Адрес отдела судебных приставов';


--
-- Name: COLUMN iplegallistcomplete.debtortaxpayeridentificationnumber; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.debtortaxpayeridentificationnumber IS 'Идентификационный номер налогоплательщика должника';


--
-- Name: COLUMN iplegallistcomplete.taxpayeridentificationnumberoforganizationcollector; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.iplegallistcomplete.taxpayeridentificationnumberoforganizationcollector IS 'Идентификационный номер налогоплательщика взыскателя-организации';


--
-- Name: rds; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.rds (
    id_decl text NOT NULL,
    reg_number text NOT NULL,
    decl_status text NOT NULL,
    decl_type text,
    date_beginning date,
    date_finish date,
    declaration_scheme text,
    product_object_type_decl text NOT NULL,
    product_type text NOT NULL,
    product_group text,
    product_name text,
    asproduct_info text,
    product_tech_reg text,
    organ_to_certification_name text,
    organ_to_certification_reg_number text,
    basis_for_decl text,
    old_basis_for_decl text,
    applicant_type text,
    person_applicant_type text,
    applicant_ogrn text,
    applicant_inn text,
    applicant_name text,
    manufacturer_type text,
    manufacturer_ogrn text,
    manufacturer_inn text,
    manufacturer_name text
);


--
-- Name: TABLE rds; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.rds IS 'Сведения из Реестра деклараций о соответствии
https://fsa.gov.ru/opendata/7736638268-rds/';


--
-- Name: COLUMN rds.id_decl; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.id_decl IS 'id';


--
-- Name: COLUMN rds.reg_number; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.reg_number IS 'Рег. Номер';


--
-- Name: COLUMN rds.decl_status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.decl_status IS 'Статус';


--
-- Name: COLUMN rds.decl_type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.decl_type IS 'Тип декларации';


--
-- Name: COLUMN rds.date_beginning; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.date_beginning IS 'Дата начала действия';


--
-- Name: COLUMN rds.date_finish; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.date_finish IS 'Дата окончания действия';


--
-- Name: COLUMN rds.declaration_scheme; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.declaration_scheme IS 'Схема декларирования';


--
-- Name: COLUMN rds.product_object_type_decl; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.product_object_type_decl IS 'Тип объекта декларирования';


--
-- Name: COLUMN rds.product_type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.product_type IS 'Вид продукции';


--
-- Name: COLUMN rds.product_group; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.product_group IS 'Группа продукции';


--
-- Name: COLUMN rds.product_name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.product_name IS 'Общее наименование продукции';


--
-- Name: COLUMN rds.asproduct_info; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.asproduct_info IS 'Информация по продукции';


--
-- Name: COLUMN rds.product_tech_reg; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.product_tech_reg IS 'Технический регламент';


--
-- Name: COLUMN rds.organ_to_certification_name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.organ_to_certification_name IS 'Наименование ОС';


--
-- Name: COLUMN rds.organ_to_certification_reg_number; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.organ_to_certification_reg_number IS 'Номер аттестата аккредитации ОС';


--
-- Name: COLUMN rds.basis_for_decl; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.basis_for_decl IS 'Основание выдачи ДС';


--
-- Name: COLUMN rds.old_basis_for_decl; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.old_basis_for_decl IS 'Основание выдачи ДС (ФГИС 1.0)';


--
-- Name: COLUMN rds.applicant_type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.applicant_type IS 'Тип заявителя';


--
-- Name: COLUMN rds.person_applicant_type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.person_applicant_type IS 'Вид заявителя';


--
-- Name: COLUMN rds.applicant_ogrn; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.applicant_ogrn IS 'ОГРН/ОГРНИП заявителя';


--
-- Name: COLUMN rds.applicant_inn; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.applicant_inn IS 'ИНН заявителя';


--
-- Name: COLUMN rds.applicant_name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.applicant_name IS 'Наименование Заявителя';


--
-- Name: COLUMN rds.manufacturer_type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.manufacturer_type IS 'Вид изготовителя';


--
-- Name: COLUMN rds.manufacturer_ogrn; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.manufacturer_ogrn IS 'ОГРН/ОГРНИП изготовителя';


--
-- Name: COLUMN rds.manufacturer_inn; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.manufacturer_inn IS 'ИНН изготовителя';


--
-- Name: COLUMN rds.manufacturer_name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.rds.manufacturer_name IS 'Полное наименование изготовителя';


--
-- Name: информация1928коап; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."информация1928коап" (
    "инн" text NOT NULL,
    "огрн" text NOT NULL,
    "кпп" text,
    "наименованиеюл" text,
    "типучастника" text,
    "суд" text,
    "номердела" text NOT NULL,
    "датавынесенияпостановления" date,
    "датавступлениявзаконнуюсилу" date
);


--
-- Name: исппроизввотнюрлиц; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."исппроизввотнюрлиц" (
    nameofdebtor text,
    addressofdebtororganization text,
    actualaddressofdebtororganization text,
    numberofenforcementproceeding text NOT NULL,
    dateofinstitutionproceeding text,
    totalnumberofenforcementproceedings text,
    executivedocumenttype text,
    dateofexecutivedocument text,
    numberofexecutivedocument text NOT NULL,
    objectofexecutivedocuments text,
    objectofexecution text,
    amountdue numeric DEFAULT 0 NOT NULL,
    debtremainingbalance numeric DEFAULT 0 NOT NULL,
    departmentsofbailiffs text,
    addressofdepartmentsofbailiff text,
    debtortaxpayeridentificationnumber text,
    taxpayeridentificationnumberoforganizationcollector text,
    repaid boolean DEFAULT false NOT NULL
);


--
-- Name: TABLE "исппроизввотнюрлиц"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public."исппроизввотнюрлиц" IS 'Исполнительные производства в отношении юридических лиц, оконченные в соответствии с пунктами 3 и 4 части 1 статьи 46 и пунктами 6 и 7 части 1 статьи 47 Федерального закона от 2 октября 2007 г. № 229-ФЗ «Об исполнительном производстве»';


--
-- Name: COLUMN "исппроизввотнюрлиц".nameofdebtor; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".nameofdebtor IS 'Наименование юридического лица';


--
-- Name: COLUMN "исппроизввотнюрлиц".addressofdebtororganization; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".addressofdebtororganization IS 'Адрес организации - должника';


--
-- Name: COLUMN "исппроизввотнюрлиц".actualaddressofdebtororganization; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".actualaddressofdebtororganization IS 'Фактический адрес организации должника';


--
-- Name: COLUMN "исппроизввотнюрлиц".numberofenforcementproceeding; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".numberofenforcementproceeding IS 'Номер исполнительного производства';


--
-- Name: COLUMN "исппроизввотнюрлиц".dateofinstitutionproceeding; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".dateofinstitutionproceeding IS 'Дата возбуждения исполнительного производства';


--
-- Name: COLUMN "исппроизввотнюрлиц".totalnumberofenforcementproceedings; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".totalnumberofenforcementproceedings IS 'Номер сводного производства по взыскателю или должнику';


--
-- Name: COLUMN "исппроизввотнюрлиц".executivedocumenttype; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".executivedocumenttype IS 'Тип исполнительного документа';


--
-- Name: COLUMN "исппроизввотнюрлиц".dateofexecutivedocument; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".dateofexecutivedocument IS 'Дата исполнительного документа';


--
-- Name: COLUMN "исппроизввотнюрлиц".numberofexecutivedocument; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".numberofexecutivedocument IS 'Номер исполнительного документа';


--
-- Name: COLUMN "исппроизввотнюрлиц".objectofexecutivedocuments; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".objectofexecutivedocuments IS 'Требования исполнительного документа';


--
-- Name: COLUMN "исппроизввотнюрлиц".objectofexecution; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".objectofexecution IS 'Предмет исполнения';


--
-- Name: COLUMN "исппроизввотнюрлиц".amountdue; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".amountdue IS 'Сумма долга';


--
-- Name: COLUMN "исппроизввотнюрлиц".debtremainingbalance; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".debtremainingbalance IS 'Остаток непогашенной задолженности';


--
-- Name: COLUMN "исппроизввотнюрлиц".departmentsofbailiffs; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".departmentsofbailiffs IS 'Наименование отдела';


--
-- Name: COLUMN "исппроизввотнюрлиц".addressofdepartmentsofbailiff; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".addressofdepartmentsofbailiff IS 'Адрес отдела судебных приставов';


--
-- Name: COLUMN "исппроизввотнюрлиц".debtortaxpayeridentificationnumber; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".debtortaxpayeridentificationnumber IS 'Идентификационный номер налогоплательщика должника';


--
-- Name: COLUMN "исппроизввотнюрлиц".taxpayeridentificationnumberoforganizationcollector; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".taxpayeridentificationnumberoforganizationcollector IS 'Идентификационный номер налогоплательщика взыскателя-организации';


--
-- Name: COLUMN "исппроизввотнюрлиц".repaid; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."исппроизввотнюрлиц".repaid IS 'погашен';


--
-- Name: налоговыеправонарушенияиштрафы; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."налоговыеправонарушенияиштрафы" (
    "датасост" date NOT NULL,
    "иннюл" text NOT NULL,
    "наиморг" text NOT NULL,
    "сумштраф" numeric DEFAULT 0 NOT NULL
);


--
-- Name: оквэд; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."оквэд" (
    "кодоквэд" text NOT NULL,
    "наимоквэд" text NOT NULL
);


--
-- Name: октао; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."октао" (
    ter text NOT NULL,
    kod1 text NOT NULL,
    kod2 text NOT NULL,
    kod3 text NOT NULL,
    razdel text NOT NULL,
    name text,
    centrum text,
    nomdescr text,
    nomakt text,
    status text,
    dateutv date,
    datevved date
);


--
-- Name: октмо; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."октмо" (
    ter text NOT NULL,
    kod1 text NOT NULL,
    kod2 text NOT NULL,
    kod3 text NOT NULL,
    razdel text NOT NULL,
    name text,
    centrum text,
    nomdescr text,
    nomakt text,
    status text,
    dateutv date,
    datevved date
);


--
-- Name: открытыйреестртоварныхзнаков; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."открытыйреестртоварныхзнаков" (
    registrationnumber text NOT NULL,
    registrationdate text,
    applicationnumber text,
    applicationdate text,
    prioritydate text,
    exhibitionprioritydate text,
    parisconventionprioritynumber text,
    parisconventionprioritydate text,
    parisconventionprioritycountrycode text,
    initialapplicationnumber text,
    initialapplicationorioritydate text,
    initialregistrationnumber text,
    initialregistrationdate text,
    internationalregistrationnumber text,
    internationalregistrationdate text,
    internationalregistrationprioritydate text,
    internationalregistrationentrydate text,
    applicationnumberforrecognitionoftrademarkfromcrimea text,
    applicationdateforrecognitionoftrademarkfromcrimea text,
    crimeantrademarkapplicationnumberforstateregistrationinukraine text,
    crimeantrademarkapplicationdateforstateregistrationinukraine text,
    crimeantrademarkcertificatenumberinukraine text,
    exclusiverightstransferagreementregistrationnumber text,
    exclusiverightstransferagreementregistrationdate text,
    legallyrelatedapplications text,
    legallyrelatedregistrations text,
    expirationdate text,
    rightholdername text,
    foreignrightholdername text,
    rightholderaddress text,
    rightholdercountrycode text,
    rightholderogrn text,
    rightholderinn text,
    correspondenceaddress text,
    collective boolean,
    collectiveusers text,
    extractionfromcharterofthecollectivetrademark text,
    colorspecification text,
    unprotectedelements text,
    kindspecification text,
    threedimensional boolean,
    threedimensionalspecification text,
    holographic boolean,
    holographicspecification text,
    sound boolean,
    soundspecification text,
    olfactory boolean,
    olfactoryspecification text,
    color boolean,
    colortrademarkspecification text,
    light boolean,
    lightspecification text,
    changing boolean,
    changingspecification text,
    positional boolean,
    positionalspecification text,
    actual boolean,
    publicationurl text
);


--
-- Name: реестрдисквалифицированныхлиц; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."реестрдисквалифицированныхлиц" (
    id bigint NOT NULL,
    fio text,
    bdate date,
    bplace text,
    orgname text,
    inn text,
    positionfl text,
    nkoap text,
    gorgname text,
    sudfio text,
    sudposition text,
    disqualificationduration text,
    disstartdate date,
    disenddate date
);


--
-- Name: реестробщеизвестныхтоварныхзнак; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."реестробщеизвестныхтоварныхзнак" (
    registrationnumber text NOT NULL,
    registrationdate date,
    wellknowntrademarkdate date,
    legallyrelatedregistrations text,
    rightholdername text,
    foreignrightholdername text,
    rightholderaddress text,
    rightholdercountrycode text,
    rightholderogrn text,
    rightholderinn text,
    correspondenceaddress text,
    collective boolean,
    collectiveusers text,
    extractionfromcharterofcollectivetrademark text,
    colorspecification text,
    unprotectedelements text,
    kindspecification text,
    threedimensional boolean,
    threedimensionalspecification text,
    holographic boolean,
    holographicspecification text,
    sound boolean,
    soundspecification text,
    olfactory boolean,
    olfactoryspecification text,
    color boolean,
    colortrademarkspecification text,
    light boolean,
    lightspecification text,
    changing boolean,
    changingspecification text,
    positional boolean,
    positionalspecification text,
    actual boolean,
    publicationurl text
);


--
-- Name: режимналогоплательщика; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."режимналогоплательщика" (
    "датадок" date NOT NULL,
    "наиморг" text NOT NULL,
    "иннюл" text NOT NULL,
    "есхн" boolean NOT NULL,
    "усн" boolean NOT NULL,
    "енвд" boolean NOT NULL,
    "срп" boolean NOT NULL
);


--
-- Name: росаккредитация; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."росаккредитация" (
    id_cert text NOT NULL,
    cert_status text,
    cert_type text,
    reg_number text,
    date_begining date,
    date_finish date,
    product_scheme text,
    product_object_type_cert text,
    product_type text,
    product_okpd2 text,
    product_tn_ved text,
    product_tech_reg text,
    product_group text,
    product_name text,
    product_info text,
    applicant_type text,
    person_applicant_type text,
    applicant_ogrn text,
    applicant_inn text,
    applicant_phone text,
    applicant_fax text,
    applicant_email text,
    applicant_website text,
    applicant_name text,
    applicant_director_name text,
    applicant_address text,
    applicant_address_actual text,
    manufacturer_type text,
    manufacturer_ogrn text,
    manufacturer_inn text,
    manufacturer_phone text,
    manufacturer_fax text,
    manufacturer_email text,
    manufacturer_website text,
    manufacturer_name text,
    manufacturer_director_name text,
    manufacturer_country text,
    manufacturer_address text,
    manufacturer_address_actual text,
    manufacturer_address_filial text,
    organ_to_certification_name text,
    organ_to_certification_reg_number text,
    organ_to_certification_head_name text,
    basis_for_certificate text,
    old_basis_for_certificate text,
    fio_expert text,
    fio_signatory text,
    product_national_standart text,
    production_analysis_for_act text,
    production_analysis_for_act_number text,
    production_analysis_for_act_date date
);


--
-- Name: сведенияобуплаченныхорганизацие; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."сведенияобуплаченныхорганизацие" (
    "датасост" date NOT NULL,
    "наиморг" text,
    "иннюл" text NOT NULL,
    "наимналог" text NOT NULL,
    "сумуплнал" numeric DEFAULT 0 NOT NULL
);


--
-- Name: сведенияобучастиивконсгруппе; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."сведенияобучастиивконсгруппе" (
    "датасост" date NOT NULL,
    "наиморг" text,
    "иннюл" text NOT NULL,
    "признучкгн" integer NOT NULL
);


--
-- Name: сведенияосреднчислработников; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."сведенияосреднчислработников" (
    "датасост" date NOT NULL,
    "наиморг" text,
    "иннюл" text NOT NULL,
    "колраб" integer DEFAULT 0 NOT NULL
);


--
-- Name: сведенияосуммахнедоимки; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."сведенияосуммахнедоимки" (
    "датадок" date NOT NULL,
    "наиморг" text,
    "иннюл" text NOT NULL,
    "наимналог" text NOT NULL,
    "сумнедналог" numeric DEFAULT 0 NOT NULL,
    "сумпени" numeric DEFAULT 0 NOT NULL,
    "сумштраф" numeric DEFAULT 0 NOT NULL,
    "общсумнедоим" numeric DEFAULT 0 NOT NULL,
    "иддок" text NOT NULL,
    "датасост" date
);


--
-- Name: COLUMN "сведенияосуммахнедоимки"."датадок"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."сведенияосуммахнедоимки"."датадок" IS 'Дата документа';


--
-- Name: COLUMN "сведенияосуммахнедоимки"."иддок"; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public."сведенияосуммахнедоимки"."иддок" IS 'ID документа';


--
-- Name: смпоквэды; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."смпоквэды" (
    "инн" text NOT NULL,
    "огрн" text NOT NULL,
    "кодоквэд" text NOT NULL,
    "основной" boolean DEFAULT false NOT NULL
);


--
-- Name: субъектымалогоисреднегопредприн; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."субъектымалогоисреднегопредприн" (
    "датасост" date NOT NULL,
    "датавклмсп" text,
    "видсубмсп" integer,
    "катсубмсп" integer,
    "призновмсп" integer,
    "сведсоцпред" integer,
    "ссчр" integer,
    "наиморг" text,
    "наиморгсокр" text,
    "иннюл" text,
    "огрнюл" text,
    "иннфл" text,
    "огрнип" text,
    "фамилия" text,
    "имя" text,
    "отчество" text,
    "номлиценз" text,
    "даталиценз" text,
    "датаначлиценз" text,
    "датаконлиценз" text,
    "датаостлиценз" text,
    "серлиценз" text,
    "видлиценз" text,
    "оргвыдлиценз" text,
    "оргостлиценз" text,
    "наимлицвд" text,
    "кодрегион" integer,
    "регионтип" text,
    "регионнаим" text,
    "районтип" text,
    "районнаим" text,
    "городтип" text,
    "городнаим" text,
    "населпункттип" text,
    "населпунктнаим" text,
    "инн" text NOT NULL,
    "огрн" text NOT NULL,
    "свпрод" jsonb DEFAULT '[]'::jsonb,
    "свпрогпарт" jsonb DEFAULT '[]'::jsonb,
    "свконтр" jsonb DEFAULT '[]'::jsonb,
    "свдог" jsonb DEFAULT '[]'::jsonb
);


--
-- Name: source_files; Type: TABLE; Schema: service_management; Owner: -
--

CREATE TABLE service_management.source_files (
    source_type text NOT NULL,
    source_link text NOT NULL,
    filename text NOT NULL,
    sha256sum text NOT NULL,
    downloaded boolean DEFAULT false NOT NULL,
    uploaded boolean DEFAULT false NOT NULL,
    task_datetime timestamp with time zone DEFAULT now() NOT NULL,
    id uuid DEFAULT service_management.uuid_generate_v4() NOT NULL
);


--
-- Name: баланс баланс_pk; Type: CONSTRAINT; Schema: accounting_statements; Owner: -
--

ALTER TABLE ONLY accounting_statements."баланс"
    ADD CONSTRAINT "баланс_pk" PRIMARY KEY ("иннюл", "год");


--
-- Name: egr_search egr_search_t_pk; Type: CONSTRAINT; Schema: egr; Owner: -
--

ALTER TABLE ONLY egr.egr_search
    ADD CONSTRAINT egr_search_t_pk PRIMARY KEY (inn, ogrn);


--
-- Name: егрип егрип_pk; Type: CONSTRAINT; Schema: egr; Owner: -
--

ALTER TABLE ONLY egr."егрип"
    ADD CONSTRAINT "егрип_pk" PRIMARY KEY ("инн", "огрн");


--
-- Name: егрюл егрюл_pk; Type: CONSTRAINT; Schema: egr; Owner: -
--

ALTER TABLE ONLY egr."егрюл"
    ADD CONSTRAINT "егрюл_pk" PRIMARY KEY ("инн", "огрн");


--
-- Name: iplegallistcomplete iplegallistcomplete_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.iplegallistcomplete
    ADD CONSTRAINT iplegallistcomplete_pk PRIMARY KEY (numberofenforcementproceeding, numberofexecutivedocument);


--
-- Name: реестробщеизвестныхтоварныхзнак otz_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."реестробщеизвестныхтоварныхзнак"
    ADD CONSTRAINT otz_pk PRIMARY KEY (registrationnumber);


--
-- Name: rds rds_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.rds
    ADD CONSTRAINT rds_pk PRIMARY KEY (id_decl);


--
-- Name: открытыйреестртоварныхзнаков tz_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."открытыйреестртоварныхзнаков"
    ADD CONSTRAINT tz_pk PRIMARY KEY (registrationnumber);


--
-- Name: hotels гостиницы_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hotels
    ADD CONSTRAINT "гостиницы_pk" PRIMARY KEY (federal_number);


--
-- Name: информация1928коап информация1928коап_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."информация1928коап"
    ADD CONSTRAINT "информация1928коап_pk" PRIMARY KEY ("инн", "огрн", "номердела");


--
-- Name: исппроизввотнюрлиц исппроизввотнюрлиц_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."исппроизввотнюрлиц"
    ADD CONSTRAINT "исппроизввотнюрлиц_pk" PRIMARY KEY (numberofenforcementproceeding, numberofexecutivedocument);


--
-- Name: налоговыеправонарушенияиштрафы налоговыеправонарушенияиштрафы_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."налоговыеправонарушенияиштрафы"
    ADD CONSTRAINT "налоговыеправонарушенияиштрафы_pk" PRIMARY KEY ("датасост", "иннюл");


--
-- Name: оквэд оквэд_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."оквэд"
    ADD CONSTRAINT "оквэд_pk" PRIMARY KEY ("кодоквэд");


--
-- Name: реестрдисквалифицированныхлиц реестрдисквалифицированныхлиц_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."реестрдисквалифицированныхлиц"
    ADD CONSTRAINT "реестрдисквалифицированныхлиц_pk" PRIMARY KEY (id);


--
-- Name: режимналогоплательщика режимналогоплательщика_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."режимналогоплательщика"
    ADD CONSTRAINT "режимналогоплательщика_pk" PRIMARY KEY ("иннюл");


--
-- Name: росаккредитация росаккредитация_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."росаккредитация"
    ADD CONSTRAINT "росаккредитация_pk" PRIMARY KEY (id_cert);


--
-- Name: сведенияобуплаченныхорганизацие сведенияобуплаченныхорганизацие_; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."сведенияобуплаченныхорганизацие"
    ADD CONSTRAINT "сведенияобуплаченныхорганизацие_" PRIMARY KEY ("иннюл", "наимналог");


--
-- Name: сведенияобучастиивконсгруппе сведенияобучастиивконсгруппе_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."сведенияобучастиивконсгруппе"
    ADD CONSTRAINT "сведенияобучастиивконсгруппе_pk" PRIMARY KEY ("иннюл");


--
-- Name: сведенияосреднчислработников сведенияосреднчислработников_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."сведенияосреднчислработников"
    ADD CONSTRAINT "сведенияосреднчислработников_pk" PRIMARY KEY ("иннюл");


--
-- Name: сведенияосуммахнедоимки сведенияосуммахнедоимки_un; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."сведенияосуммахнедоимки"
    ADD CONSTRAINT "сведенияосуммахнедоимки_un" UNIQUE ("иддок", "наимналог");


--
-- Name: смпоквэды смпоквэды_un; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."смпоквэды"
    ADD CONSTRAINT "смпоквэды_un" UNIQUE ("инн", "огрн", "кодоквэд", "основной");


--
-- Name: субъектымалогоисреднегопредприн субъектымалогоисреднегопредприн_; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."субъектымалогоисреднегопредприн"
    ADD CONSTRAINT "субъектымалогоисреднегопредприн_" PRIMARY KEY ("инн", "огрн");


--
-- Name: source_files source_files_pk; Type: CONSTRAINT; Schema: service_management; Owner: -
--

ALTER TABLE ONLY service_management.source_files
    ADD CONSTRAINT source_files_pk PRIMARY KEY (id);


--
-- Name: source_files source_files_un; Type: CONSTRAINT; Schema: service_management; Owner: -
--

ALTER TABLE ONLY service_management.source_files
    ADD CONSTRAINT source_files_un UNIQUE (source_type, filename, sha256sum);


--
-- Name: b_verion; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX b_verion ON accounting_statements."баланс" USING btree (((balance_data ->> 'ВерсФорм'::text)));


--
-- Name: balance_mat_view_accounts_payable_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_accounts_payable_idx ON accounting_statements.balance_mat_view USING btree (accounts_payable);


--
-- Name: balance_mat_view_accounts_receivable_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_accounts_receivable_idx ON accounting_statements.balance_mat_view USING btree (accounts_receivable);


--
-- Name: balance_mat_view_basic_assets_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_basic_assets_idx ON accounting_statements.balance_mat_view USING btree (basic_assets);


--
-- Name: balance_mat_view_borrowed_funds_long_term_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_borrowed_funds_long_term_idx ON accounting_statements.balance_mat_view USING btree (borrowed_funds_long_term);


--
-- Name: balance_mat_view_borrowed_funds_short_term_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_borrowed_funds_short_term_idx ON accounting_statements.balance_mat_view USING btree (borrowed_funds_short_term);


--
-- Name: balance_mat_view_capital_and_reserves_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_capital_and_reserves_idx ON accounting_statements.balance_mat_view USING btree (capital_and_reserves);


--
-- Name: balance_mat_view_cash_and_equivalents_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_cash_and_equivalents_idx ON accounting_statements.balance_mat_view USING btree (cash_and_equivalents);


--
-- Name: balance_mat_view_commercial_expenses_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_commercial_expenses_idx ON accounting_statements.balance_mat_view USING btree (commercial_expenses);


--
-- Name: balance_mat_view_cost_of_sales_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_cost_of_sales_idx ON accounting_statements.balance_mat_view USING btree (cost_of_sales);


--
-- Name: balance_mat_view_current_assets_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_current_assets_idx ON accounting_statements.balance_mat_view USING btree (current_assets);


--
-- Name: balance_mat_view_doc_date_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_doc_date_idx ON accounting_statements.balance_mat_view USING btree (doc_date DESC);


--
-- Name: balance_mat_view_expenses_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_expenses_idx ON accounting_statements.balance_mat_view USING btree (expenses);


--
-- Name: balance_mat_view_gross_profit_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_gross_profit_idx ON accounting_statements.balance_mat_view USING btree (gross_profit);


--
-- Name: balance_mat_view_income_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_income_idx ON accounting_statements.balance_mat_view USING btree (income);


--
-- Name: balance_mat_view_income_tax_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_income_tax_idx ON accounting_statements.balance_mat_view USING btree (income_tax);


--
-- Name: balance_mat_view_inn_year_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE UNIQUE INDEX balance_mat_view_inn_year_idx ON accounting_statements.balance_mat_view USING btree (year, inn);


--
-- Name: balance_mat_view_intangible_assets_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_intangible_assets_idx ON accounting_statements.balance_mat_view USING btree (intangible_assets);


--
-- Name: balance_mat_view_interest_payable_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_interest_payable_idx ON accounting_statements.balance_mat_view USING btree (interest_payable);


--
-- Name: balance_mat_view_management_expenses_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_management_expenses_idx ON accounting_statements.balance_mat_view USING btree (management_expenses);


--
-- Name: balance_mat_view_net_aassets_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_net_aassets_idx ON accounting_statements.balance_mat_view USING btree (net_aassets);


--
-- Name: balance_mat_view_net_profit_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_net_profit_idx ON accounting_statements.balance_mat_view USING btree (net_profit);


--
-- Name: balance_mat_view_non_current_assets_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_non_current_assets_idx ON accounting_statements.balance_mat_view USING btree (non_current_assets);


--
-- Name: balance_mat_view_other_expenses_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_other_expenses_idx ON accounting_statements.balance_mat_view USING btree (other_expenses);


--
-- Name: balance_mat_view_other_income_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_other_income_idx ON accounting_statements.balance_mat_view USING btree (other_income);


--
-- Name: balance_mat_view_other_non_current_assets_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_other_non_current_assets_idx ON accounting_statements.balance_mat_view USING btree (other_non_current_assets);


--
-- Name: balance_mat_view_other_short_term_liabilities_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_other_short_term_liabilities_idx ON accounting_statements.balance_mat_view USING btree (other_short_term_liabilities);


--
-- Name: balance_mat_view_profit_before_taxation_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_profit_before_taxation_idx ON accounting_statements.balance_mat_view USING btree (profit_before_taxation);


--
-- Name: balance_mat_view_profit_from_sale_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_profit_from_sale_idx ON accounting_statements.balance_mat_view USING btree (profit_from_sale);


--
-- Name: balance_mat_view_revenue_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_revenue_idx ON accounting_statements.balance_mat_view USING btree (revenue);


--
-- Name: balance_mat_view_stocks_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_stocks_idx ON accounting_statements.balance_mat_view USING btree (stocks);


--
-- Name: balance_mat_view_tax_usn_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_tax_usn_idx ON accounting_statements.balance_mat_view USING btree (tax_usn);


--
-- Name: balance_mat_view_total_assets_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_total_assets_idx ON accounting_statements.balance_mat_view USING btree (total_assets);


--
-- Name: balance_mat_view_total_liabilities_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX balance_mat_view_total_liabilities_idx ON accounting_statements.balance_mat_view USING btree (total_liabilities);


--
-- Name: balance_mat_view_сurrent_income_tax_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX "balance_mat_view_сurrent_income_tax_idx" ON accounting_statements.balance_mat_view USING btree ("сurrent_income_tax");


--
-- Name: okpo_inn_idx; Type: INDEX; Schema: accounting_statements; Owner: -
--

CREATE INDEX okpo_inn_idx ON accounting_statements.okpo USING btree (inn);


--
-- Name: egr_search_okved_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX egr_search_okved_idx ON egr.egr_search USING btree (((okved ->> 'code'::text))) WHERE (okved IS NOT NULL);


--
-- Name: egr_search_t_date_discharge_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX egr_search_t_date_discharge_idx ON egr.egr_search USING btree (date_discharge DESC);


--
-- Name: egr_search_t_date_liquidation_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX egr_search_t_date_liquidation_idx ON egr.egr_search USING btree (date_liquidation) WHERE (date_liquidation IS NOT NULL);


--
-- Name: egr_search_t_date_liquidation_null_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX egr_search_t_date_liquidation_null_idx ON egr.egr_search USING btree (date_liquidation) WHERE (date_liquidation IS NULL);


--
-- Name: egr_search_t_date_registration_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX egr_search_t_date_registration_idx ON egr.egr_search USING btree (date_registration) WHERE (date_registration IS NOT NULL);


--
-- Name: egr_search_t_inn_sort_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX egr_search_t_inn_sort_idx ON egr.egr_search USING btree (inn);


--
-- Name: egr_search_t_is_legal_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX egr_search_t_is_legal_idx ON egr.egr_search USING btree (is_legal);


--
-- Name: egr_search_t_kpp_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX egr_search_t_kpp_idx ON egr.egr_search USING btree (kpp) WHERE (kpp IS NOT NULL);


--
-- Name: egr_search_t_ogrn_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX egr_search_t_ogrn_idx ON egr.egr_search USING btree (ogrn);


--
-- Name: egr_search_t_status_notnull_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX egr_search_t_status_notnull_idx ON egr.egr_search USING hash (status) WHERE (status IS NOT NULL);


--
-- Name: egr_search_t_text_idx_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX egr_search_t_text_idx_idx ON egr.egr_search USING gin (idx_name_full);


--
-- Name: handbook_okved_trgm_code_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX handbook_okved_trgm_code_idx ON egr.handbook_okved USING gin (code);


--
-- Name: handbook_okved_trgm_title_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX handbook_okved_trgm_title_idx ON egr.handbook_okved USING gin (title);


--
-- Name: handbook_tax_authority_code_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE UNIQUE INDEX handbook_tax_authority_code_idx ON egr.handbook_tax_authority USING btree (code, name);


--
-- Name: statuses_status_full_name_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX statuses_status_full_name_idx ON egr.statuses USING hash (status_full_name);


--
-- Name: statuses_status_name_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX statuses_status_name_idx ON egr.statuses USING hash (status_name);


--
-- Name: statuses_un; Type: INDEX; Schema: egr; Owner: -
--

CREATE UNIQUE INDEX statuses_un ON egr.statuses USING btree (status_full_name);


--
-- Name: tax_authority_code_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE INDEX tax_authority_code_idx ON egr.tax_authority USING btree (code);


--
-- Name: tax_authority_inn_idx; Type: INDEX; Schema: egr; Owner: -
--

CREATE UNIQUE INDEX tax_authority_inn_idx ON egr.tax_authority USING btree (inn, ogrn);


--
-- Name: iplegallistcomplete_debtortaxpayeridentificati; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX iplegallistcomplete_debtortaxpayeridentificati ON public.iplegallistcomplete USING btree (debtortaxpayeridentificationnumber);


--
-- Name: otz_inn_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX otz_inn_idx ON public."реестробщеизвестныхтоварныхзнак" USING btree (rightholderinn);


--
-- Name: rds_applicant_inn_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX rds_applicant_inn_idx ON public.rds USING btree (applicant_inn, applicant_ogrn);


--
-- Name: гостиницы_инн_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "гостиницы_инн_idx" ON public.hotels USING btree (inn);


--
-- Name: гостиницы_огрн_огрнип_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "гостиницы_огрн_огрнип_idx" ON public.hotels USING btree (ogrn);


--
-- Name: исппроизввотнюрлиц_debtortaxpayeridentificati; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "исппроизввотнюрлиц_debtortaxpayeridentificati" ON public."исппроизввотнюрлиц" USING btree (debtortaxpayeridentificationnumber);


--
-- Name: исппроизввотнюрлиц_repaid_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "исппроизввотнюрлиц_repaid_idx" ON public."исппроизввотнюрлиц" USING btree (repaid);


--
-- Name: классификациягостиниц_номер_гост; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "классификациягостиниц_номер_гост" ON public.hotels_classification USING btree (federal_number);


--
-- Name: номерагостиниц_номер_гостиницы_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "номерагостиниц_номер_гостиницы_idx" ON public.hotels_rooms USING btree (federal_number);


--
-- Name: открытыйреестртоварныхзнаков_righth; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "открытыйреестртоварныхзнаков_righth" ON public."открытыйреестртоварныхзнаков" USING btree (rightholderinn);


--
-- Name: реестрдисквалифицированныхлиц_inn_; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "реестрдисквалифицированныхлиц_inn_" ON public."реестрдисквалифицированныхлиц" USING btree (inn);


--
-- Name: режимналогоплательщика_енвд_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "режимналогоплательщика_енвд_idx" ON public."режимналогоплательщика" USING btree ("енвд");


--
-- Name: режимналогоплательщика_есхн_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "режимналогоплательщика_есхн_idx" ON public."режимналогоплательщика" USING btree ("есхн");


--
-- Name: режимналогоплательщика_срп_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "режимналогоплательщика_срп_idx" ON public."режимналогоплательщика" USING btree ("срп");


--
-- Name: режимналогоплательщика_усн_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "режимналогоплательщика_усн_idx" ON public."режимналогоплательщика" USING btree ("усн");


--
-- Name: росаккредитация_applicant_inn_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "росаккредитация_applicant_inn_idx" ON public."росаккредитация" USING btree (applicant_inn);


--
-- Name: сведенияосуммахнедоимки_иннюл_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX "сведенияосуммахнедоимки_иннюл_idx" ON public."сведенияосуммахнедоимки" USING btree ("иннюл");


--
-- Name: source_files_task_datetime_idx; Type: INDEX; Schema: service_management; Owner: -
--

CREATE INDEX source_files_task_datetime_idx ON service_management.source_files USING btree (task_datetime DESC);


--
-- Name: егрип trigger_egrip_search; Type: TRIGGER; Schema: egr; Owner: -
--

CREATE TRIGGER trigger_egrip_search AFTER INSERT OR UPDATE ON egr."егрип" FOR EACH ROW EXECUTE FUNCTION egr.trigger_proc_egrip_search();


--
-- Name: егрюл trigger_egrul_search; Type: TRIGGER; Schema: egr; Owner: -
--

CREATE TRIGGER trigger_egrul_search AFTER INSERT OR UPDATE ON egr."егрюл" FOR EACH ROW EXECUTE FUNCTION egr.trigger_proc_egrul_search();


--
-- Name: hotels_classification hotels_classification_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hotels_classification
    ADD CONSTRAINT hotels_classification_fk FOREIGN KEY (federal_number) REFERENCES public.hotels(federal_number) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: hotels_rooms hotels_rooms_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hotels_rooms
    ADD CONSTRAINT hotels_rooms_fk FOREIGN KEY (federal_number) REFERENCES public.hotels(federal_number) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: смпоквэды кодоквэды_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."смпоквэды"
    ADD CONSTRAINT "кодоквэды_fk" FOREIGN KEY ("кодоквэд") REFERENCES public."оквэд"("кодоквэд") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: смпоквэды смпоквэды_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public."смпоквэды"
    ADD CONSTRAINT "смпоквэды_fk" FOREIGN KEY ("инн", "огрн") REFERENCES public."субъектымалогоисреднегопредприн"("инн", "огрн") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

