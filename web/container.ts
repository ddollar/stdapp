import { computed } from "vue";
import { useRoute } from "vue-router";

const container = computed(() => {
	const route = useRoute();

	switch (true) {
		case route.meta.wide:
			return "container-fluid";
		default:
			return "container";
	}
});

export default container;
