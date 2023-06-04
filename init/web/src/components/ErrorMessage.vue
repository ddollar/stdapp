<script setup>
const props = defineProps(["error"]);
</script>

<template>
	<div class="alert alert-danger" v-for="error in props.error?.graphQLErrors" :key="error">
		<div class="d-flex">
			<h6 class="mb-0 flex-grow-1">{{ error.message }}</h6>
			<a style="font-size: 0.8em" data-bs-toggle="modal" href="#stacktrace" v-if="error.extensions?.stacktrace">
				Trace
			</a>
		</div>
		<div id="stacktrace" class="modal fade" v-if="error.extensions?.stacktrace">
			<div class="modal-dialog modal-xl">
				<div class="modal-content">
					<div class="modal-body font-monospace" style="font-size: 0.7em">
						<div class="mb-2" v-for="line in error.extensions?.stacktrace" :key="line">
							{{ line }}
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
