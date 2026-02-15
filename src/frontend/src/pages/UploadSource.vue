<template>
    <div class="container-fluid">
        <div class="row mb-1 mt-5 justify-content-center">
            <div class="col-4">
                <div
                    class="alert alert-danger"
                    role="alert">
                    в разработке
                </div>
                <div class="card">
                    <div class="card-body">
                        <form @submit.prevent="onSubmit">
                            <div class="mb-3">
                                <select
                                    v-model="sourceType"
                                    class="form-select"
                                    name="sourceType">
                                    <option value="balance">
                                        Бухгалтерская отчетность
                                    </option>
                                </select>
                            </div>
                            <div class="mb-3">
                                <input
                                    class="form-control"
                                    type="file"
                                    name="sourceFile"
                                    @change="onChangeFile">
                            </div>
                            <div
                                v-if="showProgress"
                                class="mb-3">
                                <div
                                    class="progress"
                                    role="progressbar"
                                    aria-label="Basic example"
                                    :aria-valuenow="progressValue"
                                    aria-valuemin="0"
                                    aria-valuemax="100">
                                    <div
                                        class="progress-bar progress-bar-striped progress-bar-animated"
                                        :style="{ width: `${progressValue}%` }">
                                        {{ `${progressValue}%` }}
                                    </div>
                                </div>
                            </div>
                            <button
                                type="submit"
                                class="btn btn-sm btn-primary">
                                Загрузить
                            </button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

const sourceFile = ref<File | null>(null);
const sourceType = ref<string | null>(null);
const progressValue = ref(0);
const showProgress = ref(false);

function onChangeFile(e: Event) {
    const target = e.target as HTMLInputElement;
    if (target && target.files) {
        sourceFile.value = target.files[0];
    }
}

function onSubmit(e: Event | SubmitEvent) {
    let fileSize: number;
    const formData = new FormData();
    const formElement = e.target as HTMLFormElement;
    if (sourceFile.value && sourceType.value) {
        fileSize = sourceFile.value.size;
        formData.append("sourceFile", sourceFile.value);
        formData.append("sourceType", sourceType.value);
    } else {
        alert("Не все поля формы заполнены");
        return;
    }
    const xhr = new XMLHttpRequest();
    xhr.withCredentials = true;
    xhr.open("POST", "/api/upload");
    xhr.upload.onprogress = function (event) {
        progressValue.value = Math.round(event.loaded / (fileSize / 100));
    };
    xhr.onloadstart = function () {
        showProgress.value = true;
    };
    xhr.onloadend = function () {
        formElement.sourceFile.value = null;
        sourceFile.value = null;
        sourceType.value = null;
        showProgress.value = false;
        progressValue.value = 0;
        if (xhr.status == 201) {
            alert("Загружено");
        } else {
            alert(`Ошибка: ${this.status}`);
        }
    };
    xhr.send(formData);
}
</script>
