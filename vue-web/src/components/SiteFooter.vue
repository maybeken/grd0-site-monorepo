<template>
  <hr class="border-accent" />
  <div class="flex justify-center pt-2">
    <div class="text-center">
      <p class="text-secondary italic">
      <a href="https://creativecommons.org/licenses/by-nc/4.0/" target="_blank">CC BY-NC 4.0</a> @ grd0.net
      </p>
      <p class="text-secondary text-xs">
        Est. 2020 - {{ dayjs().format("YYYY") }} | <a href="http://gitlab.com/grd0.net/grd0-site-monorepo" target="_blank">Source Available @ GitLab</a>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import dayjs from 'dayjs';
</script>