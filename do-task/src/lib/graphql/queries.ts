import * as Apollo from '@apollo/client/core';
const { gql } = Apollo;

// Query all tasks
export const GET_TASKS = gql`
	query GetTasks {
		tasks {
			id
			title
			description
			status
			priority
			dueDate
			createdAt
			updatedAt
			category {
				id
				name
			}
			tags
		}
	}
`;

// Query single task
export const GET_TASK = gql`
	query GetTask($id: ID!) {
		task(id: $id) {
			id
			title
			description
			status
			priority
			dueDate
			createdAt
			updatedAt
			category {
				id
				name
			}
			tags
		}
	}
`;

// Query all categories
export const GET_CATEGORIES = gql`
	query GetCategories {
		categories {
			id
			name
		}
	}
`;

// Create task mutation
export const CREATE_TASK = gql`
	mutation CreateTask($input: CreateTaskInput!) {
		createTask(input: $input) {
			id
			title
			description
			status
			priority
			dueDate
			createdAt
			updatedAt
			category {
				id
				name
			}
			tags
		}
	}
`;

// Update task mutation
export const UPDATE_TASK = gql`
	mutation UpdateTask($id: ID!, $input: UpdateTaskInput!) {
		updateTask(id: $id, input: $input) {
			id
			title
			description
			status
			priority
			dueDate
			createdAt
			updatedAt
			category {
				id
				name
			}
			tags
		}
	}
`;

// Delete task mutation
export const DELETE_TASK = gql`
	mutation DeleteTask($id: ID!) {
		deleteTask(id: $id)
	}
`;

// Update task status mutation
export const UPDATE_TASK_STATUS = gql`
	mutation UpdateTaskStatus($id: ID!, $status: TaskStatus!) {
		updateTaskStatus(id: $id, status: $status) {
			id
			status
			updatedAt
		}
	}
`;

// Create category mutation
export const CREATE_CATEGORY = gql`
	mutation CreateCategory($name: String!) {
		createCategory(name: $name) {
			id
			name
		}
	}
`;

// Update category mutation
export const UPDATE_CATEGORY = gql`
	mutation UpdateCategory($id: ID!, $name: String!) {
		updateCategory(id: $id, name: $name) {
			id
			name
		}
	}
`;

// Delete category mutation
export const DELETE_CATEGORY = gql`
	mutation DeleteCategory($id: ID!) {
		deleteCategory(id: $id)
	}
`;

// Authentication queries and mutations

export const ME_QUERY = gql`
	query Me {
		me {
			id
			name
			email
			createdAt
			updatedAt
		}
	}
`;

export const LOGIN_MUTATION = gql`
	mutation Login($input: LoginInput!) {
		login(input: $input) {
			user {
				id
				name
				email
				createdAt
				updatedAt
			}
			token
		}
	}
`;

export const REGISTER_MUTATION = gql`
	mutation Register($input: RegisterInput!) {
		register(input: $input) {
			user {
				id
				name
				email
				createdAt
				updatedAt
			}
			token
		}
	}
`;

// Profile management mutations
export const UPDATE_PROFILE_MUTATION = gql`
	mutation UpdateProfile($input: UpdateProfileInput!) {
		updateProfile(input: $input) {
			id
			name
			email
			createdAt
			updatedAt
		}
	}
`;

export const CHANGE_PASSWORD_MUTATION = gql`
	mutation ChangePassword($input: ChangePasswordInput!) {
		changePassword(input: $input)
	}
`;
