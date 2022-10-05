<script>
	import Withdraw from './Withdraw.svelte';
  import { SvelteToast } from '@zerodevx/svelte-toast'
</script>

<main>  
  <div class="item"><Withdraw /></div>
  <div class="item"><Withdraw /></div>
</main>
<SvelteToast/>

<style>
  main {
    display: flex;
    flex-direction: column;
  }

  .item {
    justify-content: center;
    justify-items: center;
    flex-grow: 1;
  }

  @media (min-width: 1000px) {
    main {
      flex-direction: row;
    }
  }

</style>