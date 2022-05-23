# bin2png
<h2>Visualise your data using bin2png, converting bits to black or white pixels</h2>

**Usage: bin2png -f FILE -o OUTPUT [other flags]**

bin2png will read your inputted data (works recursively with directories) and will output a square png with every bit visualised as a black or white pixel.

By default a bit value of 1 will return a white pixel and a bit value of 0 will return a black pixel. You can invert this by using the --invert (-i) flag in the CLI.

Have lots of fun visualising data! <3

<h2>TODO</h2>

- [ ] Add RGB support (3 bytes per pixel)
- [ ] 16:9 support
- [ ] Add goroutines to make writing pixels faster
