<script lang="ts">
	import { onMount } from 'svelte';
	import { user } from '$lib/stores/user.js';
	import '../app.css';
	import Header from '$lib/components/layout/Header.svelte';
	import Sidebar from '$lib/components/layout/Sidebar.svelte';
	import MobileNavbar from '$lib/components/layout/MobileNavbar.svelte';
	
	let sidebarVisible = false;

	type SidebarToggleEvent = {
		detail?: {
			visible: boolean;
		};
	};

	function handleSidebarToggle(event: SidebarToggleEvent) {
		sidebarVisible = event.detail?.visible ?? !sidebarVisible;
	}

	function openSidebar() {
		sidebarVisible = true;
	}

	onMount(() => {
		let username = localStorage.getItem('username');
		if (!username) {
			username = prompt('Please enter your name:');
			if (username) {
				localStorage.setItem('username', username);
			}
		}
		user.set(username || 'Guest');
		
		// Check screen width and set sidebar visibility accordingly
		const isDesktop = window.innerWidth >= 640; // sm breakpoint is 640px in Tailwind
		sidebarVisible = isDesktop;
	});
</script>

<div class="flex flex-col h-screen">
	<Header />

	<div class="flex flex-1 overflow-hidden">
		<Sidebar bind:sidebarVisible toggleSidebar={({ visible }: { visible: boolean }) => handleSidebarToggle({detail: {visible}})} />

		<main class="flex-1 overflow-auto p-6 bg-base-100 dark:bg-gray-900">
			<slot />
		</main>
	</div>
	
	<!-- Mobile navigation bar that shows at the bottom of the screen on small devices -->
	<MobileNavbar on:openSidebar={openSidebar} />
</div>
