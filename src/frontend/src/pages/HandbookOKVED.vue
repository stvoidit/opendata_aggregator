<template>
    <div
        v-if="loaded"
        class="container-fluid">
        <div class="row">
            <div class="col p-3">
                <div class="input-group mb-2">
                    <input
                        v-model="searchStr"
                        class="form-control"
                        placeholder="поиск по названию"
                        type="text">
                    <button
                        class="btn btn-warning"
                        type="button"
                        @click="() => searchStr = ''">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="16"
                            height="16"
                            fill="currentColor"
                            class="bi bi-x"
                            viewBox="0 0 16 16">
                            <path d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708z" />
                        </svg>
                    </button>
                </div>
                <table class="table table-sm table-bordered">
                    <thead>
                        <tr>
                            <th width="7%">
                                #
                            </th>
                            <th width="8%">
                                код
                            </th>
                            <th width="80%">
                                название
                            </th>
                            <th witdh="5%">
                                версия
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr
                            v-for="(code, index) in computedData"
                            :key="`${code.code}-${code.vers}`">
                            <td>
                                {{ index+1 }}
                            </td>
                            <td class="d-flex justify-content-between">
                                <span :style="calcMargin(code.code)">{{ code.code }}</span>
                                <button
                                    class="btn btn-sm btn-link copy-btn"
                                    :style="{cursor: 'copy'}"
                                    @click="()=>copyTo(code.code)">
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        width="16"
                                        height="16"
                                        fill="currentColor"
                                        class="bi bi-subtract"
                                        viewBox="0 0 16 16">
                                        <path d="M0 2a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v2h2a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2v-2H2a2 2 0 0 1-2-2V2zm2-1a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1H2z" />
                                    </svg>
                                </button>
                            </td>
                            <td>
                                <span>{{ code.title }}</span>
                            </td>
                            <td>
                                {{ code.vers }}
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
    <SpinnerLoader v-else />
    <div class="toast-container position-fixed top-0 end-0 p-5">
        <div
            id="liveToast"
            ref="liveToast"
            class="toast align-items-center"
            role="alert"
            aria-live="assertive"
            aria-atomic="true">
            <div class="toast-body">
                copied!
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useClipboard } from "@vueuse/core";
import { Toast } from "bootstrap";
import SpinnerLoader from "@/components/SpinnerLoader.vue";
interface okvedCode {
    code: string;
    title: string;
    vers?: string;
}
const loaded = ref(false);
const liveToast = ref<HTMLElement>();
const searchStr = ref("");
const codes = ref<okvedCode[]>([]);
const fetchData = async () => {
    loaded.value = false;
    const response = await fetch("/api/handbook_okved");
    codes.value = await response.json();
    loaded.value = true;
};
onMounted(fetchData);
const calcMargin = (str: string) => {
    const count = [ ...str ].reduce((prev, cur) => cur === "." ? prev + 1 : prev, 0);
    return {
        "margin-left": `${count}rem`
    };
};
const computedData = computed(() => {
    if (!searchStr.value.length) {
        return codes.value;
    }
    const search = searchStr.value.toLocaleLowerCase();
    return codes.value.filter(v => v.title.toLowerCase().includes(search));
});
const copyTo = async (str: string) => {
    await useClipboard({ legacy: true }).copy(str);
    if (liveToast.value) {
        const toast = new Toast(liveToast.value, { delay: 750 });
        toast.show();
    }
};
</script>
<style scoped>
.copy-btn:hover {
    transform: scale(1.25);
}
</style>
