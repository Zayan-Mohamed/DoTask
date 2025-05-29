import { ApolloClient } from '@apollo/client/core';
import { InMemoryCache } from '@apollo/client/cache';
import { createHttpLink } from '@apollo/client/link/http';

// Create HTTP link to GraphQL server
const httpLink = createHttpLink({
	uri: 'http://localhost:8080/query'
});

// Create Apollo Client instance
export const client = new ApolloClient({
	link: httpLink,
	cache: new InMemoryCache(),
	defaultOptions: {
		watchQuery: {
			errorPolicy: 'all'
		},
		query: {
			errorPolicy: 'all'
		}
	}
});
