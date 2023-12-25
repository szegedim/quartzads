# QuartzAds

## Advertising Technology behind showmycard.com℠

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

You have a few options.

Create a DigitalOcean account and run the docker container `schmiedent/quartzads` as an app mapping port 7777.

Create a AWS account and run the docker container `schmiedent/quartzads` as an AWS Lambda mapping port 7777.

Launch it on a VM.

You will need changes and some security setup to go production. Create a billing page on stripe.com and set the url to `PAYMENTURL`. Update `SITETITLE` with your desired site name. For more info go to hq@schmied.us

# Warranty And License

This is a bare minimum demo.

This document is Licensed under Creative Commons CC0.
To the extent possible under law, the author(s) have dedicated all copyright and related and neighboring rights
to this document to the public domain worldwide.
This document is distributed without any warranty.
You should have received a copy of the CC0 Public Domain Dedication along with this document.
If not, see https://creativecommons.org/publicdomain/zero/1.0/legalcode.

# TODO

- Payment integration and refunds
- Search across peer sites
- Full spot bid reporting w/ impressions, clicks, uptime
- Redis integration
