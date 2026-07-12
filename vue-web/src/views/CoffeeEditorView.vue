<template>
    <TastingForm :beans="beans" :equipment-list="equipment"></TastingForm>
</template>

<script setup>
import { getCoffeeBeans, getCoffeeEquipment } from '@/services/coffee'

const { data: beans } = getCoffeeBeans();
const { data: equipment } = getCoffeeEquipment();
</script>