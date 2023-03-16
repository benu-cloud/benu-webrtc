package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/benu-cloud/benu-livestreaming-gst/pkg/pkgerrors"
)

func (r *Resolution) String() string {
	return fmt.Sprintf("%dx%d", r.Width, r.Height)
}

func (r *Resolution) Set(s string) error {
	splitRes := strings.SplitN(s, "x", 2)
	var err error
	if len(splitRes) != 2 {
		goto badFormat
	}

	r.Width, err = strconv.Atoi(splitRes[0])
	if err != nil {
		goto badFormat
	}
	r.Height, err = strconv.Atoi(splitRes[1])
	if err != nil {
		goto badFormat
	}
	if r.Height < 0 || r.Width < 0 {
		goto badFormat
	}
	return nil
badFormat:
	return pkgerrors.NewBadCommanlineArgument("Resolution", s, "[WIDTH]x[HEIGHT]")
}

func (e *VideoEncoder) String() string {
	return string(*e)
}

func (e *VideoEncoder) Set(s string) error {
	switch s {
	case "VP9":
		*e = VP9
	case "H264":
		*e = H264
	case "NVH264":
		*e = NVH264
	default:
		return pkgerrors.NewBadCommanlineArgument("VideoEncoder", s, "(VP9 / H264 / NVH264)")
	}
	return nil
}

func (p *PortNumber) String() string {
	return fmt.Sprintf("%d", uint(*p))
}

func (p *PortNumber) Set(s string) error {
	portnum, err := strconv.Atoi(s)
	if err != nil {
		goto badFormat
	}
	if portnum < 0 || portnum > 65535 {
		goto badFormat
	}
	return nil
badFormat:
	return pkgerrors.NewBadCommanlineArgument("Port", s, "0-65535")
}
