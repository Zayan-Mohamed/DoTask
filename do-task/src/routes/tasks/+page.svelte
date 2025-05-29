<script lang="ts">
    import { tasks, categories, unifiedStore, TASK_STATUS, PRIORITY, loadTasks, loadCategories } from '$lib/stores/unified-store';
    import type { TaskStatus } from '$lib/types';
    import { onMount } from 'svelte';

    // Convert these to $state if you want to make them reactive state variables
    let searchQuery = $state('');
    let statusFilter = $state('all');
    let priorityFilter = $state('all');
    let sortBy = $state('dueDate');
    let sortDirection = $state('asc');
    
    // Load tasks and categories when component mounts
    onMount(async () => {
        await unifiedStore.initializeData();
    });

    // Replace the reactive $: with $derived
    let filteredTasks = $derived(($tasks || [])
        .filter((task) => {
            // Text search
            const matchesSearch =
                searchQuery === '' ||
                task.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
                task.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
                (task.tags && task.tags.some((tag) => tag.toLowerCase().includes(searchQuery.toLowerCase())));

            // Status filter
            const matchesStatus = statusFilter === 'all' || task.status === statusFilter;

            // Priority filter
            const matchesPriority = priorityFilter === 'all' || task.priority === priorityFilter;

            return matchesSearch && matchesStatus && matchesPriority;
        })
        .sort((a, b) => {
            // Sort by selected field
            if (sortBy === 'dueDate') {
                const dateA = new Date(a.dueDate);
                const dateB = new Date(b.dueDate);
                const timeA = isNaN(dateA.getTime()) ? 0 : dateA.getTime();
                const timeB = isNaN(dateB.getTime()) ? 0 : dateB.getTime();
                return sortDirection === 'asc' ? timeA - timeB : timeB - timeA;
            } else if (sortBy === 'priority') {
                const priorityValues = {
                    [PRIORITY.HIGH]: 3,
                    [PRIORITY.MEDIUM]: 2,
                    [PRIORITY.LOW]: 1
                };
                return sortDirection === 'asc'
                    ? priorityValues[a.priority] - priorityValues[b.priority]
                    : priorityValues[b.priority] - priorityValues[a.priority];
            } else if (sortBy === 'title') {
                return sortDirection === 'asc'
                    ? a.title.localeCompare(b.title)
                    : b.title.localeCompare(a.title);
            }
            return 0;
        }));

    // Format date to human-readable string
    function formatDate(dateString: string) {
        try {
            const date = new Date(dateString);
            if (isNaN(date.getTime())) {
                return 'Invalid Date';
            }
            return date.toLocaleDateString();
        } catch (error) {
            return 'Invalid Date';
        }
    }

    // Get the appropriate status badge class
    function getStatusClass(status: string) {
        switch (status) {
            case TASK_STATUS.COMPLETED:
                return 'bg-success/20 text-success dark:bg-success/30';
            case TASK_STATUS.IN_PROGRESS:
                return 'bg-primary/20 text-primary dark:bg-primary/30';
            default:
                return 'bg-warning/20 text-warning dark:bg-warning/30';
        }
    }

    // Get the appropriate priority class
    function getPriorityClass(priority : string) {
        switch (priority) {
            case PRIORITY.HIGH:
                return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300';
            case PRIORITY.MEDIUM:
                return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300';
            default:
                return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300';
        }
    }

    // Toggle sort direction when clicking the same column
    function toggleSort(column: string) {
        if (sortBy === column) {
            sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
        } else {
            sortBy = column;
            sortDirection = 'asc';
        }
    }

    // Update task status
    async function updateTaskStatus(id: string, newStatus: TaskStatus) {
        await unifiedStore.updateTaskStatus(id, newStatus);
		await Promise.all([loadTasks(), loadCategories()]);
    }

    // Delete task
    function deleteTask(id: string) {
        if (confirm('Are you sure you want to delete this task?')) {
            unifiedStore.deleteTask(id);
        }
    }
</script>

