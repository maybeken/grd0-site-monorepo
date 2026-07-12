<template>
  <hr class="border-accent" />
  <div class="flex justify-center pt-2">
    <div class="text-center">
      <p class="text-secondary italic text-sm">
        Est. 2020 - {{ dayjs().format('YYYY') }} @ <a href="https://grd0.net">grd0.net</a>
      </p>
      <p class="text-secondary text-xs">
        <a href="https://creativecommons.org/licenses/by-nc/4.0/" target="_blank"
          >Site content licensed under CC BY-NC 4.0</a
        >
        <span class="select-none"> | </span>
        <a href="http://gitlab.com/grd0.net/grd0-site-monorepo" target="_blank"
          >Source code available under CC BY-SA 4.0 on GitLab</a
        >
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import dayjs from 'dayjs'
</script>
