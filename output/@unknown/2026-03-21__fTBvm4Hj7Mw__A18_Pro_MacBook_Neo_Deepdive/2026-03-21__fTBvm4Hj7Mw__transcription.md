---
title: "2026-03-21 A18 Pro & MacBook Neo Deep-dive (transcription)"
video_id: "fTBvm4Hj7Mw"
url: "https://www.youtube.com/watch?v=fTBvm4Hj7Mw"
channel: "@unknown"
channel_name: ""
upload_date: "2026-03-21"
duration: "14:58"
language: "en"
tags: []
categories:
  - "Science & Technology"
subtitle_type: "manual"
processed_at: "2026-03-22T13:55:53+08:00"
---

This is a MacBook Air M3 logic board. Here’s the\h
Apple M3 chip and directly next to it are two\h\h
LPDDR5x DRAM modules. That’s the shared memory.\h
A bit further away, but still really close,\h\h
are two NAND chips for the SSD. The other\h
silver chip isn’t NAND, that’s for WiFi\h\h
and Bluetooth. A tightly integrated system\h
and honestly amazing engineering from Apple.\h\h
But it’s nothing compared to an iPhone.
This is what an iPhone logic board looks\h\h
like. Suddenly, the MacBook Air\h
board seems large and clunky.
We can actually open up the iPhone logic board,\h
because it’s a sandwich. The lower part is the\h\h
RF or radio frequency board and the upper part\h
contains the SoC. In this case, it’s an Apple A15,\h\h
because this logic board was taken from an iPhone\h
14. If we turn it around, we can spot a single\h\h
NAND chip, that’s the iPhone storage.
But wait, aren’t we missing something?\h\h
We got the A15 and the SSD, but where’s\h
the DRAM? Where is the shared memory?
Every chip needs memory. But there are different\h
ways to integrate DRAM into a system. The most\h\h
common form are so called Dual Inline Memory\h
Modules, or DIMMs. Laptops usually use SO-DIMMs,\h\h
which stands for “Small Outline”.\h
They are literally just smaller.
But we can go further. The next step is called\h
“on-package memory”. That’s what the M3 we just\h\h
looked at is using. Or Intel’s Lunar Lake. Here\h
the memory modules are placed directly next to the\h\h
chip, onto the same package. The obvious advantage\h
is that you save a lot of space, because you\h\h
don’t need memory slots and large memory modules.\h
Plus, being so close to the chip has performance\h\h
and efficiency benefits. But it also means you\h
can’t upgrade or replace your RAM, if it fails.
But there’s one more step that’s closer\h
than on-package memory. It’s so close,\h\h
we can’t even see the memory. Or rather,\h
we are being tricked, because we can see\h\h
the memory, but not the SoC underneath it.
Going all the way back to 2016 with the A10,\h\h
Apple started using a very special packaging\h
technology for the A-series iPhone chips,\h\h
called InFO-PoP. It was developed by TSMC and is\h
short for “Integrated Fan-Out Package-on-Package”.\h\h
With InFO-PoP, the DRAM die is placed\h
directly on top of the SoC die. It’s\h\h
a memory sandwich. Package on package.
That’s why we can’t see the DRAM on our\h\h
iPhone logic board, because SoC and\h
DRAM are a single, closed package.
InFO-PoP is used when space is the most\h
limiting factor, like in a modern mobile\h\h
phone. Stacking DRAM on top of the SoC\h
means the logic board can be super compact,\h\h
even if there are obvious downsides. Like the\h
thermal implications of practically placing\h\h
a lid on the hottest part of your system.
That’s why the new MacBook Neo is so interesting.\h\h
It doesn’t use an M-series chip like all\h
previous Apple Silicon based MacBooks. Instead,\h\h
it’s using the A18 Pro, the same SoC that\h
powers the iPhone 16 Pro. Apple is using\h\h
an iPhone chip in a MacBook. It almost feels\h
a bit illegal. Can they do that? And how?
The first question is, how much does\h
the MacBook Neo actually share with\h\h
the iPhone 16 Pro? Is Apple really using a\h
tiny iPhone logic board inside a much larger\h\h
MacBook? The simple answer is no.
Every logic board is tailored to\h\h
a specific product. It needs to connect\h
and control dozens of different systems,\h\h
with the right control logic and connectors. The\h
SoC, the A18 Pro, is just a part of that. Take the\h\h
radio frequency board for example. An iPhone is\h
a mobile phone with a lot of connectivity aside\h\h
from just WiFi and Bluetooth. Something the\h
MacBook Neo doesn’t offer. The iPhone logic\h\h
board also needs to control a touch screen\h
and connect to multiple cameras for example.
The MacBook Neo is much simpler. One\h
screen without touch functionality,\h\h
a single 1080p webcam, keyboard plus trackpad\h
and speakers. That’s basically it. That’s all\h\h
the logic board needs to handle. Oh, and\h
the Touch-ID sensor on the 512GB model.
Thanks to teardowns, like from iFixit, we\h
can confirm this. The MacBook Neo doesn’t\h\h
use the same logic board as the iPhone\h
16 Pro. It’s a new, larger design. Not\h\h
as large as a M-series logic board, but much\h
larger than the logic boards in an iPhone.
The only thing that the iPhone 16 Pro\h
and the MacBook Neo actually share is\h\h
the A18 Pro SoC. You can see how tiny\h
it looks on the larger Neo logic board.
But isn’t an iPhone chip too weak to power a\h
MacBook and run MacOS? Aside from the obvious\h\h
fact that the MacBook Neo does so, let’s\h
take a look at what’s inside the A18 Pro.
Thanks to the guys from ChipWise we actually\h
have a great quality A18 Pro die shot. The\h\h
chip contains about 20 billion transistors on\h
a 109 square millimeter die and is produced\h\h
in TSMCs N3E. It shares process node and most\h
of its IP with the Apple M4 series. The best\h\h
way to think of it, is that M4 and A18 Pro\h
are the same chip generation, just build\h\h
with different target devices in mind.
Starting at the borders of the A18 Pro,\h\h
we can spot four 16-bit memory PHYs. They\h
connect to the LPDDR5x memory chip that’s\h\h
later stacked on top the A18 Pro using the\h
InFO-PoP packaging we just talked about. And\h\h
because four times sixteen is 64, the A18\h
Pro has a 64-bit wide memory interface.
According to Apple, the A18 Pro has a 6-core CPU,\h
with two performance and four efficiency cores.\h\h
While the names are clear, performance core\h
sound fast and efficiency core sound, well,\h\h
efficient, it can be difficult to grasp the real\h
difference. But looking at a die shot helps.
The two performance cores are located in\h
the upper left part of the chip. We can\h\h
clearly see both cores in a mirrored layout\h
side-by-side. Right below is the L2 cache with\h\h
two 8-megabyte SRAM blocks for a total of 16\h
megabytes. This cache is shared between both\h\h
performance cores. And below the cache we can\h
see a dedicated AMX unit, a special co-processor\h\h
that’s optimized to run matrix operations,\h
which are mostly used for machine learning.
As you can see, the performance cores take\h
up quiet a bit of area on the A18 Pro,\h\h
especially if you consider the large\h
L2 cache and the attached AMX unit.
In comparison, the efficiency cores are almost\h
comically small. We can spot the entire cluster\h\h
of four E-cores pretty much in the center\h
of the chip. The E-cores share a 4-megabyte\h\h
large L2 cache and also come with a dedicated\h
AMX unit. As you can see, a single P-core is\h\h
almost as large as all four E-cores combined.
Make no mistakes, both core types are great\h\h
performing cores for their size. But looking\h
at the silicon you can see that there’s a big\h\h
size difference, which also translates into a big\h
performance difference. P-cores are built to\h\h
provide maximum single thread performance and\h
run at high clock speeds. A great combination\h\h
for demanding applications. The E-cores on the\h
other hand really focus on efficiency. But not\h\h
just energy efficiency, also area efficiency.
Next is the Apple Neural Processing Unit,\h\h
or NPU. It does take up a large amount of\h
area, similar to the entire E-core cluster.\h\h
Apple calls it a 16-core NPU, but there\h
are only eight physical cores visible.\h\h
Which means they are most likely dual-cores.
Below the NPU and the E-cores we can find the\h\h
System Level Cache. This is basically Apple’s\h
version of a chip-wide shared L3 cache. It can\h\h
be used by almost all other system components.\h
On the A18 Pro, the SLC offer 24 megabytes,\h\h
divided into two 12-megabyte SRAM blocks.
Last, but certainly not least, we have\h\h
largest logic block of the entire chip.\h
The 6-core GPU is located on the lower\h\h
area of the die. Graphics are a very important\h
part of SoC performance. Not just for gaming,\h\h
but a modern GPU can accelerate many different\h
tasks, for example rendering or even AI.
If you think about it, there’s a NPU, both P-\h
and E-cores come with AMX units and the GPU can\h\h
also run AI. There’s a lot of AI performance\h
spread all over the A18 Pro. And with the A19\h\h
Pro and M5 series, Apple added so called “neural\h
accelerators” to the GPU. It seems like everything\h\h
is converging towards AI acceleration.
Now that we have labeled all major function\h\h
blocks, you can see that there’s still a lot of\h
unmarked area left. There’s some IO on the upper\h\h
shorelines of the chip and I’m pretty sure that\h
the top right part, next to the memory PHY, is the\h\h
Apple secure enclave. The other areas are used for\h
the media and display engines, image and signal\h\h
processing, and much more. But it’s really hard to\h
accurately label them without starting to guess.
One more thing I want to mention is an interesting\h
structure to the right of the P-cores. It’s very\h\h
clearly a dual-core of some sorts, I call\h
it “low-power E-cores”. It most likely is\h\h
a specialized embedded microcontroller, but\h
it does look very similar to a tiny CPU.
Now that we know what the A18 Pro actually\h
looks like, let’s compare it to the M4. Sadly,\h\h
there isn’t a similar high-res die shot\h
available like the one from ChipWise, but\h\h
we do have a low-res version from TechInsights.
Like we talked about before, the M4 uses the same\h\h
IP and is build on the same N3E node as\h
the A18 Pro. With 28 billion transistors\h\h
it has 40% more transistors but at 165 square\h
millimeters it’s about 50% larger in size.\h\h
That’s most likely because it contains more\h
I/O and analogue circuitry, something that\h\h
usually has a lower transistor density.
Looking at the die, we can see that the M4\h\h
has eight memory PHYs, instead of just four\h
on A18 Pro. Each one still offers a 16-bit\h\h
wide connection, for a total of 128-bit. The\h
GPU is also much larger at up to 10-cores.
Keep in mind that while the iPhone 16 Pro\h
uses the full six GPU cores of the A18 Pro,\h\h
the MacBook Neo comes in a binned version\h
with only 5-cores active. A fully enabled\h\h
M4 has double the GPU cores of the MacBook\h
Neo. A pretty large difference. And that’s\h\h
why it needs double the memory interface, because\h
memory bandwidth is most important for the GPU.
Next, let’s compare the CPU cores. The M4 has\h
four P-cores, again double what the A18 Pro\h\h
offers. But interestingly, both chips have\h
the same amount of shared P-core L2 cache.\h\h
It’s 16 megabytes on the A18 Pro and the M4.
And the same is true for the E-cores. While\h\h
the M4 has 50% more cores, six instead of\h
four, the shared E-core L2 cache stays at\h\h
4 megabytes. The NPU is the same too, with\h
eight dual-cores, just like on the A18 Pro.
But where it really gets interesting is the System\h
Level Cache. The M4 only has 8 megabytes of SLC,\h\h
that’s a third of what the A18 Pro\h
offers. A massive difference. But why?
First, with a 128-bit interface, the M4 already\h
offers a pretty decent memory bandwidth, there’s\h\h
less need for a large cache, which is mostly\h
used to mask memory access times. And second,\h\h
the A18 Pro has a much stronger focus on energy\h
efficiency, because it was designed as a mobile\h\h
phone SoC. Every time the SoC has to access\h
the DRAM, even if it’s packaged right on top,\h\h
it uses a lot of energy. Memory access is\h
expensive. And with a large SLC you are not only\h\h
improving performance, but also power efficiency.
That means, if you would run an application that\h\h
only uses two P-cores but needs a lot of cache,\h
the A18 Pro might actually be faster in certain\h\h
tasks, because it offers a combined 40 megabytes\h
of shared P-core L2 cache and System Level\h\h
Cache. Kind of cool if you think about it.
But make no mistake, with double the P-cores\h\h
running at higher clock speeds and 50% more\h
E-cores, the M4 handily outperforms the A18\h\h
Pro in most scenarios. Just not in all of them.
So, what does all of that mean for the MacBook\h\h
Neo. First, it shows that the A18 Pro is more\h
than capable of powering a MacBook and running\h\h
MacOS. Yes, it’s an iPhone chip, but in the end,\h
it’s still Apple Silicon. It’s basically a mini\h\h
M4. And when it comes to on-die cache,\h
the A18 Pro even has the advantage.
Building an entry level low-cost MacBook isn’t\h
easy. Apple needs to cut costs, without losing\h\h
the essence of what a MacBook is supposed to be.\h
And the A18 Pro is clearly part of that strategy.
With a 50% larger die, a single M4 wafer produces\h
only about 300 known good dies, while that same\h\h
wafer could produce 500 A18 Pros. That’s a big\h
difference. Remember, both use the same TSMC N3E\h\h
process node, which means a wafer costs the same.
As far as I know, the A18 Pro was only used in the\h\h
iPhone 16 Pro and Pro Max. And Apple also\h
only used A18 Pro’s with all six GPU-cores\h\h
enabled. Which means, chips that have\h
a production defect couldn’t be used.
As we saw, the GPU is the largest function block\h
on the A18 Pro. Which means if there’s a defect,\h\h
there a pretty good chance it’s affecting the GPU.\h
Using a binned A18 Pro is genius, because it means\h\h
Apple can now finally use all the A18 Pro’s with\h
a partial defect in the GPU area. I’m sure Apple\h\h
has been collecting these binned chips ever since\h
they started producing the A18 Pro back in 2024.
And using an iPhone SoC that is designed for\h
even lower power than the already efficient\h\h
M-series chips, Apple can get away with using\h
a cheaper and lower capacity battery while\h\h
still getting good battery life.
The A18 Pro also allows Apple to\h\h
design a smaller logic board for the\h
MacBook Neo. Not a massive difference,\h\h
but when every dollar counts, it does help.
And lastly, the A18 Pro already comes packaged\h\h
with 8 gigabytes of DRAM. Yes, in theory\h
Apple could tell TSMC to start packaging\h\h
the A18 Pro with more memory, but\h
that would not only be expensive,\h\h
but also mean they couldn’t use older A18 Pros\h
that were already produced for the iPhone 16 Pro.
The M3 MacBook Air still came in a base\h
configuration of 8 gigabytes DRAM and\h\h
256 gigabytes SSD. That was not that long ago.
In my opinion, the MacBook Neo is the smartest\h\h
way Apple could have built a low-cost MacBook. And\h
the A18 Pro is the central piece of the puzzle.
I know I’ve been missing in action a bit in\h
the last couple of months, but this is about\h\h
to change. Many projects and videos are in the\h
works, here’s a sneak peek of an amazing M1 die\h\h
shot from Fritzchens Fritz. And you can also find\h
me over on the SemiAnalysis YouTube channel.
I hope you enjoyed the video\h
and see you in the next one!
