import { writable } from 'svelte/store';

/**
 * @typedef {'todo' | 'in-progress' | 'completed'} TaskStatus
 * @typedef {'low' | 'medium' | 'high'} TaskPriority
 *
 * @typedef {Object} Task
 * @property {string} id - Unique identifier
 * @property {string} title - Task title
 * @property {string} description - Task description
 * @property {TaskStatus} status - Current status
 * @property {TaskPriority} priority - Priority level
 * @property {string} dueDate - Due date as ISO string
 * @property {string} createdAt - Creation date as ISO string
 * @property {string} category - Task category
 * @property {string[]} tags - Task tags
 * @property {string} updatedAt - Last updated date as ISO string
 */

// Task status options
/** @type {{TODO: TaskStatus, IN_PROGRESS: TaskStatus, COMPLETED: TaskStatus}} */
export const TASK_STATUS = {
	TODO: 'todo',
	IN_PROGRESS: 'in-progress',
	COMPLETED: 'completed'
};

// Priority levels
/** @type {{LOW: TaskPriority, MEDIUM: TaskPriority, HIGH: TaskPriority}} */
export const PRIORITY = {
	LOW: 'low',
	MEDIUM: 'medium',
	HIGH: 'high'
};

// Initialize sample data
/** @type {Task[]} */
const initialTasks = [
	{
		id: '1',
		title: 'Complete project setup',
		description: 'Set up the SvelteKit and Golang GraphQL project structure',
		status: TASK_STATUS.IN_PROGRESS,
		priority: PRIORITY.HIGH,
		dueDate: new Date(Date.now() + 86400000).toISOString(), // Tomorrow
		createdAt: new Date().toISOString(),
		category: 'Development',
		tags: ['setup', 'project'],
		updatedAt: new Date().toISOString()
	},
	{
		id: '2',
		title: 'Design task interface',
		description: 'Create a clean and intuitive interface for managing tasks',
		status: TASK_STATUS.TODO,
		priority: PRIORITY.MEDIUM,
		dueDate: new Date(Date.now() + 172800000).toISOString(), // Day after tomorrow
		createdAt: new Date().toISOString(),
		category: 'Design',
		tags: ['ui', 'ux'],
		updatedAt: new Date().toISOString()
	}
];

/**
 * Creates and returns a typed task store
 * @returns {import('svelte/store').Writable<Task[]> & {
 *   add: (task: Omit<Task, 'id' | 'createdAt'>) => void;
 *   updateTask: (id: string, updatedTask: Partial<Task>) => void;
 *   delete: (id: string) => void;
 *   updateStatus: (id: string, status: TaskStatus) => void;
 *   reset: () => void;
 * }}
 */
function createTaskStore() {
	// Try to load tasks from local storage
	/** @type {Task[]} */
	const storedTasks =
		typeof localStorage !== 'undefined'
			? JSON.parse(localStorage.getItem('tasks') || JSON.stringify(initialTasks))
			: initialTasks;

	const { subscribe, set, update } = writable(storedTasks);

	/**
	 * Persist changes to localStorage whenever the store is updated
	 * @param {Task[]} tasks - The tasks to persist
	 */
	const persistToLocalStorage = (tasks) => {
		if (typeof localStorage !== 'undefined') {
			localStorage.setItem('tasks', JSON.stringify(tasks));
		}
	};

	return {
		subscribe,
		set,
		update,

		/**
		 * Add a new task
		 * @param {Omit<Task, 'id' | 'createdAt'>} task - The new task data
		 */
		add: (task) =>
			update((tasks) => {
				const newTask = {
					...task,
					id: crypto.randomUUID(),
					createdAt: new Date().toISOString()
				};
				const updatedTasks = [...tasks, newTask];
				persistToLocalStorage(updatedTasks);
				return updatedTasks;
			}),

		/**
		 * Update an existing task
		 * @param {string} id - Task ID
		 * @param {Partial<Task>} updatedTask - Updated task data
		 */
		updateTask: (id, updatedTask) =>
			update((tasks) => {
				const taskIndex = tasks.findIndex((t) => t.id === id);
				if (taskIndex !== -1) {
					tasks[taskIndex] = { ...tasks[taskIndex], ...updatedTask };
				}
				persistToLocalStorage(tasks);
				return [...tasks];
			}),

		/**
		 * Delete a task
		 * @param {string} id - Task ID to delete
		 */
		delete: (id) =>
			update((tasks) => {
				const filteredTasks = tasks.filter((t) => t.id !== id);
				persistToLocalStorage(filteredTasks);
				return filteredTasks;
			}),

		/**
		 * Update task status
		 * @param {string} id - Task ID
		 * @param {TaskStatus} status - New status
		 */
		updateStatus: (id, status) =>
			update((tasks) => {
				const taskIndex = tasks.findIndex((t) => t.id === id);
				if (taskIndex !== -1) {
					tasks[taskIndex].status = status;
				}
				persistToLocalStorage(tasks);
				return [...tasks];
			}),

		// Reset to default
		reset: () => {
			set(initialTasks);
			persistToLocalStorage(initialTasks);
		}
	};
}

export const tasks = createTaskStore();

// Create a store for categories
export const categories = writable(['Development', 'Design', 'Marketing', 'Personal']);
