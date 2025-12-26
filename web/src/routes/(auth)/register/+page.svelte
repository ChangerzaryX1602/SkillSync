<script lang="ts">
  import { enhance } from "$app/forms";
  let { form } = $props();
  let loading = $state(false);
</script>

<svelte:head>
  <title>Register - SkillSync</title>
</svelte:head>

<div class="space-y-6">
  <div>
    <h3 class="text-lg font-medium leading-6 text-gray-900">
      Create a new account
    </h3>
  </div>

  {#if form?.error}
    <div class="rounded-md bg-red-50 p-4">
      <div class="flex">
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">
            {form.error}
          </h3>
        </div>
      </div>
    </div>
  {/if}

  {#if form?.success}
    <div class="rounded-md bg-green-50 p-4">
      <div class="flex">
        <div class="ml-3">
          <h3 class="text-sm font-medium text-green-800">
            Registration successful! Please <a href="/login" class="underline"
              >login</a
            >.
          </h3>
        </div>
      </div>
    </div>
  {/if}

  <form
    method="POST"
    action="?/register"
    use:enhance={() => {
      loading = true;
      return async ({ update }) => {
        loading = false;
        update();
      };
    }}
    class="space-y-6"
  >
    <div>
      <label for="username" class="block text-sm font-medium text-gray-700">
        Username
      </label>
      <div class="mt-1">
        <input
          id="username"
          name="username"
          type="text"
          required
          class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
        />
      </div>
    </div>

    <div>
      <label for="email" class="block text-sm font-medium text-gray-700">
        Email
      </label>
      <div class="mt-1">
        <input
          id="email"
          name="email"
          type="email"
          required
          class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
        />
      </div>
    </div>

    <div>
      <label for="password" class="block text-sm font-medium text-gray-700">
        Password
      </label>
      <div class="mt-1">
        <input
          id="password"
          name="password"
          type="password"
          required
          class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
        />
      </div>
    </div>

    <div>
      <button
        type="submit"
        disabled={loading}
        class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
      >
        {#if loading}
          Creating account...
        {:else}
          Register
        {/if}
      </button>
    </div>

    <div class="text-sm text-center">
      <a
        href="/login"
        class="font-medium text-indigo-600 hover:text-indigo-500"
      >
        Already have an account? Sign in
      </a>
    </div>
  </form>
</div>
