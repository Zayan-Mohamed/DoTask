<script lang="ts">
    import { onMount } from 'svelte';
    import { user } from '$lib/stores/user';
    import { 
        tasks, 
        categories, 
        loading, 
        error, 
        loadTasks, 
        loadCategories, 
        createTask,
        updateTaskStatus, 
        TASK_STATUS, 
        PRIORITY 
    } from '$lib/stores/unified-store';
    import type { Task, TaskStatus } from '$lib/types';
    import { goto } from '$app/navigation';
    
    // Using Svelte 5 runes for reactivity
    let isSubmitting = $state(false);
    
    // If the rune system isn't working, uncomment this line:
    // let isSubmitting = writable(false);
    
    onMount(async () => {
        try {
            await Promise.all([loadTasks(), loadCategories()]);
        } catch (error) {
            console.error('Failed to load data:', error);
        }
    });

    let totalTasks = $derived(($tasks || []).length);
    
    let completedTasks = $derived(($tasks || []).filter(
        (task: Task) => task.status === TASK_STATUS.COMPLETED
    ).length);
    
    let pendingTasks = $derived(($tasks || []).filter(
        (task: Task) => task.status === TASK_STATUS.TODO
    ).length);
    
    let inProgressTasks = $derived(($tasks || []).filter(
        (task: Task) => task.status === TASK_STATUS.IN_PROGRESS
    ).length);

    let urgentTasks = $derived(($tasks || []).filter(
        (task: Task) => {
            if (!task.dueDate) return false;
            const dueDate = new Date(task.dueDate);
            // Check if date is valid
            if (isNaN(dueDate.getTime())) return false;
            
            const now = new Date();
            const threeDaysFromNow = new Date();
            threeDaysFromNow.setDate(now.getDate() + 3);

            return (
                task.priority === PRIORITY.HIGH &&
                task.status !== TASK_STATUS.COMPLETED &&
                dueDate <= threeDaysFromNow
            );
        }
    ));

    let completionRate = $derived(totalTasks > 0 ? Math.round((completedTasks / totalTasks) * 100) : 0);

    function sortByCreatedAt(a: Task, b: Task) {
        return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
    }
    
    // Quick add task functionality
    let quickTask = {
        title: '',
        categoryId: '',
        status: TASK_STATUS.TODO,
        priority: PRIORITY.MEDIUM
    };

    async function handleQuickAddTask(event: Event) {
        event.preventDefault();
        if (!quickTask.title || !quickTask.categoryId) return;

        
        isSubmitting = true;
        try {
            const newTask = await createTask({
                title: quickTask.title,
                status: quickTask.status,
                priority: quickTask.priority,
                categoryId: quickTask.categoryId
            });
            
            // Reset form
            quickTask.title = '';
            
            // Optional: Navigate to the new task
            // goto(`/tasks/${newTask.id}`);
        } catch (err) {
            console.error('Failed to create task:', err);
        } finally {
            isSubmitting = false;
        }
    }
    
    // Add task status toggle functionality
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
    
    async function handleStatusToggle(task: Task) {
        try {
            const nextStatus = getNextStatus(task.status);
            await updateTaskStatus(task.id, nextStatus);
        } catch (err) {
            console.error('Failed to update task status:', err);
        }
    }
</script>

