DRAFT

# My all-star zoo, or why I hired Linus Torvalds, Kent Beck, and Rob Pike for my AI team

<x-summary>
This summer, I transformed Claude Code from a junior intern into an okay engineering team.
</x-summary>

Like everyone else, I've started with Cursor (and Windsurf, and Zed). And... it didn't work well at all.

Our production codebase is way over 200k LOC of Go code alone, plus a lot of HTML and some JS, and it was just too complex for AI to handle well. Decidedly _not_ a boilerplate CRUD app.

With that codebase, AI was only good for doing the very few boilerplate tasks like adding enums.

Over the summer, I got from _there_ to a pretty darn successful AI team, and here's my journey.


## Step 1: Claude Code

Switching to the stupid command-line `claude` tool (and porting over all rules into `CLAUDE.md`) has given me an immediate boost.

I don't know what it is that Claude Code does differently from Cursor, but it _is_ very different. It uses a _lot_ more tokens and produces _much_ better results even when using exactly the same model.

But, of course, you do not have to use exactly the same model. One of the superpowers of Claude Code is running Opus via your Claude Max subscription. Opus is much smarter than Sonnet. Claude Max costs $200/mo but gives you $250 worth of tokens every 5 hours, and that's just about enough to run your AI coding agent on Opus for hours.


## Step 2: Instructions

Claude needs a lot of instructions to work effectively, otherwise it spends way too much time re-discovering your codebase.

I got Claude writing its own instructions:

1. Before working on a new area of the codebase, I'd ask Claude to explore that area and write a summary to `_ai/` folder.

2. I then curate and rewrite that summary to emphasize the important parts.

3. When doing a task, I ask Claude to read the docs under `_ai/`. Profit!

Note: Claude really sucks at writing docs. Always did, still does. It fails to tell essentials from superflous, and tends to fixate on the weirdest things.

That's why all docs that Claude is allowed to update go under `_ai/`, while manually-written docs go into `_readme/` or `CLAUDE.md`. Letting it touch the real, human-written documentation was a disaster.


## Step 3: Planning

Following Anthropic's advice, I've added an explicit planning step that produced a plan file. The subsequent execution steps were supposed to read and update the file.

Claude supports custom commands, and here's my `/do` command from that period:

```
Your task:
$ARGUMENTS

1. Ultrathink.
2. Read any relevants docs under `_ai` and `_readme`.
3. Read a lot of code. Find related code and read it. Anticipate challenges and proactively research them too.
4. Formulate acceptance criteria for the original task.
6. Where should you put the test(s)? Work hard to find a bunch of relevant tests, read them all, find the best place.
7. Build a VERY detailed plan. Anticipate any challenges. Save it to `aiplan.txt`.
8. Execute. Iterate until all tests pass! DO NOT BE LAZY. DO NOT TAKE SHORTCUTS. Work hard. Document progress in `aiplan.txt` as you go.
9. Do not stop until all tests pass! Document progress in `aiplan.txt` as you go!
10. Review your own code, criticue it, find shortcomings. Compare to the instructions in `CLAUDE.md`. Review high-level design. Verify your work against the acceptance criteria. ALL tests must pass. Iterate to address any issues.
11. ALL TESTS MUST PASS IN THE END. Anything broken? You caused that. FIX IT.
12. DID I MENTION TO DOUBLE CHECK THAT **ALL** TESTS PASS?
13. Write down any valuable facts that are worth remembering (high bar) in one-line fact format (like `CLAUDE.md`) to the relevant file under `_ai`.
```

You can probably feel my pain; every line here is written in blood and tears.


## Intermission: What the fuck are we even trying to achieve?

There are two ways to use AI coding tools:

1. Interactive. You collaborate with AI, ask questions, refine a plan, then execute together.
2. Batch. You give it a task and switch to something else, and come back in 30-60 minutes expecting significant progress on that task.

The interactive way is <s>for suckers</s> err, for less experienced developers.

If you're a senior and you try to do that, you'd never gain any efficiency.

