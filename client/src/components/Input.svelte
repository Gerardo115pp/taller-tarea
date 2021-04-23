<script>
    export let input_id = "";
    export let input_type = "text";
    export let input_placeholder;
    export let onKeypressed;
    export let onBlur;


    const handleOutsideClickDetected = e => {
        if(e.currentTarget === e.target) {
            e.target.getElementsByTagName("input")[0].focus()
        }
    }
</script>

<style>
    .input-container {
        cursor: text;
        border: 2px solid var(--theme-color);
        border-radius: 9999px;
        padding: 1.5vh 4.2vh;
    }   

    .input-container input {
        font-family: var(--main-font);
        display: flex;
        background: none;
        border: none;
        color: var(--theme-color);
        font-size: 1.2rem;
        padding: 0;
        align-items: center;
        outline: none;
    }

    .input-container input::placeholder {
        font-family: 'Rhodium Libre';
        color: var(--dimtheme-color);
        text-transform: capitalize;
    }

</style>

<div on:click={handleOutsideClickDetected} class="input-container">
    <input id={input_id} type={input_type}
        placeholder={input_placeholder}
        on:keydown={onKeypressed}
        on:blur={onBlur}
    />
</div>