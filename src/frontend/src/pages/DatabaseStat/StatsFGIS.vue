<template>
    <div>
        <h3>ФГИС</h3>
        <table class="table table-bordered">
            <tbody v-if="!loading">
                <tr>
                    <td>Внеплановые проверки:</td>
                    <td>{{ stats_fgis?.unscheduled_count?.toLocaleString() }}</td>
                </tr>
                <tr>
                    <td>Плановые проверки:</td>
                    <td>{{ stats_fgis?.scheduled_count?.toLocaleString() }}</td>
                </tr>
            </tbody>
            <SpinnerLoader v-else />
        </table>
    </div>
</template>

<script setup>
import SpinnerLoader from "@/components/SpinnerLoader.vue";
import { ref, onMounted } from "vue";
const stats_fgis = ref({});
const loading = ref(false);
onMounted(async () => {
    loading.value = true;
    await fetch("/api/db_stats/stats_fgis")
        .then(r => r.json())
        .then((b) => {
            stats_fgis.value = b;
        })
        .finally(() => loading.value = false);
});
</script>
