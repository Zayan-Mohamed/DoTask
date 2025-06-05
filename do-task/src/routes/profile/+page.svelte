<script lang="ts">
	import { onMount } from 'svelte';
	import { user, authService, authStore } from '$lib/stores/auth';
	import { client } from '$lib/graphql/client';
	import { ME_QUERY, UPDATE_PROFILE_MUTATION, CHANGE_PASSWORD_MUTATION } from '$lib/graphql/queries';
	import type { UpdateProfileInput, ChangePasswordInput, User } from '$lib/types';
	import { goto } from '$app/navigation';

	// Profile data
	let currentUser = $state<User | null>(null);
	let isLoading = $state(false);
	let profileError = $state('');
	let profileSuccess = $state('');
	let passwordError = $state('');
	let passwordSuccess = $state('');

	// Profile editing
	let isEditingProfile = $state(false);
	let profileForm = $state<UpdateProfileInput>({
		name: '',
		email: ''
	});
	let profileValidationErrors = $state<{[key: string]: string}>({});

	// Password changing
	let isChangingPassword = $state(false);
	let passwordForm = $state<ChangePasswordInput>({
		currentPassword: '',
		newPassword: ''
	});
	let confirmPassword = $state('');
	let passwordValidationErrors = $state<{[key: string]: string}>({});
	let showCurrentPassword = $state(false);
	let showNewPassword = $state(false);
	let showConfirmPassword = $state(false);

	// Account statistics
	let accountStats = {
		tasksCreated: 0,
		tasksCompleted: 0,
		categoriesCreated: 0,
		accountAge: ''
	};

	// Profile preferences
	let preferences = $state({
		theme: 'light',
		emailNotifications: true,
		taskReminders: true,
		weeklyDigest: false
	});

	// Load user profile data
	async function loadProfile() {
		try {
			isLoading = true;
			const result = await client.query({
				query: ME_QUERY,
				fetchPolicy: 'cache-first'
			});
			
			if (result.data?.me) {
				const userData = result.data.me;
				currentUser = userData;
				profileForm = {
					name: userData.name,
					email: userData.email
				};
				
				// Calculate account age
				if (userData.createdAt) {
					const createdDate = new Date(userData.createdAt);
					const now = new Date();
					const diffTime = Math.abs(now.getTime() - createdDate.getTime());
					const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
					
					if (diffDays < 30) {
						accountStats.accountAge = `${diffDays} days`;
					} else if (diffDays < 365) {
						const months = Math.floor(diffDays / 30);
						accountStats.accountAge = `${months} ${months === 1 ? 'month' : 'months'}`;
					} else {
						const years = Math.floor(diffDays / 365);
						accountStats.accountAge = `${years} ${years === 1 ? 'year' : 'years'}`;
					}
				}
			}
		} catch (error) {
			console.error('Failed to load profile:', error);
			profileError = 'Failed to load profile. Please try again.';
		} finally {
			isLoading = false;
		}
	}

	// Validate profile form
	function validateProfileForm(): boolean {
		profileValidationErrors = {};
		
		if (!profileForm.name?.trim()) {
			profileValidationErrors.name = 'Name is required';
		} else if (profileForm.name.trim().length < 2) {
			profileValidationErrors.name = 'Name must be at least 2 characters';
		}
		
		if (!profileForm.email?.trim()) {
			profileValidationErrors.email = 'Email is required';
		} else if (!/\S+@\S+\.\S+/.test(profileForm.email)) {
			profileValidationErrors.email = 'Please enter a valid email';
		}
		
		return Object.keys(profileValidationErrors).length === 0;
	}

	// Validate password form
	function validatePasswordForm(): boolean {
		passwordValidationErrors = {};
		
		if (!passwordForm.currentPassword) {
			passwordValidationErrors.currentPassword = 'Current password is required';
		}
		
		if (!passwordForm.newPassword) {
			passwordValidationErrors.newPassword = 'New password is required';
		} else if (passwordForm.newPassword.length < 8) {
			passwordValidationErrors.newPassword = 'Password must be at least 8 characters';
		} else if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(passwordForm.newPassword)) {
			passwordValidationErrors.newPassword = 'Password must contain at least one uppercase letter, one lowercase letter, and one number';
		}
		
		if (!confirmPassword) {
			passwordValidationErrors.confirmPassword = 'Please confirm your new password';
		} else if (passwordForm.newPassword !== confirmPassword) {
			passwordValidationErrors.confirmPassword = 'Passwords do not match';
		}
		
		return Object.keys(passwordValidationErrors).length === 0;
	}

	// Update profile
	async function updateProfile(event: Event) {
		event.preventDefault();
		if (!validateProfileForm()) return;
		
		try {
			isLoading = true;
			profileError = '';
			profileSuccess = '';

			const result = await client.mutate({
				mutation: UPDATE_PROFILE_MUTATION,
				variables: {
					input: {
						name: profileForm.name?.trim(),
						email: profileForm.email?.trim()
					}
				}
			});

			if (result.data?.updateProfile) {
				currentUser = result.data.updateProfile;
				// Update auth store
				authStore.update(state => ({
					...state,
					user: currentUser
				}));
				profileSuccess = 'Profile updated successfully!';
				isEditingProfile = false;
				
				setTimeout(() => {
					profileSuccess = '';
				}, 3000);
			}
		} catch (error: any) {
			console.error('Profile update failed:', error);
			profileError = error.message || 'Failed to update profile. Please try again.';
			setTimeout(() => {
				profileError = '';
			}, 5000);
		} finally {
			isLoading = false;
		}
	}

	// Change password
	async function changePassword(event: Event) {
		event.preventDefault();
		if (!validatePasswordForm()) return;
		
		try {
			isLoading = true;
			passwordError = '';
			passwordSuccess = '';

			const result = await client.mutate({
				mutation: CHANGE_PASSWORD_MUTATION,
				variables: {
					input: {
						currentPassword: passwordForm.currentPassword,
						newPassword: passwordForm.newPassword
					}
				}
			});

			if (result.data?.changePassword) {
				passwordSuccess = 'Password changed successfully!';
				isChangingPassword = false;
				
				// Reset form
				passwordForm = {
					currentPassword: '',
					newPassword: ''
				};
				confirmPassword = '';
				
				setTimeout(() => {
					passwordSuccess = '';
				}, 3000);
			}
		} catch (error: any) {
			console.error('Password change failed:', error);
			passwordError = error.message || 'Failed to change password. Please try again.';
			setTimeout(() => {
				passwordError = '';
			}, 5000);
		} finally {
			isLoading = false;
		}
	}

	// Cancel profile editing
	function cancelProfileEdit() {
		isEditingProfile = false;
		if (currentUser) {
			profileForm = {
				name: currentUser.name,
				email: currentUser.email
			};
		}
		profileValidationErrors = {};
		profileError = '';
	}

	// Cancel password change
	function cancelPasswordChange() {
		isChangingPassword = false;
		passwordForm = {
			currentPassword: '',
			newPassword: ''
		};
		confirmPassword = '';
		passwordValidationErrors = {};
		passwordError = '';
		showCurrentPassword = false;
		showNewPassword = false;
		showConfirmPassword = false;
	}

	// Delete account (placeholder)
	async function deleteAccount() {
		if (confirm('Are you sure you want to delete your account? This action cannot be undone.')) {
			// This would require a backend implementation
			alert('Account deletion is not yet implemented. Please contact support.');
		}
	}

	// Load profile on mount
	onMount(async () => {
		if (!$user) {
			await authService.checkAuth();
			if (!$user) {
				goto('/auth/login');
				return;
			}
		}
		await loadProfile();
	});

	// Format date
	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}

	// Get avatar URL
	function getAvatarUrl(name: string): string {
		return `https://ui-avatars.com/api/?name=${encodeURIComponent(name)}&background=6366F1&color=fff&rounded=true&size=128`;
	}
