<script>
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { tasks, TASK_STATUS, PRIORITY } from '$lib/stores/unified-store';

	let { toggleSidebar, sidebarVisible = $bindable(false) } = $props();

	let isDesktop = true; 

	function handleToggleSidebar(){
		sidebarVisible = !sidebarVisible;
		toggleSidebar?.({visible: sidebarVisible});
	}

	function handleResize() {
		isDesktop = window.innerWidth >= 640;
		if (isDesktop) {
			sidebarVisible = true;
		}
	}

	let recentTasks = $derived([...$tasks]
		.sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
		.slice(0, 3));

	onMount(() => {
		handleResize();
		sidebarVisible = isDesktop;
		window.addEventListener('resize', handleResize);
		return () => window.removeEventListener('resize', handleResize);
	});
</script>

<div class="flex h-screen">	<div
		class="{sidebarVisible
			? 'translate-x-0'
			: '-translate-x-full'} sidebar fixed sm:relative sm:translate-x-0 z-30 w-64 bg-accent text-white p-4 flex flex-col dark:bg-gray-800 h-full transition-transform duration-300 ease-in-out shadow-lg"
	>
		<div class="flex items-center justify-between mb-6">
			<h2 class="text-xl font-bold">DoTask</h2>			<button aria-label="Close Sidebar" class="sm:hidden" type="button" onclick={handleToggleSidebar}>
				<svg
					class="w-6 h-6 text-white hover:text-gray-200 transition duration-150 ease-in-out"
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2"
					stroke="currentColor"
				>
					<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
				</svg>
			</button>
		</div>

		<nav class="mb-8">
			<ul class="space-y-2">
				<li>
					<a
						href="/"
						class="flex items-center p-2 rounded-lg {page.url.pathname === '/'
							? 'bg-primary dark:bg-gray-700'
							: 'hover:bg-primary/70 dark:hover:bg-gray-600'} transition-colors"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="w-5 h-5 mr-3"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
							/>
						</svg>
						Dashboard
					</a>
				</li>
				<li>
					<a
						href="/tasks"
						class="flex items-center p-2 rounded-lg {page.url.pathname.startsWith('/tasks')
							? 'bg-primary dark:bg-gray-700'
							: 'hover:bg-primary/70 dark:hover:bg-gray-600'} transition-colors"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="w-5 h-5 mr-3"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
							/>
						</svg>
						All Tasks
					</a>
				</li>
				<li>
					<a
						href="/tasks/create"
						class="flex items-center p-2 rounded-lg {page.url.pathname === '/tasks/create'
							? 'bg-primary dark:bg-gray-700'
							: 'hover:bg-primary/70 dark:hover:bg-gray-600'} transition-colors"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="w-5 h-5 mr-3"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 6v6m0 0v6m0-6h6m-6 0H6"
							/>
						</svg>
						Add Task
					</a>
				</li>
				<li>
					<a
						href="/categories"
						class="flex items-center p-2 rounded-lg {page.url.pathname.startsWith('/categories')
							? 'bg-primary dark:bg-gray-700'
							: 'hover:bg-primary/70 dark:hover:bg-gray-600'} transition-colors"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="w-5 h-5 mr-3"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
							/>
						</svg>
						Categories
					</a>
				</li>
			</ul>
		</nav>
		<div class="mt-auto">
			<h3 class="text-lg font-semibold mb-2">Recent Tasks</h3>
			<ul class="space-y-2 text-sm">
				{#if recentTasks.length === 0}
					<li class="p-2 bg-primary/20 dark:bg-gray-700 rounded-lg text-center">No tasks yet</li>
				{:else}
					{#each recentTasks as task}
						<li class="p-2 bg-primary/20 dark:bg-gray-700 rounded-lg">
							<a href="/tasks/{task.id}" class="block hover:text-primary-200 truncate">
								{task.title}
							</a>
							<div class="flex items-center text-xs mt-1">
								<span class="text-xs opacity-80">{new Date(task.dueDate).toLocaleDateString()}</span>
								<span class="ml-2 px-1.5 py-0.5 rounded-full text-xs {task.priority === PRIORITY.HIGH ? 'bg-red-500/40' : task.priority === PRIORITY.MEDIUM ? 'bg-yellow-500/40' : 'bg-green-500/40'}">
									{task.priority === PRIORITY.HIGH ? 'High' : task.priority === PRIORITY.MEDIUM ? 'Medium' : 'Low'}
								</span>
								<span class="ml-auto px-1.5 py-0.5 rounded-full text-xs {task.status === TASK_STATUS.COMPLETED ? 'bg-success/20 text-success dark:bg-success/30' : task.status === TASK_STATUS.IN_PROGRESS ? 'bg-blue-700/40 text-white dark:bg-blue-700/30 dark:text-white' : 'bg-warning/20 text-warning dark:bg-warning/30'}">
									{#if task.status === TASK_STATUS.COMPLETED}
										<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 inline-block mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
										</svg>
									{/if}
									{task.status === TASK_STATUS.TODO ? 'Todo' : task.status === TASK_STATUS.IN_PROGRESS ? 'In Progress' : 'Done'}
								</span>
							</div>
						</li>
					{/each}
				{/if}
			</ul>
		</div>
	</div>
</div>
