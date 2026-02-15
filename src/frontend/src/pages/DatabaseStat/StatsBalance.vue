<template>
    <div>
        <h4>Бухгалтерская отчетность</h4>
        <table class="table table-bordered">
            <tbody v-if="!loading">
                <tr>
                    <td>Последняя дата документа</td>
                    <td>{{ toDateLocal(stat_balance?.last_doc_date) }}</td>
                </tr>
                <tr>
                    <td>Всего записей</td>
                    <td>{{ stat_balance?.total_count?.toLocaleString() }}</td>
                </tr>
                <tr>
                    <td>Количество по годам</td>
                    <td>
                        <table class="table table-sm table-bordered">
                            <thead>
                                <tr>
                                    <th>год</th>
                                    <th>количество</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr
                                    v-for="v in stat_balance?.count_years"
                                    :key="v.year">
                                    <td>{{ v.year }}</td>
                                    <td>{{ v.count.toLocaleString() }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </td>
                </tr>
            </tbody>
            <SpinnerLoader v-else />
        </table>
    </div>
</template>

<script setup lang="ts">
import SpinnerLoader from "@/components/SpinnerLoader.vue";
import { ref, onMounted } from "vue";
const stat_balance = ref<any>({});
const loading = ref(false);
onMounted(async () => {
    loading.value = true;
    await fetch("/api/db_stats/stats_balance")
        .then(r => r.json())
        .then((b) => {
            stat_balance.value = b;
        })
        .finally(() => loading.value = false);
});
const toDateLocal = (s: string | number | Date) => {
    const date = new Date(s);
    const day = date.getDate() < 10 ? `0${date.getDate()}` : date.getDate().toString();
    const month = date.getMonth() + 1 < 10 ? `0${date.getMonth() + 1}` : (date.getMonth() + 1).toString();
    const year = date.getFullYear();
    return `${day}.${month}.${year}`;
};
</script>
