#!/bin/bash

help_text="""
Необходима переменная окружения 'PGPASSWORD'\n\n

Указать полное название таблицы, например: 'egr.egr_search'\n
./refresh_materialized_view.bash egr.egr_search 1\n\n

Вторым аргументов нужно указать '1', если таблица уже материалозована,
это добавляет к команде опцию 'CONCURENT' - обновление таблицы без блокировки.\n\n

Список материалованных представлений:\n\n
1.\tegr.egr_search\t\t\t\t- зависит от таблиц ЕГРИП и ЕГРЮЛ\n
1.2.\tegr.tax_authority\t\t\t- зависит от egr.egr_search\n
1.3.\tegr.handbook_tax_authority\t\t- зависит от мат.вью egr.tax_authority\n
2.\tegr.okved\t\t\t\t- зависит от таблиц ЕГРИП и ЕГРЮЛ\n
2.1.\tegr.handbook_okved\t\t\t- зависит от egr.handbook_okved\n
3.\taccounting_statements.balance_mat_view\t- зависит только от таблицы accounting_statements.баланс
"""

if [[ "$1" == "help" || $# -lt 1 ]]; then
    echo -e $help_text
    exit 0
fi

concurent=""
if [[ $2 -eq 1 ]]; then
    concurent=" CONCURRENTLY "
fi
command="REFRESH MATERIALIZED VIEW $concurent $1 WITH DATA"
echo $command
psql -U opendata_app -d opendata_aggregator -p 7777 -c "$command"
