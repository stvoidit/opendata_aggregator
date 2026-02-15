<template>
    <div>
        <h4>Последнее обновление источников</h4>
        <table class="table table-bordered">
            <thead>
                <tr>
                    <th>Название источника</th>
                    <th>Последнее обновление</th>
                </tr>
            </thead>
            <tbody v-if="!loading">
                <tr
                    v-for="lus in stats?.last_updates_sources"
                    :id="lus.source_type"
                    :key="lus.source_type">
                    <td>
                        <a
                            :href="stats.sources[lus.source_type]"
                            target="_blank">{{ getLabelHR(lus.source_type) }}
                        </a>
                    </td>
                    <td>{{ formatDate(lus.datetime) }}</td>
                </tr>
            </tbody>
            <SpinnerLoader v-else />
        </table>
    </div>
</template>

<script setup lang="ts">
import SpinnerLoader from "@/components/SpinnerLoader.vue";
import { ref, onMounted } from "vue";
const stats = ref<any>({});
const loading = ref(false);
onMounted(async () => {
    loading.value = true;
    await Promise.all([
        fetch("/api/db_stats/sources")
            .then(r => r.json())
            .then((b) => {
                stats.value.sources = b;
            }),
        fetch("/api/db_stats/last_updates_sources")
            .then(r => r.json())
            .then((b) => {
                stats.value.last_updates_sources = b;
            })
    ]).finally(() => loading.value = false);
});

const formatDate = (s: string | number | Date) => (new Date(s)).toLocaleString();
const ruLabes = new Map(Object.entries({
    balance: "Отчет об бухгалтерском балансе",
    iplegallist: "Исполнительные производства в отношении юридических лиц",
    iplegallistcomplete: "Оконченные производства в отношении юридических лиц",
    snr: "Сведения о специальных налоговых режимах, применяемых налогоплательщиками",
    taxoffence: "Сведения о налоговых правонарушениях и мерах ответственности за их совершение",
    okved2: "Общероссийский классификатор видов экономической деятельности (ОКВЭД2)",
    okato: "Общероссийский классификатор объектов административно-территориального деления (ОКАТО)",
    oktmo: "Общероссийский классификатор территорий муниципальных образований (ОКТМО)",
    otz: "Открытый реестр общеизвестных в Российской Федерации товарных знаков",
    tz: "Открытый реестр товарных знаков и знаков обслуживания Российской Федерации",
    zakupki: "Информация о привлечении участника закупки к административной ответственности по ст. 19.28 КоАП",
    registerdisqualified: "Реестр дисквалифицированных лиц",
    rsmp: "Единый реестр субъектов малого и среднего предпринимательства",
    rss: "Сведения из Реестра сертификатов соответствия",
    sshr: "Сведения о среднесписочной численности работников организации",
    debtam: "Сведения о суммах недоимки и задолженности по пеням и штрафам",
    paytax: "Сведения об уплаченных организацией суммах налогов и сборов",
    kgn: "Сведения об участии в консолидированной группе налогоплательщиков",
    hotels: "Федералльный реестр туристских объектов",
    egrul: "Единый государственный реестр юридических лиц",
    egrip: "Единый государственный реестр индивидуальных предпринимателей",
    rds: "Сведения из Реестра деклараций о соответствии",
    fgis_unscheduled: "ФГИС внеплановые проверки",
    fgis_plan: "ФГИС плановые проверки"
}));
const getLabelHR = (key: string) => ruLabes.get(key);
</script>
