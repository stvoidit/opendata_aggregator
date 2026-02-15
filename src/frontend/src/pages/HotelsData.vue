<template>
    <div
        v-if="loaded"
        class="container-fluid">
        <div class="row mt-2">
            <div class="col-2 text-start">
                <p><strong>Всего строк:</strong> {{ totalRows.toLocaleString() }}</p>
            </div>
            <div class="col-2">
                <div class="form-floating">
                    <select
                        id="floatingSelect"
                        v-model="selectedRegion"
                        class="form-select my-1"
                        @change="()=> page = 1">
                        <option :value="null">
                            ---
                        </option>
                        <option
                            v-for="value in regions"
                            :key="value"
                            :value="value">
                            {{ value }}
                        </option>
                    </select>
                    <label for="floatingSelect">Выбор региона</label>
                </div>
            </div>
            <div class="col-4">
                <Pagination
                    v-model="page"
                    alignment="center"
                    :per-page="limit"
                    :total="totalRows"
                    show-jump-buttons
                    :show-prev-next-button="false" />
            </div>
            <div class="col-4 text-end">
                <a
                    href="/api/download/hotels"
                    target="_blank"
                    class="btn btn-sm btn-success">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        fill="currentColor"
                        class="bi bi-filetype-xlsx"
                        viewBox="0 0 16 16">
                        <path
                            fill-rule="evenodd"
                            d="M14 4.5V11h-1V4.5h-2A1.5 1.5 0 0 1 9.5 3V1H4a1 1 0 0 0-1 1v9H2V2a2 2 0 0 1 2-2h5.5L14 4.5ZM7.86 14.841a1.13 1.13 0 0 0 .401.823c.13.108.29.192.479.252.19.061.411.091.665.091.338 0 .624-.053.858-.158.237-.105.416-.252.54-.44a1.17 1.17 0 0 0 .187-.656c0-.224-.045-.41-.135-.56a1.002 1.002 0 0 0-.375-.357 2.028 2.028 0 0 0-.565-.21l-.621-.144a.97.97 0 0 1-.405-.176.37.37 0 0 1-.143-.299c0-.156.061-.284.184-.384.125-.101.296-.152.513-.152.143 0 .266.023.37.068a.624.624 0 0 1 .245.181.56.56 0 0 1 .12.258h.75a1.093 1.093 0 0 0-.199-.566 1.21 1.21 0 0 0-.5-.41 1.813 1.813 0 0 0-.78-.152c-.293 0-.552.05-.777.15-.224.099-.4.24-.527.421-.127.182-.19.395-.19.639 0 .201.04.376.123.524.082.149.199.27.351.367.153.095.332.167.54.213l.618.144c.207.049.36.113.462.193a.387.387 0 0 1 .153.326.512.512 0 0 1-.085.29.558.558 0 0 1-.255.193c-.111.047-.25.07-.413.07-.117 0-.224-.013-.32-.04a.837.837 0 0 1-.249-.115.578.578 0 0 1-.255-.384h-.764Zm-3.726-2.909h.893l-1.274 2.007 1.254 1.992h-.908l-.85-1.415h-.035l-.853 1.415H1.5l1.24-2.016-1.228-1.983h.931l.832 1.438h.036l.823-1.438Zm1.923 3.325h1.697v.674H5.266v-3.999h.791v3.325Zm7.636-3.325h.893l-1.274 2.007 1.254 1.992h-.908l-.85-1.415h-.035l-.853 1.415h-.861l1.24-2.016-1.228-1.983h.931l.832 1.438h.036l.823-1.438Z" />
                    </svg>
                    скачать
                </a>
            </div>
        </div>
        <div class="row">
            <div class="col">
                <table class="table table-sm table-bordered shadow-sm">
                    <thead>
                        <tr>
                            <th
                                v-for="h in headers"
                                :key="h">
                                {{ h }}
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr
                            v-for="h in tableHotels"
                            :key="h.federal_number">
                            <td>
                                <table class="table-sm table table-borderless">
                                    <tbody>
                                        <tr>
                                            <td>Номер:</td>
                                            <td>{{ h.federal_number }}</td>
                                        </tr>
                                        <tr>
                                            <td>ИНН:</td>
                                            <td>{{ h.inn }}</td>
                                        </tr>
                                        <tr>
                                            <td>ОГРН:</td>
                                            <td>{{ h.ogrn }}</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </td>
                            <td>
                                <table class="table-sm table table-borderless">
                                    <tbody>
                                        <tr>
                                            <td>Название:</td>
                                            <td>{{ h.short_name }}</td>
                                        </tr>
                                        <tr>
                                            <td>Владелец:</td>
                                            <td>{{ h.owner }}</td>
                                        </tr>
                                        <tr>
                                            <td>Вид:</td>
                                            <td>{{ h.type }}</td>
                                        </tr>
                                        <tr>
                                            <td>Регион:</td>
                                            <td>{{ h.region }}</td>
                                        </tr>
                                        <tr>
                                            <td>Адрес:</td>
                                            <td>{{ h.address }}</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </td>
                            <td>
                                <ul>
                                    <li><a :href="`mailto:${h.email}`">{{ h.email }}</a></li>
                                    <li>
                                        <a
                                            :href="h.site"
                                            target="_blank">{{ h.site }}</a>
                                    </li>
                                    <li>{{ h.phone }}</li>
                                    <li>{{ h.fax }}</li>
                                </ul>
                            </td>
                            <td>
                                <table
                                    v-for="cls in h.classification"
                                    :key="cls.license_number"
                                    class="table-sm table mb-2">
                                    <tbody>
                                        <tr>
                                            <td>Номер лицензии:</td>
                                            <td>{{ cls.license_number }}</td>
                                        </tr>
                                        <tr>
                                            <td>Регистрационный номер:</td>
                                            <td>{{ cls.registration_number }}</td>
                                        </tr>
                                        <tr>
                                            <td>Категория:</td>
                                            <td>{{ cls.category }}</td>
                                        </tr>
                                        <tr>
                                            <td>Дата выдачи:</td>
                                            <td>{{ cls.date_issued }}</td>
                                        </tr>
                                        <tr>
                                            <td>Действительно до:</td>
                                            <td>{{ cls.date_end }}</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </td>
                            <td>
                                <table class="table table-sm">
                                    <thead>
                                        <tr>
                                            <th>категория</th>
                                            <th>комнаты</th>
                                            <th>места</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        <tr
                                            v-for="room in h.rooms"
                                            :key="room.category">
                                            <td>{{ room.category }}</td>
                                            <td>{{ room.rooms }}</td>
                                            <td>{{ room.seats }}</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
    <SpinnerLoader v-else />
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from "vue";
import Pagination from "@/components/Pagination.vue";
import SpinnerLoader from "@/components/SpinnerLoader.vue";
interface IHotel {
    federal_number: string;
    inn: string;
    ogrn: string;
    full_name: string;
    short_name: string;
    type: string;
    address: string;
    email: string;
    region: string;
    site: string;
    phone: string;
    fax: string;
    owner: string;
    classification: {
        federal_number: string;
        date_issued: string;
        date_end: string;
        category: string;
        license_number: string;
        registration_number: string;
    }[];
    rooms: {
        federal_number: string;
        category: string;
        rooms: number;
        seats: number;
    }[];
}
const headers = [
    "Реквизиты",
    "Название, вид и владелец",
    "Контакты",
    "Классификация",
    "Номера"
];
const selectedRegion = ref(null);
const loaded = ref(false);
const hotels = ref<IHotel[]>([]);
const page = ref(1);
const limit = ref(20);
onMounted(async () => {
    loaded.value = false;
    const response = await fetch("/api/hotels");
    hotels.value = await response.json();
    hotels.value.forEach((h) => {
        h.region = h.region.trim();
    });
    loaded.value = true;
});
const filteredHotels = computed(() => {
    if (selectedRegion.value === null) {
        return hotels.value;
    }
    return hotels.value.filter(h => h.region === selectedRegion.value);
});
const tableHotels = computed(() => {
    const start = (page.value - 1) * limit.value;
    const end = page.value * limit.value;
    const dataPart = filteredHotels.value.slice(start, end);
    dataPart.forEach((h) => {
        if (!h.site.startsWith("http") && h.site.length > 3) {
            h.site = `https://${h.site}`;
        }
    });
    return dataPart;
});
const totalRows = computed(() => filteredHotels.value.length);
const regions = computed(() => {
    const set = new Set<string>();
    for (const h of hotels.value) {
        set.add(h.region);
    }
    return Array.from(set).sort();
});
</script>
