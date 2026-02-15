<template>
    <div>
        <h4>Гостиницы</h4>
        <table class="table table-bordered">
            <tbody v-if="!loading">
                <tr>
                    <td>Всего записей</td>
                    <td>{{ stats_hotels?.total_count?.toLocaleString() }}</td>
                </tr>
                <tr>
                    <td>Количество по типам</td>
                    <td>
                        <table class="table table-sm table-bordered">
                            <thead>
                                <tr>
                                    <th width="80%">
                                        тип
                                    </th>
                                    <th>
                                        количество
                                    </th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr
                                    v-for="v in stats_hotels?.count_hotels_types"
                                    :key="v.type">
                                    <td>{{ v.type }}</td>
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
const stats_hotels = ref({});
const loading = ref(false);
onMounted(async () => {
    loading.value = true;
    await fetch("/api/db_stats/stats_hotels")
        .then(r => r.json())
        .then((b) => {
            stats_hotels.value = b;
        })
        .finally(() => loading.value = false);
});
</script>
