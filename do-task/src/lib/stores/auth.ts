import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { client } from '$lib/graphql/client';
import { LOGIN_MUTATION, REGISTER_MUTATION, ME_QUERY } from '$lib/graphql/queries';

interface User {
	id: string;
	email: string;
	name: string;
}

interface AuthState {
	user: User | null;
	isAuthenticated: boolean;
	isLoading: boolean;
	error: string | null;
}

const initialState: AuthState = {
	user: null,
	isAuthenticated: false,
	isLoading: true,
	error: null
};

export const authStore = writable<AuthState>(initialState);

// Derived stores for convenient access
export const user = derived(authStore, ($auth) => $auth.user);
export const isAuthenticated = derived(authStore, ($auth) => $auth.isAuthenticated);
export const isLoading = derived(authStore, ($auth) => $auth.isLoading);
export const authError = derived(authStore, ($auth) => $auth.error);

class AuthService {
	async login(email: string, password: string): Promise<void> {
		authStore.update((state) => ({ ...state, isLoading: true, error: null }));

		try {
			// Use GraphQL client to login instead of REST API
			const result = await client.mutate({
				mutation: LOGIN_MUTATION,
				variables: {
					input: { email, password }
				}
			});

			if (result.errors) {
				throw new Error(result.errors[0]?.message || 'Login failed');
			}

			if (!result.data?.login) {
				throw new Error('Login failed - no data returned');
			}

			const userData = result.data.login;

			authStore.update((state) => ({
				...state,
				user: userData.user,
				isAuthenticated: true,
				isLoading: false,
				error: null
			}));

			if (browser && userData.token) {
				localStorage.setItem('accessToken', userData.token);
			}

			if (browser) {
				goto('/');
			}
		} catch (error) {
			authStore.update((state) => ({
				...state,
				isLoading: false,
				error: error instanceof Error ? error.message : 'An unknown error occurred'
			}));
			throw error;
		}
	}

	async register(name: string, email: string, password: string): Promise<void> {
		authStore.update((state) => ({ ...state, isLoading: true, error: null }));

		try {
			// Use GraphQL client instead of REST API
			const result = await client.mutate({
				mutation: REGISTER_MUTATION,
				variables: {
					input: { name, email, password }
				}
			});

			if (result.errors) {
				throw new Error(result.errors[0]?.message || 'Registration failed');
			}

			if (!result.data?.register) {
				throw new Error('Registration failed - no data returned');
			}

			const userData = result.data.register;

			authStore.update((state) => ({
				...state,
				user: userData.user,
				isAuthenticated: true,
				isLoading: false,
				error: null
			}));

			if (browser && userData.token) {
				localStorage.setItem('accessToken', userData.token);
			}

			if (browser) {
				goto('/');
			}
		} catch (error) {
			authStore.update((state) => ({
				...state,
				isLoading: false,
				error: error instanceof Error ? error.message : 'An unknown error occurred'
			}));
			throw error;
		}
	}

	async logout(): Promise<void> {
		try {
			// Clear HTTP-only cookies on server
			await fetch('/api/logout', {
				method: 'POST',
				credentials: 'include'
			});
		} catch (error) {
			console.error('Logout API error:', error);
			// Continue with client-side cleanup even if server logout fails
		}

		// Always perform client-side cleanup regardless of server response
		try {
			// Clear localStorage tokens
			if (browser) {
				localStorage.removeItem('accessToken');
				localStorage.removeItem('refreshToken');
				localStorage.removeItem('username');
			}

			// Clear Apollo cache to remove any cached user data
			await client.clearStore();

			// Reset cache entirely to ensure clean state
			await client.resetStore();
		} catch (error) {
			console.error('Error during client cleanup:', error);
		}

		// Reset auth store
		authStore.set({
			...initialState,
			isLoading: false
		});

		if (browser) {
			goto('/auth/login');
		}
	}

	async checkAuth(): Promise<void> {
		if (!browser) return;

		try {
			// Use GraphQL client to check authentication
			const result = await client.query({
				query: ME_QUERY,
				fetchPolicy: 'network-only'
			});

			if (result.data?.me) {
				authStore.update((state) => ({
					...state,
					user: result.data.me,
					isAuthenticated: true,
					isLoading: false
				}));
			} else {
				authStore.update((state) => ({
					...state,
					isAuthenticated: false,
					isLoading: false
				}));
			}
		} catch (error) {
			authStore.update((state) => ({
				...state,
				isAuthenticated: false,
				isLoading: false,
				error: 'Failed to check authentication'
			}));
		}
	}
}

export const authService = new AuthService();
