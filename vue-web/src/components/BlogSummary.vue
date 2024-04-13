<template>
  <BlogCard>
    <div class="p-4">
      <div class="flex gap-6">
        <div class="flex flex-col">
          <div class="w-full">
            <ProfileIcon class="mx-auto" size="16" :src="ProfilePic"></ProfileIcon>
          </div>
          <p class="text-secondary">Author: Ken</p>
        </div>

        <div class="flex flex-col gap-2 hover:underline hover:cursor-pointer">
          <h1 class="text-4xl font-black">Title</h1>
          <h2>Sub-title</h2>
        </div>
      </div>

      <hr class="border-accent my-2" />
      
      <div class="content py-2 text-justify line-clamp-5 max-h-32 text-ellipsis">
        I recollect that my first exploit in squirrel-shooting was in a grove of tall walnut-trees that shades one side of the valley. I had wandered into it at noontime, when all nature is peculiarly quiet, and was startled by the roar of my own gun, as it broke the Sabbath stillness around and was prolonged and reverberated by the angry echoes.
        He ordered his regular breakfast. Two eggs sunnyside up, hash browns, and two strips of bacon. He continued to look at the menu wondering if this would be the day he added something new. This was also part of the routine. A few seconds of hesitation to see if something else would be added to the order before demuring and saying that would be all. It was the same exact meal that he had ordered every day for the past two years.
        Dave wasn't exactly sure how he had ended up in this predicament. He ran through all the events that had lead to this current situation and it still didn't make sense. He wanted to spend some time to try and make sense of it all, but he had higher priorities at the moment. The first was how to get out of his current situation of being naked in a tree with snow falling all around and no way for him to get down.
        He lifted the bottle to his lips and took a sip of the drink. He had tasted this before, but he couldn't quite remember the time and place it had happened. He desperately searched his mind trying to locate and remember where he had tasted this when the bicycle ran over his foot.
        Life isn't always beautiful. That was a lesson that Dan was learning. He also realized that life wasn't easy. This had come as a shock since he had lived a charmed life. He hated that this was the truth and he struggled to be happy knowing that his assumptions weren't correct. He wouldn't realize until much later in life that the difficult obstacles he was facing that were taking away the beauty in his life at this moment would ultimately make his life much more beautiful. All he knew was that at this moment was that life isn't always beautiful.
        Pink ponies and purple giraffes roamed the field. Cotton candy grew from the ground as a chocolate river meandered off to the side. What looked like stones in the pasture were actually rock candy. Everything in her dream seemed to be perfect except for the fact that she had no mouth.
        To the two friends, the treehouse was much more than a treehouse. It was a sanctuary away from the other kids where they could be themselves without being teased or bullied. It was their secret fortress hidden high in the branches of a huge oak that only they knew existed. At least that is what they thought. They were more than a little annoyed when their two younger sisters decided to turn the treehouse into a princess castle by painting the inside pink and putting glitter everywhere.
        The rain and wind abruptly stopped, but the sky still had the gray swirls of storms in the distance. Dave knew this feeling all too well. The calm before the storm. He only had a limited amount of time before all Hell broke loose, but he stopped to admire the calmness. Maybe it would be different this time, he thought, with the knowledge deep within that it wouldn't.
        No matter how hard he tried, he couldn't give her a good explanation about what had happened. It didn't even really make sense to him. All he knew was that he froze at the moment and no matter how hard he tried to react, nothing in his body allowed him to move. It was as if he had instantly become a statue and although he could see what was taking place, he couldn't move to intervene. He knew that wasn't a satisfactory explanation even though it was the truth.
        All he wanted was a candy bar. It didn't seem like a difficult request to comprehend, but the clerk remained frozen and didn't seem to want to honor the request. It might have had something to do with the gun pointed at his face.
      </div>

      <hr class="border-accent my-2" />

      <div class="flex justify-end">
        <p class="text-accent italic">Posted At: 22 Feb 2024 11:00:00 | Edited: 22 Feb 2024 11:00:00</p>
      </div>
    </div>
  </BlogCard>
</template>

<script setup lang="ts">
import ProfilePic from '@/assets/A36880CF-7EF8-49EC-A7B8-DB569A36AD3B-79271-00000BF53A8D6A23 Medium.jpeg';
</script>

<style type="poscss" scoped>
.content::first-letter {
  @apply text-2xl font-semibold;
}
</style>