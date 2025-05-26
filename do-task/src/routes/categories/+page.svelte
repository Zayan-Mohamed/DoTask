<script lang="ts">
	import { categories, tasks, unifiedStore } from '$lib/stores/unified-store';
	import type { Category } from '$lib/types';

	let newCategory = $state('');
	let editMode: string | null = $state(null);
	let editValue = $state('');
	let errorMessage = $state('');
	let successMessage = $state('');

	// Calculate tasks per category
	const categoryStats = $derived($categories.map((category) => {
		const count = $tasks.filter((task) => task.category?.id === category.id).length;
		return { name: category.name, count, category };
	}));

	function addCategory() {
		if (newCategory.trim() === '') {
			errorMessage = 'Category name cannot be empty';
			setTimeout(() => {
				errorMessage = '';
			}, 3000);
			return;
		}

		if ($categories.some(cat => cat.name === newCategory.trim())) {
			errorMessage = 'Category already exists';
			setTimeout(() => {
				errorMessage = '';
			}, 3000);
			return;
		}

		unifiedStore.createCategory(newCategory.trim());
		newCategory = '';

		successMessage = 'Category added successfully';
		setTimeout(() => {
			successMessage = '';
		}, 3000);
	}

	function deleteCategory(categoryId: string) {
		// Check if any tasks use this category
		const tasksUsingCategory = $tasks.filter((task) => task.category?.id === categoryId).length;

		if (tasksUsingCategory > 0) {
			if (
				!confirm(
					`There are ${tasksUsingCategory} task(s) using this category. Are you sure you want to delete it?`
				)
			) {
				return;
			}
		}

		unifiedStore.deleteCategory(categoryId);

		// If any tasks use this category, update them to use the first available category or null
		if (tasksUsingCategory > 0) {
			const defaultCategory = $categories.length > 0 ? $categories[0] : null;

			// Update all tasks with the deleted category
			$tasks.forEach((task) => {
				if (task.category?.id === categoryId) {
					unifiedStore.updateTask(task.id, { 
						title: task.title,
						description: task.description,
						status: task.status,
						priority: task.priority,
						dueDate: task.dueDate,
						categoryId: defaultCategory?.id || '',
						tags: task.tags
					});
				}
			});
		}

		successMessage = 'Category deleted successfully';
		setTimeout(() => {
			successMessage = '';
		}, 3000);
	}

	function startEditMode(categoryId: string, categoryName: string) {
        if (editMode) {
            errorMessage = 'Please finish editing the current category first';
            setTimeout(() => {
                errorMessage = '';
            }, 3000);
            return;
        }
		editMode = categoryId;
		editValue = categoryName;
	}

	function cancelEditMode() {
		editMode = null;
		editValue = '';
	}

	function saveEdit(categoryId: string) {
		if (editValue.trim() === '') {
			errorMessage = 'Category name cannot be empty';
			setTimeout(() => {
				errorMessage = '';
			}, 3000);
			return;
		}

		const currentCategory = $categories.find(cat => cat.id === categoryId);
		if (currentCategory && currentCategory.name === editValue.trim()) {
			// No changes
			cancelEditMode();
			return;
		}

		if ($categories.some(cat => cat.name === editValue.trim())) {
			errorMessage = 'Category already exists';
			setTimeout(() => {
				errorMessage = '';
			}, 3000);
			return;
		}

		unifiedStore.updateCategory(categoryId, editValue.trim());

		cancelEditMode();

		successMessage = 'Category updated successfully';
		setTimeout(() => {
			successMessage = '';
		}, 3000);
	}
</script>

<div class="container mx-auto py-6">
	<div class="flex justify-between items-center mb-6">
		<h1 class="text-2xl font-bold">Categories</h1>
	</div>

	<!-- Add new category form -->
	<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 mb-6">
		<h2 class="text-xl font-semibold mb-4">Add New Category</h2>

		{#if errorMessage}
			<div
				class="mb-4 bg-red-100 dark:bg-red-900/30 border-l-4 border-red-500 text-red-700 dark:text-red-300 p-4"
				role="alert"
			>
				<p>{errorMessage}</p>
			</div>
		{/if}

		{#if successMessage}
			<div
				class="mb-4 bg-green-100 dark:bg-green-900/30 border-l-4 border-green-500 text-green-700 dark:text-green-300 p-4"
				role="alert"
			>
				<p>{successMessage}</p>
			</div>
		{/if}

		<div class="flex gap-2">
			<input
				type="text"
				placeholder="Enter category name"
				bind:value={newCategory}
				class="flex-grow px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
			/>
			<button
				onclick={addCategory}
				class="px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/80 transition-colors"
			>
				Add Category
			</button>
		</div>
	</div>

	<!-- Categories list -->
	<div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
		{#if $categories.length === 0}
			<div class="p-6 text-center">
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
						d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
					/>
				</svg>
				<h3 class="mt-2 text-lg font-medium text-gray-900 dark:text-gray-100">No categories</h3>
				<p class="mt-1 text-gray-500 dark:text-gray-400">Add categories to organize your tasks.</p>
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
								Category Name
							</th>
							<th
								scope="col"
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
							>
								Tasks
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
						{#each categoryStats as { name, count, category }, i}
							<tr class="hover:bg-gray-50 dark:hover:bg-gray-700/50">
								<td class="px-6 py-4 whitespace-nowrap">
									{#if editMode === category.id}
										<div class="flex gap-2 items-center">
											<input
												type="text"
												bind:value={editValue}
												class="flex-grow px-3 py-1 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
											/>
											<button
												onclick={() => saveEdit(category.id)}
												class="p-1"
												aria-label="Save category name"
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													class="h-5 w-5 text-green-600"
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
											</button>
											<button onclick={cancelEditMode} class="p-1" aria-label="Cancel editing">
												<svg
													xmlns="http://www.w3.org/2000/svg"
													class="h-5 w-5 text-gray-500"
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
									{:else}
										<div class="flex items-center">
											<span
												class="inline-block w-3 h-3 rounded-full bg-primary mr-2"
												style="background-color: hsl({(i * 137) % 360}, 70%, 60%);"
											></span>
											<div class="text-sm font-medium text-gray-900 dark:text-gray-100">
												{name}
											</div>
										</div>
									{/if}
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<span
										class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300"
									>
										{count}
										{count === 1 ? 'task' : 'tasks'}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
									{#if editMode !== category.id}
										<div class="flex items-center justify-end space-x-2">
											<button
												onclick={() => startEditMode(category.id, name)}
												class="text-primary hover:text-primary/80"
												aria-label="Edit category"
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
											</button>

											<button
												onclick={() => deleteCategory(category.id)}
												class="text-red-600 hover:text-red-800"
												aria-label="Delete category"
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
									{/if}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>
