import { writable } from 'svelte/store';

export const settings = writable({
  theme: 'light',
});