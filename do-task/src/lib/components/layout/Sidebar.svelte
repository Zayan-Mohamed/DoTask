<script lang="ts">
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { tasks, TASK_STATUS, PRIORITY, updateTaskStatus, loadTasks, loadCategories } from '$lib/stores/unified-store';
	import type { Task, TaskStatus, TaskPriority } from '$lib/types';
	import { staggerReveal, ripple } from '$lib/utils/animations';

	let { toggleSidebar, sidebarVisible = $bindable(false) } = $props();

	let isDesktop = true; 
	let currentHour = $state(new Date().getHours());

	// Enhanced Task interface for Context7 features
	interface ContextualTask extends Task {
		contextScore: number;
		isDueToday: boolean;
		isDueTomorrow: boolean;
		isOverdue: boolean;
	}

	interface ContextualInsight {
		type: 'warning' | 'info' | 'positive' | 'neutral' | 'evening';
		message: string;
		icon: string;
	}

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

	// Task status toggle functionality
	function getNextStatus(currentStatus: TaskStatus): TaskStatus {
		switch (currentStatus) {
			case TASK_STATUS.TODO:
				return TASK_STATUS.IN_PROGRESS;
			case TASK_STATUS.IN_PROGRESS:
				return TASK_STATUS.COMPLETED;
			case TASK_STATUS.COMPLETED:
				return TASK_STATUS.TODO;
			default:
				return TASK_STATUS.TODO;
		}
	}
	
	async function handleStatusToggle(event: Event, task: ContextualTask) {
		event.preventDefault();
		event.stopPropagation();
		try {
			const nextStatus = getNextStatus(task.status);
			await updateTaskStatus(task.id, nextStatus);
			await Promise.all([loadTasks(), loadCategories()]);
		} catch (err) {
			console.error('Failed to update task status:', err);
		}
	}

	// Enhanced Context7 integration for smart task prioritization
	let contextualTasks = $derived.by((): ContextualTask[] => {
		const now = new Date();
		const today = now.toDateString();
		const currentHour = now.getHours();
		
		// Time-based context awareness
		const isWorkingHours = currentHour >= 9 && currentHour <= 17;
		const isEvening = currentHour >= 18 && currentHour <= 22;
		const isMorning = currentHour >= 6 && currentHour <= 11;
		
		// Ensure tasks is not undefined
		if (!$tasks || !Array.isArray($tasks)) {
			return [];
		}
		
		return [...$tasks]
			.filter(task => task.status !== TASK_STATUS.COMPLETED)
			.map((task): ContextualTask => {
				// Safe date parsing
				let dueDate: Date;
				try {
					dueDate = new Date(task.dueDate);
					if (isNaN(dueDate.getTime())) {
						dueDate = new Date();
					}
				} catch {
					dueDate = new Date();
				}
				
				const isDueToday = dueDate.toDateString() === today;
				const isDueTomorrow = dueDate.toDateString() === new Date(now.getTime() + 86400000).toDateString();
				const isOverdue = dueDate < now && task.status !== TASK_STATUS.COMPLETED;
				
				// Context7 smart scoring algorithm
				let contextScore = 0;
				
				// Urgency factors
				if (isOverdue) contextScore += 100;
				if (isDueToday) contextScore += 50;
				if (isDueTomorrow) contextScore += 25;
				
				// Priority weighting
				if (task.priority === PRIORITY.HIGH) contextScore += 30;
				if (task.priority === PRIORITY.MEDIUM) contextScore += 15;
				
				// Time-based context adjustments
				const taskTags = task.tags || [];
				if (isWorkingHours && Array.isArray(taskTags) && taskTags.some(tag => 
					['work', 'office', 'meeting', 'business'].includes(tag.toLowerCase())
				)) contextScore += 20;
				
				if (isEvening && Array.isArray(taskTags) && taskTags.some(tag => 
					['personal', 'home', 'family', 'hobby'].includes(tag.toLowerCase())
				)) contextScore += 15;
				
				if (isMorning && Array.isArray(taskTags) && taskTags.some(tag => 
					['planning', 'review', 'standup', 'daily'].includes(tag.toLowerCase())
				)) contextScore += 15;
				
				// Recent activity boost
				const lastUpdated = new Date(task.updatedAt);
				const hoursSinceUpdate = (now.getTime() - lastUpdated.getTime()) / (1000 * 60 * 60);
				if (hoursSinceUpdate < 2) contextScore += 10;
				
				return { ...task, contextScore, isDueToday, isDueTomorrow, isOverdue };
			})
			.sort((a, b) => b.contextScore - a.contextScore)
			.slice(0, 3);
	});

	// Get contextual greeting and insights
	let contextualInsight = $derived.by((): ContextualInsight => {
		// Ensure tasks is not undefined
		if (!$tasks || !Array.isArray($tasks)) {
			return { 
				type: 'neutral',
				message: 'Loading tasks...',
				icon: '⏳'
			};
		}
		
		const overdueTasks = $tasks.filter(task => 
			new Date(task.dueDate) < new Date() && task.status !== TASK_STATUS.COMPLETED
		).length;
		
		const todayTasks = $tasks.filter(task => 
			new Date(task.dueDate).toDateString() === new Date().toDateString() && 
			task.status !== TASK_STATUS.COMPLETED
		).length;
		
		if (overdueTasks > 0) {
			return { 
				type: 'warning',
				message: `${overdueTasks} overdue task${overdueTasks > 1 ? 's' : ''}`,
				icon: '⚠️'
			};
		}
		
		if (todayTasks > 0) {
			return { 
				type: 'info',
				message: `${todayTasks} task${todayTasks > 1 ? 's' : ''} due today`,
				icon: '📅'
			};
		}
		
		if (currentHour >= 6 && currentHour < 12) {
			return { 
				type: 'positive',
				message: 'Good morning! Ready to tackle the day?',
				icon: '🌅'
			};
		} else if (currentHour >= 12 && currentHour < 17) {
			return { 
				type: 'neutral',
				message: 'Keep up the great work!',
				icon: '💪'
			};
		} else {
			return { 
				type: 'evening',
				message: 'Wrapping up for the day?',
				icon: '🌙'
			};
		}
	});

	onMount(() => {
		handleResize();
		sidebarVisible = isDesktop;
		window.addEventListener('resize', handleResize);
		
		// Update time-based context every minute
		const timeInterval = setInterval(() => {
			currentHour = new Date().getHours();
		}, 60000);
		
		return () => {
			window.removeEventListener('resize', handleResize);
			clearInterval(timeInterval);
		};
	});
