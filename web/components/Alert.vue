<script setup>
defineProps({
	error: {
		type: Error,
		default: null,
	},
});
</script>

<template>
	<div v-if="error?.graphQLErrors">
		<div v-for="err in error?.graphQLErrors" :key="err" class="alert alert-danger">
			<div class="d-flex">
				<h6 class="mb-0 flex-grow-1">{{ err.message }}</h6>
				<a v-if="error.extensions?.stacktrace" style="font-size: 0.8em" data-bs-toggle="modal" href="#stacktrace">
					Trace
				</a>
			</div>
			<div v-if="error.extensions?.stacktrace" id="stacktrace" class="modal fade">
				<div class="modal-dialog modal-xl">
					<div class="modal-content">
						<div class="modal-body font-monospace" style="font-size: 0.7em">
							<div v-for="line in error.extensions?.stacktrace" :key="line" class="mb-2">
								{{ line }}
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
	<div v-else-if="error?.message" class="alert alert-danger">{{ error.message }}</div>
</template>
