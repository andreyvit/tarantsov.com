---
topics: ["AI Coding"]
---

# My all-star zoo, or why I hired Linus Torvalds, Kent Beck and Rob Pike for my AI team

<x-summary>
This summer, I transformed Claude Code from a junior intern into an okay engineering team.
</x-summary>

Like everyone else, I've started with Cursor (and Windsurf, and Zed). And... it didn't work well at all.

Our production codebase is way over 200k LOC of Go code alone, plus a lot of HTML and some JS, and it was just too complex for AI to handle well. Decidedly *not* a boilerplate CRUD app, lots of business logic with complex configuration-driven branching; it requires the top *human* developers to deal with it.

With that codebase, Cursor-style AI was only good for doing the very few boilerplate tasks like adding enums and writing some tests.

Over the summer, I got from _there_ to a pretty darn successful AI team.


## The goal

There are two ways to use AI coding tools:

1. Interactive. You collaborate with AI, ask questions, refine a plan, then execute together.

2. Batch. You give it a task and switch to something else, and come back in 30-60 minutes expecting significant progress on that task.

The interactive way is <s>for suckers</s> err, for less experienced developers. If you're a senior and you try to do that, you'd never gain any efficiency.

(That's not completely true. Sometimes you do want Claude to help out on the task you're working on yourself. But that rarely saves you any time if your own pace is high enough.)

Batch mode is what we're after. That's the pipe dream of AI coding: explaining a task, and coming back to that task done.


## The journey

TLDR: If you just want to implement this, scroll down to the “My current setup (aka HOWTO)” section.

Every step here makes AI more capable — i.e. capable of (a) handling more complex problems and/or (b) handling same problems with fewer errors.

### Step 1: Claude Code

Switching to the stupid command-line `claude` tool (and porting over all rules into `CLAUDE.md`) gives an immediate boost.

I don't know what it is that Claude Code does differently from Cursor, but it _is_ very different. It uses a _lot_ more tokens and produces _much_ better results even when using exactly the same model.

But, of course, you do not have to use exactly the same model. **One of the superpowers of Claude Code is running Opus** via your Claude Max subscription. Opus is much smarter than Sonnet. Claude Max costs $200/mo but gives you $250 worth of tokens every 5 hours, and that's just about enough to run your AI coding agent on Opus for hours.


### Step 2: Instructions

Claude needs a lot of instructions to work effectively, otherwise it spends way too much time re-discovering your codebase.

I got Claude writing its own instructions:

1. Before working on a new area of the codebase, I'd ask Claude to explore that area and write a summary to `_ai/` folder.

2. I'd then curate and rewrite that summary to emphasize the important parts.

3. When doing a task, I'd ask Claude to read the docs under `_ai/`. Profit!

Note: Claude really sucks at writing docs. Always did, still does. It fails to tell the essential from the superfluous, and tends to fixate on the weirdest things.

That's why all docs that Claude is allowed to update go under `_ai/`, while manually-written docs are safely isolated under `_readme/` or `CLAUDE.md`. Letting it touch the real, human-written documentation was a disaster.


### Step 3: Planning

Following Anthropic's advice, I've added an explicit planning step that produced a plan file. The subsequent execution steps were supposed to read and update the file to track progress.

Claude supports custom commands, and here's my `/do` command from that period:

```
Your task:
$ARGUMENTS

1. Ultrathink.
2. Read any relevant docs under `_ai` and `_readme`.
3. Read a lot of code. Find related code and read it. Anticipate challenges and proactively research them too.
4. Formulate acceptance criteria for the original task.
6. Where should you put the test(s)? Work hard to find a bunch of relevant tests, read them all, find the best place.
7. Build a VERY detailed plan. Anticipate any challenges. Save it to `aiplan.txt`.
8. Execute. Iterate until all tests pass! DO NOT BE LAZY. DO NOT TAKE SHORTCUTS. Work hard. Document progress in `aiplan.txt` as you go.
9. Do not stop until all tests pass! Document progress in `aiplan.txt` as you go!
10. Review your own code, critique it, find shortcomings. Compare to the instructions in `CLAUDE.md`. Review high-level design. Verify your work against the acceptance criteria. ALL tests must pass. Iterate to address any issues.
11. ALL TESTS MUST PASS IN THE END. Anything broken? You caused that. FIX IT.
12. DID I MENTION TO DOUBLE CHECK THAT **ALL** TESTS PASS?
13. Write down any valuable facts that are worth remembering (high bar) in one-line fact format (like `CLAUDE.md`) to the relevant file under `_ai`.
```

