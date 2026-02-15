<template>
    <div>
        <h3>Сумма недоимки</h3>
        <table class="table table-bordered">
            <tbody v-if="!loading">
                <tr>
                    <td>Всего записей:</td>
                    <td>{{ stats_debtam?.total_count?.toLocaleString() }}</td>
                </tr>
                <tr>
                    <td>Общая сумма:</td>
                    <td>{{ stats_debtam?.total_sum?.toLocaleString() }}</td>
                </tr>
            </tbody>
            <SpinnerLoader v-else />
        </table>
    </div>
</template>

<script setup>
import SpinnerLoader from "@/components/SpinnerLoader.vue";
import { ref, onMounted } from "vue";
const stats_debtam = ref({});
const loading = ref(false);
onMounted(async () => {
    loading.value = true;
    await fetch("/api/db_stats/stats_debtam")
        .then(r => r.json())
        .then((b) => {
            stats_debtam.value = b;
        })
        .finally(() => loading.value = false);
});
</script>
