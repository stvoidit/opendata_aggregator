### Источники

|Название|Источник|
|---|---|
| Сведения о специальных налоговых режимах, применяемых налогоплательщиками |https://www.nalog.gov.ru/opendata/7707329152-snr/|
| Сведения из Реестра сертификатов соответствия |https://fsa.gov.ru/opendata/7736638268-rss/|
| Сведения о суммах недоимки и задолженности по пеням и штрафам |https://www.nalog.gov.ru/opendata/7707329152-debtam/|
| Сведения об участии в консолидированной группе налогоплательщиков |https://www.nalog.gov.ru/opendata/7707329152-kgn/|
| Сведения о среднесписочной численности работников организации |https://www.nalog.gov.ru/opendata/7707329152-sshr2019/|
| Сведения об уплаченных организацией в календарном году, предшествующем году размещения указанных сведений в информационно-телекоммуникационной сети "Интернет" в соответствии с пунктом 1.1 статьи 102 Налогового кодекса Российской Федерации, суммах налогов и сборов (по каждому налогу и сбору) без учета сумм налогов (сборов), уплаченных в связи с ввозом товаров на таможенную территорию Евразийского экономического союза, сумм налогов, уплаченных налоговым агентом, о суммах страховых взносов |https://www.nalog.gov.ru/opendata/7707329152-paytax/|
| Исполнительные производства в отношении юридических лиц |https://opendata.fssp.gov.ru/7709576929-iplegallist|
| Открытый реестр товарных знаков и знаков обслуживания Российской Федерации |https://rospatent.gov.ru/opendata/7730176088-tz|
| Открытый реестр общеизвестных в Российской Федерации товарных знаков |https://rospatent.gov.ru/opendata/7730176088-otz|
| Единый реестр субъектов малого и среднего предпринимательства |https://www.nalog.gov.ru/opendata/7707329152-rsmp/|
| Общероссийский классификатор видов экономической деятельности (ОКВЭД2) |https://rosstat.gov.ru/opendata/7708234640-okved2|
| Сведения о налоговых правонарушениях и мерах ответственности за их совершение |https://www.nalog.gov.ru/opendata/7707329152-taxoffence/|
| Реестр дисквалифицированных лиц |https://www.nalog.gov.ru/opendata/7707329152-registerdisqualified/|
| Общероссийский классификатор территорий муниципальных образований (ОКТМО) |https://rosstat.gov.ru/opendata/7708234640-oktmo|
| Общероссийский классификатор объектов административно-территориального деления (ОКАТО) |https://rosstat.gov.ru/opendata/7708234640-okato|
| Информация о привлечении участника закупки к административной ответственности по ст. 19.28 КоАП |https://zakupki.gov.ru/epz/main/public/document/view.html?searchString=&sectionId=2369&strictEqual=false|

### Endpoints

| Метод 	| Путь                   	| Описание                                                                      	|
|-------	|------------------------	|-------------------------------------------------------------------------------	|
| GET   	| /api/search?                   	| Поиск по параметрам (возвращает список)                                       	|
| GET   	| /api/count?                   	| Кол-во строк в результате поиска для /api/search - принимает те же параметры                                       	|
| GET   	| /api/info/{inn:[0-9]+}        	| Полная карточка информации по ИНН (старый вариант, не желательно к использованию)                                             	|
| GET   	| /api/info/{inn:[0-9]+}/{ogrn:[0-9]+}        	| Полная карточка информации по ИНН                                             	|
| GET   	| /api/handbook_tax_authority    	| Возвращает полный справочник налоговых органов (имя парамерта: __ta__)        	|
| GET   	| /api/search_tax_authority?q=    	| Поиск кода НО - возвращает список, параметр q - строка, минимум 3 символа     	|
| GET   	| /api/handbook_okved            	| Возвращает полный справочник ОКВЭД (имя парамерта: __okved__)                 	|
| GET   	| /api/search_okved?q=   	        | Поиск кодов ОКВЭД - возвращает список, параметр q - строка, минимум 3 символа 	|
| GET   	| /api/handbook_categories_ip        	| Справочник категорий ИП                                             	|

