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
@import "vue-select/dist/vue-select.css";

@import "@/styles/bootstrap.scss";
@import "@/styles/stdapp.scss";
@import "@/styles/app.scss";
</style>
