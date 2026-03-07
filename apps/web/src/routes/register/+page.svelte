<script lang="ts">
	import { Button, Label } from 'bits-ui';
	import LoaderCircle from '@lucide/svelte/icons/loader-circle';
	import * as z from 'zod';
	import { toast } from 'svelte-sonner';
	import { authClient } from '$lib/client';
	import type { RegisterResponse, RegisterResult } from '$lib/gen/auth/auth_pb';

	let registerLoading = $state<boolean>(false);

	const Register = z.object({
		email: z.email(),
		password: z.string().min(6).max(24)
	});

	type RegisterData = z.infer<typeof Register>;
	let registerData = $state<RegisterData>({
		email: '',
		password: ''
	});

	function handleError() {
		const form = document.getElementById('register-form');
		const body = document.body;
		if (form) {
			form.classList.remove('animate-headshake');
			body.classList.remove('bg-red-50');
			void form.offsetWidth;
			form.classList.add('animate-headshake');
			body.classList.add('bg-red-50');

			form.addEventListener(
				'animationend',
				() => {
					setTimeout(() => {
						form.classList.remove('animate-headshake');
						body.classList.remove('bg-red-50');
					}, 1700);
				},
				{
					once: true
				}
			);
		}
	}

	function handleFail(result: z.ZodSafeParseResult<RegisterData>) {
		if (result.success) return true;
		console.log(result);

		const message = result.error.issues[0].message;
		toast.error(message, {
			duration: 2500
		});

		handleError();

		registerLoading = false;
		return false;
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		registerLoading = true;

		const result = Register.safeParse(registerData);
		if (!handleFail(result)) return;

		try {
			const data: RegisterResult = await authClient.register({
				email: registerData.email,
				password: registerData.password
			});

			if (data.result.case === 'error') {
				toast.error(data.result.value.message);
				handleError();
			} else {
				toast.success('Account created!');
			}
		} catch (err: any) {
			toast.error(err?.message || 'Registration failed');
		} finally {
			registerLoading = false;
		}
	}
</script>

<div id="register-form" class="flex h-screen flex-col items-center justify-center gap-8">
	<h1 class="px-12 font-serif text-3xl">Create a Shulker account</h1>

	<form onsubmit={handleSubmit} class="flex w-full flex-col gap-6">
		<div class="flex w-full flex-col gap-2">
			<Label.Root for="email" class="w-full">Email</Label.Root>
			<input
				bind:value={registerData.email}
				disabled={registerLoading}
				id="email"
				type="email"
				placeholder="Email Address"
				class="round h-10 w-full border px-2.5 text-sm disabled:opacity-70"
			/>
		</div>
		<div class="flex w-full flex-col gap-2">
			<Label.Root for="password" class="w-full">Password</Label.Root>
			<input
				bind:value={registerData.password}
				disabled={registerLoading}
				id="password"
				type="password"
				placeholder="Password"
				class="round h-10 w-full border px-2.5 text-sm disabled:opacity-70"
			/>
		</div>
		<Button.Root
			type="submit"
			disabled={registerLoading}
			class="round flex h-10 flex-row items-center justify-center gap-2 bg-primary text-sm text-primary-foreground hover:bg-primary/90 disabled:opacity-70"
		>
			{#if registerLoading}
				<LoaderCircle class="size-4 animate-spin text-primary-foreground/60" />
			{/if}
			Register
		</Button.Root>
		<Button.Root
			disabled={registerLoading}
			href="/login"
			class="text-center text-muted-foreground hover:text-primary hover:no-underline disabled:opacity-70"
			>Already signed up? Login here</Button.Root
		>
	</form>
</div>
