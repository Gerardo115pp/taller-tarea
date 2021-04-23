<script>
    export let options;
    export let name;
</script>

<style>
    .select-container {
        border: 2px solid var(--theme-color);
        border-radius: 9999px;
        padding: 1.5vh 4.2vh;
    }

    .select-container select {
        background: none;
        font-family: var(--main-font);
        color: var(--theme-color);
        font-size: 1.2rem;
        text-align: center;
        border: none;
        outline: none;
        appearance: none;
        -moz-appearance: none;
        -webkit-appearance: none;
    }

    .select-container option {
        background: var(--dark-color);
        text-align: center;
    }
</style>

<div class="select-container">
    <select name={name} id={`ft-${name}`}>
        {#each options as o }
            <option value={o.value}>{o.name}</option>
        {/each}
    </select>
</div>