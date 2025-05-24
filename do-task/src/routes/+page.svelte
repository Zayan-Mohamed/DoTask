<script lang="ts">
    import { user } from '$lib/stores/user.js';
    import { tasks, TASK_STATUS, PRIORITY, type Task } from '$lib/stores/tasks.js';

    let totalTasks = $derived($tasks.length);
    
    let completedTasks = $derived($tasks.filter(
        (task: Task) => task.status === TASK_STATUS.COMPLETED
    ).length);
    
    let pendingTasks = $derived($tasks.filter(
        (task: Task) => task.status === TASK_STATUS.TODO
    ).length);
    
    let inProgressTasks = $derived($tasks.filter(
        (task: Task) => task.status === TASK_STATUS.IN_PROGRESS
    ).length);

    let urgentTasks = $derived($tasks.filter(
        (task: Task) => {
            const dueDate = new Date(task.dueDate);
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
</script>

<div class="container mx-auto">
    <div class="mb-8">
        <h1 class="text-3xl font-bold mb-2">Welcome, {$user}!</h1>
        <p class="text-muted">Here's an overview of your tasks.</p>
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

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Urgent Tasks Section -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h3 class="text-lg font-semibold mb-4">Urgent Tasks</h3>

            {#if urgentTasks.length === 0}
                <p class="text-muted text-center py-4">No urgent tasks right now. Great job!</p>
            {:else}
                <ul class="space-y-3">
                    {#each urgentTasks as task}
                        <li class="border-l-4 border-red-500 bg-red-50 dark:bg-red-900/20 p-3 rounded-r-md">
                            <div class="flex justify-between">
                                <h4 class="font-medium">{task.title}</h4>
                                <span class="text-sm text-muted">
                                    {new Date(task.dueDate).toLocaleDateString()}
                                </span>
                            </div>
                            <p class="text-sm text-muted line-clamp-1">{task.description}</p>
                        </li>
                    {/each}
                </ul>
            {/if}
        </div>

        <!-- Recent Activity -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h3 class="text-lg font-semibold mb-4">Recent Activity</h3>

            {#if $tasks.length === 0}
                <p class="text-muted text-center py-4">No task activity yet. Start by adding a task!</p>
            {:else}
                <ul class="space-y-3">
                    {#each [...$tasks].sort(sortByCreatedAt).slice(0, 5) as task}
                        <li class="flex items-center p-2 border-b border-gray-100 dark:border-gray-700">
                            {#if task.status === TASK_STATUS.COMPLETED}
                                <span class="flex-shrink-0 w-3 h-3 mr-3 bg-success rounded-full"></span>
                            {:else if task.status === TASK_STATUS.IN_PROGRESS}
                                <span class="flex-shrink-0 w-3 h-3 mr-3 bg-primary rounded-full"></span>
                            {:else}
                                <span class="flex-shrink-0 w-3 h-3 mr-3 bg-warning rounded-full"></span>
                            {/if}

                            <div class="flex-1">
                                <h4 class="font-medium">{task.title}</h4>
                                <p class="text-xs text-muted">
                                    {new Date(task.createdAt).toLocaleString()}
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
</div>
