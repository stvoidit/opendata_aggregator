<template>
    <div class="container-fluid mb-3 mt-1">
        <div
            v-for="r in routes"
            :key="r.href"
            class="row my-2">
            <div class="col">
                <div class="card shadow-sm">
                    <div class="card-body">
                        <h6 class="py-3 px-2">
                            <b>{{ r.href }}</b> - {{ r.desc }}
                        </h6>
                        <p
                            v-if="r.params.length"
                            class="card-text">
                            <table class="table table-sm table-bordered mx-1">
                                <thead>
                                    <tr>
                                        <th
                                            scope="col"
                                            width="10%">
                                            параметр
                                        </th>
                                        <th
                                            scope="col"
                                            width="10%">
                                            множественный
                                        </th>
                                        <th
                                            scope="col"
                                            width="15%">
                                            пример
                                        </th>
                                        <th scope="col">
                                            описание
                                        </th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr
                                        v-for="p in r.params"
                                        :key="p.name">
                                        <td>{{ p.name }}</td>
                                        <td>{{ p.multi }}</td>
                                        <td>{{ p.example }}</td>
                                        <td>{{ p.desc }}</td>
                                    </tr>
                                </tbody>
                            </table>
                        </p>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">

const routes = [
    {
        href: "/api/search",
        desc: "Поиск по параметрам (возвращает список)",
        params: [
            {
                name: "year",
                desc: "Год бухгалтерского отчета",
                multi: true,
                example: "2022"
            },
            {
                name: "income",
                desc: "Доходы",
                multi: false,
                example: "1200"
            },
            {
                name: "expenses",
                desc: "Расходы",
                multi: false,
                example: "1200"
            },
            {
                name: "income_tax",
                desc: "Налог на прибыль",
                multi: false,
                example: "1200"
            },
            {
                name: "tax_usn",
                desc: "Налог по УСН",
                multi: false,
                example: "1200"
            },
            {
                name: "intangible_assets",
                desc: "Нематериальные активы",
                multi: false,
                example: "1200"
            },
            {
                name: "basic_assets",
                desc: "Основные средства",
                multi: false,
                example: "1200"
            },
            {
                name: "other_non_current_assets",
                desc: "Прочие внеоборотные активы",
                multi: false,
                example: "1200"
            },
            {
                name: "non_current_assets",
                desc: "Внеоборотные активы",
                multi: false,
                example: "1200"
            },
            {
                name: "stocks",
                desc: "Запасы",
                multi: false,
                example: "1200"
            },
            {
                name: "net_aassets",
                desc: "Чистые активы",
                multi: false,
                example: "1200"
            },
            {
                name: "accounts_receivable",
                desc: "Дебиторская задолженность",
                multi: false,
                example: "1200"
            },
            {
                name: "cash_and_equivalents",
                desc: "Денежные средства и денежные эквиваленты",
                multi: false,
                example: "1200"
            },
            {
                name: "current_assets",
                desc: "Оборотные активы",
                multi: false,
                example: "1200"
            },
            {
                name: "total_assets",
                desc: "Активы всего",
                multi: false,
                example: "1200"
            },
            {
                name: "capital_and_reserves",
                desc: "Капитал и резервы ",
                multi: false,
                example: "1200"
            },
            {
                name: "borrowed_funds_long_term",
                desc: "Заёмные средства (долгосрочные)",
                multi: false,
                example: "1200"
            },
            {
                name: "borrowed_funds_short_term",
                desc: "Заёмные средства (краткосрочные)",
                multi: false,
                example: "1200"
            },
            {
                name: "accounts_payable",
                desc: "Кредиторская задолженность",
                multi: false,
                example: "1200"
            },
            {
                name: "other_short_term_liabilities",
                desc: "Прочие краткосрочные обязательства",
                multi: false,
                example: "1200"
            },
            {
                name: "total_liabilities",
                desc: "Пассивы всего",
                multi: false,
                example: "1200"
            },
            {
                name: "revenue",
                desc: "Выручка",
                multi: false,
                example: "1200"
            },
            {
                name: "cost_of_sales",
                desc: "Себестоимость продаж",
                multi: false,
                example: "1200"
            },
            {
                name: "gross_profit",
                desc: "Валовая прибыль",
                multi: false,
                example: "1200"
            },
            {
                name: "commercial_expenses",
                desc: "Коммерческие расходы",
                multi: false,
                example: "1200"
            },
            {
                name: "management_expenses",
                desc: "Управленческие расходы",
                multi: false,
                example: "1200"
            },
            {
                name: "profit_from_sale",
                desc: "Прибыль (убыток) от продажи",
                multi: false,
                example: "1200"
            },
            {
                name: "interest_payable",
                desc: "Проценты к уплате",
                multi: false,
                example: "1200"
            },
            {
                name: "other_income",
                desc: "Прочие доходы",
                multi: false,
                example: "1200"
            },
            {
                name: "other_expenses",
                desc: "Прочие расходы",
                multi: false,
                example: "1200"
            },
            {
                name: "profit_before_taxation",
                desc: "Прибыль (убыток) до налогообложения",
                multi: false,
                example: "1200"
            },
            {
                name: "сurrent_income_tax",
                desc: "Текущий налог на прибыль",
                multi: false,
                example: "1200"
            },
            {
                name: "net_profit",
                desc: "Чистая прибыль (убыток)",
                multi: false,
                example: "1200"
            },
            {
                name: "okved",
                desc: "Код ОКВЭД",
                multi: true,
                example: "?okved=28.30.51&okved=10.61.3"
            },
            {
                name: "ta",
                desc: "Код Налогового органа",
                multi: false,
                example: "1200"
            },
            {
                name: "ta",
                desc: "Код Налогового органа",
                multi: false,
                example: "7715"
            },
            {
                name: "date_registration_from",
                desc: "Дата регистрации С",
                multi: false,
                example: "2000-12-31"
            },
            {
                name: "date_registration_to",
                desc: "Дата регистрации ПО",
                multi: false,
                example: "2000-12-31"
            },
            {
                name: "date_liquidation_from",
                desc: "Дата ликвидации С",
                multi: false,
                example: "2000-12-31"
            },
            {
                name: "date_liquidation_to",
                desc: "Дата ликвидации ПО",
                multi: false,
                example: "2000-12-31"
            },
            {
                name: "kpp",
                desc: "Код КПП",
                multi: false,
                example: "370201001"
            },
            {
                name: "usn",
                desc: "Налоговый режим УСН",
                multi: false,
                example: "пусто или 1"
            },
            {
                name: "eshn",
                desc: "Налоговый режим ЕСХН",
                multi: false,
                example: "пусто или 1"
            },
            {
                name: "envd",
                desc: "Налоговый режим ЕНВД",
                multi: false,
                example: "пусто или 1"
            },
            {
                name: "srp",
                desc: "Налоговый режим СРП",
                multi: false,
                example: "пусто или 1"
            },
            {
                name: "status",
                desc: "Сатутс (a - действующие (active), d - недействующие, l - ликвидация, b - банкротство, r - реорганизация )",
                multi: true,
                example: "d"
            },
            {
                name: "q",
                desc: "Текстовый поиск по названию компании или ФИО + ИНН + ОГРН",
                multi: false,
                example: "магнит"
            },
            {
                name: "is_legal",
                desc: "ИП (0) или ЮЛ (1)",
                multi: false,
                example: "пусто или 1 или 0"
            },
            {
                name: "limit",
                desc: "Кол-во элементов в списке для search, больше 0, по умолчанию 20",
                multi: false,
                example: "100"
            }
        ]
    },
    {
        href: "/api/count",
        desc: "Кол-во строк в результате поиска для /api/search - принимает те же параметры что и /api/search",
        params: []
    },
    {
        href: "/api/info/{inn:[0-9]+}/{ogrn:[0-9]+}",
        desc: "Полная карточка информации по ЮЛ или ИП",
        params: []
    },
    {
        href: "/api/handbook_tax_authority",
        desc: "Возвращает полный справочник Налоговых Органов (для парамерта ta)",
        params: []
    },
    {
        href: "/api/search_tax_authority",
        desc: "Поиск кода Налоговых Органов - возвращает список",
        params: [
            {
                name: "q",
                desc: "Поиск по названию, минимум 3 символа",
                multi: false,
                example: "ПО РЕСПУБЛИКЕ АДЫГЕЯ"
            }
        ]
    },
    {
        href: "/api/handbook_okved",
        desc: "Возвращает полный справочник ОКВЭД (для парамерта okved)",
        params: []
    },
    {
        href: "/api/search_okved",
        desc: "Поиск кодов ОКВЭД - возвращает список",
        params: [
            {
                name: "q",
                desc: "Поиск по названию, минимум 3 символа",
                multi: false,
                example: "выращивание семян"
            }
        ]
    },
    {
        href: "/api/handbook_categories_ip",
        desc: "Справочник категорий ИП",
        params: []
    },
    {
        href: "/api/legal_statuses",
        desc: "Справочник расшифровки статусов ЮЛ",
        params: []
    }
];

</script>
