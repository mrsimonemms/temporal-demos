<!--
  ~ Copyright 2025 Simon Emms <simon@simonemms.com>
  ~
  ~ Licensed under the Apache License, Version 2.0 (the "License");
  ~ you may not use this file except in compliance with the License.
  ~ You may obtain a copy of the License at
  ~
  ~     http://www.apache.org/licenses/LICENSE-2.0
  ~
  ~ Unless required by applicable law or agreed to in writing, software
  ~ distributed under the License is distributed on an "AS IS" BASIS,
  ~ WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  ~ See the License for the specific language governing permissions and
  ~ limitations under the License.
-->

<script lang="ts">
  import Navbar from '$lib/components/navbar.svelte';

  let { children } = $props();
</script>

<svelte:head>
  <title>The Grub Stop</title>
</svelte:head>

<div class="container">
  <Navbar />
</div>

<section class="section">
  <div class="container">
    <h1 class="title">Welcome to The Grub Stop</h1>
    <h2 class="subtitle">Buy yourself some delicious grub</h2>

    {@render children()}
  </div>
</section>