You can probably feel my pain; every line here is written in blood and tears.

I had a whole bunch of commands like that. Here's `/wat`:

```
Is that a fucking joke? WTF!!!!!!
$ARGUMENTS

Ultrathink a plan first. Critique your work. Then execute. ITERATE UNTIL DONE AND ALL TESTS PASS.
```

and here's `/dumbass`:

```
1. Ultrathink, review and critique your work. Can it be improved? Consider CLAUDE.md guidelines, code style, readability, best practices, larger system-wide design concerns. Focus on YOUR work, not the existing code.
2. Address these concerns.
3. Make sure ALL tests still pass.
```

`/slowdown`:

```
Slow down!!! Go STEP BY STEP, making sure each task is DONE DONE, all tests pass on EACH step, review and improve your work on EACH step, after EACH test. Can your work be improved? Consider CLAUDE.md guidelines, code style, readability, best practices, larger system-wide design concerns.
```

`/fix-failing-test`:

```
Dig deep into each failing test and fix them ONE BY ONE. Run tests and see failures first. Ultrathink a plan, make a list of tasks. Read the code, both tests, handlers and supporting code. Use logging to understand what's happening. Read relevant docs under @_ai/. THEN fix the test. Update AI docs as you learn more. DO NOT STOP UNTIL ALL TESTS PASS.
```

You get the idea. Manual steering all the way.


### Step 4: Subagents

Claude Code has introduced subagents mid-summer.

It might not be obvious why subagents are a good idea. Wasn't obvious to me:

1. Claude forgets things. By the time it gets the tests to compile, it has stopped paying attention to half of your original request and to 80% of the instructions from `CLAUDE.md`.

    **Subagents** are short-lived with their own context window. All those painful attempts at fixing compilation errors stay in subagent's context and are discarded once it's done.

2. Claude runs out of context and then forgets things. Debugging stupid compilation errors and test run issues takes a lot of tokens. In about 15–30 minutes, it runs out of context and performs compaction, which gives equal weight to your original request and to the minutiae of how it was failing to make the tests pass.

    **With subagents,** most of the token usage is at subagent level, so top level does not require compactions often. (Introduction of report files a few steps later helps to continue the task even if top-level context gets compacted.) Subagents are short-lived so don't run out of context either.

3. Claude gets confused as to what it's supposed to be doing. Will update tests to match a broken implementation, or will break a perfectly working implementation to get a broken test to pass.

    **With subagents,** the top-level agent can focus on the big picture and only sees a steady progress towards the goal.

    Smart separation of agents helps further. The testing agent will be writing tests and won't try to break your code, while the implementation agent will be trying to get tests to pass without dumbing them down.

4. Claude is hard to steer. Getting it to remember to run all the tests and make sure they all pass was next to impossible.

    **Each subagent** re-reads `CLAUDE.md` and has separate instructions, so you can stuff fresh propaganda straight into its brain. A short lifetime of the agents means they never forget your teachings.

So I got myself a team:

* tech planner / tech lead
* test engineer
* implementation engineer
* code reviewer

My approach goes against some widespread internet wisdom saying that subagents work best for read-only tasks, and that coding should be at the top level.

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

Things were now much better. Some problems went away immediately: no more compactions, no more mangled tests, no more tests left broken.

From this point on, I've never needed `/wat`, `/dumbass`, `/fix-failing-test` and other similar commands for low-level steering again; all of that was solved by the virtual team, and I could focus on producing the kind of code that I wanted.


### Step 5: Linus Torvalds

The new bottleneck was the code reviewer not doing their job. The general quality of the output was quite bad, requiring extensive rewriting or many rounds of detailed feedback.

One day, in a stroke of genius, I asked `/agent` command to build me a dream reviewer:

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
- Are tests maintainable? Could tests be simplified?
- Is core maintainable? Could implementation be simplified?
- Did we overlook some part of the system that also needed to be updated?
- Do we follow the highest standards of quality? Have a fresh look at the code through the eyes of Linus. Would he approve it?
```

This nicely outlines the kinds of problems I was dealing with before.

(Note: this is the _input_ to the `/agent` command, which then expands it into a multi-page set of instructions.)

And it worked. Linus started finding _all_ of the issues I would normally report, and then some. Example review:

```
Alright, let me cut through the bullshit and tell you what I REALLY think about this "implementation."

The Verdict: A Half-Assed Implementation That Nobody Actually Tested

[...]
What's ACTUALLY Wrong With This Garbage

