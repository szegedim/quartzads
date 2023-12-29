# QuartzAds

## Advertising Technology behind www.showmycard.com℠

# Summary

We created this solution as an alternative to Amazon, Facebook, Google, Microsoft, TikTok Ads, etc.

You can just place ads and pay for them. The specialty is that we share logs, so that you can verify results comparing the performance of designs. It may be valuable to assess the bias of AI based solutions on other ad sites.

The design is very simplistic. It is like the founding age of the internet in the late 1990s. Feel free to update it, this is why we use CC0 license.

The content can be configured at startup making it very flexible. You can add it as a wrapper on any site giving it a potential user base of the entire web.

We do not use cookies. The reason is that they require legal opt-ins that slow down prospective customers. Cookies lower the click through rate. We discourage people to use them.

This is a demo for fundraising by the way. See www.showmycard.com℠

The business logic is simple. You can collect some of the advertisers who did not get their required return at other ad sites due to a competitor with bigger marketing budget.

We encourage distributed marketing as a result. You borrow the content and put ads on it. Make sure it is super relevant. You can cluster more advertisers and pay for traffic.

The system grows until all advertisers find their audience. It may slow at that time, since manufacturers need to invest more capital, etc. that takes time.

Please make sure to buy the rights to any intellectual property that you use. This is especially true as we do remote loading of each page that you put ads on generating true traffic on the upstream site.

Example:

- You create a cluster promoting multiple advertisers. People tend to stay more and spend more in such shopping districts of furniture for example.
- You negotiate and proxy an interior design and furniture page from a generative content site.
- You pay $50 to get 6000 impressions.
- You sell ads on the page to thirty advertisers $2 each for the same period.
- They get the brand marketing and the market data that would cost $50 each otherwise.
- You cash in $10 profits.

# Usage

You have a few options to test how it works.

Create a DigitalOcean account and run the docker container `schmiedent/quartzads` as an app communicating on port 7777.

Create a AWS account and run the docker container `schmiedent/quartzads` as an AWS Lambda communicating on port 7777.

Launch it on a virtual machine, and browse `http://127.0.0.1:7777` Get full logs in `http://127.0.0.1:7777/englang`. You can just click on an ad card to buy it.

You will need to do changes and some security setup to go production.

Create a billing page at a payment provider like stripe.com or paypal.com. Note the Url.

Write an implementation configuration file like the following.
```
Set the payment url to https://buy.stripe.com/test_00gfZueca62I1BC9AB address.
Set the title to Support Wildlife text.
Proxy the https://www.worldwildlife.org site.
```

Place it to a hard to guess location like this one. Lambdas and serverless functions are ideal for this purpose.
```
https://demo.showmycard.com/bc450bc5-77ff-4492-a988-eea45bd17c12.txt
```

Update the environment variables to use this url at startup.
```bash
export IMPLEMENTATION=https://www.showmycard.com/bc450bc5-77ff-4492-a988-eea45bd17c12.txt
```

The defaults point to a WordPress site created by us showing the market potential of upstream content is almost the entire internet. Make sure you buy the copyright before mirroring any site.

We are fundraising for this project at www.showmycard.com

You can see a demo at https://demo.showmycard.com

For more info contact hq@schmied.us

# TODO List

- Payment integration and refunds. Currently, you just leave a contact and use your payment provider.
- Search across peer sites. This is what makes it powerful.
- Full spot bid reporting w/ impressions, clicks, uptime. Design is usually very characteristic and individual to each user. We will probably build a paid version with our own design. Scrape `/englang` logs with Alteryx, PowerBI, Snowflake, Splunk, Tableau.
- Finalize Redis integration. Redis can help to back up and scale with Lambdas to Big Tech levels easily. We added only stubs to streamline the codebase.

# Warranty And License

This is a bare minimum demo. We reserve patent rights. We reserve showmycard.com℠. Any other trade marks and service marks in the codebase (e.g. Stripe™) are properties of their respective owners.

This document is Licensed under Creative Commons CC0.
To the extent possible under law, the author(s) have dedicated all copyright and related and neighboring rights
to this document to the public domain worldwide.
This document is distributed without any warranty.
You should have received a copy of the CC0 Public Domain Dedication along with this document.
If not, see https://creativecommons.org/publicdomain/zero/1.0/legalcode.
