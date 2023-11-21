<script setup lang="ts">
import { GraphQLError } from "graphql";

interface ApolloError extends Error {
	graphQLErrors: GraphQLError[];
	extensions: {
		stacktrace: string[];
	};
}

defineProps<{
	error: Error | ApolloError;
}>();
</script>

<template>
	<div v-if="error && 'graphQLErrors' in error">
		<div v-for="(err, index) in error.graphQLErrors" :key="index" class="alert alert-danger">
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
