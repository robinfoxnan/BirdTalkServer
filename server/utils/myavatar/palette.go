package myavatar

import (
	"image/color"
)

var defaultPalette = []color.Color{
	color.RGBA{0xff, 0xa7, 0x26, 0xff},
	color.RGBA{0xb3, 0x9d, 0xdb, 0xff},
	color.RGBA{0xc0, 0xca, 0x33, 0xff},
	color.RGBA{0xff, 0xab, 0x00, 0xff},
	color.RGBA{0x60, 0x7d, 0x8b, 0xff},
	color.RGBA{0xd1, 0xc4, 0xe9, 0xff},
	color.RGBA{0xff, 0xc4, 0x00, 0xff},
	color.RGBA{0xff, 0xcc, 0xbc, 0xff},
	color.RGBA{0x00, 0xb0, 0xff, 0xff},
	color.RGBA{0xae, 0xd5, 0x81, 0xff},
	color.RGBA{0xbf, 0x36, 0x0c, 0xff},
	color.RGBA{0xf5, 0x00, 0x57, 0xff},
	color.RGBA{0xfd, 0xd8, 0x35, 0xff},
	color.RGBA{0x90, 0xa4, 0xae, 0xff},
	color.RGBA{0x37, 0x47, 0x4f, 0xff},
	color.RGBA{0xe1, 0xbe, 0xe7, 0xff},
	color.RGBA{0xb3, 0x88, 0xff, 0xff},
	color.RGBA{0x30, 0x3f, 0x9f, 0xff},
	color.RGBA{0xff, 0x6f, 0x00, 0xff},
	color.RGBA{0xbc, 0xaa, 0xa4, 0xff},
	color.RGBA{0x45, 0x5a, 0x64, 0xff},
	color.RGBA{0xc6, 0x28, 0x28, 0xff},
	color.RGBA{0xce, 0x93, 0xd8, 0xff},
	color.RGBA{0xae, 0xea, 0x00, 0xff},
	color.RGBA{0xff, 0xab, 0x91, 0xff},
	color.RGBA{0x00, 0x83, 0x8f, 0xff},
	color.RGBA{0xff, 0x3d, 0x00, 0xff},
	color.RGBA{0x26, 0x32, 0x38, 0xff},
	color.RGBA{0xea, 0x80, 0xfc, 0xff},
	color.RGBA{0x29, 0xb6, 0xf6, 0xff},
	color.RGBA{0x61, 0x61, 0x61, 0xff},
	color.RGBA{0x21, 0x96, 0xf3, 0xff},
	color.RGBA{0x75, 0x75, 0x75, 0xff},
	color.RGBA{0x5d, 0x40, 0x37, 0xff},
	color.RGBA{0xec, 0x40, 0x7a, 0xff},
	color.RGBA{0x02, 0x88, 0xd1, 0xff},
	color.RGBA{0x65, 0x1f, 0xff, 0xff},
	color.RGBA{0xfb, 0x8c, 0x00, 0xff},
	color.RGBA{0x54, 0x6e, 0x7a, 0xff},
	color.RGBA{0x4a, 0x14, 0x8c, 0xff},
	color.RGBA{0x4f, 0xc3, 0xf7, 0xff},
	color.RGBA{0x00, 0xb8, 0xd4, 0xff},
	color.RGBA{0xcd, 0xdc, 0x39, 0xff},
	color.RGBA{0x1a, 0x23, 0x7e, 0xff},
	color.RGBA{0x00, 0xac, 0xc1, 0xff},
	color.RGBA{0x64, 0xff, 0xda, 0xff},
	color.RGBA{0xff, 0xab, 0x40, 0xff},
	color.RGBA{0x18, 0xff, 0xff, 0xff},
	color.RGBA{0x6a, 0x1b, 0x9a, 0xff},
	color.RGBA{0x15, 0x65, 0xc0, 0xff},
	color.RGBA{0xb2, 0xff, 0x59, 0xff},
	color.RGBA{0xd5, 0x00, 0x00, 0xff},
	color.RGBA{0xf8, 0xbb, 0xd0, 0xff},
	color.RGBA{0xf5, 0x7f, 0x17, 0xff},
	color.RGBA{0xff, 0x6d, 0x00, 0xff},
	color.RGBA{0x03, 0x9b, 0xe5, 0xff},
	color.RGBA{0x00, 0xe5, 0xff, 0xff},
	color.RGBA{0x7b, 0x1f, 0xa2, 0xff},
	color.RGBA{0x1b, 0x5e, 0x20, 0xff},
	color.RGBA{0x79, 0x55, 0x48, 0xff},
	color.RGBA{0xff, 0x8a, 0x80, 0xff},
	color.RGBA{0xe9, 0x1e, 0x63, 0xff},
	color.RGBA{0xf4, 0x43, 0x36, 0xff},
	color.RGBA{0xef, 0x6c, 0x00, 0xff},
	color.RGBA{0x00, 0xc8, 0x53, 0xff},
	color.RGBA{0xe6, 0x4a, 0x19, 0xff},
	color.RGBA{0x26, 0xa6, 0x9a, 0xff},
	color.RGBA{0x2e, 0x7d, 0x32, 0xff},
	color.RGBA{0x26, 0xc6, 0xda, 0xff},
	color.RGBA{0x00, 0x60, 0x64, 0xff},
	color.RGBA{0xef, 0x53, 0x50, 0xff},
	color.RGBA{0xd3, 0x2f, 0x2f, 0xff},
	color.RGBA{0xd7, 0xcc, 0xc8, 0xff},
	color.RGBA{0x39, 0x49, 0xab, 0xff},
	color.RGBA{0x00, 0xbf, 0xa5, 0xff},
	color.RGBA{0x9c, 0x27, 0xb0, 0xff},
	color.RGBA{0x5e, 0x35, 0xb1, 0xff},
	color.RGBA{0x90, 0xca, 0xf9, 0xff},
	color.RGBA{0x55, 0x8b, 0x2f, 0xff},
	color.RGBA{0xff, 0x40, 0x81, 0xff},
	color.RGBA{0x9f, 0xa8, 0xda, 0xff},
	color.RGBA{0x5c, 0x6b, 0xc0, 0xff},
	color.RGBA{0x42, 0x42, 0x42, 0xff},
	color.RGBA{0xf4, 0x8f, 0xb1, 0xff},
	color.RGBA{0xaa, 0x00, 0xff, 0xff},
	color.RGBA{0x80, 0xde, 0xea, 0xff},
	color.RGBA{0x00, 0x97, 0xa7, 0xff},
	color.RGBA{0x82, 0x77, 0x17, 0xff},
	color.RGBA{0xf0, 0x62, 0x92, 0xff},
	color.RGBA{0x80, 0xcb, 0xc4, 0xff},
	color.RGBA{0x38, 0x8e, 0x3c, 0xff},
	color.RGBA{0xff, 0x6e, 0x40, 0xff},
	color.RGBA{0xe5, 0x39, 0x35, 0xff},
	color.RGBA{0x64, 0xb5, 0xf6, 0xff},
	color.RGBA{0xff, 0x98, 0x00, 0xff},
	color.RGBA{0xf9, 0xa8, 0x25, 0xff},
	color.RGBA{0x69, 0xf0, 0xae, 0xff},
	color.RGBA{0x00, 0x69, 0x5c, 0xff},
	color.RGBA{0x4c, 0xaf, 0x50, 0xff},
	color.RGBA{0x29, 0x79, 0xff, 0xff},
	color.RGBA{0xb2, 0xdf, 0xdb, 0xff},
	color.RGBA{0xff, 0xca, 0x28, 0xff},
	color.RGBA{0xff, 0x91, 0x00, 0xff},
	color.RGBA{0x31, 0x1b, 0x92, 0xff},
	color.RGBA{0x1e, 0x88, 0xe5, 0xff},
	color.RGBA{0x8b, 0xc3, 0x4a, 0xff},
	color.RGBA{0x7c, 0xb3, 0x42, 0xff},
	color.RGBA{0x6d, 0x4c, 0x41, 0xff},
	color.RGBA{0x21, 0x21, 0x21, 0xff},
	color.RGBA{0x45, 0x27, 0xa0, 0xff},
	color.RGBA{0x00, 0x89, 0x7b, 0xff},
	color.RGBA{0xff, 0x9e, 0x80, 0xff},
	color.RGBA{0x8e, 0x24, 0xaa, 0xff},
	color.RGBA{0x02, 0x77, 0xbd, 0xff},
	color.RGBA{0xd8, 0x1b, 0x60, 0xff},
	color.RGBA{0x00, 0x91, 0xea, 0xff},
	color.RGBA{0xb7, 0x1c, 0x1c, 0xff},
	color.RGBA{0xa5, 0xd6, 0xa7, 0xff},
	color.RGBA{0xd5, 0x00, 0xf9, 0xff},
	color.RGBA{0x03, 0xa9, 0xf4, 0xff},
	color.RGBA{0x88, 0x0e, 0x4f, 0xff},
	color.RGBA{0xbd, 0xbd, 0xbd, 0xff},
	color.RGBA{0xff, 0x17, 0x44, 0xff},
	color.RGBA{0xc2, 0x18, 0x5b, 0xff},
	color.RGBA{0xff, 0x57, 0x22, 0xff},
	color.RGBA{0x3e, 0x27, 0x23, 0xff},
	color.RGBA{0xcf, 0xd8, 0xdc, 0xff},
	color.RGBA{0x7e, 0x57, 0xc2, 0xff},
	color.RGBA{0x40, 0xc4, 0xff, 0xff},
	color.RGBA{0xff, 0xb3, 0x00, 0xff},
	color.RGBA{0xc5, 0x11, 0x62, 0xff},
	color.RGBA{0x00, 0x4d, 0x40, 0xff},
	color.RGBA{0x3f, 0x51, 0xb5, 0xff},
	color.RGBA{0xc5, 0xe1, 0xa5, 0xff},
	color.RGBA{0x67, 0x3a, 0xb7, 0xff},
	color.RGBA{0x4d, 0xd0, 0xe1, 0xff},
	color.RGBA{0x64, 0xdd, 0x17, 0xff},
	color.RGBA{0xe6, 0x51, 0x00, 0xff},
	color.RGBA{0xff, 0x70, 0x43, 0xff},
	color.RGBA{0x7c, 0x4d, 0xff, 0xff},
	color.RGBA{0xc5, 0xca, 0xe9, 0xff},
	color.RGBA{0xff, 0xd6, 0x00, 0xff},
	color.RGBA{0x81, 0xc7, 0x84, 0xff},
	color.RGBA{0x68, 0x9f, 0x38, 0xff},
	color.RGBA{0x53, 0x6d, 0xfe, 0xff},
	color.RGBA{0x44, 0x8a, 0xff, 0xff},
	color.RGBA{0xf5, 0x7c, 0x00, 0xff},
	color.RGBA{0xa1, 0x88, 0x7f, 0xff},
	color.RGBA{0x19, 0x76, 0xd2, 0xff},
	color.RGBA{0xff, 0x8f, 0x00, 0xff},
	color.RGBA{0x8c, 0x9e, 0xff, 0xff},
	color.RGBA{0x42, 0xa5, 0xf5, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0xe6, 0x76, 0xff},
	color.RGBA{0x9c, 0xcc, 0x65, 0xff},
	color.RGBA{0xdd, 0x2c, 0x00, 0xff},
	color.RGBA{0x95, 0x75, 0xcd, 0xff},
	color.RGBA{0x0d, 0x47, 0xa1, 0xff},
	color.RGBA{0x01, 0x57, 0x9b, 0xff},
	color.RGBA{0xdc, 0xe7, 0x75, 0xff},
	color.RGBA{0x00, 0xbc, 0xd4, 0xff},
	color.RGBA{0x4e, 0x34, 0x2e, 0xff},
	color.RGBA{0x82, 0xb1, 0xff, 0xff},
	color.RGBA{0xc6, 0xff, 0x00, 0xff},
	color.RGBA{0xff, 0xcc, 0x80, 0xff},
	color.RGBA{0xd4, 0xe1, 0x57, 0xff},
	color.RGBA{0xff, 0xea, 0x00, 0xff},
	color.RGBA{0x00, 0x79, 0x6b, 0xff},
	color.RGBA{0x33, 0x69, 0x1e, 0xff},
	color.RGBA{0xd8, 0x43, 0x15, 0xff},
	color.RGBA{0xe5, 0x73, 0x73, 0xff},
	color.RGBA{0xba, 0x68, 0xc8, 0xff},
	color.RGBA{0xff, 0xd7, 0x40, 0xff},
	color.RGBA{0xff, 0xd1, 0x80, 0xff},
	color.RGBA{0xf4, 0x51, 0x1e, 0xff},
	color.RGBA{0xb0, 0xbe, 0xc5, 0xff},
	color.RGBA{0x79, 0x86, 0xcb, 0xff},
	color.RGBA{0x29, 0x62, 0xff, 0xff},
	color.RGBA{0xfb, 0xc0, 0x2d, 0xff},
	color.RGBA{0x78, 0x90, 0x9c, 0xff},
	color.RGBA{0x62, 0x00, 0xea, 0xff},
	color.RGBA{0x3d, 0x5a, 0xfe, 0xff},
	color.RGBA{0x80, 0xd8, 0xff, 0xff},
	color.RGBA{0xe0, 0x40, 0xfb, 0xff},
	color.RGBA{0x30, 0x4f, 0xfe, 0xff},
	color.RGBA{0xff, 0xd5, 0x4f, 0xff},
	color.RGBA{0xff, 0xa0, 0x00, 0xff},
	color.RGBA{0xff, 0x52, 0x52, 0xff},
	color.RGBA{0x51, 0x2d, 0xa8, 0xff},
	color.RGBA{0xff, 0xb7, 0x4d, 0xff},
	color.RGBA{0x00, 0x96, 0x88, 0xff},
	color.RGBA{0x1d, 0xe9, 0xb6, 0xff},
	color.RGBA{0x43, 0xa0, 0x47, 0xff},
	color.RGBA{0xff, 0xc1, 0x07, 0xff},
	color.RGBA{0x4d, 0xb6, 0xac, 0xff},
	color.RGBA{0x66, 0xbb, 0x6a, 0xff},
	color.RGBA{0xaf, 0xb4, 0x2b, 0xff},
	color.RGBA{0xff, 0x8a, 0x65, 0xff},
	color.RGBA{0x8d, 0x6e, 0x63, 0xff},
	color.RGBA{0xef, 0x9a, 0x9a, 0xff},
	color.RGBA{0x81, 0xd4, 0xfa, 0xff},
	color.RGBA{0x9e, 0x9e, 0x9e, 0xff},
	color.RGBA{0x28, 0x35, 0x93, 0xff},
	color.RGBA{0xad, 0x14, 0x57, 0xff},
	color.RGBA{0xff, 0x80, 0xab, 0xff},
	color.RGBA{0x9e, 0x9d, 0x24, 0xff},
	color.RGBA{0xab, 0x47, 0xbc, 0xff},
	color.RGBA{0x76, 0xff, 0x03, 0xff},
}