# My all-star zoo, or why I hired Linus Torvalds, Kent Beck, and Rob Pike for my AI team

<x-summary>
This summer, I transformed **Claude Code** from a frustrating junior assistant into a high-performing engineering team.
</x-summary>

This summer, I transformed Claude Code from a frustrating junior assistant into a high-performing engineering team. The secret? I gave my AI agents the personalities of legendary developers. The results were extraordinary: Linus Torvalds became my code reviewer, delivering brutal but brilliant feedback. Kent Beck brought test-driven discipline. Rob Pike implemented with pragmatic elegance. Together, they took my AI coding from "smart junior developer" to "solid senior engineer" level—a two order-of-magnitude improvement in both task complexity and output quality.

The more experienced you are as a developer, the harder it is to get value from AI coding tools. As a CTO with 25+ years of experience, I struggled with this paradox: AI easily helps beginners, but for senior developers, the gains are elusive. We code faster, tackle more complex problems, and maintain higher quality standards. Most AI output simply doesn't meet our bar.

But here's what I discovered: for senior developers, the benefit isn't time savings—it's energy conservation. AI excels at the boring, straightforward parts that drain your willpower. More importantly, it gives you a starting point when you're tired, unfocused, or procrastinating. Even if you discard most of what it produces, you're not staring at a blank screen.

This article chronicles my journey from basic Claude Code to a fully orchestrated team of AI agents, each step bringing an order-of-magnitude improvement. If you're struggling to get value from AI coding, this blueprint will show you exactly how to level up.

## Part 1: Why This Matters — The Seven-Step Evolution

### Steps 1-2: Basic Setup Still Underwhelms

My initial Claude Code experience was immediately better than Cursor or Windsurf on my large production codebase. But "better" still meant underwhelming results. I wrote instructions, created documentation folders, even had Claude study the codebase and generate its own docs in an `_ai` folder.

The problem? AI-written documentation was terrible. It would hyperfocus on random implementation details while missing architectural patterns. It remembered things that didn't matter and forgot things that did. When I let it update human-written docs, they immediately degraded. So I quarantined AI docs in their own folder—good enough for context, but not for humans.

### Steps 3-4: The Planning Phase That Gets Forgotten

Following Claude's own recommendations, I added a planning phase. The AI would write a plan to a text file, then execute. This helped briefly—until context limitations kicked in.

The enemy was compactions. Once Claude hit its context window limit, it would compress the conversation. But compaction treats all information equally. The critical high-level plan got the same weight as debugging a typo in a method name. After compaction, quality plummeted. The AI would drift off-task, forget user intent, even forget what it was building.

I watched it burn enormous context trying to get tests to compile because it hallucinated method names. Hours of low-level debugging, while the actual task objective faded into obscurity.

### Step 5: The Subagent Breakthrough

Claude released subagent support, and initially I didn't get it. Their examples—code researcher, code reviewer—seemed trivial. Then it clicked: move the context-hungry work to subagents.

My first team:
- Planning subagent
- Test engineer
- Implementation engineer
- Code reviewer

I created a workflow: plan → test → implement → review. If review failed, loop back.

Results were immediate. No more compactions. Clear separation between test-writing and implementation modes (no more "fixing" tests to make them pass). I could give specific instructions to each agent. I told the test engineer: "Add only one test at a time. Get it to compile before moving on." This alone saved enormous context.

But the reviewer was too lenient. It would approve with "comments for the future" that never got addressed.

### Step 6: Enter Linus Torvalds

Frustrated by poor reviews passing, I added a second reviewer: Linus Torvalds, "in his signature harsh manner."

The impact was shocking. Here's an excerpt from his first review:

> **The Verdict: A Half-Assed Implementation That Nobody Actually Tested**
>
> Looking at this multi-tenant journal viewer, I see classic developer behavior: They wrote code that LOOKS like it works, slapped together some tests that LOOK comprehensive, but nobody actually RAN the damn thing end-to-end.
>
> **THE TESTS ARE FAILING, YOU MORONS:** Look at the integration test output - "corrupted journal segment file" errors everywhere! Did anyone actually RUN these tests before declaring victory? Of course not!
>
> **What Other "Reviewers" Missed**
>
> If other agents approved this garbage, they should be ASHAMED. They probably looked at the code structure, saw some tests, and rubber-stamped it.

He found issues I would miss. He'd identify future production problems, performance bottlenecks, architectural flaws. When he called out "Whoever wrote this shit needs to be fired immediately," it was oddly therapeutic—he expressed all my frustrations with AI's typical shortcuts.

The quality improvement was an order of magnitude. But implementation still struggled with his feedback.

### Step 7: Building the All-Star Team

If Linus worked so well, why not legendary developers for every role?

