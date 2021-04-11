# Recorder
Basic service to watch a folder and re-encode GIFs

# Usage
Recoder allows you to watch a folder for files ending with `.recode` and will send the output to the result directory. When a file named `.recode` is found, it expects a GIF and will remove any
disposal from the images and ColorModels to allow for easier frame extraction. It does not
automatically scan for GIFs that do not have their disposal not empty so if a user uploads a background as `custom_1.gif`, you should instead save it as `custom_1.gif.recode`

```
Usage of recoder:
  -max_quality int
        Maximum quality for Imagequant (default 100)
  -min_quality int
        Minimum quality for Imagequant
  -resulting_directory string
        Folder where to output the files. Leave default to use same folder as watch.
  -speed int
        Imagequant quantization speed. Speed 1 gives marginally better quality at significant CPU cost. Speed 10 has usu
ally 5% lower quality, but is 8 times faster than the default (default 3)
  -watch string
        Folder to watch for files ending with .recode
```