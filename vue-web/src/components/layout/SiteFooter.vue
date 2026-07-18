<template>
  <hr class="border-accent" />
  <div class="flex justify-center pt-2">
    <div class="text-center">
      <p class="text-secondary italic text-sm">
        {{ $t('footer.copyright', { year: dayjs().format('YYYY') }) }}
      </p>
      <p class="text-secondary text-xs">
        <a href="https://creativecommons.org/licenses/by-nc/4.0/" target="_blank"
          >{{ $t('footer.contentLicense') }}</a
        >
        <span class="select-none"> | </span>
        <a href="http://gitlab.com/grd0.net/grd0-site-monorepo" target="_blank"
          >{{ $t('footer.sourceLicense') }}</a
        >
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import dayjs from 'dayjs'
</script>
