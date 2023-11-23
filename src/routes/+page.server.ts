import { AUTH_EMAIL, AUTH_PASSWORD } from '$env/static/private';
import { redirect } from '@sveltejs/kit';

export const actions = {
	login: async ({ request }) => {
		const data = await request.formData();
		const email = data.get('email');
		const password = data.get('password');

		if (email === AUTH_EMAIL && password === AUTH_PASSWORD) {
			throw redirect(303, '/tools');
		} else {
			throw redirect(308, '/');
		}
	}
};
