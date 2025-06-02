<script lang="ts">
    import { page } from '$app/state';
    import { TASK_STATUS, PRIORITY, type Task } from '$lib/types.js';
    import { goto } from '$app/navigation';
    import { tasks, categories, unifiedStore } from '$lib/stores/unified-store';

    const taskId = page.params.id;

let task = $derived($tasks.find((t) => t.id === taskId));
let notFound = $derived(!task);

    let editedTask : Task = $derived(task
        ? {
            ...task,
            dueDate: (() => {
                try {
                    const date = new Date(task.dueDate);
                    return isNaN(date.getTime()) ? new Date().toISOString().split('T')[0] : date.toISOString().split('T')[0];
                } catch {
                    return new Date().toISOString().split('T')[0];
                }
            })()
        }
        : {
            id: '',
            title: '',
            description: '',
            status: TASK_STATUS.TODO,
            priority: PRIORITY.MEDIUM,
            dueDate: new Date().toISOString().split('T')[0],
            category: $categories.length > 0 ? $categories[0] : { id: 'temp-id', name: 'Create a category first' },
            tags: [],
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString()
        });

    let tagInput = $state('');
    let submitting = $state(false);
    let success = $state(false);
    let errorMessage = $state('');
    let newCategory = $state('');

    function addTag() {
        if (editedTask && tagInput.trim() !== '' && !editedTask.tags.includes(tagInput.trim())) {
            editedTask.tags = [...editedTask.tags, tagInput.trim()];
            tagInput = '';
        }
    }

    function removeTag(tag: string) {
        if (editedTask) {
            editedTask.tags = editedTask.tags.filter((t) => t !== tag);
        }
    }

    async function handleSubmit(event: SubmitEvent) {
        event.preventDefault();
        submitting = true;
        errorMessage = '';

        if (!editedTask) {
            errorMessage = 'Task not found';
            submitting = false;
            return;
        }

        if (!editedTask.title.trim()) {
            errorMessage = 'Title is required';
            submitting = false;
            return;
        }

        try {
            // Update the task in the store
            await unifiedStore.updateTask(taskId, {
                title: editedTask.title,
                description: editedTask.description,
                status: editedTask.status,
                priority: editedTask.priority,
                dueDate: new Date(`${editedTask.dueDate}T00:00:00`).toISOString(),
                categoryId: editedTask.category?.id || '',
                tags: editedTask.tags
            });

            success = true;

            // Reset success message after 3 seconds
            setTimeout(() => {
                success = false;
            }, 3000);
        } catch (error) {
            errorMessage = error instanceof Error
                ? `Error updating task: ${error.message}`
                : 'An unknown error occurred';
        } finally {
            submitting = false;
        }
    }

    async function addCategory() {
            if (newCategory.trim() !== '') {
                try {
                    // Use the createCategory function to add a new category with proper structure
                    const category = await unifiedStore.createCategory(newCategory.trim());
                    if (editedTask) {
                        editedTask.category = category;
                    }
                    newCategory = '';
                } catch (error) {
                    errorMessage = error instanceof Error
                        ? `Error creating category: ${error.message}`
                        : 'An unknown error occurred';
                }
            }
        }

    function deleteTask() {
        if (confirm('Are you sure you want to delete this task?')) {
            unifiedStore.deleteTask(taskId);
            goto('/tasks');
        }
    }
</script>