### Параметры для запроса /api/search

| Название параметра           	| Описание                                 	| Множественное 	| Пример значения 	|
|------------------------------	|------------------------------------------	|---------------	|-----------------	|
| year                       	| Год бухгалтерского отчета                	| да            	| 2021            	|
| income                       	| Доходы                                   	| нет           	| 1200            	|
| expenses                     	| Расходы                                  	| нет           	| 1200            	|
| income_tax                   	| Налог на прибыль                         	| нет           	| 1200            	|
| tax_usn                      	| Налог по УСН                             	| нет           	| 1200            	|
| intangible_assets            	| Нематериальные активы                    	| нет           	| 1200            	|
| basic_assets                 	| Основные средства                        	| нет           	| 1200            	|
| other_non_current_assets     	| Прочие внеоборотные активы               	| нет           	| 1200            	|
| non_current_assets           	| Внеоборотные активы                      	| нет           	| 1200            	|
| stocks                       	| Запасы                                   	| нет           	| 1200            	|
| net_aassets                  	| Чистые активы                            	| нет           	| 1200            	|
| accounts_receivable          	| Дебиторская задолженность                	| нет           	| 1200            	|
| cash_and_equivalents         	| Денежные средства и денежные эквиваленты 	| нет           	| 1200            	|
| current_assets               	| Оборотные активы                         	| нет           	| 1200            	|
| total_assets                 	| Активы всего                             	| нет           	| 1200            	|
| capital_and_reserves         	| Капитал и резервы                        	| нет           	| 1200            	|
| borrowed_funds_long_term     	| Заёмные средства (долгосрочные)          	| нет           	| 1200            	|
| borrowed_funds_short_term    	| Заёмные средства (краткосрочные)         	| нет           	| 1200            	|
| accounts_payable             	| Кредиторская задолженность               	| нет           	| 1200            	|
| other_short_term_liabilities 	| Прочие краткосрочные обязательства       	| нет           	| 1200            	|
| total_liabilities            	| Пассивы всего                            	| нет           	| 1200            	|
| revenue                      	| Выручка                                  	| нет           	| 1200            	|
| cost_of_sales                	| Себестоимость продаж                     	| нет           	| 1200            	|
| gross_profit                 	| Валовая прибыль                          	| нет           	| 1200            	|
| commercial_expenses          	| Коммерческие расходы                     	| нет           	| 1200            	|
| management_expenses          	| Управленческие расходы                   	| нет           	| 1200            	|
| profit_from_sale             	| Прибыль (убыток) от продажи              	| нет           	| 1200            	|
| interest_payable             	| Проценты к уплате                        	| нет           	| 1200            	|
| other_income                 	| Прочие доходы                            	| нет           	| 1200            	|
| other_expenses               	| Прочие расходы                           	| нет           	| 1200            	|
| profit_before_taxation       	| Прибыль (убыток) до налогообложения      	| нет           	| 1200            	|
| сurrent_income_tax           	| Текущий налог на прибыль                 	| нет           	| 1200            	|
| net_profit                   	| Чистая прибыль (убыток)                  	| нет           	| 1200            	|
| okved                        	| Код ОКВЭД                                	| да            	| 10.41.22        	|
| ta                           	| Код Налогового органа                    	| нет           	| 7715            	|
| date_registration_from       	| Дата регистрации С                       	| нет           	| 2000-12-31      	|
| date_registration_to         	| Дата регистрации ПО                      	| нет           	| 2022-12-31      	|
| date_liquidation_from        	| Дата ликвидации С                        	| нет           	| 2000-12-31      	|
| date_liquidation_to          	| Дата ликвидации ПО                       	| нет           	| 2022-12-31      	|
| kpp                          	| Код КПП                                  	| нет           	| 370201001       	|
| usn                          	| Налоговый режим УСН                      	| нет           	| _пусто или 1_    	|
| eshn                         	| Налоговый режим ЕСХН                     	| нет           	| _пусто или 1_    	|
| envd                         	| Налоговый режим ЕНВД                     	| нет           	| _пусто или 1_    	|
| srp                          	| Налоговый режим СРП                      	| нет           	| _пусто или 1_    	|
| status                       	| Сатутс (a - действующие (active), d - недействующие, l - ликвидация, b - банкротство )                 	| нет           	| d                	|
| q                          	| Текстовый поиск по названию компании или ФИО + ИНН + ОГРН                      	| нет           	| _магнит_    	|
| is_legal                   	| ИП (0) или ЮЛ (1)                      	| нет           	| _пусто или 1 или 0_    	|
| limit                   	| Кол-во элементов в списке для search                      	| нет           	| больше 0    	|

