<template>
    <div>
        <h4>Реестр товарных знаков</h4>
        <table class="table table-bordered">
            <tbody v-if="!loading">
                <tr>
                    <td>Открытый реестр товарных знаков:</td>
                    <td>{{ stats_register_of_trademarks?.open_registry_count?.toLocaleString() }}</td>
                </tr>
                <tr>
                    <td>Общеизвестный реестр товарных знаков:</td>
                    <td>{{ stats_register_of_trademarks?.well_known_registry_count?.toLocaleString() }}</td>
                </tr>
            </tbody>
            <SpinnerLoader v-else />
        </table>
    </div>
</template>

<script setup>
import SpinnerLoader from "@/components/SpinnerLoader.vue";
import { ref, onMounted } from "vue";
const stats_register_of_trademarks = ref({});
const loading = ref(false);
onMounted(async () => {
    loading.value = true;
    await fetch("/api/db_stats/stats_register_of_trademarks")
        .then(r => r.json())
        .then((b) => {
            stats_register_of_trademarks.value = b;
        })
        .finally(() => loading.value = false);
});
</script>
