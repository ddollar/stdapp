<script setup lang="ts">
import container from "../container";
</script>

<template>
	<nav class="navbar navbar-expand navbar-dark fixed-top bg-dark">
		<div :class="container">
			<a class="navbar-brand" href="/">
				<slot name="brand" />
			</a>
			<slot />
		</div>
	</nav>
</template>

<style scoped lang="scss">
@import "../styles/stdapp";

.navbar {
	margin-right: calc($gutter * -0.5);
}
</style>
