## Why Go (golang) ?

At the time that we (Jee-Bee and me) started the PDM in April 2022, we obviously used Python but I soon ran into some issues. For starters I did want to use PySide 6, but at that time it wasn’t there in most of the Linux repositories. Then I also ran into some other issues of Python that I didn’t like, such as v2 / v3, importing a file and the lack of types or (even worse) weak types. The PySide 6 issues and the weak types were the limit. So I went searching for an alternative programming language. Today everything runs on Python 3 and PySide 6, but weak types are still bad.

Why didn’t I use C++ or Rust? The code can be cruel, to read and write. I also just don’t like the syntax of C++ and Rust, and never will like them. Sorry. That doesn’t mean that I can’t participate in FreeCAD (C++ and Python), it’s just that I like a simple language with types and GC. And if I want to use a non-GC (or optionally GC) PL then it needs to be sane, such as Zig, V, or Nim.

Go is fast, both in compile time, a couple of seconds with [air](https://github.com/air-verse/air), and it is roughly only twice as slow as C. That is fast! Go is safe, not as safe as Rust, but a lot more safe than C / C++. Go also has concurrency, and although I try to stay away from it, concurrency is interesting. Go is also neat. The code always looks good. Everyone can read, and understand, the code that someone wrote. Go is also very stable and when it compiles, most of the times, it just runs.

Dealing with errors is necessary. That is not a matter of taste. Some Python users think Go errors are “the good, the bad and the ugly”, but errors are well thought of in Go and they work very well. Just use errors!

I can write code in the PDM server, in the file system (vault), the samba server, the authentication works and I am going to write the user interface of the vault, user and admin code. It all works, I have plenty of tests for the file system (vault), and tests runs within a second, even with the slowest machines.

Docker works pretty well with Go code. The only problem (well, three are thousands of problems...) that I am facing right now is turning samba and postgresql on and off to use docker. That is why I created a switch in the Makefile. I still don’t completely understand how I can update my files after I changed them. Should I run the command line version of FreeCAD inside docker or just change the files with the xml code? The latter also works, so I think it is only a matter of taste.

A good editor such as vscode helps a lot. I like how vscode works and I haven’t met the problems that some have. Maybe in the future I need to switch to vim or something else, but then I face the unlearning part that I don’t like.

The ecosystem of Go is also very good. There are plenty of FOSS stable libraries and the standard library is excellent.

Go compiles one massive blob for each main() function. I like that more than dealing with all the libraries that you need to put in certain directories, regularly update them, and need to use package managers. You can create a Makefile and deal with the compilation with ease. It also simplifies the compilation of the source code. A Go binary can be huge (tens of megabytes is normal), yes, that is the right, but file storage is cheap.

What I don’t like about Go that it is strict. Whenever you make a mistake you can’t compile it. Sometimes that is AARGH! But then it is also coffee time...


### Reading material:
- Effective Go
- Standard library
- Godoc
- 100 Go mistakes
