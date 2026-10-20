<template>
  <BlogCard>
    <div class="py-4 px-4 md:px-16">
      <div class="flex gap-6">
        <div class="flex flex-col">
          <div class="h-24 w-24">
            <ProfileIcon :email="$props.entry?.author?.email ?? ''"></ProfileIcon>
          </div>
          <Skeleton h="xs" w="sm" :loading="$props.loading">
            <p class="text-secondary text-center">{{ $props.entry?.author?.display_name }}</p>
          </Skeleton>
        </div>

        <div class="flex flex-col gap-2 hover:underline my-auto w-full">
          <router-link :to="`/article/${$props.entry?.uri}`">
            <Skeleton h="md" w="1/3" :loading="$props.loading">
              <h1 class="font-black">{{ $props.entry?.title || $props.entry?.uri }}</h1>
            </Skeleton>
            <Skeleton h="sm" w="2/3" :loading="$props.loading">
              <h2>{{ $props.entry?.subtitle }}</h2>
            </Skeleton>
          </router-link>
        </div>
      </div>

      <hr class="border-accent my-2" />

      <div v-if="$props.loading || $props.entry?.cover" class="px-2 py-2">
        <Skeleton h="lg" w="full" :loading="$props.loading">
          <CDNImage
            v-if="$props.entry?.cover"
            class="w-full max-h-[50vh] rounded-xl object-cover"
            :uri="$props.entry?.cover ?? ''"
            :resolution="1280"
            :quality="80"
          ></CDNImage>
        </Skeleton>
      </div>

      <div class="relative px-2 py-2">
        <div class="content text-justify max-h-[66vh] min-h-[24vh] overflow-hidden">
          <Skeleton h="md" w="full" :loading="$props.loading">
            <MarkdownDisplay
              :md="$props.entry?.summary ?? ''"
              :lazy_loading="false"
            ></MarkdownDisplay>
          </Skeleton>
        </div>

        <div
          v-if="!$props.loading"
          class="absolute bottom-0 left-0 h-28 w-full z-10 rounded-xl bg-linear-to-t from-background from-33%"
        >
          <router-link :to="`/article/${$props.entry?.uri}`">
            <div class="flex flex-col-reverse h-28 w-full hover:cursor-pointer hover:underline">
              <div class="mx-auto my-8">
                <p class="text-xl">{{ $t('article.index.continueReading') }}</p>
              </div>
            </div>
          </router-link>
        </div>
      </div>

      <hr class="border-accent my-2" />

      <div class="flex justify-end">
        <Skeleton h="sm" w="xl" :loading="$props.loading">
          <p class="text-accent italic">
            <span v-if="$props.entry?.published_at">
              {{
                $t('article.header.publishedAt', {
                  date: dayjs($props.entry?.published_at ?? '').calendar()
                })
              }}
            </span>
            <span v-if="$props.entry?.published_at && $props.entry?.updated_at"> | </span>
            <span v-if="$props.entry?.updated_at">
              {{
                $t('article.header.updatedAt', {
                  date: dayjs($props.entry?.updated_at ?? '').calendar()
                })
              }}
            </span>
          </p>
        </Skeleton>
      </div>
    </div>
  </BlogCard>
</template>

<script setup lang="ts">
import dayjs from 'dayjs'
import calendar from 'dayjs/plugin/calendar'

import type { ArticleListEntry } from '@/interfaces/Article'

dayjs.extend(calendar)

interface Props {
  entry?: ArticleListEntry
  loading: boolean
}

const $props = defineProps<Props>()
</script>

<style lang="postcss" scoped>
.content:first-letter {
  @apply text-2xl font-semibold;
}
</style>
