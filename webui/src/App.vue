<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">Wasa-Photo</a>
	</header>

	<div class="container-fluid">
		<div class="row">
			<main class="col-md-9 col-lg-10 mx-auto px-md-4">
				<RouterView />
			</main>
		</div>
	</div>
</template>

<style>
</style>
