package mpa

import (
	"io"

	"github.com/korandiz/mpa"
	"github.com/mjibson/mog/codec"
)

func init() {
	codec.RegisterCodec("MP3", "\xff\xfa", NewSongs)
	codec.RegisterCodec("MP3", "\xff\xfb", NewSongs)
	codec.RegisterCodec("MP3", "\xff\xfc", NewSongs)
	codec.RegisterCodec("MP3", "\xff\xfd", NewSongs)
	codec.RegisterCodec("MP3", "\xff\xfe", NewSongs)
	codec.RegisterCodec("MP3", "\xff\xff", NewSongs)
}

func NewSongs(rf codec.Reader) ([]codec.Song, error) {
	s, err := NewSong(rf)
	if err != nil {
		return nil, err
	}
	return []codec.Song{s}, nil
}

type Song struct {
	Reader  codec.Reader
	r       io.ReadCloser
	decoder *mpa.Decoder
	buff    [2][]float32
}

func NewSong(rf codec.Reader) (*Song, error) {
	s := &Song{Reader: rf}
	_, _, err := s.Init()
	s.Close()
	return s, err
}

func (s *Song) Init() (sampleRate, channels int, err error) {
	r, err := s.Reader()
	if err != nil {
		return 0, 0, err
	}
	s.decoder = &mpa.Decoder{Input: r}
	s.r = r
	if err := s.decoder.DecodeFrame(); err != nil {
		return 0, 0, err
	}
	return s.decoder.SamplingFrequency(), s.decoder.NChannels(), nil
}

func (s *Song) Info() (codec.SongInfo, error) {
	return codec.SongInfo{
		Time: 0, // too hard to tell without decoding
	}, nil
}

func (s *Song) Play(n int) (r []float32, err error) {
	for len(r) < n {
		if len(s.buff[0]) == 0 {
			if err = s.decoder.DecodeFrame(); err != nil {
				return
			}
			for i := 0; i < 2; i++ {
				s.buff[i] = make([]float32, s.decoder.NSamples())
				s.decoder.ReadSamples(i, s.buff[i])
			}
		}
		for len(s.buff[0]) > 0 && len(r) < n {
			r = append(r, s.buff[0][0], s.buff[1][0])
			s.buff[0], s.buff[1] = s.buff[0][1:], s.buff[1][1:]
		}
	}
	return
}

func (s *Song) Close() {
	s.r.Close()
	s.decoder, s.buff[0], s.buff[1], s.r = nil, nil, nil, nil
}
