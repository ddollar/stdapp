<script setup lang="ts">
import { computed } from "vue";
import { v4 as uuid } from "uuid";

const props = defineProps({
	time: {
		type: String,
		default: null,
	},
	zero: {
		type: String,
		default: "",
	},
});

const datetime = computed(() => Date.parse(props.time));
const id = computed(() => `timeago-${uuid()}`);
</script>

<template>
	<div>
		<timeago v-if="time" :id="id" :datetime="datetime" :auto-update="1" />
		<div v-else>{{ zero }}</div>
	</div>
</template>
