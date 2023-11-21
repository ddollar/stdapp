<script setup lang="ts">
import { computed } from "vue";
import { money } from "../format";

const props = defineProps({
	amount: {
		type: Number,
		default: () => Math.random() * 100 - 50,
	},
});

const klass = computed(() => {
	return props.amount < 0 ? "negative" : "positive";
});

// onMounted(() => {
// 	props.amount.value = Math.random() * 1000;
// });
</script>

<template>
	<div class="amount" :class="klass">{{ money(amount) }}</div>
</template>
