<template>
  <div class="pb-4">
    <table class="table-auto border-spacing-4 w-full">
      <thead class="border-b">
        <tr>
          <th class="py-2" v-for="(col, key) of columns">{{ typeof col.display_name === 'string' ? col.display_name :
            col.display_name(key) }}</th>
          <th class="py-2">Actions</th>
        </tr>
      </thead>
      <tbody class="border-b">
        <tr v-for="row of getPageItem(data || [], max_item, current_page)">
          <td class="py-2" v-for="(col, key) of columns">{{ col.formatter ? col.formatter(row[key]) : row[key] }}</td>
          <td class="flex gap-4 place-content-center py-2">
            <button v-for="action of actions" class="rounded-xl border px-4 py-px"
              @click="$emit(action.name, row[action.data_key])">{{ action.display_name }}</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>

  <div class="flex gap-4 text-sm place-content-end">
    <button class="disabled:hidden" @click="current_page = current_page - 1"
      :disabled="current_page <= 1">Previous</button>
    <span>Page</span>
    <button
      v-for="idx of getNumberOfPages(data || [])" @click="current_page = idx"
      class="disabled:hidden" :class="current_page == idx ? ['underline'] : []"
      :disabled="!(getNumberOfPages(data || []) <= 5 || idx == 1 || idx == getNumberOfPages(data || []) || current_page == idx) && (idx > current_page + 2 || idx < current_page - 2)"
    >
      {{ idx }}
    </button>
    <button class="disabled:hidden" @click="current_page = current_page + 1"
      :disabled="current_page >= getNumberOfPages(data || [])">Next</button>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import type { Columns, Action } from '@/interfaces/Editor';

interface Props {
  columns: Columns,
  data: any[],
  actions?: Action[],
}

const {
  columns,
  data,
  actions,
} = defineProps<Props>();

const $emit = defineEmits();

const max_item = ref(10);
const current_page = ref(1);

function getPageItem<T>(data: T[], max_item: number, current_page: number): T[] {
  if (!data) return [];

  const current_page_start = (current_page - 1) * max_item;
  const current_page_end = current_page * max_item;
  const current_page_min = current_page_start >= 0 ? current_page_start : 0;
  const current_page_max = data.length > current_page_end ? current_page_end : data.length;

  return data.filter((_, idx) => (idx >= current_page_min && idx <= current_page_max));
}

function getNumberOfPages(data: any[]) {
  return Math.max(1, Math.floor(data.length / max_item.value))
}
</script>