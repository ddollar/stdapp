<script setup>
import NavbarRoute from "./NavbarRoute.vue";

const props = defineProps({
	filter: {
		type: Function,
		default: () => true,
	},
	routes: {
		type: Array,
		default: null,
	},
});

const visible = props.routes.filter((route) => route.meta?.icon).filter(props.filter);
</script>

<template>
	<NavbarRoute v-for="route in visible" :key="route.name" :route="route" />
</template>
