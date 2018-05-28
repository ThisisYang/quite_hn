# concurrent_v2 branch

## Goal

In concurrent_v1, we improved performance by wrapping `GetItem` in `GetItemByChan` so that we make `GetItem` call in goroutine.

However, in concurrent_v1, it neither maintain the order of articles nor efficient (still will try to retrive all articles)

In this concurrent_v2 branch, I am trying to solve the issue to retain the original order.