<div class="container mx-auto">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-2xl font-bold">All Tasks</h1>
		<a
			href="/tasks/create"
			class="inline-flex items-center px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/80 transition-colors"
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="h-5 w-5 mr-2"
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
			Add New Task
		</a>
	</div>

	<!-- Filter and search controls -->
	<div class="bg-white dark:bg-gray-800 rounded-lg shadow mb-6 p-4">
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
			<!-- Search input -->
			<div>
				<label for="search" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
					Search
				</label>
				<input
					type="text"
					id="search"
					bind:value={searchQuery}
					placeholder="Search tasks..."
					class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
				/>
			</div>

			<!-- Status filter -->
			<div>
				<label for="status" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
					Status
				</label>
				<select
					id="status"
					bind:value={statusFilter}
					class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
				>
					<option value="all">All Statuses</option>
					<option value={TASK_STATUS.TODO}>To Do</option>
					<option value={TASK_STATUS.IN_PROGRESS}>In Progress</option>
					<option value={TASK_STATUS.COMPLETED}>Completed</option>
				</select>
			</div>

			<!-- Priority filter -->
			<div>
				<label
					for="priority"
					class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
				>
					Priority
				</label>
				<select
					id="priority"
					bind:value={priorityFilter}
					class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
				>
					<option value="all">All Priorities</option>
					<option value={PRIORITY.HIGH}>High</option>
					<option value={PRIORITY.MEDIUM}>Medium</option>
					<option value={PRIORITY.LOW}>Low</option>
				</select>
			</div>

			<!-- Sort options -->
			<div>
				<label for="sort" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
					Sort By
				</label>
				<div class="flex gap-2">
					<select
						id="sort"
						bind:value={sortBy}
						class="flex-grow px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
					>
						<option value="dueDate">Due Date</option>
						<option value="priority">Priority</option>
						<option value="title">Title</option>
					</select>
					<button
						onclick={() => (sortDirection = sortDirection === 'asc' ? 'desc' : 'asc')}
						class="px-3 py-2 bg-gray-200 dark:bg-gray-700 rounded-md"
						aria-label={sortDirection === 'asc' ? 'Sort ascending' : 'Sort descending'}
					>
						{#if sortDirection === 'asc'}
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
									d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12"
								/>
							</svg>
						{:else}
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
									d="M3 4h13M3 8h9m-9 4h9m5-4v12m0 0l-4-4m4 4l4-4"
								/>
							</svg>
						{/if}
					</button>
				</div>
			</div>
		</div>
	</div>

	<!-- Tasks list -->
	<div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
		{#if filteredTasks.length === 0}
			<div class="p-8 text-center">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mx-auto h-12 w-12 text-gray-400 dark:text-gray-500"
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
				<h3 class="mt-2 text-lg font-medium text-gray-900 dark:text-gray-100">No tasks found</h3>
				<p class="mt-1 text-gray-500 dark:text-gray-400">
					{#if searchQuery || statusFilter !== 'all' || priorityFilter !== 'all'}
						Try adjusting your filters or search query
					{:else}
						Get started by creating a new task
					{/if}
				</p>
				{#if !searchQuery && statusFilter === 'all' && priorityFilter === 'all'}
					<div class="mt-6">
						<a
							href="/tasks/create"
							class="inline-flex items-center px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/80 transition-colors"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-5 w-5 mr-2"
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
							Add New Task
						</a>
					</div>
				{/if}
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
					<thead class="bg-gray-50 dark:bg-gray-700">
						<tr>
							<th
								scope="col"
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
							>
								<button
									onclick={() => toggleSort('title')}
									class="flex items-center font-medium focus:outline-none"
								>
									Title
									{#if sortBy === 'title'}
										<span class="ml-1">
											{sortDirection === 'asc' ? '↑' : '↓'}
										</span>
									{/if}
								</button>
							</th>
							<th
								scope="col"
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider hidden md:table-cell"
							>
								Category
							</th>
							<th
								scope="col"
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
							>
								<button
									onclick={() => toggleSort('priority')}
									class="flex items-center font-medium focus:outline-none"
								>
									Priority
									{#if sortBy === 'priority'}
										<span class="ml-1">
											{sortDirection === 'asc' ? '↑' : '↓'}
										</span>
									{/if}
								</button>
							</th>
							<th
								scope="col"
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider hidden sm:table-cell"
							>
								Status
							</th>
							<th
								scope="col"
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
							>
								<button
									onclick={() => toggleSort('dueDate')}
									class="flex items-center font-medium focus:outline-none"
								>
									Due Date
									{#if sortBy === 'dueDate'}
										<span class="ml-1">
											{sortDirection === 'asc' ? '↑' : '↓'}
										</span>
									{/if}
								</button>
							</th>
							<th
								scope="col"
								class="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
							>
								Actions
							</th>
						</tr>
					</thead>
					<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
						{#each filteredTasks as task (task.id)}
							<tr>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="text-sm font-medium text-gray-900 dark:text-gray-100">
										{task.title}
									</div>
									<div class="text-xs text-gray-500 dark:text-gray-400 line-clamp-1">
										{task.description}
									</div>
									{#if task.tags && task.tags.length > 0}
										<div class="mt-1 flex flex-wrap gap-1">
											{#each task.tags as tag}
												<span
													class="inline-block px-2 py-0.5 rounded-full text-xs bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300"
												>
													{tag}
												</span>
											{/each}
										</div>
									{/if}
								</td>
								<td class="px-6 py-4 whitespace-nowrap hidden md:table-cell">
									<span
										class="inline-block px-2 py-1 rounded-md text-xs bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300"
									>
										{task.category?.name || 'Uncategorized'}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<span
										class="px-2 py-1 inline-flex text-xs rounded-full {getPriorityClass(
											task.priority
										)}"
									>
										{task.priority.charAt(0).toUpperCase() + task.priority.slice(1)}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap hidden sm:table-cell">
									<span
										class="px-2 py-1 inline-flex text-xs rounded-full {getStatusClass(task.status)}"
									>
										{task.status === TASK_STATUS.TODO
											? 'To Do'
											: task.status === TASK_STATUS.IN_PROGRESS
												? 'In Progress'
												: 'Completed'}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="text-sm text-gray-900 dark:text-gray-100">
										{formatDate(task.dueDate)}
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
									<div class="flex items-center justify-end space-x-2">
										<!-- Status dropdown -->
										<div class="relative inline-block text-left">
											<select
												class="py-1 px-2 text-xs bg-gray-100 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md"
												value={task.status}
												onchange={(e) => {
                                                const target = e.target as HTMLSelectElement;
                                                updateTaskStatus(task.id, target.value as TaskStatus);
                                            }}
											>
												<option value={TASK_STATUS.TODO}>To Do</option>
												<option value={TASK_STATUS.IN_PROGRESS}>In Progress</option>
												<option value={TASK_STATUS.COMPLETED}>Completed</option>
											</select>
										</div>

										<!-- Edit link -->
										<a
											href={`/tasks/${task.id}`}
											class="text-primary hover:text-primary/80 dark:hover:text-primary/70"
											aria-label="Edit task"
										>
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
													d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
												/>
											</svg>
										</a>

										<!-- Delete button -->
										<button
											onclick={() => deleteTask(task.id)}
											class="text-red-600 hover:text-red-800 dark:hover:text-red-400"
											aria-label="Delete task"
										>
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
													d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
												/>
											</svg>
										</button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>