Batch mode is what we're after. That's the pipe dream of AI coding: explaining a task, and coming back to that task done.


## Step 4: Subagents

Claude Code has introduced subagents mid-summer.

It might not be obvious why subagents are a good idea. Wasn't obvious to me.

So let me tell you how Claude Code without subagents sucks:

1. It forgets things. By the time it gets the tests to compile, it has stopped paying attention to half of your original request and to 80% of the instructions from `CLAUDE.md`.

2. It runs out of context and then forgets things. Debugging stupid compilation errors and test run issues takes a lot of tokens. In about 15–30 minutes, it runs out of context and performs compaction, which gives equal weight to your original request and to the minutiae of how it was failing to make the tests pass.

3. It gets confused as to what it's supposed to be doing. Will update tests to match a broken implementation, or will break a perfectly working implementation to get a broken test to pass.

4. It is hard to steer. Getting it to remember to run all the tests and make sure they all passes was next to impossible.

See that `/do` command. And a few more that I needed — here's `/wat`:

```
Is that a fucking joke? WTF!!!!!!
$ARGUMENTS

Ultrathink a plan first. Criticue your work. Then execute. ITERATE UNTIL DONE AND ALL TESTS PASS.
```

and here's `/dumbass`:

```
1. Ultrathink, review and critique your work. Can it be improved? Consider CLAUDE.md guidelines, code style, readability, best practices, larger system-wide design concerns. Focus on YOUR work, not the existing code.
2. Address these concerns.
3. Make sure ALL tests still pass.
```

and I had a dedicated `/fix-failing-test` command too.

If this doesn't scream “smooth sailing” to you...

Claude subagents are invoked with a separate context window and separate instructions. This is huge because:

1. It solves the context window problem. All those painful attempts at fixing compilation errors stay in subagent's context and are discarded once it's done.

2. It solves the steering problem. Each agent re-reads `CLAUDE.md` and has separate instructions, so you can stuff fresh propaganda straight into its brain. A short lifetime of the agents means they never forget your teachings.

3. It solves the goal drift problem if you do a smart separation of agents. The testing agent will be writing tests and won't try to break your code, while the implementation agent will be trying to get tests to pass without dumbing them down.

(If you never experienced this, you might not appreciate how big of a problem it is when instead of fixing the implementation, the agent just deletes the related test and calls it a day.)

So I created:

* a planning agent
* a test engineer
* an implementation engineer
* a code reviewer
* and a few niche ones

My approach goes against some widespread internet wisdom saying that subagents are for read-only tasks, and coding should be at the top level. That wouldn't work for me at all. Pollution of the context window with low-level coding issues was one of my main problems, and the other problem was separation between test-writing and implementation-writing activities. I needed coding to be at subagent level only.

I rewrote `/do` as:

```
Your new task:
$ARGUMENTS

Process:
- Tech lead agent first.
- Test engineer agent is always next. Even if no test changes are planned, let the test engineer confirm that.
- Loop software engineer, test engineer and problem solver agents to execute on each step until done. Manager must step in between every step to reflect the progress and ensure alignment.
- Doc writer agent if relevant (if any changes to public API made, call doc writer).
- Review via code reviewer agent.
- If there are ANY suggestions from code reviewer, go back to test engineer and software engineer to address code reviewer's suggestions. After that, call manager, and then repeat code review. LOOP UNTIL REVIEW FINDS NO SUGGESTIONS.
- Finally, librarian agent to store the accumulated knowledge in _ai.
```

Things got much better. Some problems went away immediately: no more compactions, no more mangled tests, no more tests left broken.

From this point on, I never needed `/wat`, `/dumbass`, `/fix-failing-test` and other similar commands for low-level steering; all of that was solved by the virtual team.


## Step 5: Linus Torvalds

The problems I was encountering got upgraded to a higher level, but still, code reviews were missing a lot of issues, and the team would often do something very stupid.

I got very frustrated with another self-congratulatory review one day, and asked to add another reviewer agent:

