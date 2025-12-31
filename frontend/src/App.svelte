<script lang="ts">
  import Terminal from './lib/components/Terminal.svelte';
</script>

<main class="min-h-screen flex flex-col bg-base-200">
  <div class="navbar bg-base-100">
    <a class="font-bold text-xl mx-3 hover:underline" href="/">feverpanel</a>
  </div>

  <div class="flex-1 p-4">
    <div class="container mx-auto">
      <div class="card bg-base-100 ">
        <div class="card-body">
          <h2 class="card-title">Server Console</h2>
          <Terminal />
        </div>
      </div>
    </div>
  </div>
</main>