I asked Claude who should be my test engineer. Obviously: Kent Beck. Implementation engineer? I considered John Carmack but wanted someone who reflected my pragmatic values. Rob Pike from the Go team was perfect. Project manager with attention to detail? Joel Spolsky. I added Kevlin Henney as a second reviewer, Raymond Chen for documentation, Don Melton for high-level planning.

The personalities mattered. I even "fired" Joel at one point—his personality was too focused on shipping, approving subpar work to "iterate later." In an AI context, there is no later. I rehired him for detailed technical planning under Don Melton's high-level direction.

I tried getting the agents to instill more personality into themselves: "Don't your current instructions reflect your personality well? Feel free to rewrite them." They would add more of their own style, making the simulation surprisingly effective.

The team became what I called my "circus"—watching Linus savage the other agents' work while Joel tried to ship and Don insisted on quality was genuinely entertaining.

### The Critical Innovation: Planning Iteration

The breakthrough wasn't just personalities—it was forcing planning iteration. I had Joel create a plan, then Linus review it, then Joel revise, then Linus review again. Only when Linus approved did we move to implementation.

After implementation and review, we'd go back to planning. Don, Joel, and Linus would all have to agree we were done. No more shipping half-finished work.

## Part 2: How to Build Your Own Circus

### The Current Workflow

Here's the exact three-phase system running on Claude Opus under my Claude Pro subscription:

**Phase 1 - Planning:**
1. Don Melton creates high-level plan
2. Joel Spolsky expands with technical details
3. Linus reviews both plans
4. Iterate until Linus approves

**Phase 2 - Execution:**
1. Kent Beck writes tests (one at a time, ensuring each compiles)
2. Rob Pike implements code
3. Raymond Chen updates documentation
4. Kevlin Henney and Linus review in parallel
5. Return to Phase 1 (Don→Joel→Linus) to plan next iteration

**Phase 3 - Finalization (when all three planners agree we're done):**
1. Ward Cunningham preserves learnings
2. Andy Grove (HR) updates agent instructions based on user feedback

The key commands:
- `/do [task]` - starts a new task following the workflow
- `/rev [feedback]` - requests revisions, triggering the full workflow

### Creating Your Own Legendary Agents

Here's the actual prompt I used to create Linus:

```
Linus Torvalds doing very high-level review of the changes in his signature ruthless and pragmatic style. Must run after the normal code reviewer. Focuses on high-level details only, not on the code minutea:

- Did we do the right thing? Or did we do something stupid?
- Did we cut corners to finish faster and called it a success?
- Have we implemented everything requested, or forgot something?
- Do the changes align with user intent?
- Is implementation at the right level of abstraction?
- Do our tests actually test what they claim?
- Did we overlook some part of the system that needs updating?

Be very picky and very ruthless. Use strong language. Tell other agents what you think about them. Force them into perfection.
```

Claude's `/agent` command expands this into detailed instructions. The key is capturing the essence of what makes each developer legendary:

- **Kent Beck**: Test-first discipline, one test at a time, ensuring each provides value
- **Rob Pike**: Pragmatic simplicity, clear code, Go philosophy
- **Joel Spolsky**: Technical specification detail, exhaustive planning
- **Don Melton**: High-level architecture, quality gates, shipping discipline
- **Kevlin Henney**: Code craftsmanship, naming, design patterns

Choose developers who reflect your values. I picked Go creators because their pragmatism matches mine. You might prefer different philosophies.

### Implementation Details

Store agent definitions in your project. Create an `_ai/` folder for AI-managed documentation (keep it separate from human docs). Set up the workflow commands to enforce the full cycle—no shortcuts.

The magic isn't in the tools, it's in the process:
1. Separate concerns (testing vs. implementation)
2. Force planning iteration before and after implementation
3. Require unanimous approval to finish
4. Give agents permission to be harsh
5. Let personalities clash—it improves output

## Part 3: Ongoing Challenges

### What Still Doesn't Work

**HR Agent (Andy Grove):** Should update agent instructions based on user feedback, but proposals need heavy editing. It overfocuses on minor issues while missing systematic improvements.

**Librarian (Ward Cunningham):** Hit-or-miss on preserving useful learnings. I treat all updates as proposals requiring review.

**Problem Solver (Donald Knuth):** Rarely engages. I've seen it work once—when the test engineer couldn't figure out how to test a particular component. It worked brilliantly that time, but it's too rare to evaluate properly.

These aren't failures—they're experiments in progress. The core workflow delivers massive value even with these rough edges.

The journey from basic Claude Code to this circus took me from struggling with junior-level mistakes to getting solid senior-level output. Two orders of magnitude improvement, achieved purely through workflow and personality evolution, not model upgrades.

I continue experimenting weekly and will share further results. My next frontier: trying out Codus, which I'm hearing excellent things about.

*If you enjoyed this deep dive into AI coding workflows, [subscribe to my newsletter](#) for more experiments and discoveries.*

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