```
Linus Torvalds doing very high-level review of the changes in his signature ruthless and pragmatic style. Must run after the normal code reviewer. Focuses on high-level details only, not on the code minutea:

- Did we do the right thing? Or did we do something stupid?
- Did we cut corners to finish faster and called it a success, instead of doing the right thing?
- Have we implemented everything that was requested, or forgot something?
- Do the changes align with user intent?
- Did we add a hack or a workaround where we could do the right thing?
- Is implementation at the right level of abstraction?
- Is implementation and tests in the right packages?
- Do our tests actually test what they claim they do?
- Do our tests verify the core functionalities in integration? Or do they only check superficial things?
- Did we overfocus on edge cases and forgot to test the core functionality?
- Are tests maintenable? Could tests be simplified?
- Is core maintainable? Could implementation be simplified?
- Did we overlook some part of the system that also needed to be updated?
- Do we follow the highest standards of quality? Have a fresh look at the code through the eyes of Linus. Would he approve it?
```

And, oh my god, I didn't expect it to have that much impact. Well, I did say that this is a high-level reviewer. I did say it was Linus, but the quality of those reviews, they were something else entirely. It would find things I myself would miss. It would think of production issues, future issues. Like, I would often deploy something, say we're working on a feature and it needs an internal admin area. And that internal admin area usually starts out very bare bones. And sometimes my ideas of bare bones are actually naive. So I would just do the simple thing that would work. But Linus would criticize that and say, hey, you really need, there will be performance issues as soon as there are many of these objects. And he would call these agent names. All the frustrations I felt with the coding folks, he would express them for me. That was amazing. He would say, "Whoever wrote this shit needs to be fired immediately."

So that was very therapeutic, if nothing else. But I did create a whole other layer of quality feedback. It would flag all the things I would flag myself, like code quality issues. It was amazing. It was like an order of magnitude better than what I had before.

Of course, after that, what the team struggled with is actually implementing the feedback that Linus gives. Because the workflow was not strong enough. Like I would go code review and then back to implementation engineers. I figured no. That doesn't work. I need to have a planning phase after every review.

## Step 7: Building an All-Star Team

So I've actually put the planner. And then, the planner wasn't doing as great of a job as I wanted it to do, so I figured I needed to pay more attention to detail. And then a thought struck me. So like if Linus Torvalds is so much better, if having Linus as a code reviewer is so much better than just having a generic code reviewer, would it be better to hire, to use legendary personalities for each of the agents?

So I figured, okay, who's the legendary project manager with attention to detail. Well, it's Joel Spolsky. Who's a legendary testing engineer? Well, clearly Kent Beck. Who is a legendary implementation engineer? Well, here it wasn't quite clear. Claude proposed a bunch of ideas, including John Carmack. But hey, John Carmack doesn't represent my values. And a lot of people who were great developers, they don't necessarily reflect my values, but there were some people who did reflect those values. These are the team that created the Go language because the Go language is like 100% reflection of my pragmatic outlook on development. So I figured, okay, Rob Pike. Rob is going to be my engineer.

I actually added a few more subagents before this moment. I added an HR subagent whose job would be basically updating other agents' instructions so that I could say, hey, HR, do this and this. So I would tell HR, hey, turn this generic agent into Kent Beck. And it would do it.

Anyways, so I got a team. I got Joel. I got Kent. I got Rob. I got Linus for the reviewer. The other reviewer, I don't remember who I went with initially, but I wasn't very satisfied with that. So I looked for a person who's famous for pragmatic code quality that I like, that again represents my values. I went with Kevlin Henney.

I added a documentation writer. I used Raymond Chen, the Old New Thing blog author, as the documentation engineer.

And I updated the workflow. I described the workflow in Claude. I updated the workflow to insert Joel after every iteration so that it would be Linus and Kevlin doing reviews and Joel figuring out what to do about those reviews, right? What, like, are we done? Should we address those issues?

