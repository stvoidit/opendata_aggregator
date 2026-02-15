<template>
    <div>
        <h3>ФССП</h3>
        <table class="table table-bordered">
            <tbody v-if="!loading">
                <tr>
                    <td>Исполнительные производства в отношении юридических лиц:</td>
                    <td>{{ stats_fssp?.ip_legal_list?.toLocaleString() }}</td>
                </tr>
                <tr>
                    <td>Оконченные производства в отношении юридических лиц:</td>
                    <td>{{ stats_fssp?.ip_legal_list_complite?.toLocaleString() }}</td>
                </tr>
            </tbody>
            <SpinnerLoader v-else />
        </table>
    </div>
</template>

<script setup>
import SpinnerLoader from "@/components/SpinnerLoader.vue";
import { ref, onMounted } from "vue";
const stats_fssp = ref({});
const loading = ref(false);
onMounted(async () => {
    loading.value = true;
    await fetch("/api/db_stats/stats_fssp")
        .then(r => r.json())
        .then((b) => {
            stats_fssp.value = b;
        })
        .finally(() => loading.value = false);
});
</script>