### Заметки

#### Порядок обновления материализованных представлений:

##### Схема accounting_statements:

1. __accounting_statements.balance_mat_view__ - зависит только от таблицы __accounting_statements.баланс__

##### Схема egr:

1. __egr_search__ - зависит от таблиц __егрип__ и __егрюл__
1.2. __tax_authority__ - зависит от мат.вью __egr_search__
1.3. __handbook_tax_authority__ - зависит от мат.вью __tax_authority__
2. __okved__ - зависит от таблиц __егрип__ и __егрюл__
2.1. __handbook_okved__ - зависит от мат.вью __handbook_okved__

Для обновления материализованных представлений:
 ` REFRESH MATERIALIZED VIEW CONCURRENTLY <TABLE_NAME> WITH DATA;`

При первом создании таблиц или после внесения изменений можно применять команду только без __CONCURRENTLY__

### CRON

##### Федералльный реестр туристских объектов

 ` 10 6 * * * $HOME/./opendataaggregator downloader --source=hotels 2>&1 | tee -a $HOME/logs/hotels.log && $HOME/./opendataaggregator parser --parse=hotels 2>&1 | tee -a $HOME/logs/hotels.log`

##### Единый государственный реестр индивидуальных предпринимателей

 ` 0 2 * * * $HOME/./opendataaggregator downloader --source=egrip  2>&1 | tee -a $HOME/logs/egrip .log && $HOME/./opendataaggregator parser --parse=egrip  2>&1 | tee -a $HOME/logs/egrip.log`

##### Единый государственный реестр юридических лиц

 ` 15 2 * * * $HOME/./opendataaggregator downloader --source=egrul  2>&1 | tee -a $HOME/logs/egrul .log && $HOME/./opendataaggregator parser --parse=egrul  2>&1 | tee -a $HOME/logs/egrul.log`

##### Сведения о суммах недоимки и задолженности по пеням и штрафам

 ` 0 3 * * * $HOME/./opendataaggregator downloader --source=debtam  2>&1 | tee -a $HOME/logs/debtam.log && $HOME/./opendataaggregator parser --parse=debtam  2>&1 | tee -a $HOME/logs/debtam.log`

##### Сведения об участии в консолидированной группе налогоплательщиков

 ` 20 3 * * * $HOME/./opendataaggregator downloader --source=kgn  2>&1 | tee -a $HOME/logs/kgn.log && $HOME/./opendataaggregator parser --parse=kgn  2>&1 | tee -a $HOME/logs/kgn.log`

##### Открытый реестр общеизвестных в Российской Федерации товарных знаков

 ` 30 3 * * * $HOME/./opendataaggregator downloader --source=otz  2>&1 | tee -a $HOME/logs/otz.log && $HOME/./opendataaggregator parser --parse=otz  2>&1 | tee -a $HOME/logs/otz.log`

