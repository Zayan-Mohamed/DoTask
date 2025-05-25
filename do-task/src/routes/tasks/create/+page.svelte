<script lang="ts">
	import { tasks, TASK_STATUS, PRIORITY, categories, type Task } from '$lib/stores/tasks.js';

	let newTask = $state<Task>({
		id: '',
		title: '',
		description: '',
		status: TASK_STATUS.TODO,
		priority: PRIORITY.MEDIUM,
		dueDate: new Date(Date.now() + 86400000).toISOString().split('T')[0], // Tomorrow
		createdAt: '',
		category: $categories[0],
		tags: [],
		updatedAt: ''
	});

	let tagInput = $state('');
	let submitting = $state(false);
	let success = $state(false);
	let errorMessage = $state('');

	function addTag() {
		if (tagInput.trim() !== '' && !newTask.tags.includes(tagInput.trim())) {
			newTask.tags = [...newTask.tags, tagInput.trim()];
			tagInput = '';
		}
	}

	function removeTag(tag: string) {
		newTask.tags = newTask.tags.filter((t) => t !== tag);
	}

	function handleSubmit(event: Event) {
		event.preventDefault();
		submitting = true;
		errorMessage = '';

		if (!newTask.title.trim()) {
			errorMessage = 'Title is required';
			submitting = false;
			return;
		}

		try {
			// Add the new task to the store
			tasks.add({
				title: newTask.title,
				description: newTask.description,
				status: newTask.status,
				priority: newTask.priority,
				dueDate: new Date(`${newTask.dueDate}T00:00:00`).toISOString(),
				category: newTask.category,
				tags: newTask.tags,
				updatedAt: new Date().toISOString()
			});

			// Reset form
			newTask = {
				id: '',
				title: '',
				description: '',
				status: TASK_STATUS.TODO,
				priority: PRIORITY.MEDIUM,
				dueDate: new Date(Date.now() + 86400000).toISOString().split('T')[0], // Tomorrow
				createdAt: '',
				category: $categories[0],
				tags: [],
				updatedAt: ''
			};

			success = true;

			// Reset success message after 3 seconds
			setTimeout(() => {
				success = false;
			}, 3000);
		} catch (error) {
			errorMessage = `Error adding task: ${error instanceof Error ? error.message : String(error)}`;
		} finally {
			submitting = false;
		}
	}

	// For adding custom categories
	let newCategory = $state('');

	function addCategory() {
		if (newCategory.trim() !== '' && !$categories.includes(newCategory.trim())) {
			categories.update((c) => [...c, newCategory.trim()]);
			newTask.category = newCategory.trim();
			newCategory = '';
		}
	}
</script>

