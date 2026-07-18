<template>
  <div class="flex flex-col gap-4">
    <div>
      <p class="text-2xl font-bold">{{ $t('coffee.view.title') }}<CursorBlink /></p>
      <p class="italic font-light">
        {{ $t('coffee.view.subtitle') }}
      </p>
    </div>

    <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div v-for="i in 6" :key="i" class="aspect-square rounded-xl overflow-hidden">
        <Skeleton :loading="true" h="full" w="full" />
      </div>
    </div>

    <div v-else-if="!tastings || tastings.length === 0" class="text-center py-16 opacity-60">
      <p class="text-lg">{{ $t('coffee.view.noTastings') }}</p>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <TastingCard v-for="tasting in tastings" :key="tasting.id" :tasting="tasting" />
    </div>

    <div class="hidden md:block ">
      <p class="text-xs">
        {{ $t('coffee.view.equipmentCredit', {
          author: $t('coffee.view.equipmentCreditAuthor'),
          source: $t('coffee.view.equipmentCreditSource')
        }) }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getCoffeeTastings } from '@/services/coffee'
import TastingCard from '@/components/coffee/TastingCard.vue'

const { loading, data: tastings } = getCoffeeTastings()
</script>