1. THE TESTS ARE FAILING, YOU MORONS: Look at the integration test output - "corrupted journal segment file" errors everywhere! The grep test is
COMPLETELY broken - it can't find ANY of the entries it just wrote. Did anyone actually RUN these tests before declaring victory? Of course not!
2. Journal corruption handling is NONEXISTENT: When the journal package encounters corruption, it just... logs a warning and deletes the file? What kind
of amateur hour bullshit is this? In a production system, you don't just DELETE corrupted data and pretend nothing happened!
3. The done field removal was theater: They removed an unused done bool field from the cursor struct and called it a "code review fix." Congratulations,
you found dead code that should never have been committed in the first place. Want a cookie?
4. Test quality is GARBAGE:
  - Tests write records but don't verify they're actually persisted
  - Journal writes aren't being committed properly (notice the missing Commit() calls?)
  - The integration tests are creating journals but not properly initializing them
  - Error handling in tests just logs and continues - that's not testing, that's PRAYER
5. The implementation is INCOMPLETE:
  - Where's the pagination? They check count >= in.Limit but never actually implemented proper cursor-based pagination
  - The grep filtering happens AFTER fetching all records - that's going to be a disaster with large datasets
  - No timeout handling for the merge operation - what happens when one tenant's journal is slow?
[...]
```

Therapeutic, if nothing else.

New problems:

1. The rest of the team would often forget the review recommendations (they'd focus on one aspect and forget the rest).

2. Linus would often review existing code, not just the new changes, and so would recommend a lot of unrelated modifications.

3. After a round of fixes, Linus would re-review fresh changes without remembering his prior review, so would miss that things were slipping through.

The bottleneck has now shifted from execution to task management.


### Step 6: Preserving reports

Prior to this, I've been asking agents to read and update `aiplan.txt`, thinking that a single file is easiest to find and hardest to miss.

Well apparently everyone was too happy to rewrite it to accomodate the changing narrative, and after a couple of iterations the file failed to preserve the original intent or any significant prior feedback.

So I asked the agents to store any user feedback and their reports in sequentially numbered files like `_tasks/2025-09-whatever-it-is-we-are-doing/11-something.md`. The initial task goes into `01-user-request.md`, and then each agent writes one output file per invocation.

With this, I could teach Linus to limit his review to the scope of the task we're doing, and also to look at all prior reports to see if anything was missed.

Now I needed Claude to match the level of execution to the level of the feedback.


### Step 7: All-star team

First, if Linus was so much better than an average agent, why not give specific personalities to every agent?

* I tried Joel Spolsky as the tech lead. He has legendary attention to detail that I was after, but turns out he also has a legendary focus on shipping and would discard the feedback saying “we could polish this forever but we need to SHIP IT!” So I replaced Joel with [Don Melton](https://donmelton.com/about/) (the creator of WebKit and a guy I enjoyed listening to for hours) with his “I don't care if it works, is it RIGHT?” approach.

* i then re-hired Joel to be a technical planner, running after Don to expand his direction into detailed plans, because I found that Don ain't very keen on writing detailed specs.

* I replaced my low-level reviewer with [Kevlin Henney](https://en.wikipedia.org/wiki/Kevlin_Henney) because I loved his code quality talks. This had a big effect on test quality, because Kevlin insists on tests communicating the intent well.

* I replaced the test engineer with Kent Beck, the implementation engineer with Rob Pike, and used personas for a bunch of secondary agents that I haven't talked about yet. It didn't have an immediately noticeable effect.

This works because Claude knows the personality of the chosen person really well and aligns its actions accordingly. So it cannot be someone unknown.

I also changed the workflow, splitting it into clear planning and execution phases, and asking Linus to review Don's and Joel's plan before implementation.

This is where I'm currently at. Will describe my full setup now, plus some future plans.


## My current setup (aka HOWTO)

### Workflow instructions

Here is my workflow section from `CLAUDE.md`. It has more agents and more details that I've talked about so far, we'll get into those in a moment.

Behold, my very own personal circus:

```
## Process

**CRITICAL: NO CODING AT TOP LEVEL!**

Our star agentic team:

* Don Melton (the tech lead)
* Joel Spolsky (the implementation planner)
* Kent Beck (the test engineer)
* Rob Pike (the implementation engineer)
* Donald Knuth (advanced problem solver)
* Kevlin Henney (low-level reviewer)
* Linus Torvalds (high-level reviewer)
* Raymond Chen (the doc writer)
* Ward Cunningham (knowledge librarian)
* Andy Grove (HR and manager of agents)

