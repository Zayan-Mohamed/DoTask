<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { authService, isAuthenticated, isLoading } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import '../app.css';
	import Header from '$lib/components/layout/Header.svelte';
	import Sidebar from '$lib/components/layout/Sidebar.svelte';
	import MobileNavbar from '$lib/components/layout/MobileNavbar.svelte';
	
	let { children } = $props();
	
	let sidebarVisible = $state(true);

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
		// Check authentication status
		authService.checkAuth();
		
		// Check screen width and set sidebar visibility accordingly
		const handleResize = () => {
			if (window.innerWidth >= 640) {
				sidebarVisible = true;
			}
		};
		handleResize();
		window.addEventListener('resize', handleResize);
		
		return () => {
			window.removeEventListener('resize', handleResize);
		};
	});

	$effect(() => {
		if (!$isLoading) {
			const isAuthPage = page.url.pathname.startsWith('/auth');
			
			if (!$isAuthenticated && !isAuthPage) {
				goto('/auth/login');
			} else if ($isAuthenticated && isAuthPage) {
				goto('/');
			}
		}
	});
</script>

{#if $isLoading}
	<div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
		<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary"></div>
	</div>
{:else if $isAuthenticated}
	<div class="flex h-screen bg-gray-100 dark:bg-gray-900">
		<Sidebar bind:sidebarVisible toggleSidebar={handleSidebarToggle} />
		
		<div class="flex-1 flex flex-col overflow-hidden">
			<Header />
			
			<main class="flex-1 overflow-x-hidden overflow-y-auto bg-gray-100 dark:bg-gray-900 p-4">
				{@render children()}
			</main>
		</div>
		
		<MobileNavbar onOpenSidebar={openSidebar} />
	</div>
{:else}
	<main class="min-h-screen">
		{@render children()}
	</main>
{/if}
