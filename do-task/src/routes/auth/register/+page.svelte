<script lang="ts">
	import { authService, authError, isLoading } from '$lib/stores/auth';

	let name = $state('');
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let showPassword = $state(false);
	let showConfirmPassword = $state(false);
	let validationErrors = $state<{[key: string]: string}>({});

	function validateForm(): boolean {
		validationErrors = {};
		
		if (!name.trim()) {
			validationErrors.name = 'Name is required';
		} else if (name.trim().length < 2) {
			validationErrors.name = 'Name must be at least 2 characters';
		}
		
		if (!email) {
			validationErrors.email = 'Email is required';
		} else if (!/\S+@\S+\.\S+/.test(email)) {
			validationErrors.email = 'Please enter a valid email';
		}
		
		if (!password) {
			validationErrors.password = 'Password is required';
		} else if (password.length < 8) {
			validationErrors.password = 'Password must be at least 8 characters';
		} else if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(password)) {
			validationErrors.password = 'Password must contain at least one uppercase letter, one lowercase letter, and one number';
		}
		
		if (!confirmPassword) {
			validationErrors.confirmPassword = 'Please confirm your password';
		} else if (password !== confirmPassword) {
			validationErrors.confirmPassword = 'Passwords do not match';
		}
		
		return Object.keys(validationErrors).length === 0;
	}

	async function handleRegister(event: Event) {
		event.preventDefault();
		
		if (!validateForm()) return;
		
		try {
			await authService.register(name.trim(), email, password);
		} catch (error) {
			// Error is handled in the auth store
			console.error('Registration failed:', error);
		}
	}
</script>

<svelte:head>
	<title>Register - DoTask</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-8">
		<div>
			<div class="mx-auto h-12 w-auto flex justify-center">
				<svg width="120" height="40" viewBox="0 0 160 40" fill="none" xmlns="http://www.w3.org/2000/svg">
					<rect x="2" y="8" width="24" height="24" rx="6" fill="#3B82F6" />
					<path
						d="M10 20L14 24L18 16"
						stroke="white"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					/>
					<text x="32" y="26" fill="currentColor" font-size="20" font-weight="bold" class="fill-primary">DoTask</text>
				</svg>
			</div>
			<h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900 dark:text-white">
				Create your account
			</h2>
			<p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
				Already have an account?
				<a href="/auth/login" class="font-medium text-primary hover:text-primary/80">
					Sign in here
				</a>
			</p>
		</div>

		<form class="mt-8 space-y-6" onsubmit={handleRegister}>
			{#if $authError}
				<div class="bg-red-100 dark:bg-red-900/30 border-l-4 border-red-500 text-red-700 dark:text-red-300 p-4" role="alert">
					<p>{$authError}</p>
				</div>
			{/if}

			<div class="space-y-4">
				<div>
					<label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
						Full Name
					</label>
					<input
						id="name"
						name="name"
						type="text"
						autocomplete="name"
						required
						bind:value={name}
						class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white dark:bg-gray-700 rounded-md focus:outline-none focus:ring-primary focus:border-primary focus:z-10 sm:text-sm"
						placeholder="Enter your full name"
						class:border-red-500={validationErrors.name}
					/>
					{#if validationErrors.name}
						<p class="mt-1 text-sm text-red-600 dark:text-red-400">{validationErrors.name}</p>
					{/if}
				</div>

				<div>
					<label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
						Email address
					</label>
					<input
						id="email"
						name="email"
						type="email"
						autocomplete="email"
						required
						bind:value={email}
						class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white dark:bg-gray-700 rounded-md focus:outline-none focus:ring-primary focus:border-primary focus:z-10 sm:text-sm"
						placeholder="Enter your email"
						class:border-red-500={validationErrors.email}
					/>
					{#if validationErrors.email}
						<p class="mt-1 text-sm text-red-600 dark:text-red-400">{validationErrors.email}</p>
					{/if}
				</div>

				<div>
					<label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
						Password
					</label>
					<div class="mt-1 relative">
						<input
							id="password"
							name="password"
							type={showPassword ? 'text' : 'password'}
							autocomplete="new-password"
							required
							bind:value={password}
							class="appearance-none relative block w-full px-3 py-2 pr-10 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white dark:bg-gray-700 rounded-md focus:outline-none focus:ring-primary focus:border-primary focus:z-10 sm:text-sm"
							placeholder="Enter your password"
							class:border-red-500={validationErrors.password}
						/>
						<button
							type="button"
							class="absolute inset-y-0 right-0 pr-3 flex items-center"
							onclick={() => showPassword = !showPassword}
						>
							{#if showPassword}
								<svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
								</svg>
							{:else}
								<svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21" />
								</svg>
							{/if}
						</button>
					</div>
					{#if validationErrors.password}
						<p class="mt-1 text-sm text-red-600 dark:text-red-400">{validationErrors.password}</p>
					{/if}
				</div>

				<div>
					<label for="confirmPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
						Confirm Password
					</label>
					<div class="mt-1 relative">
						<input
							id="confirmPassword"
							name="confirmPassword"
							type={showConfirmPassword ? 'text' : 'password'}
							autocomplete="new-password"
							required
							bind:value={confirmPassword}
							class="appearance-none relative block w-full px-3 py-2 pr-10 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white dark:bg-gray-700 rounded-md focus:outline-none focus:ring-primary focus:border-primary focus:z-10 sm:text-sm"
							placeholder="Confirm your password"
							class:border-red-500={validationErrors.confirmPassword}
						/>
						<button
							type="button"
							class="absolute inset-y-0 right-0 pr-3 flex items-center"
							onclick={() => showConfirmPassword = !showConfirmPassword}
						>
							{#if showConfirmPassword}
								<svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
								</svg>
							{:else}
								<svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21" />
								</svg>
							{/if}
						</button>
					</div>
					{#if validationErrors.confirmPassword}
						<p class="mt-1 text-sm text-red-600 dark:text-red-400">{validationErrors.confirmPassword}</p>
					{/if}
				</div>
			</div>

			<div>
				<button
					type="submit"
					disabled={$isLoading}
					class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-primary hover:bg-primary/80 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary disabled:opacity-50 disabled:cursor-not-allowed"
				>
					{#if $isLoading}
						<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						Creating account...
					{:else}
						Create Account
					{/if}
				</button>
			</div>

			<div class="mt-4 text-center text-sm text-gray-600 dark:text-gray-400">
				<p>By creating an account, you agree to our Terms of Service and Privacy Policy.</p>
			</div>
		</form>
	</div>
</div>
