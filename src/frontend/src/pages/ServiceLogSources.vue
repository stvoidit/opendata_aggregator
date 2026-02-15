<template>
    <div class="container-fluid">
        <div class="row mt-5">
            <div class="col">
                <table class="table table-sm table-bordered">
                    <thead>
                        <tr>
                            <th
                                v-for="h in fields"
                                :key="h">
                                {{ h }}
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr
                            v-for="log in computedLogs"
                            ref="tableBody"
                            :key="log.id">
                            <td>
                                {{ log.source_type }}
                            </td>
                            <td>
                                <a
                                    :href="log.source_link"
                                    target="_blank">{{ log.filename }}</a>
                            </td>
                            <td>
                                {{ log.sha256sum }}
                            </td>
                            <td class="text-center">
                                <input
                                    class="form-check-input"
                                    type="checkbox"
                                    disabled
                                    :checked="log.downloaded">
                            </td>
                            <td class="text-center">
                                <input
                                    class="form-check-input"
                                    type="checkbox"
                                    disabled
                                    :checked="log.uploaded">
                            </td>
                            <td>
                                {{ toDatetime(log.task_datetime) }}
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, computed } from "vue";

interface logData {
    downloaded: boolean;
    filename: string;
    id: string;
    sha256sum: string;
    source_link: string;
    source_type: string;
    task_datetime: string;
    uploaded: boolean;
}

const tableBody = ref<HTMLElement[]>([]);
const logs = ref<logData[]>([]);
const offset = ref(40);
const computedLogs = computed(() => {
    return logs.value.slice(0, offset.value);
});

/**
 * Создание обсервера, который увеличивает список отображения для логов
 *
 * @param lastEl - последний элемент на который вешается обсервер
 */
function createObserver(lastEl: HTMLElement) {
    const obs = new IntersectionObserver((entries, observer) => {
        for (const entry of entries) {
            if (entry && entry.isIntersecting) {
                offset.value += 40;
                observer.disconnect(); /** обязательно отключаем прошлый обсервер, чтобы перевесить на новый последний элемент */
                if (offset.value > logs.value.length) {
                    /** перестаем создавать новые обсерверы, если достигнут конец списка элементов для отображения */
                    return;
                }
                createObserver(tableBody.value[tableBody.value.length - 1]);
                return;
            }
        }
    }, {
        root: null,
        threshold: 0.25
    });
    obs.observe(lastEl);
}

const fetchData = async () => {
    const response = await fetch("/api/service_log_sources");
    logs.value = await response.json();
    /** создаем стартовый обсервер */
    if (logs.value.length) await nextTick(() => createObserver(tableBody.value[tableBody.value.length - 1]));
};
onMounted(fetchData);
const fields = [
    "source_type",
    "filename",
    "sha256sum",
    "downloaded",
    "uploaded",
    "task_datetime"
];
const toDatetime = (s: string) => (new Date(s)).toLocaleString();
</script>
