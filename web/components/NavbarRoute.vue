<script setup lang="ts">
import NavbarDropdown from "./NavbarDropdown.vue";
import NavbarItem from "./NavbarItem.vue";
import NavbarRoutes from "./NavbarRoutes.vue";

defineProps({
	route: {
		type: Object,
		default: null,
	},
});
</script>

<template>
	<NavbarDropdown v-if="route?.children" :name="route.name" :icon="route.meta?.icon">
		{{ route.meta?.title }}
		<template #items>
			<NavbarRoutes :routes="route?.children" />
		</template>
	</NavbarDropdown>
	<NavbarItem v-else-if="route" :name="route.name" :icon="route.meta?.icon">
		{{ route.meta?.title }}
	</NavbarItem>
</template>
