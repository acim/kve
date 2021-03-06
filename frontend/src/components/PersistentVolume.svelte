<script lang="ts">
  import store from '../persistentVolumesStore'
  import toast from '../toastStore'
  import { fade } from 'svelte/transition'

  export let i: number

  const toggleReclaimPolicy = () => {
    try {
      const name = $store[i].volume.metadata.name
      let policy = 'Retain'
      if ($store[i].volume.spec.persistentVolumeReclaimPolicy === 'Retain') {
        policy = 'Delete'
      }
      store.toggleReclaimPolicy(name, policy)
      toast.set({ message: `Reclaim policy set to ${policy.toLowerCase()}` })
    } catch (err) {
      toast.set({ message: (err as Error).message })
    }
  }
</script>

<section transition:fade={{ delay: 0, duration: 400 }}>
  <h3>{$store[i].volume.metadata.name}</h3>
  <table>
    <tr>
      <td>Storage class</td>
      <td>{$store[i].volume.spec.storageClassName}</td>
    </tr>
    <tr>
      <td>Capacity</td>
      <td>{$store[i].volume.spec.capacity.storage}</td>
    </tr>
    <tr>
      <td>Mode</td>
      <td>{$store[i].volume.spec.accessModes}</td>
    </tr>
    <tr>
      <td>Reclaim policy</td>
      <td>{$store[i].volume.spec.persistentVolumeReclaimPolicy}</td>
    </tr>
    <tr>
      <td>Status</td>
      <td>{$store[i].volume.status.phase}</td>
    </tr>
    <tr>
      <td>Reference claim kind</td>
      <td>{$store[i].volume.spec.claimRef.kind}</td>
    </tr>
    {#if $store[i].volume.spec.hasOwnProperty('claimRef')}
      <tr>
        <td>Referencing claim name</td>
        <td>
          {$store[i].volume.spec.claimRef.namespace}/{$store[i].volume.spec
            .claimRef.name}
        </td>
      </tr>
    {/if}
    {#if $store[i].volume.hasOwnProperty('claim')}
      <tr>
        <td>Associated claim name</td>
        <td>
          {$store[i].claim.metadata.namespace}/{$store[i].claim.metadata.name}
        </td>
      </tr>
      <tr>
        <td>Associated claim capacity</td>
        <td>{$store[i].claim.status.capacity.storage}</td>
      </tr>
      <tr>
        <td>Associated claim modes</td>
        <td>{$store[i].claim.status.accessModes}</td>
      </tr>
      <tr>
        <td>Associated claim status</td>
        <td>{$store[i].claim.status.phase}</td>
      </tr>
      <tr>
        {#each $store[i].pods as pod, i (pod.metadata.uid)}
          <td>Mounted by pod</td>
          <td>{$store[i].pods[i].metadata.name}</td>
        {/each}
      </tr>
    {/if}
    {#if $store[i].volume.spec.hasOwnProperty('claimRef') && !$store[i].volume.hasOwnProperty('claim')}
      <tr>
        <td>Error</td>
        <td>No associated claim (orphan volume)</td>
      </tr>
    {/if}
  </table>
  <div class="right">
    {#if $store[i].volume.spec.hasOwnProperty('claimRef') && !$store[i].volume.hasOwnProperty('claim')}
      <button
        on:click={() => {
          alert('Not implemented')
        }}
      >
        Remove orphan claim
      </button>
    {:else}
      <button
        on:click={() => {
          toggleReclaimPolicy()
        }}
      >
        Toggle reclaim policy
      </button>
    {/if}
  </div>
</section>

<style>
  h3 {
    color: var(--color-accent);
  }
</style>
