# What can I use cff for?

Here are some examples of what you can use cff for:

- Run **interdependent functions** concurrently,
  with a guarantee that a function does not run before its dependencies.

  For example, with cff you can send requests to two different APIs,
  feed results from them into requests made to five other APIs,
  some of which feed into one another, and so on,
  until it all feeds back into two structs that represent the result.

  ```mermaid
  flowchart LR
    A; B; C; D; E; F; G
    dots1[...]; dots2[...]
    X; Y;

    A & B --> C
    B --> D & E
    A & C --> F
    C & D & E --> G
    F & G --> dots1
    G & E --> dots2

    dots1 --> X
    dots2 --> Y

    style dots1 fill:none,stroke:none
    style dots2 fill:none,stroke:none
  ```

  And do all of this with **as much concurrency as possible** from the
  dependency relationships.

  Use `cff.Flow` for this.

- Run **independent functions** concurrently.

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

  You can choose whether you want to stop when the first of them fails,
  or keep going despite failures.

  Use `cff.Parallel` for this.

- Run the same function on every element of a **map** or a **slice**,
  without risk of unbounded goroutine growth.

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

  Use `cff.Parallel` with `cff.Slice` or `cff.Map` for this.

See also [What can't I use cff for?](non-use-cases.md)
