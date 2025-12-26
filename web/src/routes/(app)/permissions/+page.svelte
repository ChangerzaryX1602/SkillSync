<script lang="ts">
  import { Search } from "lucide-svelte";
  let { data } = $props();
</script>

<div class="sm:flex sm:items-center">
  <div class="sm:flex-auto">
    <h1 class="text-xl font-semibold text-gray-900">Permissions</h1>
    <p class="mt-2 text-sm text-gray-700">
      A list of all available permissions in the system.
    </p>
  </div>
</div>

<div class="mt-8 flex flex-col">
  <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
    <div class="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8">
      <div
        class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg"
      >
        <table class="min-w-full divide-y divide-gray-300">
          <thead class="bg-gray-50">
            <tr>
              <th
                scope="col"
                class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6"
                >ID</th
              >
              <th
                scope="col"
                class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
                >Group</th
              >
              <th
                scope="col"
                class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
                >Name</th
              >
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 bg-white">
            {#if data.permissions && data.permissions.length > 0}
              {#each data.permissions as perm}
                <tr>
                  <td
                    class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6"
                    >{perm.id}</td
                  >
                  <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500"
                    >{perm.group}</td
                  >
                  <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500"
                    >{perm.name}</td
                  >
                </tr>
              {/each}
            {:else}
              <tr>
                <td
                  colspan="3"
                  class="px-6 py-4 text-center text-sm text-gray-500"
                  >No permissions found.</td
                >
              </tr>
            {/if}
          </tbody>
        </table>
      </div>
    </div>
  </div>
</div>
