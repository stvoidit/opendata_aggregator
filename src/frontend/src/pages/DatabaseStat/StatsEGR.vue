<template>
    <div>
        <h4>ЕГРИП и ЕГРЮЛ</h4>
        <table class="table table-bordered">
            <tbody v-if="!loading">
                <tr>
                    <td>Дата последней выписки</td>
                    <td>{{ toDateLocal(stat_egr?.last_date_discharge) }}</td>
                </tr>
                <tr>
                    <td>Всего записей</td>
                    <td>{{ stat_egr?.total_count?.toLocaleString() }}</td>
                </tr>
                <tr>
                    <td>Количество ЕГРЮЛ</td>
                    <td>{{ stat_egr?.egrul?.toLocaleString() }}</td>
                </tr>
                <tr>
                    <td>Количество ЕГРИП</td>
                    <td>{{ stat_egr?.egrip?.toLocaleString() }}</td>
                </tr>
            </tbody>
            <SpinnerLoader v-else />
        </table>
    </div>
</template>

<script setup lang="ts">
import SpinnerLoader from "@/components/SpinnerLoader.vue";
import { ref, onMounted } from "vue";
const stat_egr = ref<any>({});
const loading = ref(false);
onMounted(async () => {
    loading.value = true;
    await fetch("/api/db_stats/stats_egr")
        .then(r => r.json())
        .then((b) => {
            stat_egr.value = b;
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
