<script setup>
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
	<li class="nav-item dropdown" @click="$emit('click', $event)">
		<a class="nav-link" href="#" role="button" data-bs-toggle="dropdown">
			<Icon v-if="icon" :name="icon" />
			<span><slot /></span>
		</a>
		<ul class="dropdown-menu dropdown-menu-end dropdown-menu-dark">
			<slot name="items" />
		</ul>
	</li>
</template>
