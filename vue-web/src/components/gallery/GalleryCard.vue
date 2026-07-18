<template>
  <div class="flex flex-col gap-2 justify-end p-1 pb-2 rounded-xl bg-shade text-secondary text-sm">
    <div class="relative">
      <CDNImage
        class="rounded-xl w-full object-cover aspect-square lg:cursor-zoom-in"
        :resolution="768"
        :quality="75"
        :uri="$props.image?.filename ? `${$props.image?.collection}/${$props.image?.filename}` : ''"
        @click="zoomImage(`${$props.image?.collection}/${$props.image?.filename}`)"
      ></CDNImage>
      <div class="absolute right-0 bottom-0 ml-auto">
        <button
          class="px-2 mr-px mb-px border border-accent bg-background/75 rounded-xl disabled:brightness-50 disabled:bg-accent/75 disabled:border-secondary"
          :title="$t('common.download')"
          :disabled="loading"
          @click="
            downloadImage(
              $props.image?.filename ? `${$props.image?.collection}/${$props.image?.filename}` : ''
            )
          "
        >
          <Icon icon="mynaui:cloud-download" height="3rem" />
        </button>
      </div>
    </div>
    <div v-if="details?.description || loading" class="px-1">
      <Skeleton h="md" w="full" :loading="loading">
        <p class="text-justify">{{ details?.description }}</p>
      </Skeleton>
    </div>
    <div class="px-1 grow">
      <Skeleton h="sm" w="full" :loading="loading">
        <p class="text-center">
          {{ $props.image?.exif?.equipment?.camera }} {{ $props.image?.exif?.equipment?.lens }}
        </p>
      </Skeleton>
      <Skeleton h="sm" w="full" :loading="loading">
        <p class="text-center"></p>
        <p class="text-center">
          ISO {{ $props.image?.exif?.iso }} | {{ $props.image?.exif?.fstop }} |
          {{ $props.image?.exif?.shutter }}s
        </p>
      </Skeleton>
    </div>
    <div class="px-1">
      <Skeleton class="ml-auto" h="sm" w="2/3" :loading="loading">
        <p class="text-center"></p>
        <p class="text-right">
          {{
            details?.tz_adjustment
              ? dayjs($props.image?.exif?.datetime).add(details?.tz_adjustment, 'h').format('LLL')
              : dayjs($props.image?.exif?.datetime).format('LLL')
          }}
        </p>
      </Skeleton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import dayjs from 'dayjs'
import LocalizedFormat from 'dayjs/plugin/localizedFormat'
dayjs.extend(LocalizedFormat)

import { getGalleryDetail } from '@/services/gallery'

import type { Asset, GalleryDetail, PaginatedResponse } from '@/interfaces/Gallery'
import Skeleton from '../shared/Skeleton.vue'

interface Props {
  image?: Asset
  loading: boolean
}

const $props = defineProps<Props>()

const details = ref<GalleryDetail>()
const response = getGalleryDetail($props.image?.collection)
const empty_response: PaginatedResponse<GalleryDetail> = {
  data: [],
  total: 0,
  page: 1,
  page_size: 40,
  total_pages: 0
}
const gallery_details = response?.data || ref(empty_response)

const ASSET_URL = import.meta.env.VITE_ASSETS_URL

const $emit = defineEmits<{
  zoom: [uri?: string]
}>()

watch(gallery_details, (newVal) => {
  if (!newVal) return

  const file_rule = newVal.data.find((item) => item.filename === $props.image?.filename)
  const collection_rule = newVal.data.find((item) => item.filename === '*')

  if (collection_rule) {
    details.value = collection_rule

    if (file_rule) {
      details.value = file_rule
    }
  } else if (file_rule) {
    details.value = file_rule
  }
})

function zoomImage(val: string) {
  $emit('zoom', val)
}

function downloadImage(val: string) {
  // create element <a> for download PDF
  const link = document.createElement('a')
  link.href = `${ASSET_URL}${val}?download`
  link.target = '_blank'

  const path_split = val.split('/').reverse()
  link.download = path_split[0] ?? ''

  // Simulate a click on the element <a>
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}
</script>
