# concurrent_v2 branch

## Goal

In concurrent_v1, we improved performance by wrapping `GetItem` in `GetItemByChan` so that we make `GetItem` call in goroutine.

However, in concurrent_v1, it neither maintain the order of articles nor efficient (still will try to retrive all articles)

In this concurrent_v2 branch, I am trying to solve the issue to retain the original order.

## `done` channel should be removed

In concurrent_v1, we have a `done` channel to notify that we got 30 articles back. However, since we are using goroutine, it is not guaranteed that the response will be back in order. In this case, we can not simple close `done` channel cause we don't know if the returns are the first 30 articles or not. We have to wait for all items to be sent back.

## Solution

In this version, I passed a `Seq int` argument to `GetItemByChan` function. The `Seq` is the sequence of the article. After retriving the article (calling `GetItem`), I include the `Seq` in response struct.

The main process received response struct, insert the response to a slice initiated previously. `Seq` is used as index (position).

After receiving all response, Another `for` loop traverse the slice, break at the 30th article and return the response.

## Con and Pro

Pro:

- article orders are maintained

Con:

- The hack-noon api would return 500 items each time. I have to retrive all 500 items first.
- number of works (goroutine) is not under control. In future, If hn api return 10 million items, this will spawn 10 million goroutines

