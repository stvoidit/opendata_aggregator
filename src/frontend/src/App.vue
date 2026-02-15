<template>
    <header class="p-1">
        <ul class="nav nav-tabs nav-justified shadow-sm">
            <router-link
                v-for="r in routes"
                :key="r.href"
                v-slot="{ isActive, href, navigate }"
                :to="r.href"
                custom>
                <li
                    class="nav-item">
                    <a
                        :href="href"
                        :class="{'active': isActive}"
                        class="nav-link fw-bold"
                        @click="navigate">
                        {{ r.label }}
                    </a>
                </li>
            </router-link>
        </ul>
    </header>
    <router-view
        v-slot="{ Component }">
        <component :is="Component" />
    </router-view>
</template>
<script setup lang="ts">
const routes = [
    {
        href: "/",
        label: "Статистика"
    },
    {
        href: "/doc",
        label: "Документация"
    },
    {
        href: "/doc/api",
        label: "Описание API"
    },
    {
        href: "/handbook_okved",
        label: "Справочник ОКВЭД"
    },
    {
        href: "/hotels",
        label: "Гостиницы"
    }
];
</script>
<style lang="scss">
@import "@/styles.scss";
</style>
<style scoped>
.active {
    background-color: rgb(81, 135, 235) !important;
    color: rgba(255, 255, 255, 0.938) !important;
}
.nav-link{
    color: rgb(34, 24, 24) ;
}
</style>
