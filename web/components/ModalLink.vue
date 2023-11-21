<script setup lang="ts">
import { Modal } from "bootstrap";
import ActionButton from "./ActionButton.vue";
import Alert from "./Alert.vue";
import { onMounted, ref } from "vue";

const emit = defineEmits(["show", "shown"]);

const props = defineProps({
	action: {
		type: Function,
		required: true,
	},
	class: {
		type: String,
		default: "",
	},
});

const error = ref();
const handler = ref();
const modal = ref();

const hide = () => {
	handler.value.hide();
	handler.value = null;
};

const show = () => {
	handler.value = new Modal(modal.value);
	handler.value.show();
	emit("show");
};

onMounted(() => {
	document.getElementsByTagName("body")[0].appendChild(modal.value);
	modal.value.addEventListener("shown.bs.modal", () => {
		modal.value.getElementsByTagName("input")[0]?.focus();
		emit("shown");
	});
});
</script>

<template>
	<a :class="props.class" @click="show">
		<slot name="link" />
	</a>
	<div ref="modal" class="modal fade" tabindex="-1">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title"><slot name="title" /></h5>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<div class="modal-body">
					<Alert :error="error" />
					<slot name="body" />
				</div>
				<div class="modal-footer">
					<ActionButton :action="action" class="btn btn-primary" @done="hide" @error="error = $event">
						<slot name="action" />
					</ActionButton>
				</div>
			</div>
		</div>
	</div>
</template>
