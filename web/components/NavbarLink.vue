<script setup>
import { RouterLink } from "vue-router";

defineEmits(["click"]);

defineProps({
	icon: {
		type: String,
		default: null,
	},
	name: {
		type: String,
		default: null,
	},
});
</script>

<template>
	<li class="nav-item" @click="$emit('click', $event)">
		<RouterLink v-if="name" :to="{ name: name }"><slot /></RouterLink>
		<a v-else><slot /></a>
		<slot name="after" />
	</li>
</template>
