<template>
    <div class="container-fluid">
        <div class="row mb-3">
            <div class="col-6 d-flex flex-wrap p-0">
                <LastUpdatesSources :class="listClasses" />
                <StatsHotels :class="listClasses" />
            </div>
            <div class="col-6 d-flex flex-wrap p-0">
                <StatsBalance :class="listClasses" />
                <StatsEGR :class="listClasses" />
                <StatsRossAccreditation :class="listClasses" />
                <StatsRegisterOffTrademarks :class="listClasses" />
                <StatsSMP :class="listClasses" />
                <StatsDEBTAM :class="listClasses" />
                <StatsTaxOffenses :class="listClasses" />
                <StatsFSSP :class="listClasses" />
                <StatsFGIS :class="listClasses" />
                <StatsTaxRegime :class="listClasses" />
            </div>
        </div>
    </div>
    <!-- <SpinnerLoader v-else /> -->
</template>

<script setup lang="ts">
import StatsEGR from "./StatsEGR.vue";
import StatsBalance from "./StatsBalance.vue";
import LastUpdatesSources from "./LastUpdatesSources.vue";
import StatsHotels from "./StatsHotels.vue";
import StatsRossAccreditation from "./StatsRossAccreditation.vue";
import StatsTaxOffenses from "./StatsTaxOffenses.vue";
import StatsRegisterOffTrademarks from "./StatsRegisterOffTrademarks.vue";
import StatsFSSP from "./StatsFSSP.vue";
import StatsSMP from "./StatsSMP.vue";
import StatsDEBTAM from "./StatsDEBTAM.vue";
import StatsFGIS from "./StatsFGIS.vue";
import StatsTaxRegime from "./StatsTaxRegime.vue";
const listClasses = [
    "p-2",
    "m-1",
    "flex-fill",
    "card",
    "shadow-sm"
];
</script>
<style>
.p-0 {
    padding: 0;
}
</style>
