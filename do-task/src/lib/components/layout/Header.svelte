<script lang="ts">
	import { onMount } from 'svelte';
	import { user } from '$lib/stores/user';
	import { settings } from '$lib/stores/settings';

	let isDark = $state(false);

	function getDate() {
		const date = new Date();
		return date.toLocaleDateString('en-US', {
			weekday: 'long',
			year: 'numeric',
			month: 'long',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit',
			hour12: true
		});
	}

	function toggleTheme() {
		isDark = !isDark;
		document.documentElement.classList.toggle('dark');
		$settings.theme = isDark ? 'dark' : 'light';
		localStorage.setItem('theme', $settings.theme);
	}

	function deleteUsername() {
		localStorage.removeItem('username');
		location.reload();
	}

	onMount(() => {
		const savedTheme = localStorage.getItem('theme');
		const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

		isDark = savedTheme === 'dark' || (!savedTheme && prefersDark);

		if (isDark) {
			document.documentElement.classList.add('dark');
			$settings.theme = 'dark';
		} else {
			document.documentElement.classList.remove('dark');
			$settings.theme = 'light';
		}
	});

	let avatarUrl = $derived(`https://ui-avatars.com/api/?name=${encodeURIComponent($user)}&background=6366F1&color=fff&rounded=true&size=64`);
	let formattedDate = $derived(getDate());
</script>

<header>
	<div
		class="w-full flex items-center justify-between p-4 bg-primary text-white dark:bg-gray-800 dark:text-gray-200"
	>
		<div class="flex items-center">
			<svg
				width="160"
				height="40"
				viewBox="0 0 160 40"
				fill="none"
				xmlns="http://www.w3.org/2000/svg"
			>
				<!-- Box with check -->
				<rect x="2" y="8" width="24" height="24" rx="6" fill="#a21caf" />
				<path
					d="M10 20L14 24L18 16"
					stroke="white"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				/>

				<!-- App Name -->
				<text x="32" y="26" fill="currentColor" font-size="20" font-weight="bold">DoTask</text>
			</svg>
		</div>

		<div class="hidden md:flex items-center text-sm">
			<span class="mr-4">{formattedDate}</span>
		</div>

		<div class="flex items-center space-x-4">
			<button
				onclick={toggleTheme}
				class="bg-none text-white hover:text-gray-200 dark:text-gray-200 dark:hover:text-gray-400 rounded-md p-2"
				aria-label="Toggle dark mode"
			>
				{#if isDark}
					<!-- Sun icon for light mode -->
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-5 w-5"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
						/>
					</svg>
				{:else}
					<!-- Moon icon for dark mode -->
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-5 w-5"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
						/>
					</svg>
				{/if}
			</button>

			<div class="relative group">
				<button
					type="button"
					class="flex items-center space-x-2 focus:outline-none"
					aria-label="User menu"
				>
					<span class="hidden md:block">{$user}</span>
					<img src={avatarUrl} alt="Profile" class="w-8 h-8 rounded-full" />
				</button>
				<div
					class="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-md shadow-lg py-1 hidden group-hover:block z-50"
				>
					<a
						href="/profile"
						class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700"
						>Your Profile</a
					>
					<a
						href="/settings"
						class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700"
						>Settings</a
					>
					<button
						onclick={deleteUsername}
						class="w-full text-left block px-4 py-2 text-sm text-red-600 hover:bg-gray-100 dark:hover:bg-gray-700"
						>Sign out</button
					>
				</div>
			</div>
		</div>
	</div>
</header>
