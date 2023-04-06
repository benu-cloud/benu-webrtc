package stream

/*
#cgo pkg-config: gstreamer-1.0 gstreamer-webrtc-1.0 gstreamer-sdp-1.0 gstreamer-rtp-1.0
#cgo CFLAGS: -I${SRCDIR}/c
#cgo LDFLAGS: -L${SRCDIR}/c -lstream
#include "stream.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"sync"

	"github.com/benu-cloud/benu-livestreaming-gst/internal/config"
	"github.com/benu-cloud/benu-livestreaming-gst/internal/message"
	"github.com/benu-cloud/benu-livestreaming-gst/pkg/pkgerrors"
)

type peer struct {
	peer_id                   string
	serverSessionDescriptions chan *message.SessionDescriptionPayload
	serverIceCandidates       chan *message.IceCandidatePayload
	serverDatachannelMessages chan *message.GenericPayload
}

type stream struct {
	mutex                 sync.Mutex
	users                 []*peer
	serverGStreamerErrors chan error
}

var instance *stream = nil

// ------------------

func checkStreamInstance() error {
	if instance == nil {
		return pkgerrors.NewStreamError(errors.New("no stream instance exists, use the pipeline setup function before doing anything else"))
	}
	return nil
}

// ------------------

func SetupPipeline(settings *config.StreamSettings) (<-chan error, error) {
	if checkStreamInstance() == nil {
		return nil, pkgerrors.NewStreamError(errors.New("pipeline setup function should be run only once"))
	}
	options := C.PipelineOptions{
		audioBaseBitrate:       (C.uint)(settings.AudioBaseBitrate),
		audioBasePacketLossPct: (C.uint)(settings.AudioBasePacketLossPct),
		videoBaseBitrate:       (C.uint)(settings.VideoBaseBitrate),
		videoBaseFramerate:     (C.uint)(settings.VideoBaseFramerate),
		videoEncoder:           (C.VideoEncoder)(settings.VideoEncoder),
		videoHeight:            (C.uint)(settings.VideoResolution.Height),
		videoWidth:             (C.uint)(settings.VideoResolution.Width),
		videoShowCursor:        (C.bool)(settings.VideoShowCursor),
	}
	result := C.SetupPipeline(options)
	if result != C.SUCCESS {
		return nil, pkgerrors.NewCStreamError(int(result))
	}
	instance = &stream{
		users:                 make([]*peer, 0),
		serverGStreamerErrors: make(chan error),
	}
	return instance.serverGStreamerErrors, nil
}

func StartPipeline() error {
	if err := checkStreamInstance(); err != nil {
		return err
	}
	result := C.StartPipeline()
	if result != C.SUCCESS {
		return pkgerrors.NewCStreamError(int(result))
	}
	return nil
}

func StopPipeline() error {
	if err := checkStreamInstance(); err != nil {
		return err
	}
	result := C.StopPipeline()
	if result != C.SUCCESS {
		return pkgerrors.NewCStreamError(int(result))
	}
	return nil
}

func AddPeerToPipeline(peerId string) (<-chan *message.SessionDescriptionPayload, <-chan *message.IceCandidatePayload, error) {
	if err := checkStreamInstance(); err != nil {
		return nil, nil, err
	}
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	// check if peer already exists
	for _, p := range instance.users {
		if p.peer_id == peerId {
			return nil, nil, pkgerrors.NewStreamError(errors.New("peer already exists"))
		}
	}
	peer := &peer{
		peer_id:                   peerId,
		serverSessionDescriptions: make(chan *message.SessionDescriptionPayload),
		serverIceCandidates:       make(chan *message.IceCandidatePayload),
		serverDatachannelMessages: make(chan *message.GenericPayload),
	}
	result := C.AddPeerToPipeline(C.CString(peerId))
	if result != C.SUCCESS {
		return nil, nil, pkgerrors.NewCStreamError(int(result))
	}
	instance.users = append(instance.users, peer)
	return peer.serverSessionDescriptions, peer.serverIceCandidates, nil
}

func RemovePeerFromPipeline(peerId string) error {
	if err := checkStreamInstance(); err != nil {
		return err
	}
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	// check if peer exists
	index := -1
	for i, p := range instance.users {
		if p.peer_id == peerId {
			index = i
			break
		}
	}
	if index == -1 {
		return pkgerrors.NewStreamError(fmt.Errorf("no peer to remove with id '%s'", peerId))
	}
	result := C.RemovePeerFromPipeline(C.CString(peerId))
	if result != C.SUCCESS {
		return pkgerrors.NewCStreamError(int(result))
	}
	// close channels
	close(instance.users[index].serverIceCandidates)
	close(instance.users[index].serverSessionDescriptions)
	// remove
	instance.users = append(instance.users[:index], instance.users[index+1:]...)
	return nil
}

func SetRemoteAnswer(peerId string, answer string) error {
	if err := checkStreamInstance(); err != nil {
		return err
	}
	result := C.SetRemoteAnswer(C.CString(peerId), C.CString(answer))
	if result != C.SUCCESS {
		return pkgerrors.NewCStreamError(int(result))
	}
	return nil
}

func AddRemoteIceCandidate(peerId string, mlineindex uint, candidate string) error {
	if err := checkStreamInstance(); err != nil {
		return err
	}
	result := C.AddRemoteIceCandidate(C.CString(peerId), C.uint(mlineindex), C.CString(candidate))
	if result != C.SUCCESS {
		return pkgerrors.NewCStreamError(int(result))
	}
	return nil
}
