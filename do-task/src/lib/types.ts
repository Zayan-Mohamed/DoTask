// filepath: /home/zayan/Documents/ProjectSvelte/DoTask/do-task/src/lib/types.ts
export type TaskStatus = 'TODO' | 'IN_PROGRESS' | 'COMPLETED';
export type TaskPriority = 'LOW' | 'MEDIUM' | 'HIGH';

export interface Category {
	id: string;
	name: string;
}

export interface Task {
	id: string;
	title: string;
	description: string;
	status: TaskStatus;
	priority: TaskPriority;
	dueDate: string;
	createdAt: string;
	updatedAt: string;
	category: Category;
	tags: string[];
}

// Input types for creating/updating tasks
export interface CreateTaskInput {
	title: string;
	description?: string;
	status: TaskStatus;
	priority: TaskPriority;
	dueDate?: string;
	categoryId: string;
	tags?: string[];
}

export interface UpdateTaskInput {
	title?: string;
	description?: string;
	status?: TaskStatus;
	priority?: TaskPriority;
	dueDate?: string;
	categoryId?: string;
	tags?: string[];
}

// Task status options - matching GraphQL enum
export const TASK_STATUS: Record<TaskStatus, TaskStatus> = {
	TODO: 'TODO',
	IN_PROGRESS: 'IN_PROGRESS',
	COMPLETED: 'COMPLETED'
} as const;

// Priority levels - matching GraphQL enum
export const PRIORITY: Record<TaskPriority, TaskPriority> = {
	LOW: 'LOW',
	MEDIUM: 'MEDIUM',
	HIGH: 'HIGH'
} as const;

// User related types
export interface User {
	id: string;
	name: string;
	email: string;
	createdAt: string;
	updatedAt: string;
}

export interface UpdateProfileInput {
	name?: string;
	email?: string;
}

export interface ChangePasswordInput {
	currentPassword: string;
	newPassword: string;
}
