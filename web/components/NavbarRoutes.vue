<script setup lang="ts">
import type { RouteLocation } from "vue-router";

import NavbarRoute from "./NavbarRoute.vue";

type FilterFunction = (route: RouteLocation) => boolean;

interface Props {
	filter: FilterFunction;
	routes: Array<RouteLocation>;
}

const props = withDefaults(defineProps<Props>(), {
	filter: () => true,
});

const visible = props.routes.filter((route) => route.meta?.icon).filter(props.filter);
</script>

<template>
	<NavbarRoute v-for="(route, index) in visible" :key="index" :route="route" />
</template>
