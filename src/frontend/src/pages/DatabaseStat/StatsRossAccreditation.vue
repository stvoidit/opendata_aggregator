<template>
    <div>
        <h4>Россакредитация</h4>
        <table class="table table-bordered">
            <tbody v-if="!loading">
                <tr>
                    <td>Всего записей</td>
                    <td>{{ stats_ross_accreditation?.total_count?.toLocaleString() }}</td>
                </tr>
                <tr>
                    <td>Количество по статусам</td>
                    <td>
                        <table class="table table-sm table-bordered">
                            <thead>
                                <tr>
                                    <th>статус сертификата</th>
                                    <th>количество</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr
                                    v-for="v in stats_ross_accreditation?.ross_accreditation_statuses"
                                    :key="v.cert_status">
                                    <td>{{ v.cert_status }}</td>
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

<script setup>
import SpinnerLoader from "@/components/SpinnerLoader.vue";
import { ref, onMounted } from "vue";
const stats_ross_accreditation = ref({});
const loading = ref(false);
onMounted(async () => {
    loading.value = true;
    await fetch("/api/db_stats/stats_ross_accreditation")
        .then(r => r.json())
        .then((b) => {
            stats_ross_accreditation.value = b;
        })
        .finally(() => loading.value = false);
});
</script>
