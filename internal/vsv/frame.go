package vsv

import (
  "fmt"
  "image"
  "image/color"
)

type Frame struct {
  Data  []byte
  
  imageRectangle image.Rectangle
  colorModel     color.Model
}

func NewFrame(data []byte, imageRectangle image.Rectangle, colorModel color.Model) (*Frame, error) {
  if colorModel != BGRColorModel &&
     colorModel != BGRAColorModel {
      return nil, fmt.Errorf("Frame color mode is not correct, pending 'BGRColorModel', 'BGRAColorModel'") 
  }

  return &Frame{
    Data: data,
    imageRectangle: imageRectangle,
    colorModel: colorModel,
  }, nil
}

func (f *Frame) ColorModel() color.Model {
  return f.colorModel
}

func (f *Frame) Bounds() image.Rectangle {
  return f.imageRectangle
}

func (f *Frame) BGRAt(x, y int) color.Color {
  width := f.imageRectangle.Max.X
  offset := (y * width + x) * 3
  b := f.Data[offset]
  g := f.Data[offset + 1]
  r := f.Data[offset + 2]
  return BGR{b,g,r}
}

func (f *Frame) BGRAAt(x, y int) color.Color {
  width := f.imageRectangle.Max.X
  offset := (y * width + x) * 4
  b := f.Data[offset]
  g := f.Data[offset + 1]
  r := f.Data[offset + 2]
  a := f.Data[offset + 2]
  return BGRA{b,g,r,a}
}

func (f *Frame) At(x, y int) color.Color {
  switch(f.colorModel) {
  case BGRColorModel:
    return f.BGRAt(x, y)
  case BGRAColorModel:
    return f.BGRAAt(x, y)
  }
  // incorrect Color mode
  panic(fmt.Sprintf("Unexpected behavior: Color model switch incorrect value of %+v", f.colorModel))
}
