import type { DocumentNode, WatchQueryOptions } from '@apollo/client/core';
import type { FetchPolicy } from '@apollo/client/core/watchQueryOptions';
import { readable, derived } from 'svelte/store';
import { client } from './client';

/**
 * Execute a GraphQL query and return the results as a Svelte readable store.
 * @param query GraphQL query document
 * @param variables Query variables
 * @param options Additional options for the query
 */
export function queryStore(
	query: DocumentNode,
	variables = {},
	options: Partial<WatchQueryOptions> = {}
) {
	const baseStore = readable({ data: undefined, loading: true, error: undefined }, (set) => {
		let active = true;

		client
			.query({
				query,
				variables,
				...options,
				fetchPolicy: (options.fetchPolicy as FetchPolicy) || ('network-only' as FetchPolicy)
			})
			.then((result) => {
				if (active) {
					set({ data: result.data, loading: false, error: undefined });
				}
			})
			.catch((error) => {
				if (active) {
					set({ data: undefined, loading: false, error });
				}
			});

		return () => {
			active = false;
		};
	});

	// Add derived stores for convenience
	const result = derived(baseStore, ($baseStore) => $baseStore.data);
	const isLoading = derived(baseStore, ($baseStore) => $baseStore.loading);
	const queryError = derived(baseStore, ($baseStore) => $baseStore.error);

	// Return the base store with additional properties
	return Object.assign(baseStore, { result, isLoading, queryError });
}

/**
 * Execute a GraphQL mutation and return a function to trigger it
 * @param mutation GraphQL mutation document
 */
export function createMutation(mutation: DocumentNode) {
	return async (variables = {}, options = {}) => {
		try {
			const result = await client.mutate({
				mutation,
				variables,
				...options
			});
			return result.data;
		} catch (error) {
			console.error('Mutation error:', error);
			throw error;
		}
	};
}
