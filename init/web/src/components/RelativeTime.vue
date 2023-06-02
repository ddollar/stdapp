<script setup>
import { computed } from "vue";
import { v4 as uuid } from "uuid";

const props = defineProps(["time", "zero"]);

const datetime = computed(() => Date.parse(props.time));
const id = computed(() => `timeago-${uuid()}`);
</script>

<template>
	<div>
		<timeago v-if="time" :id="id" :datetime="datetime" :autoUpdate="1" />
		<div v-else>{{ zero }}</div>
	</div>
</template>