That worked better. Like that team could ship stuff. Like, each step here represents an order of magnitude.

## Step 8: Iterating on the Plan

So I had a team which was Joel doing planning, then implementation, then code review. Now, very often the code review would find that the plan was not quite great. So what I did is I started calling Linus twice. I first had Joel do the plan. Then I had him call Linus to review the plan. Then back to Joel to address the issues that Linus found. Then back to Linus to review his updated plan, and only when Linus approves, we move on to the implementation. These got things much better, like surprisingly, iterating on the plan was a great idea. Like, it had bigger effect than one would expect from something like that.

## Step 9: Don Melton and the Three-Phase Workflow

I got Joel as a manager. But problem is, it turned out that Joel, for all his attention to detail, he was very much focused on shipping, just like in real life. But problem is in real life, that makes sense. When the AI manager says, hey, like, we have this great feedback, but let's do it in the future because we need to ship, ship, ship now. That's not what we want. We can spend a couple more hours checking along and implementing the feedback so that the quality is better overall because I wouldn't ship that anyway.

So I had to fire Joel. I replaced him with Don Melton. Don Melton is someone who I really admire and respect and I listen to him on podcasts for many hours and I just got him to manage the team. I quickly found, though, that he's great at insisting on quality, but he is not doing detailed technical planning. So I rehired Joel. I put him next in the workflow after Don Melton. So Don would do a high-level plan and immediately Joel would expand it into a detailed technical plan, technical spec, which is something that Joel is great at, is known to be great at. And then the implementation agents would do their work. Then Kevlin and Linus would review. And then again, Don would make the next high-level plan.

So I split the workflow into three phases. Phase one was planning. Or I call them steps. Step one was planning. And the planning step is Don, then Joel, then Linus, then back to Don, Joel, Linus, and iterate until Linus approves. Then we go into the implementation step. And the implementation step is we call Kent. Then we call Rob. Then there is actual documentation writer, Raymond. And then we do review, we do in parallel, we call Kevlin and we call Linus. And after the review is done, we go back to the planning phase. So now it's Don, then Joel, then Linus. They iterate again on the new plan. And I said, hey, all of this only finishes like we declare a task done if all three of them agree. Like if Don, and Joel, and Linus during the planning phase, they all agree that we're all done. Nothing left to do. Then we're done.

When I became an all-star team and it was Linus calling everyone incompetent, I started calling it a circus or a zoo because it was really funny to watch all of this happening. Really funny.

## The Supporting Cast

### The Librarian

There is step three, which is the finalization step. I actually had a few more things in this team that I didn't talk about. One was that when all this ends, I had a librarian agent that I then switched to a named agent, Ward Cunningham, the creator of WikiWikiWeb. There was a librarian agent that when everything, like when everyone approves and the task is done, it goes and updates those docs under _AI to remember stuff for the future. Honestly, it's a hit and miss. And all of these updates, I treat them as proposals. So like it updates a batch of files. I review what it added. And sometimes I approve, sometimes I approve only some of those changes, sometimes I discard all those changes.

### The HR Agent

And then there is my HR agent. So I said, hey, run HR agent at the end of each process. And if there were any... The way I formulated this is that if there were any revision requests from the user, try to update agent instructions to incorporate those, so that the next time user doesn't have to give those requests. So that the agent basically perform better the next time so that it wouldn't be necessary.

And when I moved to All-Star team, this became Andy Grove, the head of Intel. So Andy was in charge of updating my agent instructions. And again, these updates I treated as proposals. It would very often make changes, much more often than I would like. And often those changes would not be desirable because they would be, again, weirdly over-focusing on some small aspect that it would dedicate a lot of instruction tokens to make better. So I thought, hey, not worth it.

Ideally, I want to improve this a lot. I really want to iterate on agent definition. So like compact them, make them better, et cetera. I don't feel like I'm anywhere near though. It sometimes proposes good changes. Sometimes I accept and commit them. Very often I just rework everything that he did. So this is not a very successful experiment.

### The Problem Solver

