<script lang="ts">
	import { userClient } from '$lib/client';
	import { onMount } from 'svelte';

	let users = $state<any[]>([]);

	onMount(async () => {
		const res = await userClient.getAllUsers({});
		users = res.users ?? [];
	});
</script>

<div class="flex h-screen w-screen flex-col items-center justify-center">
	{#each users as user}
		<h1>
			{user.email}
		</h1>
	{/each}
</div>
