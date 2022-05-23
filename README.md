# bin2png
<h2>Visualise your data using bin2png, converting bits of data to images</h2>

**Usage: bin2png -f FILE -o OUTPUT [other flags]**

bin2png will read your inputted data (works recursively with directories) and will output a square png with every bit visualised as a black or white pixel.

By default a bit value of 1 will return a white pixel and a bit value of 0 will return a black pixel. You can invert this by using the `--invert (-i)` flag in the CLI.

Since we are dealing with 8 pixels per byte of the supplied data, the file sizes can get quite large, if you wish to use a different type of compression you can try out the `--compression LEVEL (-c)` flag, 0 will give you zlib's default compression, whilst 1 will not compress at all; 2 will use slightly less compression, but will be faster and 3 will give you the best compression at a slower speed.

Have lots of fun visualising data! <3

<h2>TODO</h2>

- [ ] Add RGB support (3 bytes per pixel)
- [ ] 16:9 support
- [ ] Add goroutines to make writing pixels faster
