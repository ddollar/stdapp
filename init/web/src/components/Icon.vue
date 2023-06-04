<script setup>
import { watch } from "vue";
import "@/lib/font-awesome.js";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

const props = defineProps({
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
	<FontAwesomeIcon v-if="show" :icon="name" />
</template>
