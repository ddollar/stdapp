<script setup>
import { RouterView } from "vue-router";
import Navbar from "@/components/Navbar.vue";
import container from "@/lib/container";
</script>

<template>
	<header><Navbar /></header>

	<div id="app" :class="container">
		<RouterView />
	</div>
</template>

<style lang="scss">
@import "@/styles/bootstrap.scss";
@import "bootstrap";
@import "@/styles/theme.scss";
@import "vue-select/dist/vue-select.css";
</style>
