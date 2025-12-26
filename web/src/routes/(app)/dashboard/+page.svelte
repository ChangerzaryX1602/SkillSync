<script lang="ts">
  let { data } = $props();
</script>

<div
  class="bg-white overflow-hidden shadow rounded-lg divide-y divide-gray-200"
>
  <div class="px-4 py-5 sm:px-6">
    <h1 class="text-2xl font-semibold text-gray-900">Dashboard</h1>
    <p class="mt-1 max-w-2xl text-sm text-gray-500">
      Welcome to SkillSync Operations.
    </p>
  </div>
  <div class="px-4 py-5 sm:p-6">
    <p>
      Logged in as: <strong>{data.user?.username}</strong> ({data.user?.email})
    </p>
    <div class="mt-4">
      <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
        <!-- Stats or shortcuts could go here -->
        <div class="bg-indigo-50 overflow-hidden shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <dt class="text-sm font-medium text-gray-500 truncate">
              Quick Actions
            </dt>
            <dd class="mt-1 text-3xl font-semibold text-gray-900">
              <a
                href="/users"
                class="text-sm text-indigo-600 hover:text-indigo-900"
                >Manage Users &rarr;</a
              >
            </dd>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
