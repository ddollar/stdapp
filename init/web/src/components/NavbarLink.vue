<script setup>
import { RouterLink } from "vue-router";

defineProps(["to"]);
</script>

<template>
	<li class="nav-item ms-3">
		<RouterLink :to="to" class="nav-link">
			<slot />
		</RouterLink>
	</li>
</template>
