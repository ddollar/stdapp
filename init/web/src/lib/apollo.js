import { ApolloClient, InMemoryCache } from "@apollo/client/core";
import { useSubscription } from "@vue/apollo-composable";
import gql from "graphql-tag";

const apolloCache = new InMemoryCache({
	typePolicies: {
		Client: {
			keyFields: ["mac"],
		},
	},
});

import { WebSocketLink } from "@apollo/client/link/ws";
import { watch } from "vue";

const wsLink = new WebSocketLink({
	uri: `wss://${window.location.host}${import.meta.env.BASE_URL.replace(/\/$/, "")}/api/graph`,
	options: {
		reconnect: true,
	},
});

const apolloClient = new ApolloClient({
	cache: apolloCache,
	defaultOptions: {
		watchQuery: {
			fetchPolicy: "cache-and-network",
		},
	},
	link: wsLink,
});

// function cacheDelete(query, field, id) {
// 	let data = apolloCache.readQuery({ query });
// 	let items = data[field].filter((i) => i.id !== id);
// 	data = { ...data };
// 	data[field] = items;
// 	apolloCache.writeQuery({ query: query, overwrite: true, data });
// }

// function cacheUpdate(query, field, item, sort = null) {
// 	let data = apolloCache.readQuery({ query });
// 	let items = [...data[field].filter((i) => i.id !== item.id), item];
// 	if (sort) items = items.sort((a, b) => a[sort].localeCompare(b[sort]));
// 	data = { ...data };
// 	data[field] = items;
// 	apolloCache.writeQuery({ query: query, overwrite: true, data });
// }

function watchTable(name, refetch) {
	// console.log("watching", name);
	const { result } = useSubscription(
		gql`
			subscription ($name: String!) {
				table_changed(name: $name)
			}
		`,
		{
			name: name,
		}
	);
	watch(result, () => {
		// console.log("notify", name);
		refetch();
	});
}

export { apolloClient, apolloCache, watchTable };
