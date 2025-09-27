DRAFT

# Senior devs too good for AI? Saving energy not time

<x-summary>
    If you're like me... the more experience you have, the harder it is for you to gain from coding with AI. And yet, I think I made it work.
</x-summary>

If you're like me... the more experience you have, the harder it is for you to gain from coding with AI:

1. You're doing harder tasks in the first place.

2. You're better at execution, so your baseline is much higher.

3. Your idea of what code is good is much more developed and specific, so the bar is higher, and you'll accept a much narrower range of solutions.

I am a technical director with 25 years of intense experience, so I've been suffering this conondrum ever since AI coding got traction.

Many professionals are gonna try AI agents and say “wow this is really stupid, I'd rather do it all myself”. Honestly — I was tempted myself.

And yet. I've been too curious. And I think I made it work. If you approach this thing like any other challenge in your professional life, there are lots of gains to be had.

I got three insights here.


## 1. AI saves energy, not time

When executing at the 100% of my abilities, I will outdo AI in both pace and quality.

However, I'm a human being and can't always be at 100%.

And when I am tired, sleepy, or just not feeling quite right, AI is an amazing addition to my workflow:

- It helps me start a task and get it to maybe 60-80% completion.

- It takes over those mundane straightforward low-value things that drain willpower.

- It is happy to go an extra mile and add a full suite of tests even when I am in energy saving lazy mode and would cut a corner.

So I found these good uses for AI:

- when I'm dreading to start a certain task
- when I find a task too mundane
- when my day is gonna be super hectic and a few 30-minute sessions of AI steering is all I can do
- when I don't feel quite like myself and reviewing code feels much more approachable than writing it
- when I'm busy with a big primary task but need to get a few smaller things done in parallel

Very often I'll get going and will find it comfortable to take over for that “last mile” of changes that are easier done myself than explained to AI.


## 2. AI is more capable when used smartly

I started out with a pretty naive usage with just a bunch of instruction files, but when I scaled that to a family of subagents with a defined workflow, with the right prompts, I found AI to be about two orders of magnitude more capable.

The gist of my setup is:

1. You want subagents with a clear separation of duties. At the very least, you need to split up: planning from execution; test writing from implementation; code review from everything else.

2. You want a specific workflow that iterates over planning, execution and reviews steps.

3. You want to train each subagent to adhere to the exact style and process that you like. This includes general instructions (like CLAUDE.md), agent-specific instructions, and agent personalities (naming an agent after a famous person to gain the well-known traits of that person “for free”).


## 3. Parallel execution

Just having an agent helping you out in your main worktree is very limiting (in a sense that it's hard to make it work efficiently enough to be useful).

What you want instead is 1 to 3 separate checkouts of your project, so that you can have multiple agentic tasks running in parallel. Just clone your project a few extra times, and run `claude` (or your agent of choice) in each folder, and keep each instance open in your favorite text editor/IDE.

For me, three parallel agents is about the most I can handle, because my ability to review their results and give revision requests becomes a bottleneck. An agent takes 20–60 minutes to complete a task, and it takes me about 15–30 minutes to review and give feedback, so three agents saturate my entire bandwidth.

(I know some people use Git worktrees for this, and it sounds really smart, but I just use plain old clones.)

My review workflow is quite simple: I read the changes in Sublime Merge (my Git client of choice), and go more in-depth in Zed (my editor of choice) when needed. Then I write a very detailed revision request for Claude.


## My workflow

These days I just send every task to AI, so the range varies from simple copy changes to complex new features. A good number of those are still too hard for AI to handle, but:

1. I expand each ticket with an “Impl:” section where I describe roughly how I would go about implementing it, so AI will at least partially do what I was gonna do anyway.

2. I ask AI to handle the ticket initially, and allow it a few iterations with my review and revision requests. Doing 1–3 iterations strikes the right balance for me, I rarely go beyond that.

3. Finally, I take over, finish and clean up the implementation.

That final step takes a lot of time. So, yes, I'm not saving much, or anything at all, compared to an abstract robotic version of myself who can always operate at 100%. In real life though? Feels totally worth it.


## My toolset

I believe these principles apply universally, but I'm personally using Claude Code with a manual selection of Opus model and Claude Max subscription.

(If you're wondering, Claude Max currently gives you $250 worth of tokens per 5-hour interval, for $200 per month. That's generally enough to run a single agent continuously, or three parallel agents for about 2–3 hours.)
