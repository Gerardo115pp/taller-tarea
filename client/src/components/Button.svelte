<script>
    export let label = "button";
    export let isEnabled = true;
    export let onClick = () => {};
    export let isDanger = false;


</script>

<style>

    @keyframes Enable {
        0% {filter: grayscale(90%);}
        25% {filter: grayscale(60%);}
        50% {filter: grayscale(40%);}
        75% {filter: grayscale(20%);}
        100% {filter: grayscale(0%);}
    }

    .control-btn {
        font-family: var(--main-font);
        cursor: pointer;
        background-color: var(--theme-color);
        color: var(--dark-color);
        padding: 0 5%;
        border: none;
        border-radius: 15px;
        transition: all .4s ease-in-out;
    }

    .control-btn.red {
        background-color: var(--danger) !important;
    }

    .control-btn.enable {
        animation-name: Enable;
        animation-duration: 1s;
        animation-iteration-count: 1;
    }

    .control-btn.disable {
        filter: grayscale(90%);
    }
</style>

<div on:click={isEnabled ?  onClick : ()=>{}}
    class={`control-btn ${isDanger ? "red" : ""} ${isEnabled ? "enable" : "disable"}`}>{label}</div>