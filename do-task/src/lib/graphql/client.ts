import { ApolloClient } from '@apollo/client/core';
import { InMemoryCache } from '@apollo/client/cache';
import { createHttpLink } from '@apollo/client/link/http';
import { setContext } from '@apollo/client/link/context';
import { browser } from '$app/environment';
import { env } from '$env/dynamic/public';

// Create HTTP link to GraphQL server
const httpLink = createHttpLink({
	uri: env.PUBLIC_GRAPHQL_ENDPOINT || 'http://localhost:8080/query',
	credentials: 'include'
});

// Create auth link to add authorization header
const authLink = setContext((_, { headers }) => {
	// Get token from localStorage if available
	let token = '';
	if (browser) {
		token = localStorage.getItem('accessToken') || '';
	}

	return {
		headers: {
			...headers,
			...(token && { Authorization: `Bearer ${token}` })
		}
	};
});

// Create Apollo Client instance
export const client = new ApolloClient({
	link: authLink.concat(httpLink),
	cache: new InMemoryCache({
		typePolicies: {
			Query: {
				fields: {
					tasks: {
						read(existing) {
							return existing;
						},
						merge(_existing, incoming) {
							return incoming;
						},
						keyArgs: false
					},
					categories: {
						read(existing) {
							return existing;
						},
						merge(_existing, incoming) {
							return incoming;
						},
						keyArgs: false
					},
					task: {
						keyArgs: ['id']
					},
					category: {
						keyArgs: ['id']
					}
				}
			},
			Task: {
				keyFields: ['id']
			},
			Category: {
				keyFields: ['id']
			}
		}
	}),
	defaultOptions: {
		watchQuery: {
			errorPolicy: 'all',
			fetchPolicy: 'cache-and-network'
		},
		query: {
			errorPolicy: 'all',
			fetchPolicy: 'cache-first'
		},
		mutate: {
			errorPolicy: 'all'
		}
	}
});
