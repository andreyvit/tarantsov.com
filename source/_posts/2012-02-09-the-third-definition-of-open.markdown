---
layout: post
title: "The third definition of open, or How I nearly picked GPL for my product, but ended up simply publishing the source with no license (for now)"
date: 2012-02-09 16:50
comments: true
categories: 
permalink: /blog/2012/02/the-third-definition-of-open/
---

There is an urge that every hacker has: to share whatever he creates.

When I got my first day job, the owner of the company explained that he always wanted to run a “fully open-source company” with open accounting and everything. This may seem nuts to a business person, but to me it sounded like the only natural thing to do. I still can't imagine wanting to run a business any differently.

Fast-forward about 6 years. Apple and Google battling over who's better at bending the concept of openness, John Gruber making good laughs of that, and me finishing up my first launchable product.

To give you some perspective, for the last 4 years I've been trying to get into the product business. My first idea has been to make a killer IDE for Ruby on Rails, with the beauty of TextMate and powerful code analysis and refactoring. I've gathered a great team, and it could have been an easy job, but I've been way too immature to run a product. While playing with fun stuff, we've ran our consulting business right into the ground, found ourselves arguing a lot and separated.

Then a similar thing happens with [another idea](http://www.mockko.com/). I don't have enough steam to do a big project on my own, so I get a really good partner, and we work for a little while, and then he gets drowned in other work (and me too, trying not to repeat the previous experience), splits up with his wife and leaves for another city. I continue for some time, but a proper launch seems years ahead and I'm hopelessly undermotivated.

Meanwhile, the livereload gem, which we've hacked together to make coding up our single-page web app more tolerable, gets pretty popular. I always meant to make a GUI wrapper for it, but since I already had my “startup” as a fun out-of-normal-work project, it was impossible to justify spending time on another fun project which was never going to bring any money.

What all of that boils down to, is: there are ideas I really want to implement, but doing that while also doing contract work leads both into wrecks. I also have a wife, a 7-year-old daughter and expecting a second child soon, so I need money to live on — and not just pizza-and-ramen kind of money; I feel oblidged to provide a decent lifestyle to my family.

So I decide that LiveReload GUI app could be the middleground I am looking for. A small app I have time to develop and launch, bringing in the additional income which might allow me to partly ramp down on consulting and get more awesome things into the world.

Well that, and I'm really pissed with having to baby-sit the livereload gem every time I want to fix a bunch of styles on a web page.

So LiveReload 2 is really about the money. I don't want to risk that goal, but I still believe sharing is right, so I email Richard Stallman asking to talk me into going open-source.

His response has been great and led me to carefully think through my values, but on the other hand it completely missed the point for me. See, GPL is very clear-cut when dealing with large companies (which will surface in just a few moments) and their questionable intents. They want to screw us over, we defend and screw them in return.

Where this system breaks down is a community where an absolute majority of people have good intentions, something that indie/early-stage-startup community qualifies for in my opinion. I'd certainly prefer to work on great things for free, and whatever resources and occasional luxuries we need to magically materialize themselves. The fact that people work best when they don't need to think about money is pretty well-studied.

Some people are fine with working on a day job and hacking on free software in the free time. I can only envy those guys; forcing myself to do consuting sucks up so much energy I have pretty much nothing left. And it's not like I don't like what I do; I have really awesome clients and I love the projects. But other people's projects just don't excite me for long, no matter what.

So guys like me go into indie or startup business. And we don't have the option of stuff materializing itself. The options we do have is to either sell stuff for money, or get funding, or tell our families to fuck themselves.

The reason the idea of GPL breaks down for us is that we're very motivated guys when working on our own stuff, and the world needs us to do our best work. We are also happy to share and contribute whenever we can; the better we are doing, the more awesome stuff we can make free and open-source. But we need to be making money to be banging on our stuff; it's already hard enough, and adding obstacles to that is not helping either us, or users, or other developers, or free software community, or the world in general.

I think deep down Richard Stallman believes that the kind of life we want to live is very wasteful, and our desires are unresonable. I disagree. It's ok (and desirable!) for everyone to be doing equally well, but it's not ok for everyone to be doing equally poorly. So far, our technological progress really sucks; we can't even produce enough resources to feed everyone in the world, less so to cure everyone and to enable everyone to pursue their dreams. And yet I want my children to be well-fed, healthy and enabled.

So back to LiveReload, here's chapter two in which large companies go visiting and try to get us into a tight place (also known as asshole).

While I was going rounds about releasing or not releasing the code, a certain company has contacted me with a bunch of questions on how my browser extensions operate, having done surprisingly little homework. It went like, is open-square-bracket the only requirement of your protocol? — WAAT? didn't you actually read the protocol spec? — Yes, but we don't trust specs, so we reverse-engineered your Firefox extension. — WAAT? you didn't notice the extensions are all open-source on GitHub? — etc etc.

And then it went like: Will LiveReload be mentioned in your UI? — No, sorry, this feature will really be invisible. — Will you provide a link back somewhere? — No, sorry, I'm not authorized to make such promises, but maybe the documentation guy will put it somewhere. — Can you put me through with that guy? — No, sorry, he's not involved yet and he cannot make promises either. — Can I talk to someone who can negotiate the terms of usage? — No, sorry, I'm the only one responsible for this feature. — Will the extensions themselves provide any user-visible indication of being part of LiveReload? — Sorry, I cannot to make any promises.

Initially I told them extensions and livereload.js are MIT-licensed, which btw they are (otherwise guard-livereload and others wouldn't be able to bundle them). But I happened to not publish a license on GitHub by mistake, which gave me an opportunity to think it through.

Could GPL be an answer to this shit? It could well be. You'd release your extensions under GPL, and the guys would have to fork your project, release their changes, and then maybe have their app download and install those extensions at runtime to avoid bundling them with the app. Not like this scenario really helps *you*, but there is some chance that doing bullshit like that is harder than to give a proper attribution to the original project.

As irony goes, however, the IDE that company makes *is* under GPL.

Is it a big win of the free software movement to have a bunch of douchbags create a crappy IDE which anyone can use, change and redistribute? It most certainly is. We should stop and reflect on this; 20 years ago douchebags did not share their source code.

That's only one kind of freedom the world needs, however. It is most unfortunate that the free software advocacy is largely missing out on startup boomers (can I suggest the term?) The forefront of free software, GPL, is about fighting douchebags with proprietary codebases. It does not help indie developers get the most into the hands of people (including the most free stuff!), doesn't help them share the most they could share, and it does not help against the new GPL-equipped douchebags.

So back to LiveReload once again. Ihe initial plan has been to call it licensed under “GPL with moral strings attached”. Legally it would be GPL, so you can't get into trouble for using it, but morally it would require a payment unless you have a good reason not to pay (for which being a student with little money certainly qualifies).

The approach is very fragile, though, because it heavily relies on your exact message getting through to everyone, intact.

One problem may be that the legal side does not accurately reflect the intentions. GPL says “fork, redistribute, rename, sell”, while in reality I'm not comfortable with people forking LiveReload and selling it under a different name unless I abandon the project.

However, I would say undermining the legal side can only be a good thing, given the abysmal state of copyright and software patents in US. Not having the rights at all and not having money to defend your rights in court is pretty much the same, so you might as well give them up and save the trouble.

A bigger problem is that people associate GPL with “free”. No matter how much RMS talks about freedom, the primary meaning of “free” is not the one RMS wants it to be (like no matter how much we complain, the primary meaning of “hacker” won't ever be “ingenious computer enthusiast”). Now, if I can stare into your face and explain the stuff about GPL being the legal side, and me being in it for money, and the moral strings, you will get it all right. My staring capability is limited, though, and people who simply hear GPL will get a very different idea.

So, I decided to completely avoid the legal side for now and release LiveReload on pure moral terms, described on the home page. I'm happy with people exploring it, modifying it and sharing their modifications, but by default I want to get paid. Fork the source code, but don't fork the project; I'm not happy with someone taking my app and starting to sell it under a different name. I don't want you to publicly distribute the binaries either, because people coming to my web page are as much of an asset as the money are, and both page visits and money will help me achieve my bigger goals and ultimately produce more fun stuff that everyone will benefit from.

(If you do contribute extensive modifications, I'm sure we can work out a way to make you happy about that.)

Of course, you are free to choose your own moral grounds. I don't want to extort the money from you. If you cannot pay, or you're tight on money (I know I sometimes am), or you believe payment is not morally right in your case, go ahead and use it for free. (So far, it may be a bit tricky. You can ask me for a free license, or you can get a copy from someone — there is no copy protection. I will try to work this out in a better way in the future, but I need to balance that with not sending a message of payment being a highly optional thing.)

And to top that, if I abandon the project, or die without assigning a maintainer, or just turn into Oracle, you are free to take over and run it or fork it.

So there you have it. You can define being open as “repo sync; make; make install”, which is bullshit if you then screw those who did it, or as not doing vendor lock-in and using open technologies, which is just as much a stretch with Apple, despite me being a fanboy. Or you can define “open” as sharing as much as you can share, while being honest about your intentions and NOT being a dickhead. (Oh yeah, and not stealing the address books of your users. Just sayin'.)
