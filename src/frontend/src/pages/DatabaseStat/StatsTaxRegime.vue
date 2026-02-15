<template>
    <div>
        <h3>Налоговый режим</h3>
        <table class="table table-bordered">
            <tbody v-if="!loading">
                <tr>
                    <td>Всего записей:</td>
                    <td>{{ stats_tax_regime?.count?.toLocaleString() }}</td>
                </tr>
                <tr>
                    <td>УСН:</td>
                    <td>{{ stats_tax_regime?.усн?.toLocaleString() }}</td>
                </tr>
                <tr>
                    <td>ЕСХН:</td>
                    <td>{{ stats_tax_regime?.есхн?.toLocaleString() }}</td>
                </tr>
                <tr>
                    <td>ЕНВД:</td>
                    <td>{{ stats_tax_regime?.енвд?.toLocaleString() }}</td>
                </tr>
                <tr>
                    <td>СРП:</td>
                    <td>{{ stats_tax_regime?.срп?.toLocaleString() }}</td>
                </tr>
            </tbody>
            <SpinnerLoader v-else />
        </table>
    </div>
</template>

<script setup>
import SpinnerLoader from "@/components/SpinnerLoader.vue";
import { ref, onMounted } from "vue";
const stats_tax_regime = ref({});
const loading = ref(false);
onMounted(async () => {
    loading.value = true;
    await fetch("/api/db_stats/stats_tax_regime")
        .then(r => r.json())
        .then((b) => {
            stats_tax_regime.value = b;
        })
        .finally(() => loading.value = false);
});
</script>
