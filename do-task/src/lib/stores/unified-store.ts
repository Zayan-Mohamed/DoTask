import { writable, derived } from 'svelte/store';
import { client } from '$lib/graphql/client';
import {
	GET_TASKS,
	GET_CATEGORIES,
	CREATE_TASK,
	UPDATE_TASK,
	DELETE_TASK,
	UPDATE_TASK_STATUS,
	CREATE_CATEGORY,
	UPDATE_CATEGORY,
	DELETE_CATEGORY
} from '$lib/graphql/queries';
import {
	TASK_STATUS,
	PRIORITY,
	type Task,
	type Category,
	type CreateTaskInput,
	type UpdateTaskInput,
	type TaskStatus,
	type TaskPriority
} from '$lib/types';
import type { Writable } from 'svelte/store';

// Create reactive stores
export const tasks = writable<Task[]>([]);
export const categories = writable<Category[]>([]);
export const loading = writable(false);
export const error = writable<string | null>(null);

// Re-export constants
export { TASK_STATUS, PRIORITY };

// Helper function to transform task data from GraphQL to match frontend expectations
function transformTask(task: any): Task {
	return {
		...task,
		// Ensure dates are ISO strings
		dueDate: task.dueDate,
		createdAt: task.createdAt,
		updatedAt: task.updatedAt,
		// Transform category from object to keep compatibility, but we'll migrate this
		category: task.category
	};
}

// Helper function to transform category data
function transformCategory(category: any): Category {
	return {
		id: category.id,
		name: category.name
	};
}

// Load tasks from GraphQL
export async function loadTasks() {
	loading.set(true);
	error.set(null);

	try {
		const result = await client.query({
			query: GET_TASKS,
			fetchPolicy: 'network-only' // Always fetch fresh data
		});

		const transformedTasks = result.data.tasks.map(transformTask);
		tasks.set(transformedTasks);
		return transformedTasks;
	} catch (err) {
		console.error('Error loading tasks:', err);
		const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
		error.set(errorMessage);
		throw err;
	} finally {
		loading.set(false);
	}
}

// Load categories from GraphQL
export async function loadCategories() {
	loading.set(true);
	error.set(null);

	try {
		const result = await client.query({
			query: GET_CATEGORIES,
			fetchPolicy: 'network-only'
		});

		const transformedCategories = result.data.categories.map(transformCategory);
		categories.set(transformedCategories);
		return transformedCategories;
	} catch (err) {
		console.error('Error loading categories:', err);
		const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
		error.set(errorMessage);
		throw err;
	} finally {
		loading.set(false);
	}
}

// Create a new task
export async function createTask(taskInput: CreateTaskInput): Promise<Task> {
	loading.set(true);
	error.set(null);

	try {
		const result = await client.mutate({
			mutation: CREATE_TASK,
			variables: {
				input: {
					title: taskInput.title,
					description: taskInput.description || '',
					status: taskInput.status,
					priority: taskInput.priority,
					dueDate: taskInput.dueDate,
					categoryId: taskInput.categoryId, // Note: backend expects categoryId
					tags: taskInput.tags || []
				}
			}
		});

		const newTask = transformTask(result.data.createTask);
		tasks.update((currentTasks) => [...currentTasks, newTask]);
		return newTask;
	} catch (err) {
		console.error('Error creating task:', err);
		const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
		error.set(errorMessage);
		throw err;
	} finally {
		loading.set(false);
	}
}

// Update an existing task
export async function updateTask(id: string, updates: UpdateTaskInput): Promise<Task> {
	loading.set(true);
	error.set(null);

	try {
		const result = await client.mutate({
			mutation: UPDATE_TASK,
			variables: {
				id,
				input: {
					title: updates.title,
					description: updates.description,
					status: updates.status,
					priority: updates.priority,
					dueDate: updates.dueDate,
					categoryId: updates.categoryId, // Note: backend expects categoryId
					tags: updates.tags
				}
			}
		});

		const updatedTask = transformTask(result.data.updateTask);
		tasks.update((currentTasks) =>
			currentTasks.map((task) => (task.id === id ? updatedTask : task))
		);
		return updatedTask;
	} catch (err) {
		console.error('Error updating task:', err);
		const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
		error.set(errorMessage);
		throw err;
	} finally {
		loading.set(false);
	}
}

// Delete a task
export async function deleteTask(id: string): Promise<boolean> {
	loading.set(true);
	error.set(null);

	try {
		await client.mutate({
			mutation: DELETE_TASK,
			variables: { id }
		});

		tasks.update((currentTasks) => currentTasks.filter((task) => task.id !== id));
		return true;
	} catch (err) {
		console.error('Error deleting task:', err);
		const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
		error.set(errorMessage);
		throw err;
	} finally {
		loading.set(false);
	}
}

