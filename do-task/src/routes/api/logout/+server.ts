import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ cookies }) => {
	try {
		// Clear the accessToken cookie with multiple configurations to ensure deletion
		cookies.delete('accessToken', {
			path: '/',
			httpOnly: true,
			secure: false, // Set to true in production with HTTPS
			sameSite: 'lax'
		});

		// Also try to clear with different path configurations
		cookies.delete('accessToken', {
			path: '/',
			httpOnly: true,
			secure: false,
			sameSite: 'lax'
		});

		// Also try to clear any other potential auth cookies
		cookies.delete('refreshToken', {
			path: '/',
			httpOnly: true,
			secure: false,
			sameSite: 'lax'
		});

		// Set the cookie to empty value with past expiration as backup
		cookies.set('accessToken', '', {
			path: '/',
			httpOnly: true,
			secure: false,
			sameSite: 'lax',
			maxAge: 0
		});

		return json({ success: true, message: 'Logged out successfully' });
	} catch (error) {
		console.error('Logout error:', error);
		return json({ success: false, message: 'Logout failed' }, { status: 500 });
	}
};
