import { RouteRecordRaw, createRouter, createWebHistory } from "vue-router";

import DatabaseStat from "@/pages/DatabaseStat/index.vue";
import Doc from "@/pages/Doc.vue";
import DocAPI from "@/pages/DocAPI.vue";
import HandbookOKVED from "@/pages/HandbookOKVED.vue";
import HotelsData from "@/pages/HotelsData.vue";
import ServiceLogSources from "@/pages/ServiceLogSources.vue";
import UploadSource from "@/pages/UploadSource.vue";
import { useTitle } from "@vueuse/core";

const routes: RouteRecordRaw[] = [
    {
        path: "/",
        meta: {
            subtitle: "Статистика"
        },
        component: DatabaseStat
    },
    {
        path: "/doc",
        meta: {
            subtitle: "Документация | объекты"
        },
        component: Doc
    },
    {
        path: "/doc/api",
        meta: {
            subtitle: "Документация | api"
        },
        component: DocAPI
    },
    {
        path: "/service_log_sources",
        meta: {
            subtitle: "Логи обновления"
        },
        component: ServiceLogSources
    },
    {
        path: "/handbook_okved",
        meta: {
            subtitle: "Подсказки ОКВЭД"
        },
        component: HandbookOKVED
    },
    {
        path: "/hotels",
        meta: {
            subtitle: "Гостиницы"
        },
        component: HotelsData
    },
    {
        path: "/upload",
        component: UploadSource
    }
];
const router = createRouter({
    history: createWebHistory(),
    routes: routes,
    scrollBehavior: (to) => {
        if (to.hash) {
            return new Promise((resolve) => {
                setTimeout(() => {
                    resolve({
                        el: to.hash,
                        behavior: "smooth"
                    });
                }, 500);
            });
        }
    }
});
router.afterEach((to) => {
    const baseTitle = "oda";
    const title = useTitle();
    if (to.meta?.subtitle as string) {
        title.value = `${baseTitle} | ${to.meta.subtitle as string}`;
    }
});
export default router;