##### Открытый реестр товарных знаков и знаков обслуживания Российской Федерации

 ` 40 3 * * * $HOME/./opendataaggregator downloader --source=tz  2>&1 | tee -a $HOME/logs/tz.log && $HOME/./opendataaggregator parser --parse=tz  2>&1 | tee -a $HOME/logs/tz.log`

##### Сведения об уплаченных организацией суммах налогов и сборов

 ` 50 3 * * * $HOME/./opendataaggregator downloader --source=paytax 2>&1 | tee -a $HOME/logs/paytax.log && $HOME/./opendataaggregator parser --parse=paytax 2>&1 | tee -a $HOME/logs/paytax.log`

##### Реестр дисквалифицированных лиц

 ` 0 4 * * * $HOME/./opendataaggregator downloader --source=registerdisqualified 2>&1 | tee -a $HOME/logs/registerdisqualified.log && $HOME/./opendataaggregator parser --parse=registerdisqualified 2>&1 | tee -a $HOME/logs/registerdisqualified.log`

##### Сведения о специальных налоговых режимах, применяемых налогоплательщиками

 ` 30 2 * * * $HOME/./opendataaggregator downloader --source=snr 2>&1 | tee -a $HOME/logs/snr.log && $HOME/./opendataaggregator parser --parse=snr 2>&1 | tee -a $HOME/logs/snr.log`

##### Сведения о налоговых правонарушениях и мерах ответственности за их совершение

 ` 40 2 * * * $HOME/./opendataaggregator downloader --source=taxoffence 2>&1 | tee -a $HOME/logs/taxoffence.log && $HOME/./opendataaggregator parser --parse=taxoffence 2>&1 | tee -a $HOME/logs/taxoffence.log`

##### Сведения о среднесписочной численности работников организации

 ` 50 2 * * * $HOME/./opendataaggregator downloader --source=sshr 2>&1 | tee -a $HOME/logs/sshr.log && $HOME/./opendataaggregator parser --parse=sshr 2>&1 | tee -a $HOME/logs/sshr.log`

##### Единый реестр субъектов малого и среднего предпринимательства

 ` 20 4 * * * $HOME/./opendataaggregator downloader --source=rsmp 2>&1 | tee -a $HOME/logs/rsmp.log && $HOME/./opendataaggregator parser --parse=rsmp 2>&1 | tee -a $HOME/logs/rsmp.log`

##### Информация о привлечении участника закупки к административной ответственности по ст. 19.28 КоАП

 ` 0 5 * * * $HOME/./opendataaggregator downloader --source=zakupki 2>&1 | tee -a $HOME/logs/zakupki.log && $HOME/./opendataaggregator parser --parse=zakupki 2>&1 | tee -a $HOME/logs/zakupki.log`

##### Сведения из Реестра сертификатов соответствия

 ` 40 4 * * * $HOME/./opendataaggregator downloader --source=rss 2>&1 | tee -a $HOME/logs/rss.log && $HOME/./opendataaggregator parser --parse=rss 2>&1 | tee -a $HOME/logs/rss.log`

##### Исполнительные производства в отношении юридических лиц

 ` 40 6 * * * $HOME/./opendataaggregator downloader --source=iplegallist 2>&1 | tee -a $HOME/logs/iplegallist.log && $HOME/./opendataaggregator parser --parse=iplegallist 2>&1 | tee -a $HOME/logs/iplegallist.log`

##### Оконченные производства в отношении юридических лиц

 ` 15 6 * * * $HOME/./opendataaggregator downloader --source=iplegallistcomplete 2>&1 | tee -a $HOME/logs/iplegallistcomplete.log && $HOME/./opendataaggregator parser --parse=iplegallistcomplete 2>&1 | tee -a $HOME/logs/iplegallistcomplete.log`

##### Сведения из Реестра деклараций о соответствии

 ` 15 7 * * * $HOME/./opendataaggregator downloader --source=rds 2>&1 | tee -a $HOME/logs/rds.log && $HOME/./opendataaggregator parser --parse=rds 2>&1 | tee -a $HOME/logs/rds.log`
