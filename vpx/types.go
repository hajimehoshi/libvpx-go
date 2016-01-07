// THE AUTOGENERATED LICENSE. ALL THE RIGHTS ARE RESERVED BY ROBOTS.

// WARNING: This file has automatically been generated on Thu, 07 Jan 2016 21:04:53 MSK.
// By http://git.io/cgogen. DO NOT EDIT.

package vpx

/*
#cgo pkg-config: vpx
#include <vpx/vpx_encoder.h>
#include <vpx/vpx_decoder.h>
#include <vpx/vp8.h>
#include <stdlib.h>
#include "cgo_helpers.h"
*/
import "C"
import "unsafe"

// FixedBuf as declared in vpx-1.5.0/vpx_encoder.h:110
type FixedBuf struct {
	Buf            unsafe.Pointer
	Sz             uint
	refeac28dc0    *C.vpx_fixed_buf_t
	allocseac28dc0 interface{}
}

// CodecPts type as declared in vpx-1.5.0/vpx_encoder.h:118
type CodecPts int64

// CodecFrameFlags type as declared in vpx-1.5.0/vpx_encoder.h:128
type CodecFrameFlags uint

// CodecErFlags type as declared in vpx-1.5.0/vpx_encoder.h:144
type CodecErFlags uint

// CodecCxPkt as declared in vpx-1.5.0/vpx_encoder.h:223
type CodecCxPkt struct {
	Kind           CodecCxPktKind
	Data           [128]byte
	refa671fc83    *C.vpx_codec_cx_pkt_t
	allocsa671fc83 interface{}
}

// CodecEncOutputCxPktCbFn type as declared in vpx-1.5.0/vpx_encoder.h:233
type CodecEncOutputCxPktCbFn func(pkt []CodecCxPkt, userData unsafe.Pointer)

// CodecPrivOutputCxPktCbPair as declared in vpx-1.5.0/vpx_encoder.h:240
type CodecPrivOutputCxPktCbPair struct {
	OutputCxPkt    CodecEncOutputCxPktCbFn
	UserPriv       unsafe.Pointer
	ref5727a29d    *C.vpx_codec_priv_output_cx_pkt_cb_pair_t
	allocs5727a29d interface{}
}

// Rational as declared in vpx-1.5.0/vpx_encoder.h:249
type Rational struct {
	Num            int
	Den            int
	ref48ce5779    *C.vpx_rational_t
	allocs48ce5779 interface{}
}

// EncFrameFlags type as declared in vpx-1.5.0/vpx_encoder.h:291
type EncFrameFlags int

// CodecEncCfg as declared in vpx-1.5.0/vpx_encoder.h:751
type CodecEncCfg struct {
	GUsage                  uint
	GThreads                uint
	GProfile                uint
	GW                      uint
	GH                      uint
	GBitDepth               BitDepth
	GInputBitDepth          uint
	GTimebase               Rational
	GErrorResilient         CodecErFlags
	GPass                   EncPass
	GLagInFrames            uint
	RcDropframeThresh       uint
	RcResizeAllowed         uint
	RcScaledWidth           uint
	RcScaledHeight          uint
	RcResizeUpThresh        uint
	RcResizeDownThresh      uint
	RcEndUsage              RcMode
	RcTwopassStatsIn        FixedBuf
	RcFirstpassMbStatsIn    FixedBuf
	RcTargetBitrate         uint
	RcMinQuantizer          uint
	RcMaxQuantizer          uint
	RcUndershootPct         uint
	RcOvershootPct          uint
	RcBufSz                 uint
	RcBufInitialSz          uint
	RcBufOptimalSz          uint
	Rc2passVbrBiasPct       uint
	Rc2passVbrMinsectionPct uint
	Rc2passVbrMaxsectionPct uint
	KfMode                  KfMode
	KfMinDist               uint
	KfMaxDist               uint
	SsNumberLayers          uint
	SsEnableAutoAltRef      [5]int
	SsTargetBitrate         [5]uint
	TsNumberLayers          uint
	TsTargetBitrate         [5]uint
	TsRateDecimator         [5]uint
	TsPeriodicity           uint
	TsLayerID               [16]uint
	LayerTargetBitrate      [12]uint
	TemporalLayeringMode    int
	ref37e25db9             *C.vpx_codec_enc_cfg_t
	allocs37e25db9          interface{}
}