<div class="container mx-auto">
    {#if $loading}
        <div class="flex justify-center my-12">
            <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary"></div>
        </div>
    {:else if $error}
        <div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-8" role="alert">
            <p class="font-bold">Error</p>
            <p>{$error}</p>
        </div>
    {:else}
    <div class="mb-8 flex justify-between items-center">
        <div>
            <h1 class="text-3xl font-bold mb-2">Welcome, {$user}!</h1>
            <p class="text-muted">Here's an overview of your tasks.</p>
        </div>
        <button 
            class="px-4 py-2 bg-primary text-white rounded hover:bg-primary/80 transition-colors"
            onclick={() => Promise.all([loadTasks(), loadCategories()])}
        >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 inline-block mr-1" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z" clip-rule="evenodd" />
            </svg>
            Refresh
        </button>
    </div>

    <!-- Task Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <!-- Total Tasks Card -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h3 class="text-lg font-semibold text-muted mb-2">Total Tasks</h3>
            <div class="flex items-center">
                <span class="text-3xl font-bold">{totalTasks}</span>
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-8 w-8 ml-auto text-accent"
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
        </div>

        <!-- Completed Tasks Card -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h3 class="text-lg font-semibold text-muted mb-2">Completed</h3>
            <div class="flex items-center">
                <span class="text-3xl font-bold">{completedTasks}</span>
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-8 w-8 ml-auto text-success"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M5 13l4 4L19 7"
                    />
                </svg>
            </div>
        </div>

        <!-- In Progress Tasks Card -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h3 class="text-lg font-semibold text-muted mb-2">In Progress</h3>
            <div class="flex items-center">
                <span class="text-3xl font-bold">{inProgressTasks}</span>
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-8 w-8 ml-auto text-primary"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"
                    />
                </svg>
            </div>
        </div>

        <!-- Pending Tasks Card -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h3 class="text-lg font-semibold text-muted mb-2">Pending</h3>
            <div class="flex items-center">
                <span class="text-3xl font-bold">{pendingTasks}</span>
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-8 w-8 ml-auto text-warning"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                </svg>
            </div>
        </div>
    </div>

    <!-- Progress Bar -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 mb-8">
        <div class="flex justify-between items-center mb-2">
            <h3 class="text-lg font-semibold">Completion Rate</h3>
            <span class="font-bold">{completionRate}%</span>
        </div>
        <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5">
            <div class="bg-success h-2.5 rounded-full" style="width: {completionRate}%"></div>
        </div>
    </div>
    
    <!-- Quick Add Task -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 mb-8">
        <h3 class="text-lg font-semibold mb-4">Quick Add Task</h3>
        
        {#if $categories.length === 0}
            <p class="text-muted text-center py-4">You need to create a category first before adding tasks.</p>
        {:else}
            <form onsubmit={handleQuickAddTask} class="space-y-4">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label for="title" class="block text-sm font-medium mb-1">Title</label>
                        <input
                            type="text"
                            id="title"
                            bind:value={quickTask.title}
                            class="w-full px-4 py-3 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:text-white"
                            placeholder="Task title"
                            required
                        />
                    </div>
                    <div>
                        <label for="category" class="block text-sm font-medium mb-1">Category</label>
                        <select
                            id="category"
                            bind:value={quickTask.categoryId}
                            class="w-full px-4 py-3 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:text-white"
                            required
                        >
                            <option value="" disabled>Select a category</option>
                            {#each $categories as category}
                                <option value={category.id}>{category.name}</option>
                            {/each}
                        </select>
                    </div>
                </div>
                
                <div class="flex justify-end">
                    <button
                        type="submit"
                        class="px-4 py-2 bg-primary text-white rounded hover:bg-primary/80 transition-colors"
                        disabled={isSubmitting}
                    >
                        {isSubmitting ? 'Adding...' : 'Add Task'}
                    </button>
                </div>
            </form>
        {/if}
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Urgent Tasks Section -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h3 class="text-lg font-semibold mb-4">Urgent Tasks</h3>

            {#if urgentTasks.length === 0}
                <p class="text-muted text-center py-4">No urgent tasks right now. Great job!</p>
            {:else}
                <ul class="space-y-3">
                    {#each urgentTasks as task}
                        <li class="border-l-4 border-red-500 bg-red-50 dark:bg-red-900/20 p-3 rounded-r-md hover:bg-red-100 dark:hover:bg-red-900/30 transition-colors">
                            <a href="/tasks/{task.id}" class="block">
                                <div class="flex justify-between">
                                    <h4 class="font-medium">{task.title}</h4>
                                    <span class="text-sm text-muted">
                                        {(() => {
                                            try {
                                                const date = new Date(task.dueDate);
                                                return isNaN(date.getTime()) ? 'Invalid Date' : date.toLocaleDateString();
                                            } catch {
                                                return 'Invalid Date';
                                            }
                                        })()}
                                    </span>
                                </div>
                                <p class="text-sm text-muted line-clamp-1">{task.description}</p>
                            </a>
                        </li>
                    {/each}
                </ul>
            {/if}
        </div>

        <!-- Recent Activity -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h3 class="text-lg font-semibold mb-4">Recent Activity</h3>

            {#if ($tasks || []).length === 0}
                <p class="text-muted text-center py-4">No task activity yet. Start by adding a task!</p>
            {:else}
                <ul class="space-y-3">
                    {#each [...($tasks || [])].sort(sortByCreatedAt).slice(0, 5) as task}
                        <li class="flex items-center p-2 border-b border-gray-100 dark:border-gray-700">
                            <!-- Status indicator - clickable to toggle task status -->
                            <button 
                                onclick={() => handleStatusToggle(task)}
                                class="flex-shrink-0 w-6 h-6 mr-3 flex items-center justify-center rounded-full border transition-colors"
                                class:bg-success={task.status === TASK_STATUS.COMPLETED}
                                class:bg-primary={task.status === TASK_STATUS.IN_PROGRESS}
                                class:bg-warning={task.status === TASK_STATUS.TODO}
                                title={`Mark as ${getNextStatus(task.status)}`}
                            >
                                {#if task.status === TASK_STATUS.COMPLETED}
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-white" viewBox="0 0 20 20" fill="currentColor">
                                        <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                                    </svg>
                                {:else if task.status === TASK_STATUS.IN_PROGRESS}
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-white" viewBox="0 0 20 20" fill="currentColor">
                                        <path fill-rule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clip-rule="evenodd" />
                                    </svg>
                                {:else}
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-white" viewBox="0 0 20 20" fill="currentColor">
                                        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v2H7a1 1 0 100 2h2v2a1 1 0 102 0v-2h2a1 1 0 100-2h-2V7z" clip-rule="evenodd" />
                                    </svg>
                                {/if}
                            </button>

                            <div class="flex-1">
                                <h4 class="font-medium">{task.title}</h4>
                                <p class="text-xs text-muted">
                                    {(() => {
                                        try {
                                            if (!task.createdAt) return 'Unknown date';
                                            const date = new Date(task.createdAt);
                                            return isNaN(date.getTime()) ? 'Unknown date' : date.toLocaleString();
                                        } catch {
                                            return 'Unknown date';
                                        }
                                    })()}
                                </p>
                            </div>

                            <span
                                class="px-2 py-1 text-xs rounded-full {task.priority === PRIORITY.HIGH
                                    ? 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300'
                                    : task.priority === PRIORITY.MEDIUM
                                        ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300'
                                        : 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300'}"
                            >
                                {task.priority}
                            </span>
                        </li>
                    {/each}
                </ul>
            {/if}
        </div>
    </div>
    {/if}
</div>
