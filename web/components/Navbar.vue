<script setup>
import container from "../lib/container";
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

.avatar {
	height: 24px;
}

.navbar {
	margin-right: calc($gutter * -0.5);
}
</style>
