# Recoder

Basic module to help with removing disposal from GIFs

# Usage

Recoder can be invoked with `RecodeImage(io.Reader, QuantizationAttributes)` or directly via CLI, however, CLI is only intended to demonstrate what it can do.

The purpose of recoder is to handle iterating between frames with certain libraries which do not not handle disposal to automatically apply the previous frame beforehand. It is then expected that output file sizes to be a bit larger than the source image.

````
Usage of recoder:
  -In string
        Input filename
  -Out string
        Output filename
  -Speed int
        Speed (1 slowest, 10 fastest) (default 3)```
````

# Quantization Attributes

The current QuantizationAttributes are MaxColors, MinQuality, MaxQuality and Speed which are provided by ImageQuant. Default values can be retrieved with `NewQuantizationAttributes()`.