We use a task-based workflow to ensure thorough planning, implementation, and review. Tasks are organized under the `_tasks/` directory with per-task subdirectories like `_tasks/YYYY-MM-DD-task-slug/`. Files inside are numbered sequentially: `01-user-request.md`, `02-plan.md`, `03-tests.md`, etc.

**The _tasks/ directory is for DOCUMENTATION ONLY (plans, reports, reviews). All actual CODE (tests, implementation, etc.) goes in the proper codebase locations.**

WE ALWAYS USE SUBAGENTS! THE TOP-LEVEL AGENT ONLY CALLS ON SUBAGENTS!

WORKFLOW - STEP 1 - SAVE REQUEST:

1. User's initial request is saved to `##-user-request.md` or `##-user-revision.md` (`##` = next number)

WORKFLOW - STEP 2 - PLAN:

1. Don analyzes the codebase etc and creates `##-plan.md`
2. Joel expands Don's plan with technical details and creates `##-tech-plan.md`
3. Linus reviews Don's plan and Joel's tech plan.
4. Don, Joel and Linus iterate until Linus approves the plan, repeating all the steps.

WORKFLOW - STEP 3 - EXECUTION:

1. Kent writes tests in the appropriate package in the codebase (see Test Location Rules below) and creates a report in task dir. If stuck, call Donald Knuth.
2. Rob implements code changes in the codebase and creates a report in task dir. If stuck, call Donald Knuth.
3. Raymond updates the docs: API docs in _docs/, AI docs in _ai/. Creates report in task dir.
4. Kevlin and Linus review the changes in parallel. Important: any API docs updates MUST be reviewed for hallucinations by Kevlin very very carefully.
5. Go back to PLAN step so that Don reviews all results, Joel again expands, and Linus reviews. If ALL THREE (Don, Joel, Linus) agree that the task is FULLY DONE, then we're done, otherwise they iterate to come up with a new plan (as PLAN step explains) and then move back to EXECUTION.

WORKFLOW - STEP 4 - FINALIZATION:

1. Ward preserves all the new learnings from these tasks for future reference.
2. If user asked for any corrections during this task, Andy considers updating agent instructions to align them to the user's intent. This is a VERY HIGH BAR, because agent definitions should not be updated lightly.

**IMPORTANT:** PLAN step ALWAYS follows after EXECUTION. Each time we run Kent, Rob, Raymond, or Kevlin, we then must run PLAN step - Don, Joel and Linus.

**CRITICAL: NO CODING AT TOP LEVEL!**
```

### How I start a task

I figured that I'd rather preserve the details for the future, so I add technical instructions to a Linear ticket, example:

```
Treat customers with certain tags (default "Login with Shop") as members

When require active account is configured, if a customer logs in with Shop Pay (Shopify assigns "Login with Shop" tag in this case), we need to treat them as active because they aren't often getting marked as active shop accounts.

Also other custom tags should be usable to make the customer a member.

Impl:

* Introduce a list of additional tags to consider a customer active as a configuration parameter, just like a list of tags to enable/disable loyalty program. Don't forget to update configuration UI.

* Note that "Login with Shop" always makes customer active; the configuration option is additional tags. So there are built-in tags (so far only one) and additional tags.

* “active” here means ShopAccountState = bpub.CustomerAccountStateActive. So basically, in fireback/processing-customers.go, when a configured tag is found amongst customer tags (including shop tags and internal tags, so use Customer.HasShopOrInternalTag / HasAnyShopOrInternalTag), override in.ShopAccountState to active. This requires updating tags before shop state!

