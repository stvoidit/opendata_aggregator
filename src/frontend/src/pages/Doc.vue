<template>
    <div
        v-if="loaded"
        class="container-fluid mb-3 mt-1">
        <div class="row">
            <div class="col text-start">
                <a
                    class="btn btn-primary"
                    href="/api/json_schema"
                    download="TotalResult.json">Скачать JSONSchema</a>
            </div>
            <div class="col text-end">
                <span><b>версия:</b> {{ schema.$schema }}</span>
            </div>
        </div>
    </div>
    <div
        v-if="loaded"
        class="container-fluid">
        <div
            v-for="value, key in schema.$defs"
            :key="key"
            class="row mx-1 my-2">
            <h6 :id="key.toString()">
                {{ key }}
            </h6>
            <table class="table table-bordered table-sm shadow-sm">
                <thead>
                    <tr>
                        <th
                            scope="col"
                            width="25%">
                            Имя поля
                        </th>
                        <th
                            scope="col"
                            width="25%">
                            Тип данных
                        </th>
                        <th
                            scope="col"
                            width="10%">
                            Обязательное
                        </th>
                        <th
                            scope="col"
                            width="35%">
                            Описание
                        </th>
                    </tr>
                </thead>
                <tbody>
                    <template v-if="value.properties">
                        <tr
                            v-for="propValue, propName in value.properties"
                            :key="propName">
                            <td>{{ propName }}</td>
                            <td>
                                <template v-if="propValue.type === 'array'">
                                    array <a
                                        v-if="propValue.items.$ref"
                                        :href="propValue.items.$ref.replace('/$defs/', '')">{{ propValue.items.$ref.replace("#/$defs/", "") }}</a>
                                </template>
                                <template v-else-if="propValue.type === 'object'">
                                    <pre>{{ propValue.properties }}</pre>
                                </template>
                                <template v-else-if="propValue.$ref">
                                    object <a :href="propValue.$ref.replace('/$defs/', '')">{{ propValue.$ref.replace("#/$defs/", "") }}</a>
                                </template>
                                <template v-else>
                                    {{ propValue.type ? propValue.type : propValue }}
                                </template>
                            </td>
                            <td
                                class="text-center"
                                v-html="isRequired(propName, value.required)" />
                            <td>{{ propValue.description }}</td>
                        </tr>
                    </template>
                    <template v-else>
                        <tr
                            v-for="anyRef in value.oneOf"
                            :key="anyRef.$ref">
                            <td>{{ value.title }}</td>
                            <td>
                                object <a :href="anyRef.$ref.replace('/$defs/', '')">{{ anyRef.$ref.replace("#/$defs/", "") }}</a>
                            </td>
                            <td
                                class="text-center"
                                v-html="isRequired(value.title, [value.title])" />
                            <td>
                                {{ anyRef.description }}
                            </td>
                        </tr>
                    </template>
                </tbody>
            </table>
            <br>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
const loaded = ref(false);
const schema = ref<any>({});
onMounted(async () => {
    loaded.value = false;
    const response = await fetch("/api/json_schema");
    schema.value = await response.json();
    loaded.value = true;
});
function isRequired(fieldName: any, listRequired: string[]) {
    const not = "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"16\" height=\"16\" fill=\"currentColor\" class=\"bi bi-x-lg\" viewBox=\"0 0 16 16\"><path d=\"M2.146 2.854a.5.5 0 1 1 .708-.708L8 7.293l5.146-5.147a.5.5 0 0 1 .708.708L8.707 8l5.147 5.146a.5.5 0 0 1-.708.708L8 8.707l-5.146 5.147a.5.5 0 0 1-.708-.708L7.293 8 2.146 2.854Z\"/></svg>";
    const yes = "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"16\" height=\"16\" fill=\"currentColor\" class=\"bi bi-check-lg\" viewBox=\"0 0 16 16\"><path d=\"M12.736 3.97a.733.733 0 0 1 1.047 0c.286.289.29.756.01 1.05L7.88 12.01a.733.733 0 0 1-1.065.02L3.217 8.384a.757.757 0 0 1 0-1.06.733.733 0 0 1 1.047 0l3.052 3.093 5.4-6.425a.247.247 0 0 1 .02-.022Z\"/></svg>";
    if (!listRequired) {
        return not;
    }
    if (typeof fieldName === "number") {
        fieldName = fieldName.toString();
    }
    if (listRequired.includes(fieldName as string) === true) {
        return yes;
    }
    return not;
}
</script>