Then the third agent I didn't talk about is problem solver. So sometimes I figured, sometimes the implementation agent would get stuck and wouldn't be able to make progress, especially on a harder task. So I created what I call the problem solver. I eventually proposed that I call it Donald Knuth. So Donald was my, in the all-star team, it was Donald, the problem solver, and his role would be when the agent gets stuck, he would hand off to the problem solver. And problem solver would figure out like why we got stuck. So instead of trying to do coding, he would do a thorough, thorough investigation into what's wrong, what we're doing, and he would say, "Hey, here's how to move forward."

Honestly, this was more of an idea than a practical thing because I only saw it engaged once. Maybe I missed it working like one or two times more, but it's very seldom thing that agents these days get stuck. But if they do, if they ever do, they do know whom to call. And this worked wonderfully. In those cases where I think it was a testing agent, he just couldn't figure out how to write a test for a particular area. And normally it would just veer more and more off course and give up, right? But here it figured out, hey, like I'm stuck. And I need to call the... Well, it wasn't Donald back then. It was just a generic problem solver. And it called problem solver. And problem solver did solve the problem. Like, it figured out why it couldn't progress and had enough. And it succeeded then. So, yeah. Huge success, but again, I only have one recorded case of it being useful, so maybe not that important.

## My Workflow Commands

I have two main commands. One command is called "do" and it's for starting a task. Another command is called "rev" and it's for requesting revisions on the current task. Those are really, really two most important commands. I have HR command so that I don't have to, and I have a few more weird, like I have HR instill command that tries to instill more personality into an agent that was created after some person. So it would ask an agent, hey, doesn't your current instructions reflect your personality well? Like, is there something you want to improve? Like, feel free to rewrite your own instruction file. And the agent would go and talk about himself, right? So like, add a lot more of his own style and personality into the instructions, which I feel like helps.

## The Moral of the Story

So, is there a moral to this story? Well, the moral is that I have these definitions published. Of course, these definitions are kind of specific to my projects and my values, but they're not that hard to create separately. Claude has a great subagent. Well, I think the command is called agent. So it has a great command to create new agents that you just describe the agent and it does a great job expanding it into detailed instructions. And if you create an HR agent, it would also be quite good at expanding instructions. So it's not that hard to create a team like this on your own.

And there is really not that much. I mean, there is a bunch of experimentation to do, of course. But the effect of this is amazing. Throughout the summer, without changing the models, this went, I think, at least two orders of magnitude in the complexity of tasks and the quality of the output that it could do for me. In fact, it went from feeling like a really smart junior developer to feeling like not a very solid senior developer basically. You know, the kind of senior developer that's currently on the market where you spend five years doing something and you now call yourself a senior developer. So that kind of senior developer. But still, it was a huge step up.

Each step here represents an order of magnitude. Initially, right, on the lower steps, I had to deal with just doing stupid stuff like disabling a test or not even forgetting what the task was. And at the highest level, I was dealing with things like the tests are written well, but just not exactly in the style I want, so how do I steer it to a better style? But, like, instead of doing junior level work it was doing like senior level work, just not the kind of senior engineer I'd like to have, but like a different kind of senior engineer with different values. Still a problem, but a much higher level problem, much better problem to have.

## Part 4: Resources

- [GitHub repository with all agent definitions](#) (link coming soon)
- [Complete workflow documentation](#) (link coming soon)
- [Example Claude Code sessions with the circus in action](#) (link coming soon)

### The Complete All-Star Roster

- **Don Melton** - High-level technical lead
- **Joel Spolsky** - Detailed implementation planner
- **Kent Beck** - Test engineer
- **Rob Pike** - Implementation engineer
- **Donald Knuth** - Advanced problem solver
- **Kevlin Henney** - Code quality reviewer
- **Linus Torvalds** - High-level architecture reviewer
- **Raymond Chen** - Documentation writer
- **Ward Cunningham** - Knowledge librarian
- **Andy Grove** - HR and agent instruction optimizer