</script>

<div class="flex h-screen">	<div
		class="{sidebarVisible
			? 'translate-x-0'
			: '-translate-x-full'} sidebar fixed sm:relative sm:translate-x-0 z-30 w-auto bg-accent text-white flex flex-col dark:bg-gray-800 h-full transition-transform duration-300 ease-in-out shadow-lg"
	>
		<div class="flex items-center justify-between mb-6 p-4 flex-shrink-0">
			<h2 class="text-xl font-bold">DoTask</h2>			
			<button aria-label="Close Sidebar" class="sm:hidden" type="button" onclick={handleToggleSidebar}>
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

		<!-- Scrollable content area -->
		<div class="flex-1 overflow-y-auto px-4">
			<nav class="mb-8">
			<ul class="space-y-2">
				<li>
					<a
						href="/"
						use:ripple
						class="flex items-center p-3 rounded-lg group transition-all duration-200 
							   {page.url.pathname === '/'
								? 'bg-primary dark:bg-gray-700 shadow-lg transform scale-[1.02]'
								: 'hover:bg-primary/70 dark:hover:bg-gray-600 hover:transform hover:scale-[1.02]'}"
					>
						<div class="flex items-center justify-center w-5 h-5 mr-3 transition-transform duration-200 group-hover:scale-110">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="w-5 h-5"
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
						</div>
						<span class="font-medium">Dashboard</span>
					</a>
				</li>
				<li>
					<a
						href="/tasks"
						use:ripple
						class="flex items-center p-3 rounded-lg group transition-all duration-200
							   {page.url.pathname.startsWith('/tasks')
								? 'bg-primary dark:bg-gray-700 shadow-lg transform scale-[1.02]'
								: 'hover:bg-primary/70 dark:hover:bg-gray-600 hover:transform hover:scale-[1.02]'}"
					>
						<div class="flex items-center justify-center w-5 h-5 mr-3 transition-transform duration-200 group-hover:scale-110">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="w-5 h-5"
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
						</div>
						<span class="font-medium">All Tasks</span>
					</a>
				</li>
				<li>
					<a
						href="/tasks/create"
						use:ripple
						class="flex items-center p-3 rounded-lg group transition-all duration-200
							   {page.url.pathname === '/tasks/create'
								? 'bg-primary dark:bg-gray-700 shadow-lg transform scale-[1.02]'
								: 'hover:bg-primary/70 dark:hover:bg-gray-600 hover:transform hover:scale-[1.02]'}"
					>
						<div class="flex items-center justify-center w-5 h-5 mr-3 transition-transform duration-200 group-hover:scale-110 group-hover:rotate-90">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="w-5 h-5"
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
						</div>
						<span class="font-medium">Add Task</span>
					</a>
				</li>
				<li>
					<a
						href="/categories"
						use:ripple
						class="flex items-center p-3 rounded-lg group transition-all duration-200
							   {page.url.pathname.startsWith('/categories')
								? 'bg-primary dark:bg-gray-700 shadow-lg transform scale-[1.02]'
								: 'hover:bg-primary/70 dark:hover:bg-gray-600 hover:transform hover:scale-[1.02]'}"
					>
						<div class="flex items-center justify-center w-5 h-5 mr-3 transition-transform duration-200 group-hover:scale-110">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="w-5 h-5"
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
						</div>
						<span class="font-medium">Categories</span>
					</a>
				</li>
			</ul>
		</nav>
		
		<!-- Context7 Smart Insights -->
		<div class="mb-4 p-3 rounded-lg bg-white/5 border border-white/10">
			<div class="flex items-center gap-2 text-sm">
				<span class="text-lg">{contextualInsight.icon}</span>
				<span class="text-white/80 font-medium">{contextualInsight.message}</span>
			</div>
			{#if contextualInsight.type === 'warning'}
				<div class="mt-2 w-full bg-red-900/30 rounded-full h-1">
					<div class="bg-red-400 h-1 rounded-full w-full animate-pulse"></div>
				</div>
			{:else if contextualInsight.type === 'info'}
				<div class="mt-2 w-full bg-blue-900/30 rounded-full h-1">
					<div class="bg-blue-400 h-1 rounded-full w-2/3"></div>
				</div>
			{/if}
		</div>

		<div class="flex items-center justify-between mb-3">
			<h3 class="text-lg font-semibold text-white/90 flex items-center gap-2">
				<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-primary-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
				</svg>
				Smart Focus
			</h3>
			<a 
				href="/tasks" 
				class="text-xs text-white/70 hover:text-white/90 transition-colors duration-200 hover:underline"
			>
				View all
			</a>
		</div>

		<div class="space-y-2 pb-4">
			{#if contextualTasks.length === 0}
				<div class="p-4 bg-white/5 dark:bg-gray-800/50 rounded-lg border border-white/10 text-center">
					<div class="flex flex-col items-center gap-2">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 text-white/40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
						</svg>
						<span class="text-sm text-white/60">All caught up!</span>
						<span class="text-xs text-white/50">No urgent tasks right now</span>
						<a 
							href="/tasks/create" 
							class="text-xs text-primary-200 hover:text-primary-100 transition-colors duration-200"
						>
							Create a new task
						</a>
					</div>
				</div>
			{:else}
				{#each contextualTasks as task, index}
					<div 
						in:staggerReveal={{ delay: index * 100 }}
						class="group relative bg-white/5 dark:bg-gray-800/30 rounded-lg border border-white/10 
							   hover:bg-white/10 dark:hover:bg-gray-700/50 hover:border-white/20 
							   transition-all duration-300 ease-out transform hover:scale-[1.02] 
							   hover:shadow-lg hover:shadow-black/20
							   {task.isOverdue ? 'border-red-400/50 bg-red-900/20' : ''}
							   {task.isDueToday ? 'border-yellow-400/50 bg-yellow-900/10' : ''}"
					>
						<!-- Task Content -->
						<div class="flex items-start gap-3 p-3">
							<!-- Status Toggle Button -->
							<button 
								onclick={(e) => handleStatusToggle(e, task)}
								class="flex-shrink-0 mt-1 w-6 h-6 flex items-center justify-center rounded-full border-2 transition-all duration-200 hover:scale-110
									{task.status === TASK_STATUS.COMPLETED 
										? 'bg-green-500/20 border-green-500/50 text-green-300' 
										: task.status === TASK_STATUS.IN_PROGRESS 
											? 'bg-blue-500/20 border-blue-500/50 text-blue-300' 
											: 'bg-gray-500/20 border-gray-500/50 text-gray-300'
									} hover:bg-white/20"
								title={`Mark as ${getNextStatus(task.status)}`}
							>
								{#if task.status === TASK_STATUS.COMPLETED}
									<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
									</svg>
								{:else if task.status === TASK_STATUS.IN_PROGRESS}
									<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
									</svg>
								{:else}
									<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
									</svg>
								{/if}
							</button>

							<!-- Task Details (Clickable for navigation) -->
							<a 
								href="/tasks/{task.id}" 
								use:ripple={{ color: 'rgba(255, 255, 255, 0.1)' }}
								class="flex-1 block rounded-lg transition-colors duration-200"
							>
								<!-- Context7 Enhanced Header -->
								<div class="flex items-start justify-between gap-2 mb-2">
									<div class="flex items-start gap-2 flex-1">
										<!-- Context Priority Indicator -->
										<div class="flex-shrink-0 mt-0.5">
											{#if task.isOverdue}
												<div class="w-2 h-2 bg-red-400 rounded-full animate-pulse" title="Overdue"></div>
											{:else if task.isDueToday}
												<div class="w-2 h-2 bg-yellow-400 rounded-full animate-pulse" title="Due Today"></div>
											{:else if task.isDueTomorrow}
												<div class="w-2 h-2 bg-blue-400 rounded-full animate-pulse" title="Due Tomorrow"></div>
											{:else if task.status === TASK_STATUS.COMPLETED}
												<div class="w-2 h-2 bg-green-400 rounded-full" title="Completed"></div>
											{:else if task.status === TASK_STATUS.IN_PROGRESS}
												<div class="w-2 h-2 bg-blue-400 rounded-full animate-pulse" title="In Progress"></div>
											{:else}
												<div class="w-2 h-2 bg-gray-400 rounded-full" title="To Do"></div>
											{/if}
										</div>
										
										<h4 class="font-medium text-white/90 group-hover:text-white line-clamp-2 text-sm leading-tight flex-1">
											{task.title}
										</h4>
									</div>
									
									<!-- Context Score Badge -->
									{#if task.contextScore > 50}
										<div class="flex-shrink-0">
											<span class="px-1.5 py-0.5 rounded text-xs font-bold
												{task.contextScore > 80 ? 'bg-red-500/30 text-red-200' : 
												task.contextScore > 60 ? 'bg-yellow-500/30 text-yellow-200' : 
												'bg-blue-500/30 text-blue-200'}">
												{task.isOverdue ? '🚨' : task.isDueToday ? '⏰' : '⚡'}
											</span>
										</div>
									{/if}
								</div>
								
								<!-- Enhanced Task Meta Information -->
								<div class="flex items-center justify-between text-xs">
									<div class="flex items-center gap-2">
										<!-- Smart Due Date Display -->
										<span class="text-white/60 flex items-center gap-1 {task.isOverdue ? 'text-red-300' : task.isDueToday ? 'text-yellow-300' : ''}">
											<svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
											</svg>
											{#if task.isOverdue}
												Overdue
											{:else if task.isDueToday}
												Today
											{:else if task.isDueTomorrow}
												Tomorrow
											{:else}
												{new Date(task.dueDate).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
											{/if}
										</span>
										
										<!-- Enhanced Priority Badge -->
										<span class="px-1.5 py-0.5 rounded-full text-xs font-medium
											{task.priority === PRIORITY.HIGH 
												? 'bg-red-500/20 text-red-300 border border-red-500/30' 
												: task.priority === PRIORITY.MEDIUM 
													? 'bg-yellow-500/20 text-yellow-300 border border-yellow-500/30' 
													: 'bg-green-500/20 text-green-300 border border-green-500/30'
											}">
											{task.priority === PRIORITY.HIGH ? 'High' : task.priority === PRIORITY.MEDIUM ? 'Medium' : 'Low'}
										</span>
									</div>
									
									<!-- Enhanced Status Badge with Context -->
									<span class="px-2 py-0.5 rounded-full text-xs font-medium flex items-center gap-1
										{task.status === TASK_STATUS.COMPLETED 
											? 'bg-green-500/20 text-green-300 border border-green-500/30' 
											: task.status === TASK_STATUS.IN_PROGRESS 
												? 'bg-blue-500/20 text-blue-300 border border-blue-500/30' 
												: 'bg-gray-500/20 text-gray-300 border border-gray-500/30'
										}">
										{#if task.status === TASK_STATUS.COMPLETED}
											<svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
											</svg>
										{:else if task.status === TASK_STATUS.IN_PROGRESS}
											<svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
											</svg>
										{:else}
											<svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
											</svg>
										{/if}
										{task.status === TASK_STATUS.TODO ? 'Todo' : task.status === TASK_STATUS.IN_PROGRESS ? 'Active' : 'Done'}
									</span>
								</div>
							</a>
						</div>
						
						<!-- Enhanced Hover Gradient Effect with Context -->
						<div class="absolute inset-0 rounded-lg pointer-events-none transition-opacity duration-300
							{task.isOverdue ? 'bg-gradient-to-r from-red-500/10 to-transparent opacity-0 group-hover:opacity-100' :
							task.isDueToday ? 'bg-gradient-to-r from-yellow-500/10 to-transparent opacity-0 group-hover:opacity-100' :
							'bg-gradient-to-r from-primary/5 to-transparent opacity-0 group-hover:opacity-100'}"></div>
					</div>
				{/each}
			{/if}
		</div>
	</div>
	</div>
</div>

<style>
	/* Enhanced text truncation for task titles */
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}

	/* Smooth animations for status indicators */
	@keyframes pulse-glow {
		0%, 100% {
			box-shadow: 0 0 0 0 rgba(59, 130, 246, 0.7);
		}
		50% {
			box-shadow: 0 0 0 4px rgba(59, 130, 246, 0);
		}
	}

	/* Custom scrollbar for sidebar content */
	.sidebar .overflow-y-auto::-webkit-scrollbar {
		width: 4px;
	}

	.sidebar .overflow-y-auto::-webkit-scrollbar-track {
		background: rgba(255, 255, 255, 0.1);
		border-radius: 2px;
	}

	.sidebar .overflow-y-auto::-webkit-scrollbar-thumb {
		background: rgba(255, 255, 255, 0.3);
		border-radius: 2px;
	}

	.sidebar .overflow-y-auto::-webkit-scrollbar-thumb:hover {
		background: rgba(255, 255, 255, 0.5);
	}

	/* Custom scrollbar for sidebar */
	.sidebar::-webkit-scrollbar {
		width: 4px;
	}

	.sidebar::-webkit-scrollbar-track {
		background: rgba(255, 255, 255, 0.1);
		border-radius: 2px;
	}

	.sidebar::-webkit-scrollbar-thumb {
		background: rgba(255, 255, 255, 0.3);
		border-radius: 2px;
	}

	.sidebar::-webkit-scrollbar-thumb:hover {
		background: rgba(255, 255, 255, 0.5);
	}

	/* Backdrop blur effect for enhanced glassmorphism */
	.sidebar {
		backdrop-filter: blur(10px);
		-webkit-backdrop-filter: blur(10px);
	}

	/* Enhanced hover effects for navigation items */
	nav a:hover {
		background: linear-gradient(135deg, rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0.05));
		backdrop-filter: blur(10px);
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	/* Task card hover animation */
	.group:hover {
		box-shadow: 
			0 10px 25px -5px rgba(0, 0, 0, 0.1),
			0 10px 10px -5px rgba(0, 0, 0, 0.04),
			inset 0 1px 0 rgba(255, 255, 255, 0.1);
	}

	/* Status indicator pulse animation */
	.animate-pulse {
		animation: pulse-glow 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
	}

	/* Smooth gradient overlay */
	.group::before {
		content: '';
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: linear-gradient(135deg, transparent, rgba(255, 255, 255, 0.02));
		border-radius: inherit;
		opacity: 0;
		transition: opacity 300ms ease;
		pointer-events: none;
	}

	.group:hover::before {
		opacity: 1;
	}
</style>
