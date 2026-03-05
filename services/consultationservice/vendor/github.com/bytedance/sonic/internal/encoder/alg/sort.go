

package alg



func radixQsort(kvs []_MapPair, d, maxDepth int) {
    for len(kvs) > 11 {
        
        
        
        if maxDepth == 0 {
            heapSort(kvs, 0, len(kvs))
            return
        }
        maxDepth--

        p := pivot(kvs, d)
        lt, i, gt := 0, 0, len(kvs)
        for i < gt {
            c := byteAt(kvs[i].k, d)
            if c < p {
                swap(kvs, lt, i)
                i++
                lt++
            } else if c > p {
                gt--
                swap(kvs, i, gt)
            } else {
                i++
            }
        }

        
        
        
        
        
        
        
        
        
        if p == -1 {
            if lt > len(kvs) - gt {
                radixQsort(kvs[gt:], d, maxDepth)
                kvs = kvs[:lt]
            } else {
                radixQsort(kvs[:lt], d, maxDepth)
                kvs = kvs[gt:]
            }
        } else {
            ml := maxThree(lt, gt-lt, len(kvs)-gt)
            if ml == lt {
                radixQsort(kvs[lt:gt], d+1, maxDepth)
                radixQsort(kvs[gt:], d, maxDepth)
                kvs = kvs[:lt]
            } else if ml == gt-lt {
                radixQsort(kvs[:lt], d, maxDepth)
                radixQsort(kvs[gt:], d, maxDepth)
                kvs = kvs[lt:gt]
                d += 1
            } else {
                radixQsort(kvs[:lt], d, maxDepth)
                radixQsort(kvs[lt:gt], d+1, maxDepth)
                kvs = kvs[gt:] 
            }
        }
    }
    insertRadixSort(kvs, d)
}

func insertRadixSort(kvs []_MapPair, d int) {
    for i := 1; i < len(kvs); i++ {
        for j := i; j > 0 && lessFrom(kvs[j].k, kvs[j-1].k, d); j-- {
            swap(kvs, j, j-1)
        }
    }
}

func pivot(kvs []_MapPair, d int) int {
    m := len(kvs) >> 1
    if len(kvs) > 40 {
        
        t := len(kvs) / 8
        return medianThree(
            medianThree(byteAt(kvs[0].k, d), byteAt(kvs[t].k, d), byteAt(kvs[2*t].k, d)),
            medianThree(byteAt(kvs[m].k, d), byteAt(kvs[m-t].k, d), byteAt(kvs[m+t].k, d)),
            medianThree(byteAt(kvs[len(kvs)-1].k, d),
                byteAt(kvs[len(kvs)-1-t].k, d),
                byteAt(kvs[len(kvs)-1-2*t].k, d)))
    }
    return medianThree(byteAt(kvs[0].k, d), byteAt(kvs[m].k, d), byteAt(kvs[len(kvs)-1].k, d))
}

func medianThree(i, j, k int) int {
    if i > j {
        i, j = j, i
    } 
    if k < i {
        return i
    }
    if k > j {
        return j
    }
    return k
}

func maxThree(i, j, k int) int {
    max := i
    if max < j {
        max = j
    }
    if max < k {
        max = k
    }
    return max
}



func maxDepth(n int) int {
    var depth int
    for i := n; i > 0; i >>= 1 {
        depth++
    }
    return depth * 2
}



func siftDown(kvs []_MapPair, lo, hi, first int) {
    root := lo
    for {
        child := 2*root + 1
        if child >= hi {
            break
        }
        if child+1 < hi && kvs[first+child].k < kvs[first+child+1].k {
            child++
        }
        if kvs[first+root].k >= kvs[first+child].k {
            return
        }
        swap(kvs, first+root, first+child)
        root = child
    }
}

func heapSort(kvs []_MapPair, a, b int) {
    first := a
    lo := 0
    hi := b - a

    
    for i := (hi - 1) / 2; i >= 0; i-- {
        siftDown(kvs, i, hi, first)
    }

    
    for i := hi - 1; i >= 0; i-- {
        swap(kvs, first, first+i)
        siftDown(kvs, lo, i, first)
    }
}


func swap(kvs []_MapPair, a, b int) {
    kvs[a].k, kvs[b].k = kvs[b].k, kvs[a].k
    kvs[a].v, kvs[b].v = kvs[b].v, kvs[a].v
}


func lessFrom(a, b string, d int) bool {
    l := len(a)
    if l > len(b) {
        l = len(b)
    }
    for i := d; i < l; i++ {
        if a[i] == b[i] {
            continue
        }
        return a[i] < b[i]
    }
    return len(a) < len(b)
}

func byteAt(b string, p int) int {
    if p < len(b) {
        return int(b[p])
    }
    return -1
}
