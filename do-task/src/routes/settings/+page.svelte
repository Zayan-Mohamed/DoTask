<script lang="ts">
	import { onMount } from 'svelte';
	import { user } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	// Settings state using Svelte 5 syntax
	let notifications = $state({
		email: true,
		push: false,
		taskReminders: true,
		weeklyDigest: false,
		marketingEmails: false
	});

	let privacy = $state({
		profileVisibility: 'friends',
		taskVisibility: 'private',
		activityStatus: true,
		dataCollection: false
	});

	let preferences = $state({
		theme: 'system',
		language: 'en',
		timezone: 'UTC',
		dateFormat: 'MM/DD/YYYY',
		startWeekOn: 'monday'
	});

	let account = $state({
		twoFactorEnabled: false,
		sessionTimeout: 30,
		downloadData: false,
		deleteAccount: false
	});

	let activeTab = $state('notifications');
	let isLoading = $state(false);
	let message = $state('');
	let messageType = $state<'success' | 'error'>('success');

	onMount(async () => {
		if (!$user) {
			goto('/auth/login');
		}
		// Load user settings from backend
		await loadSettings();
	});

	async function loadSettings() {
		isLoading = true;
		try {
			// TODO: Implement actual API call to load user settings
			// For now, use default values
			console.log('Loading user settings...');
		} catch (error) {
			console.error('Failed to load settings:', error);
		} finally {
			isLoading = false;
		}
	}

	async function saveSettings() {
		isLoading = true;
		message = '';
		
		try {
			// TODO: Implement actual API call to save settings
			await new Promise(resolve => setTimeout(resolve, 1000)); // Simulate API call
			
			message = 'Settings saved successfully!';
			messageType = 'success';
			setTimeout(() => message = '', 3000);
		} catch (error) {
			console.error('Failed to save settings:', error);
			message = 'Failed to save settings. Please try again.';
			messageType = 'error';
		} finally {
			isLoading = false;
		}
	}

	function handleTabChange(tab: string) {
		activeTab = tab;
	}
</script>

