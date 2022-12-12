# cff: a concurrency toolkit for Go

cff (pronounce *caff* as in caffeine) is a library and code generator for Go
that makes it easy to write concurrent code in Go.

It gives you:
* Bounded resource consumption: cff uses a pool of goroutines to run all operations, preventing issues arising from unbounded goroutine growth in your application.
* Panic-safety: cff prevents panics in your concurrent code from crashing your application in a predictable manner.

cff can be useful when you are trying to:

1. Run interdependent functions concurrently, with a guarantee that a function does not run before its dependencies.
  ```mermaid
  flowchart TD
    A; B; dots[...]; H

    done(( Done ))

    A --> done
    B --Error--x done
    dots -.-> done
    H --> done

    style done fill:none,stroke:none
    style dots fill:none,stroke:none
  ```

2. Run independent functions concurrently.
  ```mermaid
  flowchart TD
    A; B; dots[...]; H

    done(( Done ))

    A --> done
    B --Error--x done
    dots -.-> done
    H --> done

    style done fill:none,stroke:none
    style dots fill:none,stroke:none
  ```

3. Run the same function on every element of a map or a slice, without risk of unbounded goroutine growth.
  ```mermaid
  flowchart RL
    subgraph Slice ["[]T"]
      i0["x1"]; i1["x2"]; dots1[...]; iN["xN"]
      style dots1 fill:none,stroke:none
    end

    subgraph Map ["map[K]V"]
      m1["(k1, v1)"]; m2["(k2, v2)"]; dots2[...]; mN["(kN, vN)"]
      style dots2 fill:none,stroke:none
    end

    subgraph Workers
      direction LR
      1; 2
    end

    Slice & Map -.-> Workers
```

See our documentation at https://uber-go.github.io/cff for more information.

## Installation

```bash
go get -u go.uber.org/cff
```

## Project status

At Uber, we've been using cff in production for several years.
We're confident in the stability of its core functionality.

Although its APIs have satisfied a majority of our needs,
we expect to add or modify some of these once the project is public.

That said, we intend to make these changes in compliance with
[Semantic Versioning](https://semver.org/).
