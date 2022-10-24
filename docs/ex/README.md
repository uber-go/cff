This directory holds example code for documentation to pull from.
Place code and tests inside subdirectories here.

Then to pull code snippets in documentation,
add a code block with an `mdoc-exec` attribute.

```markdown
```go mdox-exec='cat ex/foo/bar.go'
// ...
```

This will run the specified command and use its output as the code snippet,
when you run `make fmt` on the documentation.

Use the `region` shell script as the command to pull marked sections
instead of the whole file.

```plain mdox-exec='region' mdox-expect-exit-code='1'
USAGE: region FILE REGION1 REGION2 ...

Extracts text from FILE marked by "// region" blocks.
```

Regions are marked by `// region` and `// endregion`.

```
foo
// region myregion
bar
// endregion myregion
baz
```

`region $FILE myregion` above will return only `bar`.
