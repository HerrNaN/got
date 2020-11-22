# Got
This is supposed to be a subset of the famous version control system
Git, but written in Go.

I've always wanted to understand Git a little better. Since it is such
an integral part of my development workflow I thought I'd try to implement
it myself to understand it better.

## Currently supported actions/features
Porcelain:
- `got init`
- `got add <filespec>...`
- `got restore [--staged] <filespec>...`
- `got status`
- `got commit -m <message>`
- `got branch {-d <branchname> | --list | <newbranch>}`
- `got checkout`

Plumbing:
- `got hash-object [-w] <file>...`
- `got cat-file { -t | -p } <object>`
- `got read-tree <object>`
- `got write-tree`
- `got update-index [--add] <file>`