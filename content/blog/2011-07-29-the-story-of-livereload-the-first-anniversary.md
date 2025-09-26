{
  "permalink": "/blog/2011/07/the-story-of-livereload-the-first-anniversary/",
  "topics": ["LiveReload"],
}

# The Story of LiveReload: The First Anniversary

<x-summary>
    How LiveReload 2 was born out of a failed startup and a brief moment in the spotlight.
</x-summary>

In the summer of 2010 we were hard at work implementing [a startup idea](http://www.mockko.com/). That did not go very far, but it did give birth to an [unrelated worthy tool](http://livereload.com/) used and praised by many people.

That summer I've started using [LESS](http://lesscss.org/) and [CoffeeScript](http://jashkenas.github.com/coffee-script/) for the first time. LESS was fantastic, but it also presented a problem: as a long-time loving user of [CSSEdit](http://macrabbit.com/cssedit/), I would never agree to a less-than-instant editing experience.

Mockko (our “startup”) was WebKit-only, so using the Firefox-based [XRefresh](http://xrefresh.binaryage.com/) was not an option. After stuggling for a few days, I set off to write a similar tool for Safari.

It was a Sunday night when I finished putting the final touches on LiveReload 1.0, published the command-line tool and went to sleep, fully expecting maybe 10 downloads over the next week.

I woke up to find some unusual Twitter activity. After a quick check, it seemed like the whole Internet started talking about LiveReload. Envylabs [made a great screencast](http://blog.envylabs.com/2010/07/livereload-screencast/), misattributing the authorship to my cofounder-at-the-time who happened to fix a few lines in README while I've been sleeping. We even got into the list of trending repositories on GitHub.

**1.1–1.4.** The active development continued over the course of another month: more configurability, Chrome extension, [our own file system monitoring gem](https://github.com/mockko/em-dir-watcher) with powerful filtering facilities and all 3 platforms supported.

Soon, however, that progress came to an end. The market was recovering from the crisis and I got a lot more load in my consulting business. I've also parted with my co-founder, and poured the little time I had left into the development of Mockko itself and not the side projects.

**1.5–1.6.** Luckily, LiveReload has attracted [an active contributor](https://github.com/NV) who continued to improve the existing browser extensions and added one for Firefox. With his help we released versions 1.5 and 1.6.

Other people have made [alternative command-line tools](http://rubygems.org/gems/guard-livereload) for our extensions or even [included LiveReload into their apps](http://compass.handlino.com/).

By the spring of 2011, it was clear that I would never continue the development of LiveReload. Mockko has been stagnating for months and I had many ideas for new projects, so spending time improving an already-working open-source tool did not seem reasonable.

On the other hand, I've been using LiveReload a lot in my consulting work, and really hated to run 4 or 5 monitoring tools on the command line. I badly wanted a UI version that would also take care of running the compilers and minifiers.

**2.0.** A LiveReload GUI could help many users and potentially bring in a nice cash bonus. That, however, did not get me excited enough. The ultimate motivator (competition!) arose when I learned that someone else started working on a similar tool. Two weekends later, I had the initial prototype ready, crafted a pure-HTML (meaning no-CSS) web site and offered a private alpha version to several trusted people, calling it LiveReload 2.0.

The feedback has been positive. In fact, I've never heard anything negative about any version of LiveReload. Even those who couldn't install or run it considered the idea somewhere between very useful and incredible. This was encouraging but also troubling, and I've always felt like I was overlooking some catch.

After living through several failed-to-deliver-anything “startups”, I was intending to launch LiveReload 2 very quickly, with 2.0 being on par with 1.6 feature-wise. Based on the feedback, however, I've changed my mind and included compilation support into the roadmap. The design side of things lagged behind the development, so launching quickly was not an option anyway.

It is now the end of July, which means LiveReload has been out for a year. It has gathered [199 Twitter followers](http://twitter.com/#!/livereload/followers), [792 watchers on GitHub](https://github.com/mockko/livereload/watchers), [10788 gem downloads](http://rubygems.org/gems/livereload), 5 alpha builds released and [460 users trying them at least once](http://livereload.com/stats.php). We have a strong team which has been working together for many months, with [NV](http://twitter.com/#!/ELV1S) going from a major LR1.x contributor to a LR2.x partner, focusing on the browser side of LiveReload. We've just [got an icon designed](http://99designs.com/buttons-icons/contests/mac-app-icon-livereload-86859) and preparing to submit to the App Store soon.

What's really exciting is that several competitors are rumored to launch soon too, and even FogCreek has released a [kinda-similarly-targeted tool](http://www.webputty.net/). Stepping into an unknown territory of marketing is equally thrilling: after reading dozens of articles on launching products and getting coverage, would be nice to finally proceed.

So tell me: are you still hitting Refresh?
