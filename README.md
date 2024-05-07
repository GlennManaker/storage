Storage:
1) I decided to use AVL tree inside storage to guarantee O(logN) in any case
2) to insert we should provide with key, value and expired at date
3) to delete automatically expired data, we should start storage and then  it spawns goroutine that inside has select case expression to wait for a ticker with the earliest expiration date