<div class="container mx-auto py-6">
	<div class="max-w-2xl mx-auto bg-white dark:bg-gray-800 rounded-lg shadow p-6">
		<h1 class="text-2xl font-bold mb-6">Create New Task</h1>

		{#if success}
			<div
				class="bg-green-100 dark:bg-green-900/30 border-l-4 border-green-500 text-green-700 dark:text-green-300 p-4 mb-6"
				role="alert"
			>
				<p>Task created successfully!</p>
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
				<label for="title" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
					>Title<span class="text-red-500" >*</span></label
				>
				<input
					type="text"
					id="title"
					bind:value={newTask.title}
					class="w-full px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:text-white"
					placeholder="Enter task title"
					required
				/>
			</div>

			<!-- Description -->
			<div>
				<label
					for="description"
					class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Description (optional)</label
				>
				<textarea
					id="description"
					bind:value={newTask.description}
					class="w-full px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:text-white"
					rows="4"
					placeholder="Enter task description"
				></textarea>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<!-- Status -->
				<div>
					<label
						for="status"
						class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Status</label
					>
					<select
						id="status"
						bind:value={newTask.status}
						class="w-full px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:text-white"
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
						class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Priority</label
					>
					<select
						id="priority"
						bind:value={newTask.priority}
						class="w-full px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:text-white"
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
						class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Due Date</label
					>
					<input
						type="date"
						id="dueDate"
						bind:value={newTask.dueDate}
						min={new Date().toISOString().split('T')[0]}
						class="w-full px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:text-white"
					/>
				</div>

				<!-- Category -->
				<div>
					<label
						for="category"
						class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Category</label
					>
					<div class="flex space-x-2">
						<select
							id="category"
							bind:value={newTask.category}
							class="w-full px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:text-white"
						>
							{#each $categories as category}
								<option value={category}>{category}</option>
							{/each}
						</select>

						<button
                            aria-label="Add Category"
							type="button"
							class="px-3 py-2 bg-accent text-white rounded-md hover:bg-accent/80"
							onclick={() => {
								document.getElementById('categoryModal')?.classList.remove('hidden');
							}}
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-5 w-5"
								viewBox="0 0 20 20"
								fill="currentColor"
							>
								<path
									fill-rule="evenodd"
									d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z"
									clip-rule="evenodd"
								/>
							</svg>
						</button>
					</div>
				</div>
			</div>

			<!-- Tags -->
			<div>
				<label for="tags" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Tags</label>
				<div class="flex space-x-2">
					<input
						type="text"
						bind:value={tagInput}
						class="w-full px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:text-white"
						placeholder="Add tags"
					/>
					<button
						type="button"
						class="px-4 py-2 bg-accent text-white rounded-md hover:bg-accent/80 focus:outline-none focus:ring-2 focus:ring-primary"
						onclick={addTag}
					>
						Add
					</button>
				</div>

				{#if newTask.tags.length > 0}
					<div class="mt-2 flex flex-wrap gap-2">
						{#each newTask.tags as tag}
							<span
								class="bg-primary/10 text-primary px-3 py-1 rounded-full text-sm flex items-center"
							>
								{tag}
								<button
                                    aria-label="Remove Tag"
									type="button"
									class="ml-2 focus:outline-none"
									onclick={() => removeTag(tag)}
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										class="h-4 w-4"
										viewBox="0 0 20 20"
										fill="currentColor"
									>
										<path
											fill-rule="evenodd"
											d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
											clip-rule="evenodd"
										/>
									</svg>
								</button>
							</span>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Submit Button -->
			<div class="flex justify-end">
				<button
					type="submit"
					disabled={submitting}
					class="px-6 py-2 bg-primary text-white rounded-md hover:bg-primary/80 focus:outline-none focus:ring-2 focus:ring-primary disabled:opacity-50 disabled:cursor-not-allowed"
				>
					{#if submitting}
						<span class="flex items-center">
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
								></circle>
								<path
									class="opacity-75"
									fill="currentColor"
									d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
								></path>
							</svg>
							Creating...
						</span>
					{:else}
						Create Task
					{/if}
				</button>
			</div>
		</form>
	</div>

	<!-- Category Modal -->
	<div
		id="categoryModal"
		class="fixed inset-0 bg-black bg-opacity-50 items-center justify-center z-50 hidden"
	>
		<div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-lg w-full max-w-md">
			<h3 class="text-lg font-bold mb-4">Add New Category</h3>

			<div class="mb-4">
				<label
					for="newCategory"
					class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
					>Category Name</label
				>
				<input
					type="text"
					id="newCategory"
					bind:value={newCategory}
					class="w-full px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-primary dark:bg-gray-700 dark:text-white"
					placeholder="Enter category name"
				/>
			</div>

			<div class="flex justify-end space-x-3">
				<button
					type="button"
					class="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-white rounded-md hover:bg-gray-300 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-400"
					onclick={() => {
						document.getElementById('categoryModal')?.classList.add('hidden');
					}}
				>
					Cancel
				</button>
				<button
					type="button"
					class="px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/80 focus:outline-none focus:ring-2 focus:ring-primary"
					onclick={() => {
						addCategory();
						document.getElementById('categoryModal')?.classList.add('hidden');
					}}
				>
					Add
				</button>
			</div>
		</div>
	</div>
</div>
