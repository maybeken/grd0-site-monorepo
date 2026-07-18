<template>
  <div class="flex flex-col gap-4">
    <h1 class="text-4xl">{{ $t('notFound.title') }}</h1>
    <p class="text-lg">
      {{ $t('notFound.message') }}
    </p>
    <p class="text-xs text-secondary pt-8">{{ $t('notFound.hint') }}</p>
  </div>
</template>
