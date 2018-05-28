# hn

In hn package, we wrap `GetItem` within `GetItemByChan` so that we can utilize golang's `chan` primitive.

Pro:

- get items will be running in concurrently.

Con:

- Order is not maintained. It render the articles in random order.
- Slow. We first get all itmes' id, and wait for all items back (controlled by `wg.Wait`).