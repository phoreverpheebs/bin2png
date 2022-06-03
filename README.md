# bin2png
<h2>Features</h2>

- Visualise any data form as an image
- Output to PNG and/or JPEG
  - Adjust compression with PNG files and quality for JPEG files
- Invert pixel values
- Explore huge image files
- Explore your data visually

<h2>Usage</h2>

**`bin2png -f FILE [-o OUTPUT] [other flags]`**

- `-f FILE` (required) -> the file or directory you'd like to read and convert to an image.
- `-o OUTPUT` -> filename of output.
- `-v` -> verbose output.
- `-c LEVEL` -> amount of compression for PNG files, default is 0.
    - for speed use 2
    - for smallest size use 3 (very very slow)
    - for no compression use 1
- `-q QUALITY` -> quality of JPEG output, default is 100.
    - for smallest size and lowest quality use 1
    - for largest size and highest quality use 100
- `--jpeg` -> output a JPEG file.
- `--png` -> output a PNG file.
- `--invert` -> use this flag to invert the black and white pixels. 
    By default a binary bit 1 represents a white pixel and 0 represents a black pixel. This flag will invert those values.
    
<h2>Example</h2>

Our beloved golang source code visualised:

![bin2pngexample](https://user-images.githubusercontent.com/96285600/170034972-6982816e-72f0-4b23-bec4-84ffe7547dee.png)
    
<h2>Thank you for using bin2png <3</h2>

Have lots of fun visualising data! <3

<h2>Credits</h2>

This code uses a partial rewrite of the image/png encoder, so credits go to the Go Authors for writing the package in the first place!

<h2>TODO</h2>

- [ ] Add RGB support (3 bytes per pixel)
- [ ] 16:9 support
- [ ] Add goroutines to make writing pixels faster
- [ ] Optional recursion depth