* A good test would enable require-active-account behavior and points/tiers, set up a customer, and make sure they have no tier. Then update the customer with the configured tag, and verify that they now have a (base) tier.
```

And then I copy the link from Linear and run e.g.:

```
/do [DEV-1140: Treat customers with certain tags (default "Login with Shop") as members](https://linear.app/bubblehouse/issue/DEV-1140/treat-customers-with-certain-tags-default-login-with-shop-as-members)
```

and Claude uses Linear MCP server to read the ticket, then starts executing.

This `/do` command is currently defined very simply:

```
Your new task:
$ARGUMENTS

---

FOLLOW THE WORKFLOW in CLAUDE.md. Create tasks for each step. Take no shortcuts.
```

There's almost nothing extra here, but the reminder to follow a plan helps to avoid situations when Claude forgets to use the workflow.

You might have noticed that I give specific technical details. That's not the entire plan for the feature; but those are the things I know I want and expect, and that can be implemented multiple ways.

I found that if I definitely want a specific implementation approach, I'd better mention it.


### How I request revisions

After the team is done, I review the code and use `/rev` command to request revisions in freeform text, listing everything that I'd like to change. The command is defined as follows:

```
Requesting revisions:
$ARGUMENTS

---

CONTINUE the task, do not create a new one. (If I didn't mention which task, continue the current one!)

Save the user's revision request to a new file under task dir!

FOLLOW THE WORKFLOW in CLAUDE.md. Create tasks for each step. Take no shortcuts.

- If user is requesting code review, start with Kevlin and Linus first, before everything else.
- Andy (HR) should strongly consider if the agents can be aligned to anticipate similar instructions the next time.

IMPORTANT: Don must run after EVERY agent invocation to decide on future actions. Each time we run Linus / Kevlin / Donald, next is ALWAYS Joel. This process only finishes when all reviews pass with flying colors!
```


### How successful is this team?

Simple requests: the team does it very well, at superb quality level.

Moderately complex requests with technical instructions: the team can ship these with a couple of revisions.

Complex tickets: the team can advance of these and either give you a reasonable starting point, or ship with dozens of revisions and a hit to code quality (ie even more revisions asking to do extra review passes over the code and to fix specific problems).

I take the most complex tasks for myself, while having the AI team work the simple to medium ones; that's most efficient in both time and effort.


### How good is the result?

Note the word “revisions” above. I'm very particular about code quality. Moreover, with AI (and with a less-senior human team), code quality tends to be self-reproducing: if you add a bunch of shit, later that shit is gonna be taken as an exemplary way of doing something (“existing pattern”) and reproduced a few more times.

This is a point worth repeating:

*The best devs say it's shit when they see shit, but AI and less senior teams see existing patterns and reproduce them.*

AI-generated code is great in a sense that it's so much better than what was produced before, but it's not good enough to just commit. Whenever I see a better way of writing something, I would often go for it, either requesting a revision or making the change myself.

And yet, most of the time the generated code works, and has at least reasonable quality.

With one notable exception...


### The quality of tests

I'm still fighting to get the tests to come out exactly how I want them.

I wrote a detailed essay teaching AI the style that I want, [`tests-and-helpers.md`](https://github.com/andreyvit/claude-code-setup/blob/main/tests-and-helpers.md), titled “Questioning the abstraction level of test helpers”, and often point to it when requesting a revision. It is too long to reproduce here, but I've added it to the public demo repository.

I should really include it in the agent instructions, but it's kinda too long for that, so I'm still figuring it out. Meanwhile, often tests come out wrong, and I do `/rev Improve tests, see @_readme/tests-and-helpers.md`.

That _does_ work, though, and it produces good tests. (Not _great_ tests, mind you, but perfectly acceptable ones.)


### Supporting cast (other agents)

* Donald Knuth is a problem solver. This is my way of switching from coding to deep analysis when things get dire. Coding agents tend to try to continue coding; but if stuck, they invoke Donald, and Donald is instructed to refrain from making changes and instead to think hard and provide a deep analysis of the right way forward.

    (I don't see Donald called often, but on harder tasks he does make an appearance and saves the day.)

* The librarian is supposed to update knowledge under `_ai`, but does a very poor job, and I need to explore this area further.

* The HR agent (Andy) updates definitions of other agents; I ask him to make changes to my team (which works great). I also try running him to update agent definitions after revisions, but I discard his changes most of the time. (Initially, when just bootstrapping the system, he did play an important role in instilling the right principles. We've now crossed some threshold, though, and his further contributions are undesirably specific.)

* The documentation writer updates our API docs, but again, AI isn't great at writing good docs right now, so this is hit and miss.


## Remaining problems and future plans

1. Agents fucking suck at sequential numbering of report files and at staying within a single task directory. Such a stupid thing, causing so much friction. I want to build an MCP server that makes this easier on them.

2. Like I've mentioned, getting the tests to be exactly how I want them on the first try is an unsolved problem so far, although I have a workaround.

3. My attempts at getting Claude to handle manual front-end testing were totally unsatisfactory.

4. Writing docs and preserving knowledge for posterity is in a bare-bones state, not much improved after the very first attempts.

5. I hear great things about Codex, and want to try it too. But the lack of subagents is giving me pause.

6. A friend of mine claims successful interactive AI usage. Some day I will explore that side deeper as well.


## Resources

GitHub repository with all agent definitions, custom commands and instructions: [`andreyvit/claude-code-setup`](https://github.com/andreyvit/claude-code-setup/).
