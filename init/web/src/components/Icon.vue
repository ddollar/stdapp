<script setup>
import "@/lib/font-awesome.js";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

const props = defineProps({
	name: {
		type: String,
		required: true,
	},
});
</script>

<template>
	<FontAwesomeIcon :icon="name" />
</template>
