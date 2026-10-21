<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

interface Point {
  x: number
  y: number
}

interface Particle {
  x: number
  y: number
  vx: number
  vy: number
  radius: number
  alpha: number
  phase: number
  twinkle: number
  delay: number
  trail: Point[]
  trailRemaining: number
  active: boolean
}

const MIN_PARTICLES = 3
const MAX_PARTICLES = 12
const SPAWN_DELAY = 5
const EDGE_MARGIN = 8
const TRAIL_LIFETIME = 5

const canvas = ref<HTMLCanvasElement | null>(null)
const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)')

let context: CanvasRenderingContext2D | null = null
let particles: Particle[] = []
let frame = 0
let resizeFrame = 0
let width = 0
let height = 0
let lastTime = 0

function createParticle(): Particle {
  return {
    x: 0,
    y: 0,
    vx: 0,
    vy: 0,
    radius: 1,
    alpha: 0.3,
    phase: Math.random() * Math.PI * 2,
    twinkle: 0.4 + Math.random() * 1.2,
    delay: Math.random() * SPAWN_DELAY,
    trail: [],
    trailRemaining: 0,
    active: false
  }
}

function activate(particle: Particle) {
  const depth = Math.random()
  const speed = 100 + depth * 400
  const edge = Math.floor(Math.random() * 4)

  let angle = 0

  if (edge === 0) {
    particle.x = -EDGE_MARGIN
    particle.y = Math.random() * height
    angle = Math.random() * Math.PI - Math.PI / 2
  } else if (edge === 1) {
    particle.x = width + EDGE_MARGIN
    particle.y = Math.random() * height
    angle = Math.PI / 2 + Math.random() * Math.PI
  } else if (edge === 2) {
    particle.x = Math.random() * width
    particle.y = -EDGE_MARGIN
    angle = Math.random() * Math.PI
  } else {
    particle.x = Math.random() * width
    particle.y = height + EDGE_MARGIN
    angle = Math.PI + Math.random() * Math.PI
  }

  particle.vx = Math.cos(angle) * speed
  particle.vy = Math.sin(angle) * speed
  particle.radius = 0.5 + depth * 1.7
  particle.alpha = 0.12 + depth * 0.4
  particle.delay = 0
  particle.trail = [{ x: particle.x, y: particle.y }]
  particle.trailRemaining = 0
  particle.active = true
}

function drawTrail(particle: Particle, alpha: number) {
  if (!context || particle.trail.length < 2) return

  context.globalAlpha = alpha
  context.lineWidth = particle.radius * 2
  context.beginPath()
  context.moveTo(particle.trail[0].x, particle.trail[0].y)

  for (let i = 1; i < particle.trail.length; i++) {
    context.lineTo(particle.trail[i].x, particle.trail[i].y)
  }

  context.stroke()
}

function draw(time: number, delta = 0) {
  if (!context) return

  if (reducedMotion.matches) {
    context.clearRect(0, 0, width, height)
    context.fillStyle = 'rgb(245 245 245)'

    for (const particle of particles) {
      if (!particle.active) continue

      context.globalAlpha = particle.alpha
      context.beginPath()
      context.arc(particle.x, particle.y, particle.radius, 0, Math.PI * 2)
      context.fill()
    }

    context.globalAlpha = 1
    return
  }

  context.clearRect(0, 0, width, height)
  context.strokeStyle = 'rgb(245 245 245)'
  context.lineCap = 'round'

  for (const particle of particles) {
    if (particle.trailRemaining > 0) {
      particle.trailRemaining = Math.max(0, particle.trailRemaining - delta)

      if (particle.trailRemaining === 0) particle.trail = []
    }

    if (particle.active) {
      particle.x += particle.vx * delta
      particle.y += particle.vy * delta

      if (
        particle.x < -EDGE_MARGIN ||
        particle.x > width + EDGE_MARGIN ||
        particle.y < -EDGE_MARGIN ||
        particle.y > height + EDGE_MARGIN
      ) {
        particle.active = false
        particle.delay = Math.random() * SPAWN_DELAY
        particle.trailRemaining = TRAIL_LIFETIME
      } else {
        particle.trail.push({ x: particle.x, y: particle.y })
      }
    } else if (particle.trailRemaining <= 0) {
      if (particle.delay > 0) particle.delay -= delta
      else activate(particle)
    }

    if (particle.trail.length > 1) {
      const fade = particle.active ? 1 : particle.trailRemaining / TRAIL_LIFETIME
      const flicker = 0.6 + 0.4 * Math.sin(time * particle.twinkle + particle.phase)

      drawTrail(particle, particle.alpha * flicker * fade)
    }
  }

  context.globalAlpha = 1
}

function tick(now: number) {
  const delta = Math.min((now - lastTime) / 1000, 0.05)
  lastTime = now

  draw(now / 1000, delta)
  frame = requestAnimationFrame(tick)
}

function start() {
  if (frame || reducedMotion.matches) return

  lastTime = performance.now()
  frame = requestAnimationFrame(tick)
}

function stop() {
  cancelAnimationFrame(frame)
  frame = 0
}

function resize() {
  const element = canvas.value
  if (!element) return

  const dpr = Math.min(window.devicePixelRatio || 1, 2)
  width = element.clientWidth
  height = element.clientHeight
  element.width = Math.round(width * dpr)
  element.height = Math.round(height * dpr)

  context = element.getContext('2d')
  if (!context) return

  context.setTransform(dpr, 0, 0, dpr, 0, 0)

  const count =
    MIN_PARTICLES + Math.floor(Math.random() * (MAX_PARTICLES - MIN_PARTICLES + 1))
  particles = Array.from({ length: count }, createParticle)

  if (reducedMotion.matches) {
    particles.forEach(activate)
    draw(0)
  }
}

function handleResize() {
  cancelAnimationFrame(resizeFrame)
  resizeFrame = requestAnimationFrame(resize)
}

function handleVisibilityChange() {
  if (document.visibilityState === 'hidden') stop()
  else start()
}

onMounted(() => {
  resize()
  window.addEventListener('resize', handleResize)
  document.addEventListener('visibilitychange', handleVisibilityChange)
  start()
})

onBeforeUnmount(() => {
  stop()
  cancelAnimationFrame(resizeFrame)
  window.removeEventListener('resize', handleResize)
  document.removeEventListener('visibilitychange', handleVisibilityChange)
})
</script>

<template>
  <canvas
    ref="canvas"
    class="pointer-events-none fixed inset-0 -z-10 h-full w-full"
    aria-hidden="true"
  ></canvas>
  <div
    class="pointer-events-none fixed inset-0 -z-10 bg-black/20 backdrop-blur-[2px]"
    aria-hidden="true"
  ></div>
</template>