// Update task status
export async function updateTaskStatus(id: string, status: TaskStatus): Promise<Task> {
	loading.set(true);
	error.set(null);

	try {
		const result = await client.mutate({
			mutation: UPDATE_TASK_STATUS,
			variables: { id, status }
		});

		const updatedTask = transformTask(result.data.updateTaskStatus);

		tasks.update((currentTasks) =>
			currentTasks.map((task) =>
				task.id === id ? { ...task, ...updatedTask, updatedAt: new Date().toISOString() } : task
			)
		);
		return updatedTask;
	} catch (err) {
		console.error('Error updating task status:', err);
		const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
		error.set(errorMessage);
		throw err;
	} finally {
		loading.set(false);
	}
}

// Create a new category
export async function createCategory(name: string): Promise<Category> {
	loading.set(true);
	error.set(null);

	try {
		const result = await client.mutate({
			mutation: CREATE_CATEGORY,
			variables: { name }
		});

		const newCategory = transformCategory(result.data.createCategory);
		categories.update((currentCategories) => [...currentCategories, newCategory]);
		return newCategory;
	} catch (err) {
		console.error('Error creating category:', err);
		const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
		error.set(errorMessage);
		throw err;
	} finally {
		loading.set(false);
	}
}

// Update a category
export async function updateCategory(id: string, name: string): Promise<Category> {
	loading.set(true);
	error.set(null);

	try {
		const result = await client.mutate({
			mutation: UPDATE_CATEGORY,
			variables: { id, name }
		});

		const updatedCategory = transformCategory(result.data.updateCategory);
		categories.update((currentCategories) =>
			currentCategories.map((category) => (category.id === id ? updatedCategory : category))
		);
		return updatedCategory;
	} catch (err) {
		console.error('Error updating category:', err);
		const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
		error.set(errorMessage);
		throw err;
	} finally {
		loading.set(false);
	}
}

// Delete a category
export async function deleteCategory(id: string): Promise<boolean> {
	loading.set(true);
	error.set(null);

	try {
		await client.mutate({
			mutation: DELETE_CATEGORY,
			variables: { id }
		});

		categories.update((currentCategories) =>
			currentCategories.filter((category) => category.id !== id)
		);
		return true;
	} catch (err) {
		console.error('Error deleting category:', err);
		const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
		error.set(errorMessage);
		throw err;
	} finally {
		loading.set(false);
	}
}

// Initialize data when the store is first used
export async function initializeData() {
	try {
		await Promise.all([loadTasks(), loadCategories()]);
	} catch (err) {
		console.error('Error initializing data:', err);
		error.set('Failed to load initial data');
	}
}

// Helper functions to maintain compatibility with existing components
export const tasksStore = {
	subscribe: tasks.subscribe,
	// Add compatibility methods that mirror the old store API
	add: createTask,
	updateTask,
	delete: deleteTask,
	updateStatus: updateTaskStatus,
	reload: loadTasks
};

export const categoriesStore = {
	subscribe: categories.subscribe,
	// Add compatibility methods
	add: createCategory,
	update: updateCategory,
	delete: deleteCategory,
	reload: loadCategories
};

// Create derived store for unified state
const unifiedState = derived(
	[tasks, categories, loading, error],
	([$tasks, $categories, $loading, $error]) => ({
		tasks: $tasks,
		categories: $categories,
		loading: $loading,
		error: $error
	})
);

// Main unified store object with subscribe method for use with $ prefix
export const unifiedStore = {
	// Add subscribe method from the derived store
	subscribe: unifiedState.subscribe,

	// Store subscriptions
	tasks,
	categories,
	loading,
	error,

	// Task methods
	loadTasks,
	createTask,
	updateTask,
	deleteTask,
	updateTaskStatus,

	// Category methods
	loadCategories,
	createCategory,
	updateCategory,
	deleteCategory,

	// Utility methods
	initializeData,

	// Compatibility methods (same as individual functions but organized)
	add: createTask,
	updateStatus: updateTaskStatus,
	delete: deleteTask
};

// Define the type for the extended task store
export type TaskStore = Writable<Task[]> & {
	add: (task: Omit<Task, 'id' | 'createdAt'>) => void;
	updateTask: (id: string, updatedTask: Partial<Task>) => void;
	delete: (id: string) => void;
	updateStatus: (id: string, status: TaskStatus) => void;
	reset: () => void;
};