<svelte:head>
	<title>Settings - DoTask</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 dark:bg-gray-900 py-8">
	<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
		<!-- Header -->
		<div class="bg-white dark:bg-gray-800 shadow rounded-lg mb-8">
			<div class="px-6 py-8">
				<div class="flex items-center justify-between">
					<div>
						<h1 class="text-3xl font-bold text-gray-900 dark:text-white">
							Settings
						</h1>
						<p class="text-lg text-gray-600 dark:text-gray-300 mt-2">
							Manage your application preferences and account settings
						</p>
					</div>
					<a
						href="/profile"
						class="inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 shadow-sm text-sm font-medium rounded-md text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
					>
						<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
						</svg>
						Profile
					</a>
				</div>
			</div>
		</div>

		<!-- Message -->
		{#if message}
			<div class="mb-6 rounded-md p-4 {messageType === 'success' ? 'bg-green-50 dark:bg-green-900' : 'bg-red-50 dark:bg-red-900'}">
				<div class="flex">
					<div class="flex-shrink-0">
						{#if messageType === 'success'}
							<svg class="h-5 w-5 text-green-400" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
							</svg>
						{:else}
							<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
							</svg>
						{/if}
					</div>
					<div class="ml-3">
						<p class="text-sm font-medium {messageType === 'success' ? 'text-green-800 dark:text-green-200' : 'text-red-800 dark:text-red-200'}">
							{message}
						</p>
					</div>
				</div>
			</div>
		{/if}

		<!-- Main Content -->
		<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
			<div class="flex">
				<!-- Sidebar -->
				<div class="w-1/4 border-r border-gray-200 dark:border-gray-700">
					<nav class="space-y-1 p-6">
						<button
							class="w-full text-left px-3 py-2 rounded-md text-sm font-medium transition-colors {activeTab === 'notifications' ? 'bg-primary text-white' : 'text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700'}"
							onclick={() => handleTabChange('notifications')}
						>
							<svg class="w-4 h-4 inline mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5 5-5-5h5V3h0z" />
							</svg>
							Notifications
						</button>
						<button
							class="w-full text-left px-3 py-2 rounded-md text-sm font-medium transition-colors {activeTab === 'privacy' ? 'bg-primary text-white' : 'text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700'}"
							onclick={() => handleTabChange('privacy')}
						>
							<svg class="w-4 h-4 inline mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
							</svg>
							Privacy
						</button>
						<button
							class="w-full text-left px-3 py-2 rounded-md text-sm font-medium transition-colors {activeTab === 'preferences' ? 'bg-primary text-white' : 'text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700'}"
							onclick={() => handleTabChange('preferences')}
						>
							<svg class="w-4 h-4 inline mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4" />
							</svg>
							Preferences
						</button>
						<button
							class="w-full text-left px-3 py-2 rounded-md text-sm font-medium transition-colors {activeTab === 'account' ? 'bg-primary text-white' : 'text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700'}"
							onclick={() => handleTabChange('account')}
						>
							<svg class="w-4 h-4 inline mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
							</svg>
							Account
						</button>
					</nav>
				</div>

				<!-- Content -->
				<div class="flex-1 p-6">
					{#if activeTab === 'notifications'}
						<div class="space-y-6">
							<div>
								<h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Notification Preferences</h3>
								<div class="space-y-4">
									<div class="flex items-center justify-between">
										<div>
											<div class="text-sm font-medium text-gray-900 dark:text-white">Email Notifications</div>
											<div class="text-sm text-gray-500 dark:text-gray-400">Receive notifications via email</div>
										</div>
										<input
											type="checkbox"
											bind:checked={notifications.email}
											class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded"
										/>
									</div>
									<div class="flex items-center justify-between">
										<div>
											<div class="text-sm font-medium text-gray-900 dark:text-white">Push Notifications</div>
											<div class="text-sm text-gray-500 dark:text-gray-400">Receive push notifications in your browser</div>
										</div>
										<input
											type="checkbox"
											bind:checked={notifications.push}
											class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded"
										/>
									</div>
									<div class="flex items-center justify-between">
										<div>
											<div class="text-sm font-medium text-gray-900 dark:text-white">Task Reminders</div>
											<div class="text-sm text-gray-500 dark:text-gray-400">Get reminded about upcoming tasks</div>
										</div>
										<input
											type="checkbox"
											bind:checked={notifications.taskReminders}
											class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded"
										/>
									</div>
									<div class="flex items-center justify-between">
										<div>
											<div class="text-sm font-medium text-gray-900 dark:text-white">Weekly Digest</div>
											<div class="text-sm text-gray-500 dark:text-gray-400">Weekly summary of your activity</div>
										</div>
										<input
											type="checkbox"
											bind:checked={notifications.weeklyDigest}
											class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded"
										/>
									</div>
									<div class="flex items-center justify-between">
										<div>
											<div class="text-sm font-medium text-gray-900 dark:text-white">Marketing Emails</div>
											<div class="text-sm text-gray-500 dark:text-gray-400">Receive updates about new features</div>
										</div>
										<input
											type="checkbox"
											bind:checked={notifications.marketingEmails}
											class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded"
										/>
									</div>
								</div>
							</div>
						</div>
					{:else if activeTab === 'privacy'}
						<div class="space-y-6">
							<div>
								<h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Privacy Settings</h3>
								<div class="space-y-4">
									<div>
										<label for="profileVisibility" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Profile Visibility</label>
										<select
											id="profileVisibility"
											bind:value={privacy.profileVisibility}
											class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-primary focus:border-primary sm:text-sm rounded-md dark:bg-gray-700 dark:text-white"
										>
											<option value="public">Public</option>
											<option value="friends">Friends Only</option>
											<option value="private">Private</option>
										</select>
									</div>
									<div>
										<label for="taskVisibility" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Task Visibility</label>
										<select
											id="taskVisibility"
											bind:value={privacy.taskVisibility}
											class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-primary focus:border-primary sm:text-sm rounded-md dark:bg-gray-700 dark:text-white"
										>
											<option value="public">Public</option>
											<option value="friends">Friends Only</option>
											<option value="private">Private</option>
										</select>
									</div>
									<div class="flex items-center justify-between">
										<div>
											<div class="text-sm font-medium text-gray-900 dark:text-white">Show Activity Status</div>
											<div class="text-sm text-gray-500 dark:text-gray-400">Let others see when you're active</div>
										</div>
										<input
											type="checkbox"
											bind:checked={privacy.activityStatus}
											class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded"
										/>
									</div>
									<div class="flex items-center justify-between">
										<div>
											<div class="text-sm font-medium text-gray-900 dark:text-white">Data Collection</div>
											<div class="text-sm text-gray-500 dark:text-gray-400">Allow analytics and usage data collection</div>
										</div>
										<input
											type="checkbox"
											bind:checked={privacy.dataCollection}
											class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded"
										/>
									</div>
								</div>
							</div>
						</div>
					{:else if activeTab === 'preferences'}
						<div class="space-y-6">
							<div>
								<h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Application Preferences</h3>
								<div class="space-y-4">
									<div>
										<label for="theme" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Theme</label>
										<select
											id="theme"
											bind:value={preferences.theme}
											class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-primary focus:border-primary sm:text-sm rounded-md dark:bg-gray-700 dark:text-white"
										>
											<option value="light">Light</option>
											<option value="dark">Dark</option>
											<option value="system">System</option>
										</select>
									</div>
									<div>
										<label for="language" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Language</label>
										<select
											id="language"
											bind:value={preferences.language}
											class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-primary focus:border-primary sm:text-sm rounded-md dark:bg-gray-700 dark:text-white"
										>
											<option value="en">English</option>
											<option value="es">Spanish</option>
											<option value="fr">French</option>
											<option value="de">German</option>
										</select>
									</div>
									<div>
										<label for="timezone" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Timezone</label>
										<select
											id="timezone"
											bind:value={preferences.timezone}
											class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-primary focus:border-primary sm:text-sm rounded-md dark:bg-gray-700 dark:text-white"
										>
											<option value="UTC">UTC</option>
											<option value="America/New_York">Eastern Time</option>
											<option value="America/Chicago">Central Time</option>
											<option value="America/Denver">Mountain Time</option>
											<option value="America/Los_Angeles">Pacific Time</option>
										</select>
									</div>
									<div>
										<label for="dateFormat" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Date Format</label>
										<select
											id="dateFormat"
											bind:value={preferences.dateFormat}
											class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-primary focus:border-primary sm:text-sm rounded-md dark:bg-gray-700 dark:text-white"
										>
											<option value="MM/DD/YYYY">MM/DD/YYYY</option>
											<option value="DD/MM/YYYY">DD/MM/YYYY</option>
											<option value="YYYY-MM-DD">YYYY-MM-DD</option>
										</select>
									</div>
									<div>
										<label for="startWeekOn" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Start Week On</label>
										<select
											id="startWeekOn"
											bind:value={preferences.startWeekOn}
											class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-primary focus:border-primary sm:text-sm rounded-md dark:bg-gray-700 dark:text-white"
										>
											<option value="sunday">Sunday</option>
											<option value="monday">Monday</option>
										</select>
									</div>
								</div>
							</div>
						</div>
					{:else if activeTab === 'account'}
						<div class="space-y-6">
							<div>
								<h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Account Management</h3>
								<div class="space-y-4">
									<div class="flex items-center justify-between">
										<div>
											<div class="text-sm font-medium text-gray-900 dark:text-white">Two-Factor Authentication</div>
											<div class="text-sm text-gray-500 dark:text-gray-400">Add an extra layer of security to your account</div>
										</div>
										<input
											type="checkbox"
											bind:checked={account.twoFactorEnabled}
											class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded"
										/>
									</div>
									<div>
										<label for="sessionTimeout" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Session Timeout (minutes)</label>
										<input
											type="number"
											id="sessionTimeout"
											bind:value={account.sessionTimeout}
											min="5"
											max="480"
											class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-primary focus:border-primary sm:text-sm dark:bg-gray-700 dark:text-white"
										/>
									</div>
									<div class="border-t border-gray-200 dark:border-gray-700 pt-4">
										<h4 class="text-sm font-medium text-gray-900 dark:text-white mb-2">Data Management</h4>
										<div class="space-y-2">
											<button
												type="button"
												class="text-sm text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300"
												onclick={() => alert('Data download functionality will be implemented')}
											>
												Download My Data
											</button>
											<button
												type="button"
												class="block text-sm text-red-600 dark:text-red-400 hover:text-red-500 dark:hover:text-red-300"
												onclick={() => alert('Account deletion functionality will be implemented')}
											>
												Delete Account
											</button>
										</div>
									</div>
								</div>
							</div>
						</div>
					{/if}

					<!-- Save Button -->
					<div class="flex justify-end pt-6 border-t border-gray-200 dark:border-gray-700">
						<button
							type="button"
							onclick={saveSettings}
							disabled={isLoading}
							class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary hover:bg-primary/80 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary disabled:opacity-50 disabled:cursor-not-allowed"
						>
							{#if isLoading}
								<svg class="animate-spin -ml-1 mr-3 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
									<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
									<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
								</svg>
								Saving...
							{:else}
								Save Changes
							{/if}
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
