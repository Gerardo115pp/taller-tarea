<script>
    import searhbar_logo from '../vectors/boolean-search.svg';

    export let searchCallback;
    export let name = "generic";

    let counter =  null;

    const search = () => {
        const input_search = document.getElementById(`input-${name}`);
        if (searchCallback !== undefined) {
            searchCallback(input_search.value);
        }

    }

    const awaitForPause = () => {    
        if (counter !== null) {
            clearTimeout(counter);
        }
        counter = window.setTimeout(search, 800);
    }
</script>

<style>
    .gradiant-searchbar {
        display: flex;
        background: var(--theme-gradiant);
        width: 80%;
        height: 2.5rem;
        justify-content: space-between;
        align-items: center;
        padding: 0 2%;
        border-radius: 15px;
    }

    .gradiant-searchbar input {
        font-family: var(--main-font);
        background: none;
        width: 80%;
        color: var(--dark-color);
        font-size: 1.3rem;
        border: none;
        outline: none;
    }

    #searchbar-logo-container {
        cursor: pointer;
        display: flex;
        align-items: center;
    }

    :global(#searchbar-logo-container svg) {
        fill: var(--dark-color);
        width: 1.7rem;
    }
</style>

<div id={`searchbar-${name}`} class="gradiant-searchbar">
    <input id={`input-${name}`} on:keypress={awaitForPause} type="text">
    <div id="searchbar-logo-container">
        {@html searhbar_logo}
    </div>
</div>