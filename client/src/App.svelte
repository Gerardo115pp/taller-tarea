<script>
    import Router from 'svelte-spa-router'
    import { routes } from './routes';

</script>

<style>
    @import url('https://fonts.googleapis.com/css2?family=Rhodium+Libre&display=swap');

    :global(:root) {
        --theme-color: #86E84A;
        --dimtheme-color: #b5ff87;
        --theme-gradiant: linear-gradient(90.58deg, #86E84A 7.24%, #57B61D 93.27%);;
        --dark-color: #001023;
        --danger:#fd1b1b;
        --main-font: 'Rhodium Libre';
    }

    :global(body) {
        font-family: 'Rhodium Libre';
        overflow-y: hidden;
        background-color: var(--dark-color);
        color: var(--theme-color);
    }

    :global(h1, h2, h3, h4, h5, h6) {
        margin: 0;
    }
</style>

<div id="main-app-container">
    <Router {routes}/>
</div>