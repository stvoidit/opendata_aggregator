<template>
    <div>
        <h3>СМП</h3>
        <table class="table table-bordered">
            <tbody v-if="!loading">
                <tr>
                    <td>Всего записей:</td>
                    <td>{{ stats_smp.total_count?.toLocaleString() }}</td>
                </tr>
            </tbody>
            <SpinnerLoader v-else />
        </table>
    </div>
</template>

<script setup>
import SpinnerLoader from "@/components/SpinnerLoader.vue";
import { ref, onMounted } from "vue";
const stats_smp = ref({});
const loading = ref(false);
onMounted(async () => {
    loading.value = true;
    await fetch("/api/db_stats/stats_smp")
        .then(r => r.json())
        .then((b) => {
            stats_smp.value = b;
        })
        .finally(() => loading.value = false);
});
</script>
