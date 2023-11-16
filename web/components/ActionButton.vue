<script setup>
import { computed, ref } from "vue";

import Icon from "./Icon.vue";

const emit = defineEmits(["done", "error"]);

const props = defineProps({
	action: {
		type: Function,
		required: true,
	},
});

const mutate = async () => {
	try {
		state.value = "running";
		emit("error", null);
		await props.action();
		state.value = "done";
		await new Promise((resolve) => setTimeout(resolve, 1000));
		emit("done");
		await new Promise((resolve) => setTimeout(resolve, 500));
		state.value = "ready";
	} catch (err) {
		state.value = "error";
		emit("error", err);
	}
};

const state = ref("ready");

const klass = computed(() => {
	return state.value == "done" ? "btn-success" : "btn-primary";
});
</script>
<template>
	<button type="button" class="btn" :class="klass" @click="mutate">
		<slot />
		<Icon name="spinner" class="fa-spin ms-2" :show="state == 'running'" />
		<Icon name="circle-check" class="ms-2" :show="state == 'done'" />
	</button>
</template>
