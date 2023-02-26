package vsv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
)

var (
  prefix = [4]byte{0x56, 0x53, 0x56, 0x32}
  
  color_rgb = [4]byte{0x52, 0x47, 0x42, 0x00}
  color_rgba = [4]byte{0x52, 0x47, 0x42, 0x41}
  color_yuyv = [4]byte{0x59, 0x55, 0x59, 0x56}
)

type Header struct {
  Prefix     [4]byte       // 4 byte
  Color      [4]byte       // 4 byte
  Width      uint32        // 4 byte
  Height     uint32        // 4 byte
  PixelSize  uint32        // 4 byte
  FrameRate  uint32        // 4 byte
                           // -------
  //_ [CHUNK_SIZE - 24]byte  // 24 byte
}

func NewHeader(data []byte) (*Header, error) {
  header := Header{}

  if err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &header); err != nil {
    return nil, fmt.Errorf("Header parse error %w", err)
  }

  if bytes.Equal(header.Prefix[:], prefix[:]) == false {
    return nil, fmt.Errorf("Header prefix '%x' is not correct '%x'", header.Prefix, prefix)
  }

  if bytes.Equal(header.Color[:], color_rgb[:]) == false &&
     bytes.Equal(header.Color[:], color_rgba[:]) == false {
    return nil, fmt.Errorf("Header color mode '%x' is not correct, pending '%x', '%x'", header.Color, color_rgb, color_rgba)
  }

  return &header, nil
}

func (h *Header) ColorMode() color.Model {
  if bytes.Equal(h.Color[:], color_rgb[:]) {
    return BGRColorModel
  }

  if bytes.Equal(h.Color[:], color_rgba[:]) {
    return BGRAColorModel
  }

  // incorrect Color mode
  panic(fmt.Sprintf("Unexpected behavior: Color model incorrect value of %+v", h.Color))
}

func (h *Header) Clone() *Header {
  return &Header{
    Prefix: h.Prefix,
    Color: h.Color,
    Width: h.Width,
    Height: h.Height,
    PixelSize: h.PixelSize,
    FrameRate: h.FrameRate,
  }
}

func (h *Header) Bounds() image.Rectangle {
  return image.Rectangle{image.Point{0, 0}, image.Point{int(h.Width), int(h.Height)}}
}

func (h *Header) FrameChunks() uint64 {
  size := uint64(h.Width) * uint64(h.Height) * uint64(h.PixelSize)

  chunkCount := size / CHUNK_SIZE
  if (size % CHUNK_SIZE) != 0 {
    chunkCount += 1
  }

  return chunkCount
}
