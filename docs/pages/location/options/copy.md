# Copy

Instead of specifying multiple `to` backends for a given `location` you can also use the `copy` option. Instead of recalculating the backup multiple times, you can copy the freshly copied snapshot from one backend to the other, avoiding recomputation.

###### Example

```yaml | .autorestic.yml
locations:
  my-location:
    from: /data
    to:
      - a #Fast
      - b #Fast
      - c #Slow
```

Becomes

```yaml | .autorestic.yml
locations:
  my-location:
    from: /data
    to:
      - a
      - b
    copy:
      a:
        - c
```

Instead of backing up to each backend separately, you can choose that the snapshot created to `a` will be copied over to `c`, avoiding heavy computation on `c`.
