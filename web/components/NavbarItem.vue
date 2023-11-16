<script setup>
import { RouterLink } from "vue-router";

import Icon from "./Icon.vue";

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
		<RouterLink v-if="name" :to="{ name: name }">
			<Icon v-if="icon" :name="icon" />
			<slot />
		</RouterLink>
		<a v-else><slot /></a>
	</li>
</template>
