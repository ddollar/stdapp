<script setup lang="ts">
import { ref } from "vue";

defineEmits(["update:modelValue"]);

defineExpose({
	focus: () => {
		field?.value?.focus();
	},
});

defineProps({
	label: {
		type: String,
		default: null,
	},
	modelValue: {
		type: String,
		default: null,
	},
	type: {
		type: String,
		default: "text",
	},
});

const field = ref();
</script>

<template>
	<div class="form-group mb-3">
		<label v-if="label" :for="field?.id" class="form-label">{{ label }}</label>
		<div v-if="type === 'custom'"><slot /></div>
		<div v-else-if="type === 'money'" class="input-group">
			<span class="input-group-text">$</span>
			<input
				ref="field"
				class="form-control"
				type="number"
				inputmode="decimal"
				min="0.01"
				step="0.01"
				:value="modelValue"
				@change="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
			/>
		</div>
		<div v-else-if="type === 'text'">
			<input
				ref="field"
				class="form-control"
				:type="type"
				:value="modelValue"
				@change="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
			/>
		</div>
	</div>
</template>
