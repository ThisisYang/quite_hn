# concurrent_v5

In this version, we implemented local cache.

The way the cache is implemented is different from README. 
In each request, I get the `TopStories` first and read/write cache on each story level.

In the video, it actually cached all stories as single  item. This can avoid even calling `TopStories` and achieve A/B caching.

The downside of this is that if cache duration is too long, it means `TopStories` id list might be stale. This can be solved by set a relatively short expiration duration and use A/B cache.