</script>

<svelte:head>
	<title>Profile - DoTask</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 dark:bg-gray-900 py-8">
	<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
		<!-- Header -->
		<div class="bg-white dark:bg-gray-800 shadow rounded-lg mb-8">
			<div class="px-6 py-8">
				<div class="flex items-center space-x-6">
					{#if currentUser}
						<img 
							src={getAvatarUrl(currentUser.name)} 
							alt="Profile" 
							class="w-24 h-24 rounded-full shadow-lg"
						/>
						<div class="flex-1">
							<h1 class="text-3xl font-bold text-gray-900 dark:text-white">
								{currentUser.name}
							</h1>
							<p class="text-lg text-gray-600 dark:text-gray-300">
								{currentUser.email}
							</p>
							<p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
								Member since {formatDate(currentUser.createdAt)}
							</p>
						</div>
						<div class="text-right space-y-3">
							<div class="bg-primary/10 dark:bg-primary/20 px-4 py-2 rounded-lg">
								<p class="text-sm font-medium text-primary dark:text-primary-light">
									Account Age
								</p>
								<p class="text-lg font-bold text-primary dark:text-primary-light">
									{accountStats.accountAge}
								</p>
							</div>
							<a
								href="/settings"
								class="inline-flex items-center px-3 py-2 border border-gray-300 dark:border-gray-600 shadow-sm text-sm font-medium rounded-md text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
							>
								<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
								</svg>
								Settings
							</a>
						</div>
					{:else}
						<div class="animate-pulse flex items-center space-x-6">
							<div class="w-24 h-24 bg-gray-300 dark:bg-gray-600 rounded-full"></div>
							<div class="flex-1">
								<div class="h-8 bg-gray-300 dark:bg-gray-600 rounded w-48 mb-2"></div>
								<div class="h-6 bg-gray-300 dark:bg-gray-600 rounded w-64 mb-1"></div>
								<div class="h-4 bg-gray-300 dark:bg-gray-600 rounded w-32"></div>
							</div>
						</div>
					{/if}
				</div>
			</div>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
			<!-- Main Content -->
			<div class="lg:col-span-2 space-y-8">
				<!-- Profile Information -->
				<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
					<div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
						<div class="flex justify-between items-center">
							<h2 class="text-xl font-semibold text-gray-900 dark:text-white">
								Profile Information
							</h2>
							{#if !isEditingProfile}
								<button
									onclick={() => isEditingProfile = true}
									class="inline-flex items-center px-3 py-2 border border-transparent text-sm font-medium rounded-md text-primary bg-primary/10 hover:bg-primary/20 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
								>
									<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
									</svg>
									Edit Profile
								</button>
							{/if}
						</div>
					</div>

					<div class="px-6 py-6">
						{#if profileError}
							<div class="mb-4 bg-red-100 dark:bg-red-900/30 border-l-4 border-red-500 text-red-700 dark:text-red-300 p-4" role="alert">
								{profileError}
							</div>
						{/if}

						{#if profileSuccess}
							<div class="mb-4 bg-green-100 dark:bg-green-900/30 border-l-4 border-green-500 text-green-700 dark:text-green-300 p-4" role="alert">
								{profileSuccess}
							</div>
						{/if}

						{#if isEditingProfile}
							<form onsubmit={updateProfile} class="space-y-4">
								<div>
									<label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
										Full Name
									</label>
									<input
										id="name"
										type="text"
										bind:value={profileForm.name}
										class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
										class:border-red-500={profileValidationErrors.name}
									/>
									{#if profileValidationErrors.name}
										<p class="mt-1 text-sm text-red-600 dark:text-red-400">
											{profileValidationErrors.name}
										</p>
									{/if}
								</div>

								<div>
									<label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
										Email Address
									</label>
									<input
										id="email"
										type="email"
										bind:value={profileForm.email}
										class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
										class:border-red-500={profileValidationErrors.email}
									/>
									{#if profileValidationErrors.email}
										<p class="mt-1 text-sm text-red-600 dark:text-red-400">
											{profileValidationErrors.email}
										</p>
									{/if}
								</div>

								<div class="flex justify-end space-x-3 pt-4">
									<button
										type="button"
										onclick={cancelProfileEdit}
										class="px-4 py-2 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
									>
										Cancel
									</button>
									<button
										type="submit"
										disabled={isLoading}
										class="px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/80 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary disabled:opacity-50 disabled:cursor-not-allowed"
									>
										{isLoading ? 'Saving...' : 'Save Changes'}
									</button>
								</div>
							</form>
						{:else if currentUser}
							<div class="space-y-4">
								<div>
									<div class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
										Full Name
									</div>
									<p class="text-lg text-gray-900 dark:text-white">
										{currentUser.name}
									</p>
								</div>

								<div>
									<div class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
										Email Address
									</div>
									<p class="text-lg text-gray-900 dark:text-white">
										{currentUser.email}
									</p>
								</div>

								<div>
									<div class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
										Last Updated
									</div>
									<p class="text-lg text-gray-900 dark:text-white">
										{formatDate(currentUser.updatedAt)}
									</p>
								</div>
							</div>
						{/if}
					</div>
				</div>

				<!-- Password Change -->
				<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
					<div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
						<div class="flex justify-between items-center">
							<h2 class="text-xl font-semibold text-gray-900 dark:text-white">
								Security Settings
							</h2>
							{#if !isChangingPassword}
								<button
									onclick={() => isChangingPassword = true}
									class="inline-flex items-center px-3 py-2 border border-transparent text-sm font-medium rounded-md text-primary bg-primary/10 hover:bg-primary/20 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
								>
									<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
									</svg>
									Change Password
								</button>
							{/if}
						</div>
					</div>

					<div class="px-6 py-6">
						{#if passwordError}
							<div class="mb-4 bg-red-100 dark:bg-red-900/30 border-l-4 border-red-500 text-red-700 dark:text-red-300 p-4" role="alert">
								{passwordError}
							</div>
						{/if}

						{#if passwordSuccess}
							<div class="mb-4 bg-green-100 dark:bg-green-900/30 border-l-4 border-green-500 text-green-700 dark:text-green-300 p-4" role="alert">
								{passwordSuccess}
							</div>
						{/if}

						{#if isChangingPassword}
							<form onsubmit={changePassword} class="space-y-4">
								<div>
									<label for="currentPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
										Current Password
									</label>
									<div class="relative">
										<input
											id="currentPassword"
											type={showCurrentPassword ? 'text' : 'password'}
											bind:value={passwordForm.currentPassword}
											class="w-full px-3 py-2 pr-10 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
											class:border-red-500={passwordValidationErrors.currentPassword}
										/>
										<button
											type="button"
											class="absolute inset-y-0 right-0 pr-3 flex items-center"
											onclick={() => showCurrentPassword = !showCurrentPassword}
										>
											{#if showCurrentPassword}
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
									{#if passwordValidationErrors.currentPassword}
										<p class="mt-1 text-sm text-red-600 dark:text-red-400">
											{passwordValidationErrors.currentPassword}
										</p>
									{/if}
								</div>

								<div>
									<label for="newPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
										New Password
									</label>
									<div class="relative">
										<input
											id="newPassword"
											type={showNewPassword ? 'text' : 'password'}
											bind:value={passwordForm.newPassword}
											class="w-full px-3 py-2 pr-10 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
											class:border-red-500={passwordValidationErrors.newPassword}
										/>
										<button
											type="button"
											class="absolute inset-y-0 right-0 pr-3 flex items-center"
											onclick={() => showNewPassword = !showNewPassword}
										>
											{#if showNewPassword}
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
									{#if passwordValidationErrors.newPassword}
										<p class="mt-1 text-sm text-red-600 dark:text-red-400">
											{passwordValidationErrors.newPassword}
										</p>
									{/if}
								</div>

								<div>
									<label for="confirmPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
										Confirm New Password
									</label>
									<div class="relative">
										<input
											id="confirmPassword"
											type={showConfirmPassword ? 'text' : 'password'}
											bind:value={confirmPassword}
											class="w-full px-3 py-2 pr-10 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
											class:border-red-500={passwordValidationErrors.confirmPassword}
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
									{#if passwordValidationErrors.confirmPassword}
										<p class="mt-1 text-sm text-red-600 dark:text-red-400">
											{passwordValidationErrors.confirmPassword}
										</p>
									{/if}
								</div>

								<div class="flex justify-end space-x-3 pt-4">
									<button
										type="button"
										onclick={cancelPasswordChange}
										class="px-4 py-2 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
									>
										Cancel
									</button>
									<button
										type="submit"
										disabled={isLoading}
										class="px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/80 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary disabled:opacity-50 disabled:cursor-not-allowed"
									>
										{isLoading ? 'Changing...' : 'Change Password'}
									</button>
								</div>
							</form>
						{:else}
							<div class="space-y-4">
								<div class="flex items-center justify-between">
									<div>
										<h3 class="text-lg font-medium text-gray-900 dark:text-white">Password</h3>
										<p class="text-sm text-gray-600 dark:text-gray-400">
											Last changed on {currentUser ? formatDate(currentUser.updatedAt) : 'Unknown'}
										</p>
									</div>
									<div class="flex items-center">
										<svg class="w-5 h-5 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
										</svg>
										<span class="text-sm text-green-600 dark:text-green-400">Secure</span>
									</div>
								</div>
							</div>
						{/if}
					</div>
				</div>
			</div>

			<!-- Sidebar -->
			<div class="space-y-8">
				<!-- Account Statistics -->
				<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
					<div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
						<h2 class="text-xl font-semibold text-gray-900 dark:text-white">
							Account Statistics
						</h2>
					</div>
					<div class="px-6 py-6">
						<div class="space-y-4">
							<div class="flex justify-between items-center">
								<span class="text-sm font-medium text-gray-600 dark:text-gray-400">Tasks Created</span>
								<span class="text-lg font-bold text-primary">{accountStats.tasksCreated}</span>
							</div>
							<div class="flex justify-between items-center">
								<span class="text-sm font-medium text-gray-600 dark:text-gray-400">Tasks Completed</span>
								<span class="text-lg font-bold text-green-600">{accountStats.tasksCompleted}</span>
							</div>
							<div class="flex justify-between items-center">
								<span class="text-sm font-medium text-gray-600 dark:text-gray-400">Categories</span>
								<span class="text-lg font-bold text-blue-600">{accountStats.categoriesCreated}</span>
							</div>
							<div class="flex justify-between items-center">
								<span class="text-sm font-medium text-gray-600 dark:text-gray-400">Member Since</span>
								<span class="text-lg font-bold text-purple-600">{accountStats.accountAge}</span>
							</div>
						</div>
					</div>
				</div>

				<!-- Preferences -->
				<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
					<div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
						<h2 class="text-xl font-semibold text-gray-900 dark:text-white">
							Preferences
						</h2>
					</div>
					<div class="px-6 py-6">
						<div class="space-y-4">
							<div class="flex items-center justify-between">
								<div>
									<h3 class="text-sm font-medium text-gray-900 dark:text-white">Theme</h3>
									<p class="text-xs text-gray-600 dark:text-gray-400">Choose your preferred theme</p>
								</div>
								<select 
									bind:value={preferences.theme}
									class="px-3 py-1 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary"
								>
									<option value="light">Light</option>
									<option value="dark">Dark</option>
									<option value="auto">Auto</option>
								</select>
							</div>

							<div class="flex items-center justify-between">
								<div>
									<h3 class="text-sm font-medium text-gray-900 dark:text-white">Email Notifications</h3>
									<p class="text-xs text-gray-600 dark:text-gray-400">Receive email updates</p>
								</div>
								<label class="relative inline-flex items-center cursor-pointer">
									<input 
										type="checkbox" 
										bind:checked={preferences.emailNotifications}
										class="sr-only peer"
									/>
									<div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/25 dark:peer-focus:ring-primary/50 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-primary"></div>
								</label>
							</div>

							<div class="flex items-center justify-between">
								<div>
									<h3 class="text-sm font-medium text-gray-900 dark:text-white">Task Reminders</h3>
									<p class="text-xs text-gray-600 dark:text-gray-400">Get reminded about due tasks</p>
								</div>
								<label class="relative inline-flex items-center cursor-pointer">
									<input 
										type="checkbox" 
										bind:checked={preferences.taskReminders}
										class="sr-only peer"
									/>
									<div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/25 dark:peer-focus:ring-primary/50 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-primary"></div>
								</label>
							</div>

							<div class="flex items-center justify-between">
								<div>
									<h3 class="text-sm font-medium text-gray-900 dark:text-white">Weekly Digest</h3>
									<p class="text-xs text-gray-600 dark:text-gray-400">Weekly summary emails</p>
								</div>
								<label class="relative inline-flex items-center cursor-pointer">
									<input 
										type="checkbox" 
										bind:checked={preferences.weeklyDigest}
										class="sr-only peer"
									/>
									<div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/25 dark:peer-focus:ring-primary/50 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-primary"></div>
								</label>
							</div>
						</div>
					</div>
				</div>

				<!-- Quick Actions -->
				<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
					<div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
						<h2 class="text-xl font-semibold text-gray-900 dark:text-white">
							Quick Actions
						</h2>
					</div>
					<div class="px-6 py-6">
						<div class="space-y-3">
							<a
								href="/tasks"
								class="flex items-center p-3 text-sm font-medium text-gray-900 dark:text-white rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
							>
								<svg class="w-5 h-5 mr-3 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
								</svg>
								View All Tasks
							</a>

							<a
								href="/tasks/create"
								class="flex items-center p-3 text-sm font-medium text-gray-900 dark:text-white rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
							>
								<svg class="w-5 h-5 mr-3 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
								</svg>
								Create New Task
							</a>

							<a
								href="/categories"
								class="flex items-center p-3 text-sm font-medium text-gray-900 dark:text-white rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
							>
								<svg class="w-5 h-5 mr-3 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
								</svg>
								Manage Categories
							</a>
						</div>
					</div>
				</div>

				<!-- Danger Zone -->
				<div class="bg-white dark:bg-gray-800 shadow rounded-lg border border-red-200 dark:border-red-800">
					<div class="px-6 py-4 border-b border-red-200 dark:border-red-800">
						<h2 class="text-xl font-semibold text-red-600 dark:text-red-400">
							Danger Zone
						</h2>
					</div>
					<div class="px-6 py-6">
						<div class="space-y-4">
							<div>
								<h3 class="text-sm font-medium text-gray-900 dark:text-white mb-2">
									Delete Account
								</h3>
								<p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
									Once you delete your account, there is no going back. Please be certain.
								</p>
								<button
									onclick={deleteAccount}
									class="px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
								>
									Delete Account
								</button>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