// SvcExtraCfg as declared in vpx-1.5.0/vpx_encoder.h:764
type SvcExtraCfg struct {
	MaxQuantizers        [12]int
	MinQuantizers        [12]int
	ScalingFactorNum     [12]int
	ScalingFactorDen     [12]int
	TemporalLayeringMode int
	ref7a0d6872          *C.vpx_svc_extra_cfg_t
	allocs7a0d6872       interface{}
}

// CodecCaps type as declared in vpx-1.5.0/vpx_codec.h:153
type CodecCaps int

// CodecFlags type as declared in vpx-1.5.0/vpx_codec.h:165
type CodecFlags int

// CodecIface as declared in vpx-1.5.0/vpx_codec.h:173
type CodecIface C.vpx_codec_iface_t

// CodecPriv as declared in vpx-1.5.0/vpx_codec.h:181
type CodecPriv C.vpx_codec_priv_t

// CodecIter type as declared in vpx-1.5.0/vpx_codec.h:188
type CodecIter unsafe.Pointer

// CodecCtx as declared in vpx-1.5.0/vpx_codec.h:213
type CodecCtx struct {
	Name           string
	Iface          []CodecIface
	Err            CodecErr
	ErrDetail      string
	InitFlags      CodecFlags
	Config         [8]byte
	Priv           []CodecPriv
	ref8abc1e81    *C.vpx_codec_ctx_t
	allocs8abc1e81 interface{}
}

// Image as declared in vpx-1.5.0/vpx_image.h:133
type Image struct {
	Fmt            ImageFormat
	Cs             ColorSpace
	Range          ColorRange
	W              uint
	H              uint
	BitDepth       uint
	DW             uint
	DH             uint
	RW             uint
	RH             uint
	XChromaShift   uint
	YChromaShift   uint
	Planes         [4][]byte
	Stride         [4]int
	Bps            int
	UserPriv       unsafe.Pointer
	ImgData        []byte
	ImgDataOwner   int
	SelfAllocd     int
	FbPriv         unsafe.Pointer
	refc09455e3    *C.vpx_image_t
	allocsc09455e3 interface{}
}

// ImageRect as declared in vpx-1.5.0/vpx_image.h:141
type ImageRect struct {
	X              uint
	Y              uint
	W              uint
	H              uint
	reff3ce051f    *C.vpx_image_rect_t
	allocsf3ce051f interface{}
}

// CodecStreamInfo as declared in vpx-1.5.0/vpx_decoder.h:93
type CodecStreamInfo struct {
	Sz             uint
	W              uint
	H              uint
	IsKf           uint
	ref342546e4    *C.vpx_codec_stream_info_t
	allocs342546e4 interface{}
}

// CodecDecCfg as declared in vpx-1.5.0/vpx_decoder.h:111
type CodecDecCfg struct {
	Threads        uint
	W              uint
	H              uint
	ref7df355ac    *C.vpx_codec_dec_cfg_t
	allocs7df355ac interface{}
}

// CodecPutFrameCbFn type as declared in vpx-1.5.0/vpx_decoder.h:260
type CodecPutFrameCbFn func(userPriv unsafe.Pointer, img []Image)

// CodecPutSliceCbFn type as declared in vpx-1.5.0/vpx_decoder.h:300
type CodecPutSliceCbFn func(userPriv unsafe.Pointer, img []Image, valid []ImageRect, update []ImageRect)

// CodecFrameBuffer as declared in vpx-1.5.0/vpx_frame_buffer.h:43
type CodecFrameBuffer struct {
	Data           []byte
	Size           uint
	Priv           unsafe.Pointer
	refd319b8f1    *C.vpx_codec_frame_buffer_t
	allocsd319b8f1 interface{}
}

// GetFrameBufferCbFn type as declared in vpx-1.5.0/vpx_frame_buffer.h:63
type GetFrameBufferCbFn func(priv unsafe.Pointer, minSize uint, fb []CodecFrameBuffer) int

// ReleaseFrameBufferCbFn type as declared in vpx-1.5.0/vpx_frame_buffer.h:76
type ReleaseFrameBufferCbFn func(priv unsafe.Pointer, fb []CodecFrameBuffer) int
