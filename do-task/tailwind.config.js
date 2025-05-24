/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	darkMode: 'class',
	theme: {
		extend: {
			colors: {
				primary: '#3B82F6',
        accent: '#6366F1',
        'base-100': '#F9FAFB',
        'text-main': '#111827',
        muted: '#6B7280',
        success: '#10B981',
        warning: '#F59E0B',
        error: '#EF4444',
        info: '#0EA5E9',
        disabled: '#E5E7EB',
			},
			fontFamily: {
				sans: ['Inter', 'sans-serif'],
				heading: ['Poppins', 'sans-serif'],  
			}
		}
	},
	plugins: []
};
