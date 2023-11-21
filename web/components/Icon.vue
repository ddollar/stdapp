<script setup lang="ts">
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

defineProps({
	name: {
		type: String,
		required: true,
	},
	show: {
		type: Boolean,
		default: true,
	},
});
</script>

<template>
	<FontAwesomeIcon v-if="show" :icon="name" class="icon fa-fw" />
</template>
