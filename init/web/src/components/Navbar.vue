<script setup>
import container from "@/lib/container";
import NavbarLink from "@/components/NavbarLink.vue";
</script>

<template>
	<nav class="navbar navbar-expand-md navbar-dark fixed-top bg-dark">
		<div :class="container">
			<a class="navbar-brand" href="/">stdapp-init</a>
			<ul class="navbar-nav ml-auto">
				<NavbarLink to="/">Home</NavbarLink>
			</ul>
		</div>
	</nav>
</template>
