# Go counter

Count your tokens of your Go programs for code golf purposes.

There are a few reasons for why there is a custom program for this:

1. To discourage people from making their code super ugly by changing all of
   their variables to single letters. Variables of all lengths will count as an
   equal number of tokens.

2. This program adds to this anti-uglify rule in the sense that package imports
   and package declarations are not counted towards the token count. I.e.

   ```go
   package main // not counted

   import ( // not counted
       "fmt" // not counted
       _ "regexp" // not counted
   ) // not counted

   import "fmt" // this style is also not counted

   func main() {} // yes counted!
   ```

   The reason for this is that this allows you to split up your files without
   receiving a penalty for having to add an additional `package main` at the
   top of the second (or n-th) file, and it also allows you to re-import stuff
   you've already imported in the first file without it contributing to the
   count.

3. Comments are also not counted, so if you produce an especially hacky piece
   of code in very few tokens, you can add an explanation for how it works free
   of charge.

   ```go
   // not counted
   ```