<div class="container mx-auto py-6">
    {#if notFound}
        <div class="max-w-2xl mx-auto bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h1 class="text-2xl font-bold mb-6">Task Not Found</h1>
            <p class="mb-6">The task you are looking for does not exist or has been deleted.</p>
            <a
                href="/tasks"
                class="inline-block px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/80 transition-colors"
            >
                Return to Task List
            </a>
        </div>
    {:else}
        <div class="flex justify-between items-center mb-6">
            <a
                href="/tasks"
                class="flex items-center text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5 mr-1"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M15 19l-7-7 7-7"
                    />
                </svg>
                Back to Tasks
            </a>
            <button
                onclick={deleteTask}
                class="inline-flex items-center px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition-colors"
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
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                    />
                </svg>
                Delete Task
            </button>
        </div>

        <div class="max-w-2xl mx-auto bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h1 class="text-2xl font-bold mb-6">Edit Task</h1>

            {#if success}
                <div
                    class="bg-green-100 dark:bg-green-900/30 border-l-4 border-green-500 text-green-700 dark:text-green-300 p-4 mb-6"
                    role="alert"
                >
                    <p>Task updated successfully!</p>
                </div>
            {/if}

            {#if errorMessage}
                <div
                    class="bg-red-100 dark:bg-red-900/30 border-l-4 border-red-500 text-red-700 dark:text-red-300 p-4 mb-6"
                    role="alert"
                >
                    <p>{errorMessage}</p>
                </div>
            {/if}

            <form onsubmit={handleSubmit} class="space-y-6">
                <!-- Title -->
                <div>
                    <label
                        for="title"
                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                    >
                        Title <span class="text-red-600">*</span>
                    </label>
                    <input
                        type="text"
                        id="title"
                        bind:value={editedTask.title}
                        required
                        class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
                    />
                </div>

                <!-- Description -->
                <div>
                    <label
                        for="description"
                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                    >
                        Description
                    </label>
                    <textarea
                        id="description"
                        bind:value={editedTask.description}
                        rows="3"
                        class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
                    ></textarea>
                </div>

                <!-- Status -->
                <div>
                    <label
                        for="status"
                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                    >
                        Status
                    </label>
                    <select
                        id="status"
                        bind:value={editedTask.status}
                        class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
                    >
                        <option value={TASK_STATUS.TODO}>To Do</option>
                        <option value={TASK_STATUS.IN_PROGRESS}>In Progress</option>
                        <option value={TASK_STATUS.COMPLETED}>Completed</option>
                    </select>
                </div>

                <!-- Priority -->
                <div>
                    <label
                        for="priority"
                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                    >
                        Priority
                    </label>
                    <select
                        id="priority"
                        bind:value={editedTask.priority}
                        class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
                    >
                        <option value={PRIORITY.LOW}>Low</option>
                        <option value={PRIORITY.MEDIUM}>Medium</option>
                        <option value={PRIORITY.HIGH}>High</option>
                    </select>
                </div>

                <!-- Due Date -->
                <div>
                    <label
                        for="dueDate"
                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                    >
                        Due Date
                    </label>
                    <input
                        type="date"
                        id="dueDate"
                        bind:value={editedTask.dueDate}
                        class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
                    />
                </div>

                <!-- Category -->
                <div>
                    <label
                        for="category"
                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                    >
                        Category
                    </label>                        <div class="flex gap-2">
                        <select
                            id="category"
                            class="flex-grow px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
                            onchange={(e) => {
                                const target = e.target as HTMLSelectElement;
                                const selectedId = target.value;
                                const selectedCategory = $categories.find(cat => cat.id === selectedId);
                                if (selectedCategory) {
                                    editedTask.category = selectedCategory;
                                }
                            }}
                        >
                            {#each $categories as category}
                                <option 
                                    value={category.id} 
                                    selected={editedTask.category?.id === category.id}
                                >
                                    {category.name}
                                </option>
                            {/each}
                        </select>
                        <div class="flex items-center gap-2">
                            <input
                                type="text"
                                placeholder="New category"
                                bind:value={newCategory}
                                class="px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
                            />
                            <button
                                type="button"
                                onclick={addCategory}
                                class="p-2 bg-accent text-white rounded-md hover:bg-accent/80"
                                disabled={!newCategory.trim()}
                                aria-label="Add category"
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
                                        d="M12 6v6m0 0v6m0-6h6m-6 0H6"
                                    />
                                </svg>
                            </button>
                        </div>
                    </div>
                </div>

                <!-- Tags -->
                <div>
                    <label for="tags" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Tags
                    </label>
                    <div class="flex flex-wrap gap-2 mb-2">
                        {#each editedTask.tags as tag}
                            <div
                                class="inline-flex items-center bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300 px-2 py-1 rounded-full text-sm"
                            >
                                {tag}
                                <button
                                    aria-label="Remove tag"
                                    type="button"
                                    onclick={() => removeTag(tag)}
                                    class="ml-1 text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
                                >
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        class="h-4 w-4"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                        stroke="currentColor"
                                    >
                                        <path
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            stroke-width="2"
                                            d="M6 18L18 6M6 6l12 12"
                                        />
                                    </svg>
                                </button>
                            </div>
                        {/each}
                    </div>
                    <div class="flex gap-2">
                        <input
                            type="text"
                            id="tags"
                            bind:value={tagInput}
                            placeholder="Add a tag"
                            class="flex-grow px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
                        />
                        <button
                            type="button"
                            onclick={addTag}
                            class="p-2 bg-accent text-white rounded-md hover:bg-accent/80"
                            disabled={!tagInput.trim()}
                            aria-label="Add tag"
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
                                    d="M12 6v6m0 0v6m0-6h6m-6 0H6"
                                />
                            </svg>
                        </button>
                    </div>
                </div>

                <!-- Submit button -->
                <div class="flex justify-end">
                    <button
                        type="submit"
                        disabled={submitting}
                        class="inline-flex items-center px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/80 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {#if submitting}
                            <svg
                                class="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                            >
                                <circle
                                    class="opacity-25"
                                    cx="12"
                                    cy="12"
                                    r="10"
                                    stroke="currentColor"
                                    stroke-width="4"
                                />
                                <path
                                    class="opacity-75"
                                    fill="currentColor"
                                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                />
                            </svg>
                        {/if}
                        {submitting ? 'Updating...' : 'Update Task'}
                    </button>
                </div>
            </form>
        </div>
    {/if}
</div>
