<template>
    <div>
        <h4>Налоговые правонарушения и штрафы</h4>
        <table class="table table-bordered">
            <tbody v-if="!loading">
                <tr>
                    <td>Всего записей</td>
                    <td>{{ stats_tax_offenses?.total_count?.toLocaleString() }}</td>
                </tr>
                <tr>
                    <td>Сумма всего ₽</td>
                    <td>{{ stats_tax_offenses?.total_sum?.toLocaleString() }}</td>
                </tr>
                <tr>
                    <td>Количество по годам</td>
                    <td>
                        <table class="table table-sm table-bordered">
                            <thead>
                                <tr>
                                    <th>год</th>
                                    <th>количество</th>
                                    <th>сумма ₽</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr
                                    v-for="v in stats_tax_offenses?.sums_by_years"
                                    :key="v.year">
                                    <td>{{ v.year }}</td>
                                    <td>{{ v.count.toLocaleString() }}</td>
                                    <td>{{ v.sum.toLocaleString() }}</td>
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

<script setup>
import SpinnerLoader from "@/components/SpinnerLoader.vue";
import { ref, onMounted } from "vue";
const stats_tax_offenses = ref({});
const loading = ref(false);
onMounted(async () => {
    loading.value = true;
    await fetch("/api/db_stats/stats_tax_offenses")
        .then(r => r.json())
        .then((b) => {
            stats_tax_offenses.value = b;
        })
        .finally(() => loading.value = false);
});
</script>
